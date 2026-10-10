// Package gateway terminates WebSocket connections: authentication, character
// management and forwarding of in-world messages to the game world.
package gateway

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/coder/websocket"

	"mmorpg/internal/auth"
	"mmorpg/internal/game"
	"mmorpg/internal/gamedata"
	"mmorpg/internal/protocol"
	"mmorpg/internal/storage"
)

var nameRe = regexp.MustCompile(`^[A-Za-zÄÖÜäöüß]{2,16}$`)

type Gateway struct {
	store storage.Store
	auth  *auth.Service
	world *game.World
}

func New(store storage.Store, a *auth.Service, w *game.World) *Gateway {
	return &Gateway{store: store, auth: a, world: w}
}

type conn struct {
	ws     *websocket.Conn
	out    chan []byte
	cancel context.CancelFunc
	acct   *storage.Account
	entID  int
	closed atomic.Bool
}

// Send implements game.Sender. It never blocks the simulation: slow clients are dropped.
func (c *conn) Send(v any) {
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	select {
	case c.out <- b:
	default:
		if c.closed.CompareAndSwap(false, true) {
			log.Printf("dropping slow client")
			c.cancel()
		}
	}
}

func (g *Gateway) ServeWS(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return
	}
	ws.SetReadLimit(1 << 16)
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	c := &conn{ws: ws, out: make(chan []byte, 256), cancel: cancel}

	go func() { // writer
		for {
			select {
			case b := <-c.out:
				wctx, wcancel := context.WithTimeout(ctx, 5*time.Second)
				err := ws.Write(wctx, websocket.MessageText, b)
				wcancel()
				if err != nil {
					cancel()
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	defer func() {
		if c.entID != 0 {
			g.world.Leave(c.entID)
		}
		_ = ws.CloseNow()
	}()

	for {
		_, data, err := ws.Read(ctx)
		if err != nil {
			return
		}
		var env protocol.Envelope
		if json.Unmarshal(data, &env) != nil {
			continue
		}
		g.handle(c, env.T, data)
	}
}

func (g *Gateway) fail(c *conn, t string, err error) {
	c.Send(protocol.M{"t": t, "error": err.Error()})
}

func (g *Gateway) handle(c *conn, t string, data []byte) {
	switch t {
	case protocol.CRegister, protocol.CLogin:
		if c.acct != nil {
			return
		}
		var cr protocol.Credentials
		if json.Unmarshal(data, &cr) != nil {
			return
		}
		cr.User = strings.TrimSpace(cr.User)
		var acct *storage.Account
		var err error
		if t == protocol.CRegister {
			acct, err = g.auth.Register(cr.User, cr.Pass)
		} else {
			acct, err = g.auth.Login(cr.User, cr.Pass)
		}
		if err != nil {
			g.fail(c, "authError", err)
			return
		}
		c.acct = acct
		c.Send(protocol.M{
			"t": "auth", "user": acct.User,
			"classes": gamedata.Classes, "classOrder": gamedata.ClassOrder,
			"races": gamedata.Races, "raceOrder": gamedata.RaceOrder,
			"skills": gamedata.Skills,
		})
		g.sendChars(c)
	case protocol.CCreate:
		if c.acct == nil || c.entID != 0 {
			return
		}
		g.createChar(c, data)
	case protocol.CDelete:
		if c.acct == nil || c.entID != 0 {
			return
		}
		var m protocol.CharID
		if json.Unmarshal(data, &m) == nil {
			_ = g.store.DeleteCharacter(c.acct.ID, m.ID)
		}
		g.sendChars(c)
	case protocol.CEnter:
		if c.acct == nil || c.entID != 0 {
			return
		}
		var m protocol.CharID
		if json.Unmarshal(data, &m) != nil {
			return
		}
		ch, ok := g.store.CharacterByID(m.ID)
		if !ok || ch.AccountID != c.acct.ID {
			g.fail(c, "enterError", storage.ErrNotFound)
			return
		}
		e := g.world.Enter(ch, c)
		c.entID = e.ID
	case protocol.CLogout:
		if c.entID != 0 {
			g.world.Leave(c.entID)
			c.entID = 0
			g.sendChars(c)
		}
	default:
		if c.entID != 0 {
			g.world.Dispatch(c.entID, t, data)
		}
	}
}

func (g *Gateway) sendChars(c *conn) {
	list := []protocol.M{}
	for _, ch := range g.store.CharactersOf(c.acct.ID) {
		list = append(list, protocol.M{"id": ch.ID, "name": ch.Name, "race": ch.Race, "class": ch.Class, "level": ch.Level, "look": ch.Look})
	}
	c.Send(protocol.M{"t": "chars", "chars": list})
}

func (g *Gateway) createChar(c *conn, data []byte) {
	var m protocol.CreateChar
	if json.Unmarshal(data, &m) != nil {
		return
	}
	m.Name = strings.TrimSpace(m.Name)
	if !nameRe.MatchString(m.Name) {
		g.fail(c, "createError", errString("Name: 2-16 Buchstaben"))
		return
	}
	m.Name = strings.ToUpper(m.Name[:1]) + strings.ToLower(m.Name[1:])
	if gamedata.Races[m.Race] == nil {
		g.fail(c, "createError", errString("Unbekannte Rasse"))
		return
	}
	if gamedata.Classes[m.Class] == nil {
		g.fail(c, "createError", errString("Unbekannte Klasse"))
		return
	}
	if m.Look.Gender != "m" && m.Look.Gender != "f" {
		m.Look.Gender = "m"
	}
	m.Look.Hair = max(0, min(5, m.Look.Hair))
	m.Look.Skin = max(0, min(4, m.Look.Skin))
	ch := &storage.Character{
		AccountID: c.acct.ID, Name: m.Name, Race: m.Race, Class: m.Class, Look: m.Look,
		Level: 1, Quests: map[string]*storage.QuestState{},
	}
	if err := g.store.CreateCharacter(ch); err != nil {
		g.fail(c, "createError", err)
		return
	}
	c.Send(protocol.M{"t": "created", "id": ch.ID})
	g.sendChars(c)
}

type errString string

func (e errString) Error() string { return string(e) }

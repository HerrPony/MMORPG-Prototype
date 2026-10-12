package game

import (
	"math"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/loot"
	"mmorpg/internal/storage"
)

// Sender delivers messages to one client. Implemented by the gateway.
type Sender interface {
	Send(v any)
}

const (
	KindPlayer = "player"
	KindNPC    = "npc"
	KindMob    = "mob"
	KindPet    = "pet"
	KindObj    = "obj"
	KindShip   = "ship"

	MaxInventory = 28
	InteractDist = 8.0
	MaxMoveSpeed = 8.5
	InterestDist = 130.0 // entities further away are not sent to a client
	MaxParty     = 5
)

type Dot struct {
	Src   int
	Per   float64
	Until float64
	Next  float64
}

// Entity is anything in the world: players, NPCs, mobs, pets and objects.
type Entity struct {
	ID    int
	Kind  string
	Name  string
	Title string
	Model string
	Level int
	X, Y  float64
	Z, R  float64
	HP    int
	MaxHP int
	Dead  bool
	Inst  int     // dungeon instance id, 0 = open world
	Stuck float64 // seconds a mover has been blocked by an obstacle
	Nav   navState

	// temporary effects
	Dots        []Dot
	SlowMult    float64
	SlowUntil   float64
	BuffMult    float64
	BuffUntil   float64
	BuffTaken   float64 // extra damage taken while the buff lasts
	Shield      float64
	ShieldUntil float64
	StunUntil   float64
	PolyUntil   float64 // polymorph: ends when damaged
	HPBoostEnd  float64 // "Letztes Gefecht"

	P    *Player
	M    *Mob
	N    *gamedata.NPC
	O    *Obj
	Pet  *Pet
	Ship *Ferry
}

type CastState struct {
	Skill  string
	End    float64
	Target int
	SX, SZ float64
}

type Player struct {
	Char       *storage.Character
	Out        Sender
	Class      *gamedata.Class
	Res        float64
	ResMax     int
	Stats      Stats
	PotCD      float64
	Target     int
	Auto       bool
	CD         map[string]float64
	GCD        float64
	Cast       *CastState
	LastCombat float64
	LastMove   float64
	HPAcc      float64
	Marks      map[int]string
	Stance     string // warriors: battle | defensive

	PartyID        int
	Invite         *Invite
	Flight         *Flight
	ZoneKey        string
	PoiID          string
	PetEnt         int // entity id of the summoned pet (0 = none)
	ElemEnt        int // entity id of the water elemental
	LastQuestCheck float64
	Aboard         int     // entity id of the ferry the player rides on (0 = none)
	DeckX, DeckZ   float64 // position on the deck in ship coordinates
}

type Invite struct {
	From   int
	Expire float64
}

// Flight is a griffin ride.
type Flight struct {
	Dest   string
	SX, SZ float64
	DX, DZ float64
	Dist   float64
	Done   float64
	Speed  float64
}

type Mob struct {
	Def        *gamedata.Mob
	SX, SZ     float64
	Target     int
	Threat     map[int]float64
	TauntUntil float64
	TauntBy    int
	Returning  bool
	NextAtk    float64
	WX, WZ     float64
	WanderWait float64
	Wandering  bool
	CorpseEnd  float64
	RespawnAt  float64
	Gone       bool
	Loot       map[int]*loot.Bag // personal loot by player entity id
	Hit        map[int]bool      // players that damaged the mob (they get credit)
	AbNext     []float64
	AbDone     []bool
	Enraged    bool
	Adds       []int // summoned helpers of a boss
}

// Obj is an interactive world object.
type Obj struct {
	Def       *gamedata.ObjectDef
	Gone      bool
	RespawnAt float64
}

// Pet is a hunter companion or a summoned elemental.
type Pet struct {
	Owner    int
	Kind     string // beast | elemental
	MobID    string
	ExpireAt float64 // 0 = permanent
	NextAtk  float64
	Target   int
	HomeX    float64
}

func (e *Entity) speedMult(now float64) float64 {
	if now < e.SlowUntil && e.SlowMult > 0 {
		return math.Max(0, 1-e.SlowMult)
	}
	return 1
}

func (e *Entity) dmgMult(now float64) float64 {
	if now < e.BuffUntil && e.BuffMult > 0 {
		return e.BuffMult
	}
	return 1
}

// takenMult is the multiplier on damage the entity receives.
func (e *Entity) takenMult(now float64) float64 {
	m := 1.0
	if now < e.BuffUntil && e.BuffTaken > 0 {
		m += e.BuffTaken
	}
	if e.P != nil && e.P.Stance == "defensive" {
		m *= 0.6
	}
	return m
}

func (e *Entity) stunned(now float64) bool { return now < e.StunUntil }

func dist(a, b *Entity) float64 { return math.Hypot(a.X-b.X, a.Z-b.Z) }

func clampf(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }

func round2(v float64) float64 { return math.Round(v*100) / 100 }

// navState is a walker's current route around obstacles (see World.walkTo).
type navState struct {
	Path   [][2]float64 // waypoints still to visit
	GX, GZ float64      // the target the path was computed for
	At     float64      // when it was computed
	Fail   float64      // no new search before this time after a failed one
}

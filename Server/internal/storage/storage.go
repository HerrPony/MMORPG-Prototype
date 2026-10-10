// Package storage defines persistent models and the Store interface.
// The JSON file implementation can be swapped for a SQL store later.
package storage

import (
	"errors"
	"time"
)

var (
	ErrUserExists = errors.New("Benutzername bereits vergeben")
	ErrNameTaken  = errors.New("Charaktername bereits vergeben")
	ErrTooMany    = errors.New("Maximal 5 Charaktere pro Account")
	ErrNotFound   = errors.New("nicht gefunden")
)

const MaxCharsPerAccount = 5

type Account struct {
	ID      string    `json:"id"`
	User    string    `json:"user"`
	Hash    []byte    `json:"hash"`
	Created time.Time `json:"created"`
}

type Look struct {
	Gender string `json:"gender"` // "m" | "f"
	Hair   int    `json:"hair"`
	Skin   int    `json:"skin"`
}

type ItemStack struct {
	Item  string `json:"i"`
	Count int    `json:"c"`
}

type QuestState struct {
	Status  string `json:"status"` // "active" | "complete" | "done"
	Counts  []int  `json:"counts"`
	Expires int64  `json:"expires,omitempty"` // unix seconds, timed quests
}

// PetData is the hunter's tamed companion.
type PetData struct {
	Mob  string `json:"mob"`
	Name string `json:"name"`
	Dead bool   `json:"dead"`
}

type Character struct {
	ID           string                       `json:"id"`
	AccountID    string                       `json:"accountId"`
	Name         string                       `json:"name"`
	Race         string                       `json:"race"`
	Class        string                       `json:"class"`
	Look         Look                         `json:"look"`
	Level        int                          `json:"level"`
	XP           int                          `json:"xp"`
	X            float64                      `json:"x"`
	Y            float64                      `json:"y"`
	Z            float64                      `json:"z"`
	R            float64                      `json:"r"`
	Copper       int                          `json:"copper"`
	Skills       []string                     `json:"skills"` // learned skills in learning order
	Equipment    map[string]string            `json:"equipment"`
	Sets         map[string]map[string]string `json:"sets"` // saved equipment sets
	Inventory    []ItemStack                  `json:"inventory"`
	Quests       map[string]*QuestState       `json:"quests"`
	Discovered   []string                     `json:"discovered"`   // visited zones and places
	FlightPoints []string                     `json:"flightPoints"` // unlocked griffin landing points
	Pet          *PetData                     `json:"pet,omitempty"`
	ReturnX      float64                      `json:"returnX"` // open world position while inside a dungeon
	ReturnZ      float64                      `json:"returnZ"`
	Created      time.Time                    `json:"created"`
}

// Clone returns a deep copy.
func (c *Character) Clone() *Character {
	n := *c
	n.Inventory = append([]ItemStack(nil), c.Inventory...)
	n.Skills = append([]string(nil), c.Skills...)
	n.Discovered = append([]string(nil), c.Discovered...)
	n.FlightPoints = append([]string(nil), c.FlightPoints...)
	if c.Pet != nil {
		p := *c.Pet
		n.Pet = &p
	}
	n.Equipment = make(map[string]string, len(c.Equipment))
	for k, v := range c.Equipment {
		n.Equipment[k] = v
	}
	n.Sets = make(map[string]map[string]string, len(c.Sets))
	for k, set := range c.Sets {
		cp := make(map[string]string, len(set))
		for s, i := range set {
			cp[s] = i
		}
		n.Sets[k] = cp
	}
	n.Quests = make(map[string]*QuestState, len(c.Quests))
	for k, v := range c.Quests {
		q := *v
		q.Counts = append([]int(nil), v.Counts...)
		n.Quests[k] = &q
	}
	return &n
}

type Store interface {
	CreateAccount(user string, hash []byte) (*Account, error)
	AccountByUser(user string) (*Account, bool)
	CharactersOf(accountID string) []*Character
	CharacterByID(id string) (*Character, bool)
	CreateCharacter(c *Character) error
	SaveCharacter(c *Character) error
	DeleteCharacter(accountID, id string) error
	Close() error
}

package storage

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// JSONStore keeps everything in memory and flushes to JSON files.
type JSONStore struct {
	mu       sync.Mutex
	dir      string
	accounts map[string]*Account   // by id
	chars    map[string]*Character // by id
	dirty    bool
	stop     chan struct{}
	done     chan struct{}
}

func NewID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func OpenJSON(dir string) (*JSONStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	s := &JSONStore{
		dir:      dir,
		accounts: map[string]*Account{},
		chars:    map[string]*Character{},
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}
	if err := s.load("accounts.json", &s.accounts); err != nil {
		return nil, err
	}
	if err := s.load("characters.json", &s.chars); err != nil {
		return nil, err
	}
	go s.flushLoop()
	return s, nil
}

func (s *JSONStore) load(name string, v any) error {
	b, err := os.ReadFile(filepath.Join(s.dir, name))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func (s *JSONStore) write(name string, v any) error {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}
	tmp := filepath.Join(s.dir, name+".tmp")
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(s.dir, name))
}

func (s *JSONStore) flushLocked() error {
	if !s.dirty {
		return nil
	}
	if err := s.write("accounts.json", s.accounts); err != nil {
		return err
	}
	if err := s.write("characters.json", s.chars); err != nil {
		return err
	}
	s.dirty = false
	return nil
}

func (s *JSONStore) flushLoop() {
	defer close(s.done)
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			s.mu.Lock()
			_ = s.flushLocked()
			s.mu.Unlock()
		case <-s.stop:
			return
		}
	}
}

func (s *JSONStore) Close() error {
	close(s.stop)
	<-s.done
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.flushLocked()
}

func (s *JSONStore) CreateAccount(user string, hash []byte) (*Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, a := range s.accounts {
		if strings.EqualFold(a.User, user) {
			return nil, ErrUserExists
		}
	}
	a := &Account{ID: NewID(), User: user, Hash: hash, Created: time.Now()}
	s.accounts[a.ID] = a
	s.dirty = true
	return a, nil
}

func (s *JSONStore) AccountByUser(user string) (*Account, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, a := range s.accounts {
		if strings.EqualFold(a.User, user) {
			c := *a
			return &c, true
		}
	}
	return nil, false
}

func (s *JSONStore) CharactersOf(accountID string) []*Character {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []*Character
	for _, c := range s.chars {
		if c.AccountID == accountID {
			out = append(out, c.Clone())
		}
	}
	// stable order: oldest first
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].Created.Before(out[j-1].Created); j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

func (s *JSONStore) CharacterByID(id string) (*Character, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.chars[id]
	if !ok {
		return nil, false
	}
	return c.Clone(), true
}

func (s *JSONStore) CreateCharacter(c *Character) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, o := range s.chars {
		if strings.EqualFold(o.Name, c.Name) {
			return ErrNameTaken
		}
		if o.AccountID == c.AccountID {
			n++
		}
	}
	if n >= MaxCharsPerAccount {
		return ErrTooMany
	}
	if c.ID == "" {
		c.ID = NewID()
	}
	c.Created = time.Now()
	s.chars[c.ID] = c.Clone()
	s.dirty = true
	return nil
}

func (s *JSONStore) SaveCharacter(c *Character) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.chars[c.ID]; !ok {
		return ErrNotFound
	}
	s.chars[c.ID] = c.Clone()
	s.dirty = true
	return nil
}

func (s *JSONStore) DeleteCharacter(accountID, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.chars[id]
	if !ok || c.AccountID != accountID {
		return ErrNotFound
	}
	delete(s.chars, id)
	s.dirty = true
	return nil
}

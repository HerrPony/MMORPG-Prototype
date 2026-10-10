// Package auth handles account registration and login.
package auth

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"

	"mmorpg/internal/storage"
)

var (
	ErrBadUser  = errors.New("Benutzername: 3-20 Zeichen (Buchstaben, Zahlen, _ -)")
	ErrBadPass  = errors.New("Passwort: mindestens 4 Zeichen")
	ErrBadLogin = errors.New("Benutzername oder Passwort falsch")
	userRe      = regexp.MustCompile(`^[A-Za-z0-9_-]{3,20}$`)
)

type Service struct{ store storage.Store }

func New(s storage.Store) *Service { return &Service{store: s} }

func (s *Service) Register(user, pass string) (*storage.Account, error) {
	if !userRe.MatchString(user) {
		return nil, ErrBadUser
	}
	if len(pass) < 4 || len(pass) > 72 {
		return nil, ErrBadPass
	}
	h, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return s.store.CreateAccount(user, h)
}

func (s *Service) Login(user, pass string) (*storage.Account, error) {
	a, ok := s.store.AccountByUser(user)
	if !ok || bcrypt.CompareHashAndPassword(a.Hash, []byte(pass)) != nil {
		return nil, ErrBadLogin
	}
	return a, nil
}

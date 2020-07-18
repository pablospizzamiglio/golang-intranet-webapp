package main

import "errors"

// MemStore is a in-memory session storage
type MemStore struct {
	sessions map[string]Session
}

// NewMemoryStore creates a new memory storage
func NewMemoryStore() Store {
	return &MemStore{
		sessions: make(map[string]Session),
	}
}

// Get the session from the store
func (s MemStore) Get(id string) (Session, error) {
	session, ok := s.sessions[id]
	if !ok {
		return Session{}, errors.New("Session not found")
	}
	return session, nil
}

// Set the session in the store
func (s *MemStore) Set(id string, session Session) error {
	s.sessions[id] = session
	return nil
}

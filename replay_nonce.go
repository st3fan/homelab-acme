package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

type ReplayNonceService interface {
	New() (string, error)
	Use(nonce string) error
	Cleanup()
}

type nonceEntry struct {
	created time.Time
	used    bool
}

type InMemoryReplayNonceService struct {
	sync.Mutex
	nonces   map[string]*nonceEntry
	ttl      time.Duration
	cleanup  time.Duration
	stopChan chan struct{}
}

func NewInMemoryReplayNonceService(ttl, cleanupInterval time.Duration) *InMemoryReplayNonceService {
	s := &InMemoryReplayNonceService{
		nonces:   make(map[string]*nonceEntry),
		ttl:      ttl,
		cleanup:  cleanupInterval,
		stopChan: make(chan struct{}),
	}
	go s.runCleanup()
	return s
}

func (s *InMemoryReplayNonceService) New() (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	encoded := base64.RawURLEncoding.EncodeToString(nonce)

	s.Lock()
	{
		s.nonces[encoded] = &nonceEntry{created: time.Now()}
	}
	s.Unlock()

	return encoded, nil
}

func (s *InMemoryReplayNonceService) Use(nonce string) error {
	s.Lock()
	defer s.Unlock()

	entry, ok := s.nonces[nonce]
	if !ok {
		return errors.New("badNonce")
	}

	if time.Since(entry.created) > s.ttl {
		delete(s.nonces, nonce)
		return errors.New("badNonce")
	}

	if entry.used {
		return errors.New("badNonce")
	}

	entry.used = true

	return nil
}

func (s *InMemoryReplayNonceService) Cleanup() {
	s.Lock()
	s.Unlock()

	now := time.Now()

	for k, v := range s.nonces {
		if v.used || now.Sub(v.created) > s.ttl {
			delete(s.nonces, k)
		}
	}
}

func (s *InMemoryReplayNonceService) runCleanup() {
	ticker := time.NewTicker(s.cleanup)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.Cleanup()
		case <-s.stopChan:
			return
		}
	}
}

func (s *InMemoryReplayNonceService) Stop() {
	close(s.stopChan)
}

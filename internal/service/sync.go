package service

import (
	"log"
	"sync"
	"time"

	"github.com/pauloborszcz/tics/internal/config"
	"github.com/pauloborszcz/tics/internal/glpi"
)

type SyncService struct {
	client        *glpi.Client
	cfg           *config.Config
	mu            sync.Mutex
	tickets       []glpi.Ticket
	knownIDs      map[int]bool
	listeners     []func([]glpi.Ticket)
	onNew         []func([]glpi.Ticket)
	stopCh        chan struct{}
	autoFollowup  bool
}

func NewSyncService(client *glpi.Client, cfg *config.Config) *SyncService {
	return &SyncService{
		client:   client,
		cfg:      cfg,
		knownIDs: make(map[int]bool),
		stopCh:   make(chan struct{}),
	}
}

// OnUpdate registers a callback for when the ticket list is refreshed.
func (s *SyncService) OnUpdate(fn func([]glpi.Ticket)) {
	s.listeners = append(s.listeners, fn)
}

// OnNewTickets registers a callback for when new tickets are detected.
func (s *SyncService) OnNewTickets(fn func([]glpi.Ticket)) {
	s.onNew = append(s.onNew, fn)
}

// Tickets returns the current cached list of tickets.
func (s *SyncService) Tickets() []glpi.Ticket {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]glpi.Ticket, len(s.tickets))
	copy(result, s.tickets)
	return result
}

// Refresh does a single sync, fetching tickets and notifying listeners.
func (s *SyncService) Refresh() {
	log.Println("sync: fetching tickets...")
	tickets, err := s.client.SearchMyTickets()
	if err != nil {
		log.Printf("sync: error fetching tickets: %v", err)
		return
	}
	log.Printf("sync: found %d tickets", len(tickets))

	var newTickets []glpi.Ticket
	s.mu.Lock()
	for _, t := range tickets {
		if !s.knownIDs[t.ID] {
			newTickets = append(newTickets, t)
		}
	}
	// Update known IDs
	s.knownIDs = make(map[int]bool, len(tickets))
	for _, t := range tickets {
		s.knownIDs[t.ID] = true
	}
	s.tickets = tickets
	s.mu.Unlock()

	for _, fn := range s.listeners {
		fn(tickets)
	}
	if len(newTickets) > 0 {
		for _, fn := range s.onNew {
			fn(newTickets)
		}
	}

	// Auto-followup: send "ATENDIMENTO ABERTO" to processing tickets (only if enabled)
	s.mu.Lock()
	auto := s.autoFollowup
	s.mu.Unlock()
	if auto {
		go AutoFollowup(s.client, s.cfg, tickets)
	}
}

// Start begins the periodic sync loop.
func (s *SyncService) Start() {
	// Do an initial sync immediately
	s.Refresh()

	go func() {
		ticker := time.NewTicker(s.cfg.SyncInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.Refresh()
			case <-s.stopCh:
				return
			}
		}
	}()
}

// Stop ends the sync loop.
func (s *SyncService) Stop() {
	close(s.stopCh)
}

// SetAutoFollowup enables or disables automatic followup sending.
func (s *SyncService) SetAutoFollowup(enabled bool) {
	s.mu.Lock()
	s.autoFollowup = enabled
	s.mu.Unlock()
}

// AutoFollowupEnabled returns whether auto-followup is enabled.
func (s *SyncService) AutoFollowupEnabled() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.autoFollowup
}

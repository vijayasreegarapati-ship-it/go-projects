package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Desk represents a physical workspace.
type Desk struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Floor           int    `json:"floor"`
	IsAvailable     bool   `json:"is_available"`
	CurrentOccupant string `json:"current_occupant,omitempty"`
}

type BookingRequest struct {
	DeskID       string `json:"desk_id"`
	EmployeeName string `json:"employee_name"`
}

// InMemoryStore handles thread-safe desk operations.
type InMemoryStore struct {
	mu    sync.RWMutex
	desks map[string]*Desk
}

func NewInMemoryStore() *InMemoryStore {
	store := &InMemoryStore{
		desks: make(map[string]*Desk),
	}

	// mock data
	store.desks["desk-1"] = &Desk{ID: "desk-1", Name: "Window Seat A", Floor: 4, IsAvailable: true}
	store.desks["desk-2"] = &Desk{ID: "desk-2", Name: "Quiet Zone B", Floor: 4, IsAvailable: true}
	store.desks["desk-3"] = &Desk{ID: "desk-3", Name: "Standing Desk C", Floor: 5, IsAvailable: false, CurrentOccupant: "Ada Lovelace"}

	return store
}

func (s *InMemoryStore) GetAvailableDesks() []*Desk {
	s.mu.RLock()
	defer s.mu.RUnlock()

	available := make([]*Desk, 0, len(s.desks))
	for _, desk := range s.desks {
		if desk.IsAvailable {
			available = append(available, desk)
		}
	}
	return available
}

func (s *InMemoryStore) BookDesk(deskID, employeeName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	desk, exists := s.desks[deskID]
	if !exists {
		return errors.New("desk not found")
	}

	if !desk.IsAvailable {
		return errors.New("desk already booked")
	}

	desk.IsAvailable = false
	desk.CurrentOccupant = employeeName
	return nil
}

// WorkplaceServer holds dependencies for HTTP handlers.
type WorkplaceServer struct {
	store *InMemoryStore
}

// loggingMiddleware intercepts requests to log method, path, and duration.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// pass control back to the actual route handler
		next.ServeHTTP(w, r)

		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// writeJSON is a utility for sending JSON responses.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// handlers
func (s *WorkplaceServer) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func (s *WorkplaceServer) handleGetDesks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.GetAvailableDesks())
}

func (s *WorkplaceServer) handleBookDesk(w http.ResponseWriter, r *http.Request) {
	var req BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	if req.DeskID == "" || req.EmployeeName == "" {
		sendError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	if err := s.store.BookDesk(req.DeskID, req.EmployeeName); err != nil {
		sendError(w, http.StatusConflict, err.Error())
		return
	}

	// async sync to third-party calendar
	go func(emp, deskID string) {
		log.Printf("syncing calendar for %s...", emp)
		time.Sleep(2 * time.Second) // simulate network latency
		log.Printf("calendar sync complete for %s", emp)
	}(req.EmployeeName, req.DeskID)

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": fmt.Sprintf("booked %s", req.DeskID),
	})
}

func (s *WorkplaceServer) handleGetHeatmap(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	resultChan := make(chan string)

	go func() {
		// simulate heavy spatial query
		time.Sleep(1 * time.Second)
		resultChan <- "Heatmap: Floor 4 is 85% full."
	}()

	select {
	case <-ctx.Done():
		sendError(w, http.StatusGatewayTimeout, "analytics timeout")
	case result := <-resultChan:
		writeJSON(w, http.StatusOK, map[string]string{"data": result})
	}
}

func main() {
	server := &WorkplaceServer{
		store: NewInMemoryStore(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", server.handleHealthCheck)
	mux.HandleFunc("GET /desks", server.handleGetDesks)
	mux.HandleFunc("POST /bookings", server.handleBookDesk)
	mux.HandleFunc("GET /analytics/heatmap", server.handleGetHeatmap)

	// wrap the entire mux router with our logging middleware
	handler := loggingMiddleware(mux)

	port := ":8080"
	log.Printf("server listening on %s", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("server fatal error: %v", err)
	}
}

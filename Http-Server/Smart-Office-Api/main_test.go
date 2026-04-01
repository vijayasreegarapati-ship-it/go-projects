package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// setupTestEnv acts as our Composition Root for testing.
// It gives every single test a brand new, clean in-memory database and router.
func setupTestEnv() http.Handler {
	server := &WorkplaceServer{
		store: NewInMemoryStore(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", server.handleHealthCheck)
	mux.HandleFunc("GET /desks", server.handleGetDesks)
	mux.HandleFunc("POST /bookings", server.handleBookDesk)
	mux.HandleFunc("GET /analytics/heatmap", server.handleGetHeatmap)

	// Wrap it in our middleware exactly like main.go does
	return loggingMiddleware(mux)
}

// --- 1. Basic Endpoint Tests ---

func TestHandleHealthCheck(t *testing.T) {
	handler := setupTestEnv()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestHandleGetDesks(t *testing.T) {
	handler := setupTestEnv()

	req := httptest.NewRequest(http.MethodGet, "/desks", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// We expect exactly 2 desks to be available based on our mock data
	var desks []Desk
	if err := json.NewDecoder(rr.Body).Decode(&desks); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(desks) != 2 {
		t.Errorf("expected 2 available desks, got %d", len(desks))
	}
}

// --- 2. Table-Driven Test (for booking a desk) ---

func TestHandleBookDesk(t *testing.T) {
	handler := setupTestEnv()

	// The Table: Define all possible scenarios, inputs, and expected outputs
	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success - Valid Booking",
			payload:        `{"desk_id": "desk-1", "employee_name": "Alan Turing"}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   "booked desk-1",
		},
		{
			name:           "Error - Missing Required Fields",
			payload:        `{"desk_id": "desk-1"}`, // missing employee_name
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing required fields",
		},
		{
			name:           "Error - Invalid JSON format",
			payload:        `{"desk_id": "desk-1", oops formatting}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid payload",
		},
		{
			name:           "Error - Desk Already Booked",
			payload:        `{"desk_id": "desk-3", "employee_name": "Grace Hopper"}`, // desk-3 is hardcoded as taken
			expectedStatus: http.StatusConflict,
			expectedBody:   "desk already booked",
		},
		{
			name:           "Error - Desk Does Not Exist",
			payload:        `{"desk_id": "desk-999", "employee_name": "Grace Hopper"}`,
			expectedStatus: http.StatusConflict,
			expectedBody:   "desk not found",
		},
	}

	// Loop over the table and run each scenario
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body := bytes.NewBuffer([]byte(tc.payload))
			req := httptest.NewRequest(http.MethodPost, "/bookings", body)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// 1. Assert Status Code
			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			// 2. Assert Error/Success Message
			if !strings.Contains(rr.Body.String(), tc.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tc.expectedBody, rr.Body.String())
			}
		})
	}
}

// --- 3. Context & Timeout Test ---

func TestHandleGetHeatmap(t *testing.T) {
	handler := setupTestEnv()

	req := httptest.NewRequest(http.MethodGet, "/analytics/heatmap", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Since our main.go code has the AI sleeping for 1 second, and the Context
	// timeout is 2 seconds, this should succeed with a 200 OK.
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Heatmap: Floor 4") {
		t.Errorf("expected body to contain heatmap data, got %s", rr.Body.String())
	}
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChain_SingleMiddleware(t *testing.T) {
	called := false

	// Final handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("handler"))
	})

	// Single middleware
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			next.ServeHTTP(w, r)
		})
	}

	wrapped := Chain(handler, middleware)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if !called {
		t.Error("Middleware was not called")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestChain_MultipleMiddlewares(t *testing.T) {
	order := []string{}

	// Final handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(http.StatusOK)
	})

	// Middlewares
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m1-before")
			next.ServeHTTP(w, r)
			order = append(order, "m1-after")
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m2-before")
			next.ServeHTTP(w, r)
			order = append(order, "m2-after")
		})
	}

	middleware3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m3-before")
			next.ServeHTTP(w, r)
			order = append(order, "m3-after")
		})
	}

	// Chain: m1, m2, m3
	wrapped := Chain(handler, middleware1, middleware2, middleware3)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	// Expected order: m1 -> m2 -> m3 -> handler -> m3 -> m2 -> m1
	expectedOrder := []string{
		"m1-before", "m2-before", "m3-before", "handler", "m3-after", "m2-after", "m1-after",
	}

	if len(order) != len(expectedOrder) {
		t.Fatalf("Expected %d calls, got %d", len(expectedOrder), len(order))
	}

	for i, expected := range expectedOrder {
		if order[i] != expected {
			t.Errorf("Step %d: expected %s, got %s", i, expected, order[i])
		}
	}
}

func TestChain_NoMiddlewares(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("handler"))
	})

	// Chain with no middlewares should return the handler as-is
	wrapped := Chain(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "handler" {
		t.Errorf("Expected body 'handler', got %s", w.Body.String())
	}
}

func TestChain_MiddlewareShortCircuit(t *testing.T) {
	handlerCalled := false

	// Final handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Middleware that short-circuits
	blockingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			// Don't call next.ServeHTTP - short circuit!
		})
	}

	wrapped := Chain(handler, blockingMiddleware)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if handlerCalled {
		t.Error("Handler should not have been called (middleware short-circuited)")
	}

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

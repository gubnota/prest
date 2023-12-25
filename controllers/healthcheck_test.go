package controllers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/prest/prest/adapters"
	"github.com/prest/prest/testutils"
	"github.com/stretchr/testify/require"
)

func TestCheckDBHealth(t *testing.T) {
	require.Nil(t, CheckDBHealth(context.Background(), nil))
}

func healthyDB(context.Context, adapters.Adapter) error { return nil }
func unhealthyDB(context.Context, adapters.Adapter) error {
	return errors.New("could not connect to the database")
}

func TestHealthStatus(t *testing.T) {
	for _, tc := range []struct {
		checkDBHealth func(context.Context, adapters.Adapter) error
		desc          string
		expected      int
	}{
		{healthyDB, "healthy database", http.StatusOK},
		{unhealthyDB, "unhealthy database", http.StatusServiceUnavailable},
	} {
		router := mux.NewRouter()
		checks := CheckList{tc.checkDBHealth}
		cfg := Config{adapter: nil}
		router.HandleFunc("/_health", cfg.WrappedHealthCheck(checks)).Methods("GET")
		server := httptest.NewServer(router)
		defer server.Close()
		testutils.DoRequest(t, server.URL+"/_health", nil, "GET", tc.expected, "")
	}
}

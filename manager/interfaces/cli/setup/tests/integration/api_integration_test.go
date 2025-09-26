package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockAPIServer creates a mock API server for testing
func createMockAPIServer() *httptest.Server {
	mux := http.NewServeMux()

	// Mock setup endpoint
	mux.HandleFunc("/api/v1/setup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		response := map[string]interface{}{
			"success": true,
			"config": map[string]interface{}{
				"manager": map[string]interface{}{
					"home_dir":     "/tmp/syntropy",
					"log_level":    "info",
					"api_endpoint": "https://api.syntropy.com",
				},
				"owner_key": map[string]interface{}{
					"type": "Ed25519",
					"path": "/tmp/syntropy/keys/owner.key",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Mock validation endpoint
	mux.HandleFunc("/api/v1/validate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		response := map[string]interface{}{
			"valid": true,
			"environment": map[string]interface{}{
				"os":              "linux",
				"architecture":    "amd64",
				"has_admin":       true,
				"available_disk":  100.0,
				"has_internet":    true,
				"home_dir":        "/home/test",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Mock status endpoint
	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		response := map[string]interface{}{
			"status":      "active",
			"last_check":  time.Now().Unix(),
			"version":     "1.0.0",
			"environment": "linux",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	return httptest.NewServer(mux)
}

// TestAPIIntegration tests the API integration functionality
func TestAPIIntegration(t *testing.T) {
	server := createMockAPIServer()
	defer server.Close()

	t.Run("Setup API Integration", func(t *testing.T) {
		t.Run("Should successfully call setup API", func(t *testing.T) {
			client := &http.Client{Timeout: 10 * time.Second}
			
			req, err := http.NewRequest("POST", server.URL+"/api/v1/setup", nil)
			require.NoError(t, err)
			
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			
			assert.True(t, response["success"].(bool))
			assert.NotNil(t, response["config"])
		})

		t.Run("Should handle API timeout gracefully", func(t *testing.T) {
			client := &http.Client{Timeout: 1 * time.Nanosecond}
			
			req, err := http.NewRequest("POST", server.URL+"/api/v1/setup", nil)
			require.NoError(t, err)
			
			_, err = client.Do(req)
			assert.Error(t, err)
			// Check for timeout-related error messages
			errorMsg := err.Error()
			assert.True(t, 
				strings.Contains(errorMsg, "timeout") || 
				strings.Contains(errorMsg, "deadline exceeded") || 
				strings.Contains(errorMsg, "Client.Timeout exceeded"),
				"Expected timeout-related error, got: %s", errorMsg)
		})

		t.Run("Should handle invalid API endpoint", func(t *testing.T) {
			client := &http.Client{Timeout: 5 * time.Second}
			
			req, err := http.NewRequest("POST", "http://invalid-endpoint/api/v1/setup", nil)
			require.NoError(t, err)
			
			_, err = client.Do(req)
			assert.Error(t, err)
		})
	})

	t.Run("Validation API Integration", func(t *testing.T) {
		t.Run("Should successfully call validation API", func(t *testing.T) {
			client := &http.Client{Timeout: 10 * time.Second}
			
			req, err := http.NewRequest("POST", server.URL+"/api/v1/validate", nil)
			require.NoError(t, err)
			
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			
			assert.True(t, response["valid"].(bool))
			assert.NotNil(t, response["environment"])
		})

		t.Run("Should handle validation API errors", func(t *testing.T) {
			client := &http.Client{Timeout: 10 * time.Second}
			
			req, err := http.NewRequest("GET", server.URL+"/api/v1/validate", nil)
			require.NoError(t, err)
			
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		})
	})

	t.Run("Status API Integration", func(t *testing.T) {
		t.Run("Should successfully call status API", func(t *testing.T) {
			client := &http.Client{Timeout: 10 * time.Second}
			
			req, err := http.NewRequest("GET", server.URL+"/api/v1/status", nil)
			require.NoError(t, err)
			
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			
			assert.Equal(t, "active", response["status"])
			assert.NotNil(t, response["last_check"])
		})
	})
}

// TestAPIIntegrationWithContext tests API integration with context
func TestAPIIntegrationWithContext(t *testing.T) {
	server := createMockAPIServer()
	defer server.Close()

	t.Run("Should respect context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		client := &http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequestWithContext(ctx, "POST", server.URL+"/api/v1/setup", nil)
		require.NoError(t, err)
		
		_, err = client.Do(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("Should respect context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		client := &http.Client{}
		req, err := http.NewRequestWithContext(ctx, "POST", server.URL+"/api/v1/setup", nil)
		require.NoError(t, err)
		
		_, err = client.Do(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context deadline exceeded")
	})
}

// TestAPIIntegrationRetry tests retry logic for API calls
func TestAPIIntegrationRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		
		response := map[string]interface{}{
			"success": true,
			"attempt": attempts,
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	t.Run("Should retry on server errors", func(t *testing.T) {
		maxRetries := 3
		client := &http.Client{Timeout: 10 * time.Second}
		
		var lastErr error
		for i := 0; i < maxRetries; i++ {
			req, err := http.NewRequest("POST", server.URL, nil)
			require.NoError(t, err)
			
			resp, err := client.Do(req)
			if err != nil {
				lastErr = err
				continue
			}
			
			if resp.StatusCode == http.StatusOK {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				resp.Body.Close()
				require.NoError(t, err)
				
				assert.True(t, response["success"].(bool))
				assert.Equal(t, float64(3), response["attempt"].(float64))
				return
			}
			
			resp.Body.Close()
			time.Sleep(100 * time.Millisecond) // Brief delay between retries
		}
		
		if lastErr != nil {
			t.Fatalf("All retries failed: %v", lastErr)
		}
	})
}

// BenchmarkAPIIntegration benchmarks API integration performance
func BenchmarkAPIIntegration(b *testing.B) {
	server := createMockAPIServer()
	defer server.Close()

	client := &http.Client{Timeout: 10 * time.Second}

	b.Run("Setup API Call", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("POST", server.URL+"/api/v1/setup", nil)
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})

	b.Run("Status API Call", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", server.URL+"/api/v1/status", nil)
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})
}
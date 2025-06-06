// Package bot provides the main service for running CronAI in bot mode.
package bot

import (
	"net"
	"os"
	"strings"
	"testing"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: Config{
				Port:      "8080",
				Secret:    "test-secret",
				Model:     "openai",
				Processor: "",
			},
			wantErr: false,
		},
		{
			name:    "empty config",
			cfg:     Config{},
			wantErr: false, // Should use defaults
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewService(tt.cfg)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if service == nil {
					t.Error("NewService() returned nil service")
					return
				}

				if service.server == nil {
					t.Error("NewService() server is nil")
				}

				if service.router == nil {
					t.Error("NewService() router is nil")
				}

				if service.model == nil {
					t.Error("NewService() model is nil")
				}

				if service.logger == nil {
					t.Error("NewService() logger is nil")
				}

				// Check that handlers are registered (only if router is not nil)
				if service.router != nil {
					types := service.router.GetRegisteredTypes()
					if len(types) == 0 {
						t.Error("NewService() no handlers registered")
					}
				}
			}
		})
	}
}

func TestStartService(t *testing.T) {
	// Save original env vars
	originalPort := os.Getenv("CRONAI_BOT_PORT")
	originalSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	defer func() {
		// Restore original values
		if originalPort != "" {
			_ = os.Setenv("CRONAI_BOT_PORT", originalPort) //nolint:errcheck // Test cleanup
		} else {
			_ = os.Unsetenv("CRONAI_BOT_PORT") //nolint:errcheck // Test cleanup
		}
		if originalSecret != "" {
			_ = os.Setenv("GITHUB_WEBHOOK_SECRET", originalSecret) //nolint:errcheck // Test cleanup
		} else {
			_ = os.Unsetenv("GITHUB_WEBHOOK_SECRET") //nolint:errcheck // Test cleanup
		}
	}()

	tests := []struct {
		name          string
		port          string
		webhookSecret string
		wantErr       bool
		errorContains string
	}{
		{
			name:          "short webhook secret",
			port:          "0",
			webhookSecret: "short",
			wantErr:       true,
			errorContains: "webhook secret",
		},
		{
			name:          "invalid port",
			port:          "invalid",
			webhookSecret: "valid-secret-123",
			wantErr:       true,
			errorContains: "port",
		},
		{
			name:          "port already in use",
			port:          "8080",
			webhookSecret: "valid-secret-123",
			wantErr:       true,
			errorContains: "port",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test environment variables
			_ = os.Setenv("CRONAI_BOT_PORT", tt.port)                //nolint:errcheck // Test setup
			_ = os.Setenv("GITHUB_WEBHOOK_SECRET", tt.webhookSecret) //nolint:errcheck // Test setup

			// For "port already in use" test, create a listener on that port
			var listener net.Listener
			if tt.name == "port already in use" {
				var err error
				listener, err = net.Listen("tcp", ":8080")
				if err != nil {
					t.Skipf("Cannot test port in use: %v", err)
				}
				defer func() { _ = listener.Close() }() //nolint:errcheck // Test cleanup
			}

			err := StartService("")

			if (err != nil) != tt.wantErr {
				t.Errorf("StartService() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && tt.errorContains != "" {
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("StartService() error = %v, want error containing %s", err, tt.errorContains)
				}
			}
		})
	}
}

func TestServiceStartStop(t *testing.T) {
	service, err := NewService(Config{
		Port:      "0", // Use port 0 for testing
		Secret:    "test",
		Model:     "openai",
		Processor: "",
	})
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	// Test that we can create the service
	if service == nil {
		t.Fatal("NewService() returned nil")
	}

	// We could test Start/Stop but it would require more complex setup
	// For now, just test that the service was created successfully
}

package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		main()
		return
	}

	tests := []struct {
		name        string
		setupEnv    func() (map[string]string, func())
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name: "invalid credentials",
			setupEnv: func() (map[string]string, func()) {
				return map[string]string{
					"SECURE_ENV_INITIAL_USER": "test_user",
					"SECURE_ENV_INITIAL_PASS": "test_pass",
				}, func() {}
			},
			args:        []string{"wrong_user", "wrong_pass"},
			wantErr:     true,
			errContains: "Credenciais inválidas",
		},
		{
			name: "missing environment variables",
			setupEnv: func() (map[string]string, func()) {
				return map[string]string{}, func() {}
			},
			args:        []string{"test_user", "test_pass"},
			wantErr:     true,
			errContains: "SECURE_ENV_INITIAL_USER and SECURE_ENV_INITIAL_PASS must be set",
		},
		{
			name: "missing command line arguments",
			setupEnv: func() (map[string]string, func()) {
				return map[string]string{
					"SECURE_ENV_INITIAL_USER": "test_user",
					"SECURE_ENV_INITIAL_PASS": "test_pass",
				}, func() {}
			},
			args:        []string{},
			wantErr:     true,
			errContains: "Usage: secure-env <user> <pass>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build test binary
			tmpDir, err := os.MkdirTemp("", "test-secure-env")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			testBin := filepath.Join(tmpDir, "secure-env.test")
			build := exec.Command("go", "build", "-o", testBin)
			if err := build.Run(); err != nil {
				t.Fatalf("Failed to build test binary: %v", err)
			}

			// Prepare command
			cmd := exec.Command(testBin, tt.args...)
			cmd.Env = []string{"TEST_SUBPROCESS=1"}

			// Setup environment
			env, cleanup := tt.setupEnv()
			defer cleanup()
			for k, v := range env {
				cmd.Env = append(cmd.Env, k+"="+v)
			}

			// Capture output
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			// Run command with timeout
			done := make(chan error, 1)
			go func() {
				done <- cmd.Run()
			}()

			// Wait for command to finish or timeout
			var cmdErr error
			select {
			case cmdErr = <-done:
				// Command completed
			case <-time.After(5 * time.Second):
				cmd.Process.Kill()
				t.Fatal("Command timed out")
			}

			// Check error
			if (cmdErr != nil) != tt.wantErr {
				t.Errorf("Process exited with error = %v, wantErr %v", cmdErr, tt.wantErr)
				return
			}

			if tt.wantErr {
				errOutput := stderr.String() + stdout.String()
				if tt.errContains != "" && !bytes.Contains([]byte(errOutput), []byte(tt.errContains)) {
					t.Errorf("Expected error containing %q, got %q", tt.errContains, errOutput)
				}
			}
		})
	}
}

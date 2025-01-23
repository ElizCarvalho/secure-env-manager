package config

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	// Backup das variáveis de ambiente originais
	originalUser := os.Getenv("SECURE_ENV_INITIAL_USER")
	originalPass := os.Getenv("SECURE_ENV_INITIAL_PASS")
	defer func() {
		os.Setenv("SECURE_ENV_INITIAL_USER", originalUser)
		os.Setenv("SECURE_ENV_INITIAL_PASS", originalPass)
	}()

	tests := []struct {
		name        string
		setupEnv    func()
		wantErr     bool
		errContains string
	}{
		{
			name: "Sucesso - Variáveis configuradas corretamente",
			setupEnv: func() {
				os.Setenv("SECURE_ENV_INITIAL_USER", "test_user")
				os.Setenv("SECURE_ENV_INITIAL_PASS", "test_pass")
			},
			wantErr: false,
		},
		{
			name: "Erro - Usuário não configurado",
			setupEnv: func() {
				os.Unsetenv("SECURE_ENV_INITIAL_USER")
				os.Setenv("SECURE_ENV_INITIAL_PASS", "test_pass")
			},
			wantErr:     true,
			errContains: "SECURE_ENV_INITIAL_USER and SECURE_ENV_INITIAL_PASS must be set",
		},
		{
			name: "Erro - Senha não configurada",
			setupEnv: func() {
				os.Setenv("SECURE_ENV_INITIAL_USER", "test_user")
				os.Unsetenv("SECURE_ENV_INITIAL_PASS")
			},
			wantErr:     true,
			errContains: "SECURE_ENV_INITIAL_USER and SECURE_ENV_INITIAL_PASS must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			cfg, err := New()

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				if err.Error() != tt.errContains {
					t.Errorf("Expected error '%s' but got '%s'", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if cfg == nil {
				t.Error("Expected config but got nil")
			}
		})
	}
}

func TestValidateCredentials(t *testing.T) {
	cfg := &Config{
		InitialUser: "test_user",
		InitialPass: "test_pass",
	}

	tests := []struct {
		name        string
		user        string
		pass        string
		wantErr     bool
		errContains string
	}{
		{
			name:    "Sucesso - Credenciais válidas",
			user:    "test_user",
			pass:    "test_pass",
			wantErr: false,
		},
		{
			name:        "Erro - Usuário inválido",
			user:        "wrong_user",
			pass:        "test_pass",
			wantErr:     true,
			errContains: "invalid credentials",
		},
		{
			name:        "Erro - Senha inválida",
			user:        "test_user",
			pass:        "wrong_pass",
			wantErr:     true,
			errContains: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cfg.ValidateCredentials(tt.user, tt.pass)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				if err.Error() != tt.errContains {
					t.Errorf("Expected error '%s' but got '%s'", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

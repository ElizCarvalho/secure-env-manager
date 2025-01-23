package handler

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"secure-env-manager/internal/crypto"
)

// mockStdin creates a mock stdin for testing
func mockStdin(input string) (*os.File, error) {
	tmpfile, err := os.CreateTemp("", "stdin")
	if err != nil {
		return nil, err
	}

	if _, err := tmpfile.Write([]byte(input)); err != nil {
		return nil, err
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		return nil, err
	}

	return tmpfile, nil
}

func TestHandleEncrypt(t *testing.T) {
	h := New("test-password")

	// Criar diretório temporário
	tmpDir, err := os.MkdirTemp("", "test-project")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Criar arquivo .env
	envFile := filepath.Join(tmpDir, "test.env")
	if err := os.WriteFile(envFile, []byte("TEST=value"), 0644); err != nil {
		t.Fatalf("failed to create env file: %v", err)
	}

	// Criar diretório do projeto
	projectDir := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	// Executar função
	if err := h.HandleEncrypt(projectDir, envFile); err != nil {
		t.Errorf("HandleEncrypt() error = %v", err)
	}

	// Verificar se arquivo foi criado
	encPath := filepath.Join(projectDir, ".env.enc")
	if _, err := os.Stat(encPath); os.IsNotExist(err) {
		t.Errorf("encrypted file was not created")
	}
}

func TestHandleDecrypt(t *testing.T) {
	h := New("test-password")

	// Criar diretório temporário
	tmpDir, err := os.MkdirTemp("", "test-project")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Criar diretório do projeto e arquivo criptografado
	projectDir := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	// Criar arquivo criptografado
	encrypted, err := crypto.Encrypt([]byte("TEST=value"), "test-password")
	if err != nil {
		t.Fatalf("failed to encrypt data: %v", err)
	}
	encPath := filepath.Join(projectDir, ".env.enc")
	if err := os.WriteFile(encPath, []byte(encrypted), 0644); err != nil {
		t.Fatalf("failed to create encrypted file: %v", err)
	}

	// Definir arquivo de saída
	outputFile := filepath.Join(tmpDir, "output.env")

	// Executar função
	if err := h.HandleDecrypt(projectDir, outputFile); err != nil {
		t.Errorf("HandleDecrypt() error = %v", err)
	}

	// Verificar se arquivo foi criado
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("decrypted file was not created")
	}
}

func TestReadOption(t *testing.T) {
	h := New("test-password")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid option",
			input:    "1\n",
			expected: "1",
		},
		{
			name:     "option with whitespace",
			input:    "  2  \n",
			expected: "2",
		},
		{
			name:     "invalid option",
			input:    "invalid\n",
			expected: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare test input
			oldStdin := os.Stdin
			tmpStdin, err := mockStdin(tt.input)
			if err != nil {
				t.Fatalf("failed to mock stdin: %v", err)
			}
			defer tmpStdin.Close()
			os.Stdin = tmpStdin
			defer func() { os.Stdin = oldStdin }()

			// Run the function
			result := h.ReadOption()

			// Verify result
			if result != tt.expected {
				t.Errorf("ReadOption() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShowMenu(t *testing.T) {
	h := New("test-password")

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	// Run the function
	h.ShowMenu()

	// Restore stdout
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify menu items are present
	expectedItems := []string{
		"=== Secure ENV Manager ===",
		"1. Encrypt .env file",
		"2. Decrypt file",
		"3. Exit",
		"Choose an option:",
	}

	for _, item := range expectedItems {
		if !strings.Contains(output, item) {
			t.Errorf("menu missing item: %s", item)
		}
	}
}

func TestHandleEncryptErrors(t *testing.T) {
	h := New("test-password")

	tests := []struct {
		name        string
		setupFunc   func() (string, string, func())
		wantErr     bool
		errContains string
	}{
		{
			name: "non-existent env file",
			setupFunc: func() (string, string, func()) {
				tmpDir, _ := os.MkdirTemp("", "test-project")
				nonExistentFile := filepath.Join(tmpDir, "nonexistent.env")
				return tmpDir, nonExistentFile, func() { os.RemoveAll(tmpDir) }
			},
			wantErr:     true,
			errContains: "error reading file",
		},
		{
			name: "permission denied on project dir",
			setupFunc: func() (string, string, func()) {
				// Criar diretório temporário
				tmpDir, _ := os.MkdirTemp("", "test-project")

				// Criar arquivo .env
				envFile := filepath.Join(tmpDir, "test.env")
				os.WriteFile(envFile, []byte("TEST=value"), 0644)

				// Criar diretório pai sem permissão de escrita
				readOnlyParent := filepath.Join(tmpDir, "readonly-parent")
				os.MkdirAll(readOnlyParent, 0755)
				os.Chmod(readOnlyParent, 0555) // r-xr-xr-x

				// Tentar criar diretório do projeto dentro do pai somente leitura
				projectDir := filepath.Join(readOnlyParent, "project")

				return projectDir, envFile, func() {
					os.Chmod(readOnlyParent, 0755)
					os.RemoveAll(tmpDir)
				}
			},
			wantErr:     true,
			errContains: "error creating directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectDir, envFile, cleanup := tt.setupFunc()
			defer cleanup()

			err := h.HandleEncrypt(projectDir, envFile)

			if (err != nil) != tt.wantErr {
				t.Errorf("HandleEncrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("HandleEncrypt() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
}

func TestHandleDecryptErrors(t *testing.T) {
	h := New("test-password")

	tests := []struct {
		name        string
		setupFunc   func() (string, string, func())
		wantErr     bool
		errContains string
	}{
		{
			name: "non-existent project",
			setupFunc: func() (string, string, func()) {
				tmpDir, _ := os.MkdirTemp("", "test-project")
				nonExistentProject := filepath.Join(tmpDir, "nonexistent")
				outputFile := filepath.Join(tmpDir, "output.env")
				return nonExistentProject, outputFile, func() { os.RemoveAll(tmpDir) }
			},
			wantErr:     true,
			errContains: "error reading file",
		},
		{
			name: "invalid encrypted data",
			setupFunc: func() (string, string, func()) {
				tmpDir, _ := os.MkdirTemp("", "test-project")
				projectDir := filepath.Join(tmpDir, "test-project")
				os.MkdirAll(projectDir, 0755)
				encPath := filepath.Join(projectDir, ".env.enc")
				os.WriteFile(encPath, []byte("invalid-data"), 0644)
				outputFile := filepath.Join(tmpDir, "output.env")
				return projectDir, outputFile, func() { os.RemoveAll(tmpDir) }
			},
			wantErr:     true,
			errContains: "error decrypting",
		},
		{
			name: "permission denied on output file",
			setupFunc: func() (string, string, func()) {
				// Criar diretório temporário
				tmpDir, _ := os.MkdirTemp("", "test-project")

				// Criar diretório do projeto e arquivo criptografado
				projectDir := filepath.Join(tmpDir, "test-project")
				os.MkdirAll(projectDir, 0755)
				encrypted, _ := crypto.Encrypt([]byte("TEST=value"), "test-password")
				encPath := filepath.Join(projectDir, ".env.enc")
				os.WriteFile(encPath, []byte(encrypted), 0644)

				// Criar diretório pai sem permissão de escrita
				readOnlyParent := filepath.Join(tmpDir, "readonly-parent")
				os.MkdirAll(readOnlyParent, 0755)
				os.Chmod(readOnlyParent, 0555) // r-xr-xr-x

				// Tentar criar arquivo dentro do diretório somente leitura
				outputFile := filepath.Join(readOnlyParent, "output.env")

				return projectDir, outputFile, func() {
					os.Chmod(readOnlyParent, 0755)
					os.RemoveAll(tmpDir)
				}
			},
			wantErr:     true,
			errContains: "error creating file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectDir, outputFile, cleanup := tt.setupFunc()
			defer cleanup()

			err := h.HandleDecrypt(projectDir, outputFile)

			if (err != nil) != tt.wantErr {
				t.Errorf("HandleDecrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("HandleDecrypt() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
}

package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		password string
		wantErr  bool
	}{
		{
			name:     "valid encryption/decryption",
			data:     []byte("test data"),
			password: "test-password",
			wantErr:  false,
		},
		{
			name:     "empty data",
			data:     []byte{},
			password: "test-password",
			wantErr:  false,
		},
		{
			name:     "empty password",
			data:     []byte("test data"),
			password: "",
			wantErr:  true,
		},
		{
			name:     "very long data",
			data:     bytes.Repeat([]byte("a"), 1024*1024), // 1MB
			password: "test-password",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encryption
			encrypted, err := Encrypt(tt.data, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Test decryption
			decrypted, err := Decrypt(encrypted, tt.password)
			if err != nil {
				t.Errorf("Decrypt() error = %v", err)
				return
			}

			if !bytes.Equal(decrypted, tt.data) {
				t.Errorf("Decrypt() = %v, want %v", decrypted, tt.data)
			}
		})
	}
}

func TestDecryptErrors(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		password string
		wantErr  bool
	}{
		{
			name:     "invalid base64",
			data:     "invalid base64!@#$",
			password: "test-password",
			wantErr:  true,
		},
		{
			name:     "wrong password",
			data:     "dmFsaWQgYmFzZTY0", // valid base64 but not encrypted data
			password: "wrong-password",
			wantErr:  true,
		},
		{
			name:     "corrupted data",
			data:     "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo=", // valid base64 but not encrypted data
			password: "test-password",
			wantErr:  true,
		},
		{
			name:     "empty encrypted data",
			data:     "",
			password: "test-password",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(tt.data, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeriveKey(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantLen  int
	}{
		{
			name:     "normal password",
			password: "test-password",
			wantLen:  32,
		},
		{
			name:     "empty password",
			password: "",
			wantLen:  32,
		},
		{
			name:     "very long password",
			password: string(bytes.Repeat([]byte("a"), 1024)),
			wantLen:  32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := DeriveKey(tt.password)
			if len(key) != tt.wantLen {
				t.Errorf("DeriveKey() len = %v, want %v", len(key), tt.wantLen)
			}
		})
	}
}

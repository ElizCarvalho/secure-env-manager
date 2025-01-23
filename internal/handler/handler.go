package handler

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"secure-env-manager/internal/crypto"
)

// Handler handles the encryption/decryption operations
type Handler struct {
	password string
}

// New creates a new Handler instance
func New(password string) *Handler {
	return &Handler{
		password: password,
	}
}

// ShowMenu displays the main menu
func (h *Handler) ShowMenu() {
	fmt.Println("\n=== Secure ENV Manager ===")
	fmt.Println("1. Encrypt .env file")
	fmt.Println("2. Decrypt file")
	fmt.Println("3. Exit")
	fmt.Print("\nChoose an option: ")
}

// ReadOption reads a menu option from stdin
func (h *Handler) ReadOption() string {
	reader := bufio.NewReader(os.Stdin)
	option, _ := reader.ReadString('\n')
	return strings.TrimSpace(option)
}

// HandleEncrypt handles the encryption of a .env file
func (h *Handler) HandleEncrypt(project, envPath string) error {
	// Ler arquivo .env
	data, err := os.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Criar diretório do projeto se não existir
	if err := os.MkdirAll(project, 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	// Criptografar dados
	encrypted, err := crypto.Encrypt(data, h.password)
	if err != nil {
		return fmt.Errorf("error encrypting: %v", err)
	}

	// Salvar arquivo criptografado
	encPath := filepath.Join(project, ".env.enc")
	if err := os.WriteFile(encPath, []byte(encrypted), 0644); err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	fmt.Printf("\nEncrypted file saved at: %s\n", encPath)
	return nil
}

// HandleDecrypt handles the decryption of an encrypted file
func (h *Handler) HandleDecrypt(project, outputPath string) error {
	// Ler arquivo criptografado
	encPath := filepath.Join(project, ".env.enc")
	data, err := os.ReadFile(encPath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Descriptografar dados
	decrypted, err := crypto.Decrypt(string(data), h.password)
	if err != nil {
		return fmt.Errorf("error decrypting: %v", err)
	}

	// Salvar arquivo descriptografado
	if err := os.WriteFile(outputPath, decrypted, 0644); err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	fmt.Printf("Decrypted file saved at: %s\n", outputPath)
	return nil
}

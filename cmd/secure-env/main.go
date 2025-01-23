package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"secure-env-manager/internal/handler"
)

func main() {
	// Verificar credenciais
	if len(os.Args) != 3 {
		fmt.Println("Usage: secure-env <user> <pass>")
		os.Exit(1)
	}

	// Obter credenciais do ambiente
	initialUser := os.Getenv("SECURE_ENV_INITIAL_USER")
	initialPass := os.Getenv("SECURE_ENV_INITIAL_PASS")

	if initialUser == "" || initialPass == "" {
		fmt.Println("SECURE_ENV_INITIAL_USER and SECURE_ENV_INITIAL_PASS must be set")
		os.Exit(1)
	}

	// Validar credenciais
	if os.Args[1] != initialUser || os.Args[2] != initialPass {
		fmt.Println("Credenciais inválidas")
		os.Exit(1)
	}

	// Criar handler
	h := handler.New(initialPass)

	// Menu principal
	for {
		fmt.Println("\n=== Secure ENV Manager ===")
		fmt.Println("1. Encrypt .env file")
		fmt.Println("2. Decrypt .env file")
		fmt.Println("3. Exit")
		fmt.Print("\nOption: ")

		reader := bufio.NewReader(os.Stdin)
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			fmt.Print("\nProject name: ")
			project, _ := reader.ReadString('\n')
			project = strings.TrimSpace(project)

			fmt.Print("Path to .env file: ")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSpace(path)

			if err := h.HandleEncrypt(project, path); err != nil {
				fmt.Printf("Error: %v\n", err)
			}

		case "2":
			fmt.Print("\nProject name: ")
			project, _ := reader.ReadString('\n')
			project = strings.TrimSpace(project)

			fmt.Print("Path to save decrypted file: ")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSpace(path)

			if err := h.HandleDecrypt(project, path); err != nil {
				fmt.Printf("Error: %v\n", err)
			}

		case "3":
			fmt.Println("\nGoodbye!")
			os.Exit(0)

		default:
			fmt.Println("\nInvalid option")
		}
	}
}

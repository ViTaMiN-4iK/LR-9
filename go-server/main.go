package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("🚀 Go TCP Server starting on :8080...")

	// Создаем TCP-сервер, слушающий порт 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("❌ Error starting server: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("✅ Server is listening on :8080")

	// Принимаем одно соединение
	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("❌ Error accepting connection: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("✅ Accepted connection from %s\n", conn.RemoteAddr())

	// Читаем данные от клиента
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("❌ Error reading: %v\n", err)
		return
	}

	message := strings.TrimSpace(string(buffer[:n]))
	fmt.Printf("📩 Received: %q\n", message)

	// Формируем и отправляем ответ
	response := "Hello from Go"
	_, err = conn.Write([]byte(response + "\n"))
	if err != nil {
		fmt.Printf("❌ Error writing: %v\n", err)
		return
	}

	fmt.Printf("📤 Sent: %q\n", response)
	fmt.Println("👋 Connection closed, server shutting down...")
}

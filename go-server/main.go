package main

import (
	"fmt"
	"net"
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
	fmt.Println("🔄 Server will handle multiple connections CONCURRENTLY with goroutines")

	// Бесконечный цикл для обработки соединений
	for {
		// Принимаем соединение
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("❌ Error accepting connection: %v\n", err)
			continue // Продолжаем слушать дальше, даже если была ошибка
		}

		// ЗАПУСКАЕМ КАЖДОЕ СОЕДИНЕНИЕ В ОТДЕЛЬНОЙ ГОРУТИНЕ!
		go handleConnection(conn)
	}
}

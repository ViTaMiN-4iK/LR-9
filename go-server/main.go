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

// handleConnection обрабатывает одно клиентское соединение
// Теперь эта функция будет вызываться в горутинах
func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("✅ Accepted connection from %s (goroutine started)\n", conn.RemoteAddr())

	// Читаем данные от клиента
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("❌ Error reading from %s: %v\n", conn.RemoteAddr(), err)
		return
	}

	message := CleanMessage(buffer[:n])
	fmt.Printf("📩 Received from %s: %q\n", conn.RemoteAddr(), message)

	// Формируем и отправляем ответ
	response := ProcessMessage(message)
	_, err = conn.Write([]byte(response + "\n"))
	if err != nil {
		fmt.Printf("❌ Error writing to %s: %v\n", conn.RemoteAddr(), err)
		return
	}

	fmt.Printf("📤 Sent to %s: %q\n", conn.RemoteAddr(), response)
}

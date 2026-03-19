package main

import (
	"net"
	"testing"
	"time"
)

// TestHandleConnection проверяет логику обработки соединения
func TestHandleConnection(t *testing.T) {
	// Создаем пару соединений (клиент-сервер) в памяти
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	// Канал для получения результата
	done := make(chan bool)

	// Запускаем handleConnection в горутине
	go func() {
		handleConnection(server)
		done <- true
	}()

	// Отправляем тестовое сообщение от клиента
	testMessage := "ping"
	_, err := client.Write([]byte(testMessage + "\n"))
	if err != nil {
		t.Fatalf("Failed to write to connection: %v", err)
	}

	// Читаем ответ
	response := make([]byte, 1024)
	n, err := client.Read(response)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	expected := "Hello from Go\n"
	if string(response[:n]) != expected {
		t.Errorf("Expected %q, got %q", expected, string(response[:n]))
	}

	// Даем время на завершение handleConnection
	select {
	case <-done:
		// Все хорошо
	case <-time.After(time.Second):
		t.Fatal("handleConnection didn't close properly")
	}
}

// BenchmarkHandleConnection замеряет производительность обработки
func BenchmarkHandleConnection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client, server := net.Pipe()

		// Запускаем обработчик
		go handleConnection(server)

		// Отправляем и получаем данные
		client.Write([]byte("ping\n"))
		response := make([]byte, 1024)
		client.Read(response)

		client.Close()
		server.Close()
	}
}

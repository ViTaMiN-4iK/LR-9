package main

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

// TestConcurrentHandling проверяет, что сервер может обрабатывать
// несколько соединений одновременно
func TestConcurrentHandling(t *testing.T) {
	// Создаем слушатель на случайном порту
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port

	// Запускаем серверную часть в горутине
	serverDone := make(chan bool)
	go func() {
		for i := 0; i < 5; i++ { // Примем 5 соединений
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			// Здесь критически важно ЗАПУСТИТЬ ГОРУТИНУ!
			// Если убрать "go", тест упадет
			go handleConnection(conn)
		}
		serverDone <- true
	}()

	// Создаем WaitGroup для синхронизации клиентов
	var wg sync.WaitGroup

	// Канал для сбора результатов
	results := make(chan string, 5)

	// Запускаем 5 параллельных клиентов
	start := time.Now()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			// Каждый клиент подключается и отправляет сообщение
			conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
			if err != nil {
				t.Errorf("Client %d failed to connect: %v", clientID, err)
				return
			}
			defer conn.Close()

			// Отправляем сообщение
			message := fmt.Sprintf("ping from client %d", clientID)
			_, err = conn.Write([]byte(message + "\n"))
			if err != nil {
				t.Errorf("Client %d failed to write: %v", clientID, err)
				return
			}

			// Читаем ответ
			response := make([]byte, 1024)
			n, err := conn.Read(response)
			if err != nil {
				t.Errorf("Client %d failed to read: %v", clientID, err)
				return
			}

			results <- string(response[:n])
		}(i)
	}

	// Ждем завершения всех клиентов
	wg.Wait()
	close(results)
	elapsed := time.Since(start)

	// Проверяем, что все ответы получены
	responseCount := 0
	for range results {
		responseCount++
	}

	if responseCount != 5 {
		t.Errorf("Expected 5 responses, got %d", responseCount)
	}

	t.Logf("All 5 concurrent connections handled in %v", elapsed)
}

// BenchmarkConcurrentHandling замеряет производительность параллельной обработки
func BenchmarkConcurrentHandling(b *testing.B) {
	for i := 0; i < b.N; i++ {
		listener, _ := net.Listen("tcp", ":0")
		port := listener.Addr().(*net.TCPAddr).Port

		// Сервер
		go func() {
			for j := 0; j < 10; j++ {
				conn, _ := listener.Accept()
				go handleConnection(conn)
			}
		}()

		// 10 параллельных клиентов
		var wg sync.WaitGroup
		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				conn, _ := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
				conn.Write([]byte("ping\n"))
				response := make([]byte, 1024)
				conn.Read(response)
				conn.Close()
			}()
		}
		wg.Wait()
		listener.Close()
	}
}

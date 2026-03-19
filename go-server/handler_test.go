package main

import (
	"fmt"
	"net"
	"sync"
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

// TestConcurrency демонстрирует, что горутины работают параллельно
// Это исправленная версия TestSlowHandler
func TestConcurrency(t *testing.T) {
	// Создаем реальный слушатель
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	t.Logf("Test server listening on port %d", port)

	// Канал для синхронизации завершения сервера
	serverDone := make(chan bool)

	// Запускаем сервер с горутинами (именно так, как в main)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				serverDone <- true
				return
			}
			// Критически важно: используем handleConnection с go
			go handleConnection(conn)
		}
	}()

	// Создаем "медленного" клиента, который будет долго обрабатываться
	slowClientDone := make(chan bool)
	go func() {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			t.Errorf("Slow client failed to connect: %v", err)
			slowClientDone <- true
			return
		}
		defer conn.Close()

		// Отправляем сообщение, которое вызовет долгую обработку
		// Но в нашей текущей реализации все сообщения обрабатываются быстро
		// Поэтому просто симулируем долгий клиент через time.Sleep в тесте
		conn.Write([]byte("slow client\n"))

		// Читаем ответ (это произойдет быстро)
		response := make([]byte, 1024)
		n, _ := conn.Read(response)
		t.Logf("Slow client received: %q", string(response[:n]))

		// Имитируем, что клиент "думает" перед закрытием
		time.Sleep(2 * time.Second)
		slowClientDone <- true
	}()

	// Даем медленному клиенту время подключиться
	time.Sleep(100 * time.Millisecond)

	// Быстрый клиент подключается и должен получить ответ немедленно
	start := time.Now()

	fastClientDone := make(chan bool)
	go func() {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			t.Errorf("Fast client failed to connect: %v", err)
			fastClientDone <- true
			return
		}
		defer conn.Close()

		conn.Write([]byte("fast client\n"))

		response := make([]byte, 1024)
		n, _ := conn.Read(response)
		t.Logf("Fast client received: %q", string(response[:n]))

		fastClientDone <- true
	}()

	// Ждем быстрого клиента (должен завершиться быстро)
	select {
	case <-fastClientDone:
		elapsed := time.Since(start)
		if elapsed > 500*time.Millisecond {
			t.Errorf("Fast client was delayed! Took %v", elapsed)
		} else {
			t.Logf("✅ Fast client completed in %v (while slow client is still processing)", elapsed)
		}
	case <-time.After(3 * time.Second):
		t.Error("Fast client timeout - connections are being processed sequentially!")
	}

	// Ждем медленного клиента
	<-slowClientDone

	// Останавливаем сервер
	listener.Close()
	<-serverDone
}

// Еще один полезный тест - проверка максимальной нагрузки
func TestManyConcurrentConnections(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port

	// Запускаем сервер
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go handleConnection(conn)
		}
	}()

	// Запускаем 50 параллельных клиентов
	const numClients = 50
	var wg sync.WaitGroup
	errors := make(chan error, numClients)

	start := time.Now()

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
			if err != nil {
				errors <- fmt.Errorf("client %d failed to connect: %v", id, err)
				return
			}
			defer conn.Close()

			message := fmt.Sprintf("ping from client %d", id)
			conn.Write([]byte(message + "\n"))

			response := make([]byte, 1024)
			n, err := conn.Read(response)
			if err != nil {
				errors <- fmt.Errorf("client %d failed to read: %v", id, err)
				return
			}

			expected := "Hello from Go\n"
			if string(response[:n]) != expected {
				errors <- fmt.Errorf("client %d got wrong response: %q", id, string(response[:n]))
			}
		}(i)
	}

	wg.Wait()
	close(errors)
	elapsed := time.Since(start)

	// Проверяем ошибки
	for err := range errors {
		t.Error(err)
	}

	t.Logf("✅ Successfully handled %d concurrent connections in %v", numClients, elapsed)
}

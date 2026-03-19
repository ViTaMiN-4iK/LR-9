package main

import (
	"os"
	"testing"
	"time"
)

func TestMainDoesNotPanic(t *testing.T) {
	// Запускаем main в горутине, так как он теперь блокирующий
	done := make(chan bool)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("main panicked: %v", r)
			}
			done <- true
		}()

		// Временно перенаправляем stdout, чтобы не захламлять вывод тестов
		oldStdout := os.Stdout
		_, w, _ := os.Pipe()
		os.Stdout = w

		main()

		w.Close()
		os.Stdout = oldStdout
	}()

	// Даем main время запуститься, но не ждем соединения
	// Через 100мс закрываем тест (main все еще висит на Accept)
	select {
	case <-done:
		// main завершился сам (например, с ошибкой)
	case <-time.After(100 * time.Millisecond):
		// main все еще работает на Accept - это нормально для данного теста
		t.Log("main is running and waiting for connections")
	}
}

package main

import (
	"os"
	"testing"
	"time"
)

func TestMainDoesNotPanic(t *testing.T) {
	// Запускаем main в горутине, так как он теперь бесконечный
	done := make(chan bool)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("main panicked: %v", r)
			}
			done <- true
		}()

		// Временно перенаправляем stdout
		oldStdout := os.Stdout
		_, w, _ := os.Pipe()
		os.Stdout = w

		main()

		w.Close()
		os.Stdout = oldStdout
	}()

	// Даем main время запуститься
	time.Sleep(500 * time.Millisecond)

	// Проверяем, что main все еще работает (не завершился)
	select {
	case <-done:
		t.Error("main terminated unexpectedly")
	case <-time.After(100 * time.Millisecond):
		// main все еще работает - это хорошо!
		t.Log("main is running and waiting for connections")
	}
}

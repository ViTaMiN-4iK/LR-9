package main

import (
	"os"
	"testing"
)

// TestMainExists проверяет, что функция main объявлена
// Это базовый тест, который гарантирует, что пакет компилируется
func TestMainExists(t *testing.T) {
	// Просто проверяем, что main функция существует
	// Если бы main не было, тест бы даже не скомпилировался
	t.Log("✅ Пакет main успешно компилируется")
}

// TestMainOutput проверяет, что main выводит ожидаемое сообщение
// Но мы не можем легко перехватить вывод main(), поэтому
// пока просто проверяем, что main не паникует
func TestMainDoesNotPanic(t *testing.T) {
	// Сохраняем оригинальный stdout
	oldStdout := os.Stdout

	// Временный файл для перехвата вывода
	_, w, _ := os.Pipe()
	os.Stdout = w

	// Функция для восстановления stdout
	defer func() {
		w.Close()
		os.Stdout = oldStdout
	}()

	// Проверяем, что main не паникует
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main panicked: %v", r)
		}
	}()

	// Запускаем main
	main()
}

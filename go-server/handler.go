package main

import "strings"

// ProcessMessage обрабатывает входящее сообщение и возвращает ответ
func ProcessMessage(message string) string {
	// Пока просто возвращаем приветствие
	// Позже здесь будет более сложная логика
	return "Hello from Go"
}

// CleanMessage удаляет пробельные символы по краям
func CleanMessage(raw []byte) string {
	return strings.TrimSpace(string(raw))
}

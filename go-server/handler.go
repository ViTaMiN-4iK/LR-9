package main

import (
	"fmt"
	"net" // <- Добавь этот импорт!
	"strconv"
	"strings"
)

// ProcessMessage обрабатывает входящее сообщение и возвращает ответ
func ProcessMessage(message string) string {
	// Проверяем, не команда ли это на обработку данных
	if strings.HasPrefix(message, "process:") {
		// Убираем префикс "process:"
		dataStr := strings.TrimPrefix(message, "process:")

		// Разбираем числа
		numberStrings := strings.Split(dataStr, ",")
		result := make([]string, 0, len(numberStrings))

		for _, numStr := range numberStrings {
			numStr = strings.TrimSpace(numStr)
			if numStr == "" {
				continue
			}

			num, err := strconv.Atoi(numStr)
			if err != nil {
				fmt.Printf("Error parsing number: %v\n", err)
				continue
			}

			result = append(result, strconv.Itoa(num*2))
		}

		return "result:" + strings.Join(result, ",")
	}

	return "Hello from Go"
}

// CleanMessage удаляет пробельные символы по краям
func CleanMessage(raw []byte) string {
	return strings.TrimSpace(string(raw))
}

// handleConnection обрабатывает одно клиентское соединение
func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("✅ Accepted connection from %s (goroutine started)\n", conn.RemoteAddr())

	// Увеличиваем буфер для больших сообщений (64KB)
	buffer := make([]byte, 64*1024) // 64KB буфер
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("❌ Error reading from %s: %v\n", conn.RemoteAddr(), err)
		return
	}

	message := CleanMessage(buffer[:n])
	fmt.Printf("📩 Received %d bytes from %s\n", n, conn.RemoteAddr())

	// Показываем первые 50 символов сообщения для отладки
	preview := message
	if len(preview) > 50 {
		preview = preview[:50] + "..."
	}
	fmt.Printf("📩 Message preview: %q\n", preview)

	response := ProcessMessage(message)
	_, err = conn.Write([]byte(response + "\n"))
	if err != nil {
		fmt.Printf("❌ Error writing to %s: %v\n", conn.RemoteAddr(), err)
		return
	}

	// Показываем первые 50 символов ответа
	responsePreview := response
	if len(responsePreview) > 50 {
		responsePreview = responsePreview[:50] + "..."
	}
	fmt.Printf("📤 Sent to %s: %q\n", conn.RemoteAddr(), responsePreview)
}

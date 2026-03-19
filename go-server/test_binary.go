// Временный файл для тестирования бинарника
// Можно запустить как: go run test_binary.go
package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("Testing compiled binary...")

	// Компилируем
	cmd := exec.Command("go", "build", "-o", "test_server.exe")
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Build failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Build successful")

	// Запускаем бинарник
	server := exec.Command("./test_server.exe")
	server.Stdout = os.Stdout
	server.Stderr = os.Stderr
	if err := server.Start(); err != nil {
		fmt.Printf("❌ Failed to start server: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Server started")

	// Даем серверу время запуститься
	time.Sleep(1 * time.Second)

	// Тестируем подключение
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("❌ Failed to connect: %v\n", err)
		server.Process.Kill()
		os.Exit(1)
	}

	// Отправляем тестовое сообщение
	conn.Write([]byte("test from binary\n"))

	// Читаем ответ
	response := make([]byte, 1024)
	n, _ := conn.Read(response)
	fmt.Printf("✅ Got response: %q\n", string(response[:n]))

	conn.Close()
	server.Process.Kill()

	// Удаляем тестовый бинарник
	os.Remove("test_server.exe")

	fmt.Println("✅ Binary test passed!")
}

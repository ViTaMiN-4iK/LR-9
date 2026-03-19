#!/usr/bin/env python3
"""
Python client for Go TCP Server
Лабораторная работа №9: Мультиязычное программирование
"""

import socket
import logging
from typing import Optional

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


class GoTCPClient:
    """Клиент для взаимодействия с Go TCP сервером"""
    
    def __init__(self, host: str = 'localhost', port: int = 8080, timeout: int = 5):
        """
        Инициализация клиента
        
        Args:
            host: хост сервера
            port: порт сервера
            timeout: таймаут подключения в секундах
        """
        self.host = host
        self.port = port
        self.timeout = timeout
        self.socket: Optional[socket.socket] = None
    
    def connect(self) -> bool:
        """
        Подключение к серверу
        
        Returns:
            True если подключение успешно, иначе False
        """
        try:
            self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.socket.settimeout(self.timeout)
            self.socket.connect((self.host, self.port))
            logger.info(f"✅ Connected to {self.host}:{self.port}")
            return True
        except socket.error as e:
            logger.error(f"❌ Connection failed: {e}")
            self.socket = None
            return False
    
    def send_message(self, message: str) -> Optional[str]:
        """
        Отправка сообщения серверу и получение ответа
        
        Args:
            message: сообщение для отправки
            
        Returns:
            Ответ сервера или None в случае ошибки
        """
        if not self.socket:
            logger.error("Not connected to server")
            return None
        
        try:
            # Отправляем сообщение
            self.socket.send(message.encode('utf-8'))
            logger.info(f"📤 Sent: {message}")
            
            # Получаем ответ
            response = self.socket.recv(1024).decode('utf-8').strip()
            logger.info(f"📩 Received: {response}")
            
            return response
        except socket.error as e:
            logger.error(f"❌ Communication error: {e}")
            return None
    
    def close(self):
        """Закрытие соединения"""
        if self.socket:
            self.socket.close()
            self.socket = None
            logger.info("🔒 Connection closed")
    
    def __enter__(self):
        """Контекстный менеджер для использования with"""
        self.connect()
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        """Выход из контекстного менеджера"""
        self.close()


def main():
    """Пример использования клиента"""
    # Использование с контекстным менеджером
    with GoTCPClient() as client:
        if client.socket:
            # Отправляем несколько сообщений
            messages = ["ping", "hello from python", "test message"]
            
            for msg in messages:
                response = client.send_message(msg)
                if response:
                    print(f"Server response: {response}")
                else:
                    print(f"Failed to send: {msg}")
        else:
            print("Failed to connect to server")


if __name__ == "__main__":
    main()
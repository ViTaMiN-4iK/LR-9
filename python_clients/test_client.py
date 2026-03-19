#!/usr/bin/env python3
"""
Тесты для Python клиента Go TCP сервера
Используем моки, чтобы не зависеть от реального сервера
"""

import unittest
from unittest.mock import Mock, patch, MagicMock
import socket
from go_tcp_client import GoTCPClient


class TestGoTCPClient(unittest.TestCase):
    """Тесты для GoTCPClient"""
    
    def setUp(self):
        """Подготовка перед каждым тестом"""
        self.client = GoTCPClient('localhost', 8080, timeout=1)
    
    def tearDown(self):
        """Очистка после каждого теста"""
        self.client.close()
    
    @patch('socket.socket')
    def test_connect_success(self, mock_socket):
        """Тест успешного подключения"""
        # Настраиваем мок
        mock_socket_instance = MagicMock()
        mock_socket.return_value = mock_socket_instance
        
        # Вызываем метод
        result = self.client.connect()
        
        # Проверяем
        self.assertTrue(result)
        mock_socket_instance.connect.assert_called_once_with(('localhost', 8080))
        mock_socket_instance.settimeout.assert_called_once_with(1)
    
    @patch('socket.socket')
    def test_connect_failure(self, mock_socket):
        """Тест ошибки подключения"""
        # Настраиваем мок на ошибку
        mock_socket_instance = MagicMock()
        mock_socket_instance.connect.side_effect = socket.error("Connection refused")
        mock_socket.return_value = mock_socket_instance
        
        # Вызываем метод
        result = self.client.connect()
        
        # Проверяем
        self.assertFalse(result)
        self.assertIsNone(self.client.socket)
    
    @patch('socket.socket')
    def test_send_message_success(self, mock_socket):
        """Тест успешной отправки сообщения"""
        # Настраиваем мок
        mock_socket_instance = MagicMock()
        mock_socket_instance.recv.return_value = b"Hello from Go\n"
        self.client.socket = mock_socket_instance
        
        # Вызываем метод
        response = self.client.send_message("ping")
        
        # Проверяем
        self.assertEqual(response, "Hello from Go")
        mock_socket_instance.send.assert_called_once_with(b"ping")
        mock_socket_instance.recv.assert_called_once_with(1024)
    
    @patch('socket.socket')
    def test_send_message_not_connected(self, mock_socket):
        """Тест отправки без подключения"""
        self.client.socket = None
        response = self.client.send_message("ping")
        self.assertIsNone(response)
    
    @patch('socket.socket')
    def test_send_message_error(self, mock_socket):
        """Тест ошибки при отправке"""
        # Настраиваем мок на ошибку
        mock_socket_instance = MagicMock()
        mock_socket_instance.send.side_effect = socket.error("Connection lost")
        self.client.socket = mock_socket_instance
        
        # Вызываем метод
        response = self.client.send_message("ping")
        
        # Проверяем
        self.assertIsNone(response)
    
    def test_context_manager(self):
        """Тест контекстного менеджера"""
        with patch('socket.socket') as mock_socket:
            mock_socket_instance = MagicMock()
            mock_socket.return_value = mock_socket_instance
            
            # Используем контекстный менеджер
            with GoTCPClient() as client:
                self.assertIsNotNone(client.socket)
            
            # Проверяем, что соединение закрыто
            mock_socket_instance.close.assert_called_once()
    
    @patch('socket.socket')
    def test_multiple_messages(self, mock_socket):
        """Тест нескольких сообщений подряд"""
        # Настраиваем мок
        mock_socket_instance = MagicMock()
        mock_socket_instance.recv.side_effect = [
            b"Hello from Go\n",
            b"Hello from Go\n",
            b"Hello from Go\n"
        ]
        self.client.socket = mock_socket_instance
        
        # Отправляем несколько сообщений
        messages = ["ping1", "ping2", "ping3"]
        for msg in messages:
            response = self.client.send_message(msg)
            self.assertEqual(response, "Hello from Go")
        
        # Проверяем, что send вызывался 3 раза
        self.assertEqual(mock_socket_instance.send.call_count, 3)
        self.assertEqual(mock_socket_instance.recv.call_count, 3)


if __name__ == '__main__':
    unittest.main()
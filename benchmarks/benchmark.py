#!/usr/bin/env python3
"""
Сравнительный анализ производительности:
- Чистый Python
- Python + Rust (PyO3)
- Python + Go (TCP-сервер)
"""

import time
import statistics
import socket
import json
import sys
import os
from typing import List, Callable
import matplotlib.pyplot as plt
import numpy as np

# Добавляем путь для импорта наших модулей
sys.path.insert(0, os.path.abspath('.'))

try:
    from rust_lib import Analyzer
    RUST_AVAILABLE = True
except ImportError:
    RUST_AVAILABLE = False
    print("⚠️  Rust module not available")

# ============= ЧИСТЫЙ PYTHON =============
class PythonAnalyzer:
    """Чистая Python-реализация для сравнения"""
    
    def __init__(self, multiplier: float):
        self.multiplier = multiplier
    
    def get_multiplier(self) -> float:
        return self.multiplier
    
    def process_data(self, data: List[int]) -> List[int]:
        return [int(x * self.multiplier) for x in data]
    
    def process_data_f64(self, data: List[float]) -> List[float]:
        return [x * self.multiplier for x in data]

# ============= GO КЛИЕНТ =============
class GoAnalyzer:
    """Клиент для Go-сервера"""
    
    def __init__(self, host='localhost', port=8080):
        self.host = host
        self.port = port
        self.multiplier = 2.0  # Заглушка
    
    def process_data(self, data: List[int]) -> List[int]:
        """Отправляет данные в Go-сервер и получает результат"""
        try:
            with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
                sock.connect((self.host, self.port))
                
                # Упаковываем данные в простой протокол
                message = f"process:{','.join(map(str, data))}\n"
                sock.send(message.encode())
                
                # Получаем ответ
                response = sock.recv(65536).decode().strip()
                
                # Распаковываем результат
                if response.startswith("result:"):
                    return list(map(int, response[7:].split(',')))
                return []
        except Exception as e:
            print(f"⚠️  Go connection error: {e}")
            return []

# ============= БЕНЧМАРКИ =============
def benchmark_function(func: Callable, data: List, iterations: int = 10) -> dict:
    """Замеряет время выполнения функции"""
    times = []
    
    for i in range(iterations):
        start = time.perf_counter()
        result = func(data)
        end = time.perf_counter()
        times.append(end - start)
        
        # Проверяем корректность результата (только первый раз)
        if i == 0:
            expected = [x * 2 for x in data[:5]]  # multiplier = 2
            actual = result[:5]
            assert actual == expected, f"Wrong result: {actual} != {expected}"
    
    return {
        'min': min(times),
        'max': max(times),
        'mean': statistics.mean(times),
        'median': statistics.median(times),
        'stdev': statistics.stdev(times) if len(times) > 1 else 0,
        'total_items': len(data),
        'items_per_second': len(data) / statistics.mean(times)
    }

def run_benchmarks(sizes: List[int] = [1000, 10000, 100000, 500000, 1000000]):
    """Запускает бенчмарки для разных размеров данных"""
    
    results = {
        'python': {'sizes': [], 'times': [], 'items_per_sec': []},
        'rust': {'sizes': [], 'times': [], 'items_per_sec': []},
        'go': {'sizes': [], 'times': [], 'items_per_sec': []}
    }
    
    # Создаём экземпляры
    py_analyzer = PythonAnalyzer(2.0)
    rs_analyzer = Analyzer(2.0) if RUST_AVAILABLE else None
    go_analyzer = GoAnalyzer()
    
    print("\n🔬 Запуск бенчмарков...")
    print("=" * 60)
    
    for size in sizes:
        print(f"\n📊 Размер данных: {size:,} элементов")
        data = list(range(size))
        
        # Python
        print("  🐍 Python...", end="", flush=True)
        py_results = benchmark_function(
            lambda d: py_analyzer.process_data(d), 
            data
        )
        results['python']['sizes'].append(size)
        results['python']['times'].append(py_results['mean'])
        results['python']['items_per_sec'].append(py_results['items_per_second'])
        print(f" {py_results['mean']*1000:.2f} ms ({py_results['items_per_second']:,.0f} items/sec)")
        
        # Rust
        if rs_analyzer:
            print("  🦀 Rust...", end="", flush=True)
            rs_results = benchmark_function(
                lambda d: rs_analyzer.process_data(d),
                data
            )
            results['rust']['sizes'].append(size)
            results['rust']['times'].append(rs_results['mean'])
            results['rust']['items_per_sec'].append(rs_results['items_per_second'])
            print(f" {rs_results['mean']*1000:.2f} ms ({rs_results['items_per_second']:,.0f} items/sec)")
        
        # Go (если сервер запущен)
        try:
            print("  🦫 Go...", end="", flush=True)
            go_results = benchmark_function(
                lambda d: go_analyzer.process_data(d),
                data,
                iterations=3  # Меньше итераций для Go
            )
            results['go']['sizes'].append(size)
            results['go']['times'].append(go_results['mean'])
            results['go']['items_per_sec'].append(go_results['items_per_second'])
            print(f" {go_results['mean']*1000:.2f} ms ({go_results['items_per_second']:,.0f} items/sec)")
        except:
            print("  ⚠️  Go server not available")
    
    return results

def plot_results(results: dict):
    """Строит графики результатов"""
    
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 5))
    
    # График времени выполнения
    ax1.set_title('Время выполнения (меньше = лучше)')
    ax1.set_xlabel('Размер данных (элементов)')
    ax1.set_ylabel('Время (мс)')
    ax1.set_xscale('log')
    ax1.set_yscale('log')
    
    colors = {'python': 'blue', 'rust': 'orange', 'go': 'green'}
    labels = {'python': 'Python', 'rust': 'Python+Rust', 'go': 'Python+Go'}
    
    for impl in ['python', 'rust', 'go']:
        if results[impl]['sizes']:
            ax1.plot(
                results[impl]['sizes'], 
                [t*1000 for t in results[impl]['times']], 
                marker='o', 
                color=colors[impl],
                label=labels[impl]
            )
    
    ax1.legend()
    ax1.grid(True, alpha=0.3)
    
    # График производительности (элементов в секунду)
    ax2.set_title('Производительность (больше = лучше)')
    ax2.set_xlabel('Размер данных (элементов)')
    ax2.set_ylabel('Элементов в секунду')
    ax2.set_xscale('log')
    ax2.set_yscale('log')
    
    for impl in ['python', 'rust', 'go']:
        if results[impl]['sizes']:
            ax2.plot(
                results[impl]['sizes'], 
                results[impl]['items_per_sec'], 
                marker='s', 
                color=colors[impl],
                label=labels[impl]
            )
    
    ax2.legend()
    ax2.grid(True, alpha=0.3)
    
    plt.tight_layout()
    plt.savefig('benchmarks/benchmark_results.png', dpi=150)
    plt.show()
    
    print("\n📈 Графики сохранены в benchmarks/benchmark_results.png")

def print_summary(results: dict):
    """Выводит сводную таблицу результатов"""
    
    print("\n" + "="*80)
    print("📊 СВОДНАЯ ТАБЛИЦА ПРОИЗВОДИТЕЛЬНОСТИ")
    print("="*80)
    
    print(f"{'Размер':>12} | {'Python':>15} | {'Rust':>15} | {'Go':>15} | {'Rust vs Python':>15}")
    print("-"*80)
    
    for i, size in enumerate(results['python']['sizes']):
        py_time = results['python']['times'][i] * 1000
        py_items = results['python']['items_per_sec'][i]
        
        rust_time = 0
        rust_items = 0
        if results['rust']['sizes'] and i < len(results['rust']['times']):
            rust_time = results['rust']['times'][i] * 1000
            rust_items = results['rust']['items_per_sec'][i]
        
        go_time = 0
        go_items = 0
        if results['go']['sizes'] and i < len(results['go']['times']):
            go_time = results['go']['times'][i] * 1000
            go_items = results['go']['items_per_sec'][i]
        
        speedup = py_time / rust_time if rust_time > 0 else 0
        
        print(f"{size:12,} | {py_time:8.2f} ms | {rust_time:8.2f} ms | {go_time:8.2f} ms | {speedup:>14.2f}x")
    
    print("="*80)

def main():
    """Главная функция"""
    
    print("\n" + "="*60)
    print("🚀 БЕНЧМАРК: Python vs Rust vs Go")
    print("="*60)
    
    # Размеры данных для тестирования
    sizes = [1000, 10000, 100000, 500000, 1000000]
    
    # Запускаем бенчмарки
    results = run_benchmarks(sizes)
    
    # Выводим результаты
    print_summary(results)
    
    # Строим графики
    try:
        plot_results(results)
    except ImportError:
        print("\n⚠️  matplotlib не установлен. Установите: pip install matplotlib")
    
    # Сохраняем результаты в JSON
    with open('benchmarks/results.json', 'w') as f:
        json.dump(results, f, indent=2)
    print("\n📁 Результаты сохранены в benchmarks/results.json")

if __name__ == "__main__":
    main()
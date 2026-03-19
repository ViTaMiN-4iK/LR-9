#!/usr/bin/env python3
"""
Тесты для Rust модуля через Python
Запуск: pytest tests/test_rust_module.py -v
"""

import pytest
from rust_lib import Analyzer, sum_as_string

# Тесты для Analyzer
def test_analyzer_creation():
    """Тест создания Analyzer"""
    a = Analyzer(2.5)
    assert a.get_multiplier() == 2.5

def test_analyzer_zero():
    """Тест с нулевым множителем"""
    a = Analyzer(0.0)
    assert a.get_multiplier() == 0.0

def test_analyzer_negative():
    """Тест с отрицательным множителем"""
    a = Analyzer(-3.14)
    assert a.get_multiplier() == -3.14

def test_analyzer_multiple_instances():
    """Тест нескольких экземпляров"""
    a1 = Analyzer(1.5)
    a2 = Analyzer(2.5)
    assert a1.get_multiplier() == 1.5
    assert a2.get_multiplier() == 2.5

# Тесты для process_data
def test_process_data_basic():
    """Базовый тест process_data"""
    a = Analyzer(2.0)
    result = a.process_data([1, 2, 3, 4, 5])
    assert result == [2, 4, 6, 8, 10]

def test_process_data_float_multiplier():
    """Тест с дробным множителем"""
    a = Analyzer(1.5)
    result = a.process_data([2, 4, 6, 8])
    assert result == [3, 6, 9, 12]

def test_process_data_empty():
    """Тест с пустым списком"""
    a = Analyzer(2.0)
    result = a.process_data([])
    assert result == []

def test_process_data_negative():
    """Тест с отрицательным множителем"""
    a = Analyzer(-2.0)
    result = a.process_data([1, 2, 3])
    assert result == [-2, -4, -6]

def test_process_data_large():
    """Тест с большим списком"""
    a = Analyzer(2.0)
    input_data = list(range(1000))
    result = a.process_data(input_data)
    expected = [x * 2 for x in input_data]
    assert result == expected

def test_process_data_f64():
    """Тест версии для float (если доступна)"""
    a = Analyzer(2.0)
    if hasattr(a, 'process_data_f64'):
        result = a.process_data_f64([1.5, 2.5, 3.5])
        assert result == [3.0, 5.0, 7.0]

# Тесты для sum_as_string
def test_sum_as_string():
    """Тест функции суммирования"""
    assert sum_as_string(5, 3) == "8"
    assert sum_as_string(0, 0) == "0"
    assert sum_as_string(100, 200) == "300"

def test_sum_as_string_large():
    """Тест с большими числами"""
    assert sum_as_string(1000000, 2000000) == "3000000"

# Тест производительности (опционально)
def test_process_data_performance():
    """Простой тест производительности"""
    import time
    
    a = Analyzer(2.0)
    data = list(range(100000))  # 100k элементов
    
    start = time.time()
    result = a.process_data(data)
    elapsed = time.time() - start
    
    print(f"\n⏱️  Processed 100k elements in {elapsed:.4f} seconds")
    assert len(result) == len(data)
    assert result[0] == 0
    assert result[-1] == (len(data)-1) * 2

if __name__ == "__main__":
    pytest.main([__file__, "-v", "--capture=no"])
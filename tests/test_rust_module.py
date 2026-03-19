#!/usr/bin/env python3
"""
Тесты для Rust модуля через Python
Запуск: pytest tests/test_rust_module.py -v
"""

import pytest
from rust_lib import Analyzer, sum_as_string

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

def test_sum_as_string():
    """Тест функции суммирования"""
    assert sum_as_string(5, 3) == "8"
    assert sum_as_string(0, 0) == "0"
    assert sum_as_string(100, 200) == "300"

def test_sum_as_string_large():
    """Тест с большими числами"""
    assert sum_as_string(1000000, 2000000) == "3000000"

if __name__ == "__main__":
    pytest.main([__file__, "-v"])
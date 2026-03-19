use pyo3::prelude::*;
use pyo3::types::PyModule;

/// Простая структура для демонстрации экспорта в Python
/// Хранит коэффициент, который можно использовать для вычислений
#[pyclass]
pub struct Analyzer {
    multiplier: f64,
}

#[pymethods]
impl Analyzer {
    /// Конструктор структуры
    #[new]
    pub fn new(multiplier: f64) -> Self {
        Analyzer { multiplier }
    }
    
    /// Геттер для поля multiplier
    pub fn get_multiplier(&self) -> f64 {
        self.multiplier
    }
}

/// Функция для суммирования чисел (оставляем для примера)
#[pyfunction]
fn sum_as_string(a: usize, b: usize) -> PyResult<String> {
    Ok((a + b).to_string())
}

/// Модуль Python
#[pymodule]
fn rust_lib(m: &Bound<'_, PyModule>) -> PyResult<()> {
    m.add_function(wrap_pyfunction!(sum_as_string, m)?)?;
    m.add_class::<Analyzer>()?;
    Ok(())
}

// ============= ЮНИТ-ТЕСТЫ =============
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_analyzer_new() {
        let analyzer = Analyzer::new(2.5);
        assert_eq!(analyzer.multiplier, 2.5);
    }

    #[test]
    fn test_analyzer_get_multiplier() {
        let analyzer = Analyzer::new(3.14);
        assert_eq!(analyzer.get_multiplier(), 3.14);
    }

    #[test]
    fn test_analyzer_with_zero() {
        let analyzer = Analyzer::new(0.0);
        assert_eq!(analyzer.get_multiplier(), 0.0);
    }

    #[test]
    fn test_analyzer_with_negative() {
        let analyzer = Analyzer::new(-5.0);
        assert_eq!(analyzer.get_multiplier(), -5.0);
    }

    #[test]
    fn test_sum_as_string() {
        assert_eq!(sum_as_string(5, 3).unwrap(), "8");
        assert_eq!(sum_as_string(0, 0).unwrap(), "0");
        assert_eq!(sum_as_string(999, 1).unwrap(), "1000");
    }
}
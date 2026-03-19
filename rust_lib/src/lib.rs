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
    
    /// Обрабатывает список целых чисел, умножая каждое на multiplier
    /// Принимает Vec<i32>, возвращает Vec<i32>
    pub fn process_data(&self, data: Vec<i32>) -> Vec<i32> {
        data.iter()
            .map(|&x| (x as f64 * self.multiplier) as i32)
            .collect()
    }
    
    /// Версия для чисел с плавающей точкой (демонстрация перегрузки)
    pub fn process_data_f64(&self, data: Vec<f64>) -> Vec<f64> {
        data.iter()
            .map(|&x| x * self.multiplier)
            .collect()
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
    fn test_process_data_integer() {
        let analyzer = Analyzer::new(2.0);
        let input = vec![1, 2, 3, 4, 5];
        let expected = vec![2, 4, 6, 8, 10];
        assert_eq!(analyzer.process_data(input), expected);
    }

    #[test]
    fn test_process_data_float_multiplier() {
        let analyzer = Analyzer::new(1.5);
        let input = vec![2, 4, 6, 8];
        let expected = vec![3, 6, 9, 12]; // 2*1.5=3, 4*1.5=6, ...
        assert_eq!(analyzer.process_data(input), expected);
    }

    #[test]
    fn test_process_data_empty() {
        let analyzer = Analyzer::new(2.0);
        let input: Vec<i32> = vec![];
        let expected: Vec<i32> = vec![];
        assert_eq!(analyzer.process_data(input), expected);
    }

    #[test]
    fn test_process_data_negative() {
        let analyzer = Analyzer::new(-2.0);
        let input = vec![1, 2, 3];
        let expected = vec![-2, -4, -6];
        assert_eq!(analyzer.process_data(input), expected);
    }

    #[test]
    fn test_process_data_f64() {
        let analyzer = Analyzer::new(2.0);
        let input = vec![1.5, 2.5, 3.5];
        let expected = vec![3.0, 5.0, 7.0];
        assert_eq!(analyzer.process_data_f64(input), expected);
    }

    #[test]
    fn test_sum_as_string() {
        assert_eq!(sum_as_string(5, 3).unwrap(), "8");
        assert_eq!(sum_as_string(0, 0).unwrap(), "0");
        assert_eq!(sum_as_string(999, 1).unwrap(), "1000");
    }
}
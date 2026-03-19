// Интеграционные тесты для PyO3 модуля
// Используем актуальный API для PyO3 0.27

use pyo3::prelude::*;
use pyo3::types::PyModule;

#[test]
fn test_analyzer_from_python() -> PyResult<()> {
    // Инициализируем Python
    pyo3::prepare_freethreaded_python();
    
    Python::with_gil(|py| {
        // Импортируем наш модуль
        let module = PyModule::import(py, "rust_lib")?;
        
        // Получаем класс Analyzer
        let analyzer_class = module.getattr("Analyzer")?;
        
        // Создаем экземпляр
        let analyzer = analyzer_class.call1((2.5,))?;
        
        // Вызываем метод get_multiplier
        let multiplier: f64 = analyzer
            .getattr("get_multiplier")?
            .call0()?
            .extract()?;
        
        assert_eq!(multiplier, 2.5);
        Ok(())
    })
}

#[test]
fn test_sum_as_string_from_python() -> PyResult<()> {
    pyo3::prepare_freethreaded_python();
    
    Python::with_gil(|py| {
        let module = PyModule::import(py, "rust_lib")?;
        
        let result: String = module
            .getattr("sum_as_string")?
            .call1((5, 3))?
            .extract()?;
        
        assert_eq!(result, "8");
        Ok(())
    })
}

#[test]
fn test_analyzer_multiple_instances() -> PyResult<()> {
    pyo3::prepare_freethreaded_python();
    
    Python::with_gil(|py| {
        let module = PyModule::import(py, "rust_lib")?;
        let analyzer_class = module.getattr("Analyzer")?;
        
        let a1 = analyzer_class.call1((1.5,))?;
        let a2 = analyzer_class.call1((3.0,))?;
        
        let m1: f64 = a1.getattr("get_multiplier")?.call0()?.extract()?;
        let m2: f64 = a2.getattr("get_multiplier")?.call0()?.extract()?;
        
        assert_eq!(m1, 1.5);
        assert_eq!(m2, 3.0);
        Ok(())
    })
}
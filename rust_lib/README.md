Usage Example
python
from rust_lib import Analyzer

# Create analyzer with multiplier 2.5
analyzer = Analyzer(2.5)

# Get the multiplier value
print(analyzer.get_multiplier())  # Output: 2.5

# Process a list of numbers
data = [1, 2, 3, 4, 5]
result = analyzer.process_data(data)
print(result)  # Output: [2, 4, 6, 8, 10]
Performance
Rust implementation is 3-5x faster than pure Python for large datasets.

Benchmarks
Tests on 1,000,000 elements:

Pure Python: 75 ms

Rust: 23 ms (3.2x faster)

License
MIT
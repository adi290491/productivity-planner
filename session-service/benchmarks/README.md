# Performance Benchmarks

This directory contains performance test results for the session-service.

## Files

- `BASELINE_YYYYMMDD.txt` - Baseline performance metrics
- `load_test_YYYYMMDD_HHMMSS.txt` - Individual load test results
- `load_test_latest.txt` - Symlink to most recent load test
- `scalability_test_YYYYMMDD_HHMMSS.txt` - Scalability test results
- `scalability_test_latest.txt` - Symlink to most recent scalability test

## Running Tests
```bash
# Load test (50 concurrent users)
./load_test_with_metrics.sh

# Scalability test (10-500 concurrent users)
./scalability_test.sh
```

## Performance Targets (SLOs)

| Metric | Target | Current |
|--------|--------|---------|
| P95 Latency | < 50ms | 4-11ms ✅ |
| P99 Latency | < 100ms | 12ms ✅ |
| Throughput | > 100 req/s | 238 req/s ✅ |
| Error Rate | < 1% | 0% ✅ |
| Availability | > 99.9% | - |

## Baseline Performance

See `BASELINE_20260124.txt` for initial production-ready metrics.

## Historical Tracking

To compare performance over time:
```bash
# View latest results
cat load_test_latest.txt

# Compare to baseline
diff BASELINE_20260124.txt load_test_latest.txt
```

## CI/CD Integration

These benchmarks can be used in CI/CD to detect performance regressions:
```bash
# Example: Fail if P95 latency > 50ms
P95=$(grep "P95 Latency" load_test_latest.txt | awk '{print $3}' | sed 's/ms//')
if [ "$P95" -gt 50 ]; then
    echo "Performance regression detected!"
    exit 1
fi
```

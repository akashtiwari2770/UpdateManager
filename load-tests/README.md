# Load Testing with Artillery

This directory contains Artillery load test configurations for the Update Manager API.

## Prerequisites

Install Artillery:

```bash
npm install -g artillery@latest
```

Or using yarn:

```bash
yarn global add artillery
```

## Test Configurations

### 1. `artillery-config.yml` - Mixed Load Test
Comprehensive test with a mix of read and write operations:
- **Phases**: Warm-up, ramp-up, sustained load, spike test, cool down
- **Duration**: ~5 minutes
- **Scenarios**: Health checks, product operations, versions, notifications, audit logs, compatibility

### 2. `artillery-read-heavy.yml` - Read-Heavy Load Test
Focuses on read operations (GET requests):
- **Phases**: Warm-up, read-heavy load, cool down
- **Duration**: ~5 minutes
- **Scenarios**: Product listings, active products, notifications, audit logs
- **Use case**: Testing read performance and caching effectiveness

### 3. `artillery-write-heavy.yml` - Write-Heavy Load Test
Focuses on write operations (POST requests):
- **Phases**: Warm-up, write-heavy load, cool down
- **Duration**: ~3 minutes
- **Scenarios**: Creating products, versions, notifications, update detections, upgrade paths
- **Use case**: Testing write performance and database load

### 4. `artillery-spike-test.yml` - Spike Test
Tests system behavior under sudden traffic spikes:
- **Phases**: Normal load, sudden spike (200 req/s), recovery
- **Duration**: ~1.5 minutes
- **Scenarios**: Health checks and product listings
- **Use case**: Testing system resilience and auto-scaling

## Running Load Tests

### Using Makefile

```bash
# Run mixed load test
make load-test

# Run read-heavy load test
make load-test-read

# Run write-heavy load test
make load-test-write

# Run spike test
make load-test-spike
```

### Direct Artillery Commands

```bash
# Mixed load test
artillery run load-tests/artillery-config.yml

# Read-heavy test
artillery run load-tests/artillery-read-heavy.yml

# Write-heavy test
artillery run load-tests/artillery-write-heavy.yml

# Spike test
artillery run load-tests/artillery-spike-test.yml
```

### Custom Target

To test against a different server:

```bash
artillery run --target http://localhost:3000 load-tests/artillery-config.yml
```

## Understanding Results

Artillery provides detailed metrics including:

- **Request rate**: Requests per second
- **Response times**: Min, max, median, p95, p99
- **Status codes**: Distribution of HTTP status codes
- **Errors**: Error rate and types
- **Throughput**: Requests completed per second

### Example Output

```
Summary report @ 14:30:15(+0000) 2025-01-15
  Scenarios launched:  1000
  Scenarios completed: 998
  Requests completed:  5000
  Mean response/sec: 83.33
  Response time (msec):
    min: 12
    max: 450
    median: 45
    p95: 120
    p99: 250
  Scenario counts:
    Health Check: 100 (10%)
    List Products: 300 (30%)
    ...
  Codes:
    200: 4800
    201: 150
    404: 50
  Errors:
    ETIMEDOUT: 2
```

## Performance Targets

Based on the test configurations, here are suggested performance targets:

### Response Times
- **Health check**: < 50ms (p95)
- **GET operations**: < 200ms (p95)
- **POST operations**: < 500ms (p95)
- **Complex queries**: < 1000ms (p95)

### Error Rates
- **Target**: < 0.1% error rate
- **Acceptable**: < 1% error rate
- **Critical**: > 5% error rate

### Throughput
- **Normal load**: 50-100 req/s
- **Peak load**: 200+ req/s
- **Sustained**: 100+ req/s for 5+ minutes

## Customization

### Adjusting Load

Edit the `phases` section in any config file:

```yaml
phases:
  - duration: 60      # Duration in seconds
    arrivalRate: 20   # Requests per second
    rampTo: 50        # Ramp up to this rate
```

### Adding Scenarios

Add new scenarios to the `scenarios` section:

```yaml
scenarios:
  - name: "My Custom Test"
    weight: 10  # Relative weight (higher = more frequent)
    flow:
      - get:
          url: "/api/v1/my-endpoint"
          expect:
            - statusCode: 200
```

### Custom Headers

Modify the `defaults.headers` section:

```yaml
defaults:
  headers:
    Content-Type: "application/json"
    X-User-ID: "custom-user"
    Authorization: "Bearer token"
```

## Best Practices

1. **Start with low load**: Begin with small arrival rates and gradually increase
2. **Monitor resources**: Watch CPU, memory, and database during tests
3. **Test in isolation**: Run tests against dedicated test environments
4. **Clean test data**: Ensure test data doesn't interfere with results
5. **Baseline first**: Establish baseline performance before optimization
6. **Iterate**: Run tests multiple times to identify patterns

## Troubleshooting

### Connection Errors
- Verify the server is running: `curl http://localhost:8080/health`
- Check firewall/network settings
- Ensure MongoDB is running and accessible

### High Error Rates
- Check server logs for errors
- Verify database connection
- Monitor resource usage (CPU, memory, disk)

### Slow Response Times
- Check database query performance
- Review server resource usage
- Consider database indexing
- Check for connection pool exhaustion

## Continuous Integration

You can integrate Artillery tests into CI/CD:

```yaml
# Example GitHub Actions
- name: Run Load Tests
  run: |
    npm install -g artillery
    make load-test-read
```

## Additional Resources

- [Artillery Documentation](https://www.artillery.io/docs)
- [Artillery GitHub](https://github.com/artilleryio/artillery)
- [Load Testing Best Practices](https://www.artillery.io/docs/guides/load-testing-best-practices)


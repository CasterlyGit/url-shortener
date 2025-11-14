# Load Testing Results

## Test Configuration
- **Tool**: k6
- **Scenario**: Gradual ramp-up from 50 to 200 users
- **Duration**: ~5 minutes
- **Endpoints Tested**: Redirect service (port 8081)

## Expected Metrics to Track
- Requests per second (RPS)
- 95th percentile response time
- Error rate
- Cache hit rate (inferred from response times)

## Performance Goals
- ✅ < 1% error rate
- ✅ p95 response time < 500ms
- ✅ Handle 100+ RPS consistently
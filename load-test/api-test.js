import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 10 },  // Ramp up to 10 users
    { duration: '1m', target: 10 },   // Stay at 10 users
    { duration: '30s', target: 20 },  // Ramp up to 20 users
    { duration: '1m', target: 20 },   // Stay at 20 users
    { duration: '30s', target: 0 },   // Ramp down to 0 users
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],   // Error rate < 1%
    http_req_duration: ['p(95)<1000'], // 95% of requests < 1000ms
  },
};

const API_BASE = 'http://localhost:8080';

export default function () {
  // Test URL creation endpoint
  const payload = JSON.stringify({
    long_url: `https://example.com/loadtest/${Math.random()}`
  });
  
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  
  const res = http.post(`${API_BASE}/api/shorten`, payload, params);
  
  check(res, {
    'status is 201': (r) => r.status === 201,
    'response has short_url': (r) => JSON.parse(r.body).short_url !== undefined,
    'response time < 1000ms': (r) => r.timings.duration < 1000,
  });
  
  sleep(1); // Longer delay for write operations
}
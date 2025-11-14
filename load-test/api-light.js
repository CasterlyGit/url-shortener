import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 5 },   // Very conservative
    { duration: '1m', target: 5 },    
    { duration: '30s', target: 10 },  
    { duration: '1m', target: 10 },   
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],   
    http_req_duration: ['p(95)<1000'], 
  },
};

const API_BASE = 'http://localhost:8080';

export default function () {
  const payload = JSON.stringify({
    long_url: `https://example.com/test/${__VU}-${Date.now()}`
  });
  
  const params = {
    headers: { 'Content-Type': 'application/json' },
    timeout: '5s',
  };
  
  const res = http.post(`${API_BASE}/api/shorten`, payload, params);
  
  check(res, {
    'status is 201': (r) => r.status === 201,
    'response time < 1s': (r) => r.timings.duration < 1000,
  });
  
  sleep(1); // Conservative delay
}
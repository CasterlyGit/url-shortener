import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 50 },  // Ramp up to 50 users
    { duration: '1m', target: 50 },   // Stay at 50 users
    { duration: '30s', target: 100 }, // Ramp up to 100 users  
    { duration: '1m', target: 100 },  // Stay at 100 users
    { duration: '30s', target: 200 }, // Ramp up to 200 users
    { duration: '1m', target: 200 },  // Stay at 200 users
    { duration: '30s', target: 0 },   // Ramp down to 0 users
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],   // Error rate < 1%
    http_req_duration: ['p(95)<500'], // 95% of requests < 500ms
  },
};

// First, let's create some short URLs to test with
const API_BASE = 'http://localhost:8080';
const REDIRECT_BASE = 'http://localhost:8081';

export function setup() {
  // Create test URLs before the load test
  const urls = [];
  for (let i = 0; i < 100; i++) {
    const payload = JSON.stringify({
      long_url: `https://example.com/test/${i}?loadtest=true`
    });
    
    const params = {
      headers: {
        'Content-Type': 'application/json',
      },
    };
    
    const res = http.post(`${API_BASE}/api/shorten`, payload, params);
    
    if (res.status === 201) {
      const shortUrl = JSON.parse(res.body).short_url;
      const shortCode = shortUrl.split('/').pop();
      urls.push(shortCode);
    }
  }
  
  console.log(`Created ${urls.length} test URLs for load testing`);
  return { shortCodes: urls };
}

export default function (data) {
  // Pick a random short code from our test data
  const shortCode = data.shortCodes[Math.floor(Math.random() * data.shortCodes.length)];
  
  // Test the redirect endpoint (this is our performance-critical path)
  const res = http.get(`${REDIRECT_BASE}/${shortCode}`, {
    redirects: 0, // Don't follow redirects, just test our service response
  });
  
  // Check if response is correct (should be 302 redirect)
  check(res, {
    'status is 302': (r) => r.status === 302,
    'has location header': (r) => r.headers.Location !== undefined,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });
  
  sleep(0.1); // Small delay between requests
}
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 1, // 1 virtual user
  duration: '30s',
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% requests harus < 500ms
    http_req_failed: ['rate<0.01'],   // Error rate < 1%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Test home endpoint
  let res = http.get(`${BASE_URL}/`);
  check(res, {
    'home status is 200': (r) => r.status === 200,
    'home has welcome message': (r) => r.json('message').includes('Welcome'),
  });

  // Test health endpoint
  res = http.get(`${BASE_URL}/health`);
  check(res, {
    'health status is 200': (r) => r.status === 200,
    'database is connected': (r) => r.json('database') === 'connected',
  });

  // Test get all products
  res = http.get(`${BASE_URL}/api/produk`);
  check(res, {
    'products status is 200': (r) => r.status === 200,
    'products response is array': (r) => Array.isArray(r.json('data')),
  });

  // Test get all categories
  res = http.get(`${BASE_URL}/categories`);
  check(res, {
    'categories status is 200': (r) => r.status === 200,
    'categories response has success': (r) => r.json('success') === true,
  });

  sleep(1);
}
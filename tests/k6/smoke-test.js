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

// Test data - sesuai dengan models/product.go dan models/category.go
let createdProductId;
let createdCategoryId;

export default function () {
  // ===== HOME & HEALTH ENDPOINTS =====
  let res = http.get(`${BASE_URL}/`);
  check(res, {
    'home status is 200': (r) => r.status === 200,
    'home has welcome message': (r) => r.json('message').includes('Welcome'),
  });

  sleep(0.5);

  res = http.get(`${BASE_URL}/health`);
  check(res, {
    'health status is 200': (r) => r.status === 200,
    'database is connected': (r) => r.json('database') === 'connected',
  });

  sleep(0.5);

  // ===== CATEGORY ENDPOINTS =====
  
  // GET all categories
  res = http.get(`${BASE_URL}/categories`);
  check(res, {
    'GET categories status is 200': (r) => r.status === 200,
    'categories response is array': (r) => Array.isArray(r.body) || r.body !== null,
  });

  sleep(0.5);

  // POST create category - sesuai models.Category {ID, Name, Description}
  const categoryPayload = JSON.stringify({
    name: `Test Category ${Date.now()}`,
    description: 'K6 Smoke Test Category'
  });

  res = http.post(`${BASE_URL}/categories`, categoryPayload, {
    headers: { 'Content-Type': 'application/json' },
  });
  check(res, {
    'POST category status is 201': (r) => r.status === 201,
    'POST category returns name': (r) => r.json('name') !== undefined,
  });
  
  if (res.status === 201) {
    createdCategoryId = res.json('id');
  }

  sleep(0.5);

  // GET category by ID
  if (createdCategoryId) {
    res = http.get(`${BASE_URL}/categories/${createdCategoryId}`);
    check(res, {
      'GET category by ID status is 200': (r) => r.status === 200,
      'GET category returns id': (r) => r.json('id') === createdCategoryId,
      'GET category returns name': (r) => r.json('name') !== undefined,
      'GET category returns description': (r) => r.json('description') !== undefined,
    });
  }

  sleep(0.5);

  // PUT update category
  if (createdCategoryId) {
    const updateCategoryPayload = JSON.stringify({
      name: `Updated Category ${Date.now()}`,
      description: 'Updated by K6'
    });

    res = http.put(`${BASE_URL}/categories/${createdCategoryId}`, updateCategoryPayload, {
      headers: { 'Content-Type': 'application/json' },
    });
    check(res, {
      'PUT category status is 200': (r) => r.status === 200,
      'PUT category returns updated name': (r) => r.json('name').includes('Updated'),
    });
  }

  sleep(0.5);

  // ===== PRODUCT ENDPOINTS =====
  
  // GET all products
  res = http.get(`${BASE_URL}/api/produk`);
  check(res, {
    'GET products status is 200': (r) => r.status === 200,
    'products response is array': (r) => Array.isArray(r.body) || r.body !== null,
  });

  sleep(0.5);

  // POST create product - sesuai models.Product {ID, Name, Price, Stock, CategoryID, CategoryName}
  const productPayload = JSON.stringify({
    name: `Test Product ${Date.now()}`,
    price: 10000,
    stock: 50,
    category_id: createdCategoryId || 1
  });

  res = http.post(`${BASE_URL}/api/produk`, productPayload, {
    headers: { 'Content-Type': 'application/json' },
  });
  check(res, {
    'POST product status is 201': (r) => r.status === 201,
    'POST product returns name': (r) => r.json('name') !== undefined,
    'POST product returns price': (r) => r.json('price') !== undefined,
    'POST product returns stock': (r) => r.json('stock') !== undefined,
  });

  if (res.status === 201) {
    createdProductId = res.json('id');
  }

  sleep(0.5);

  // GET product by ID (dengan category_name jika ada JOIN)
  if (createdProductId) {
    res = http.get(`${BASE_URL}/api/produk/${createdProductId}`);
    check(res, {
      'GET product by ID status is 200': (r) => r.status === 200,
      'GET product returns id': (r) => r.json('id') === createdProductId,
      'GET product returns name': (r) => r.json('name') !== undefined,
      'GET product returns price': (r) => r.json('price') !== undefined,
      'GET product returns stock': (r) => r.json('stock') !== undefined,
      'GET product returns category_id': (r) => r.json('category_id') !== undefined,
    });
  }

  sleep(0.5);

  // PUT update product
  if (createdProductId) {
    const updateProductPayload = JSON.stringify({
      name: `Updated Product ${Date.now()}`,
      price: 15000,
      stock: 100,
      category_id: createdCategoryId || 1
    });

    res = http.put(`${BASE_URL}/api/produk/${createdProductId}`, updateProductPayload, {
      headers: { 'Content-Type': 'application/json' },
    });
    check(res, {
      'PUT product status is 200': (r) => r.status === 200,
      'PUT product returns updated name': (r) => r.json('name').includes('Updated'),
      'PUT product returns updated price': (r) => r.json('price') === 15000,
    });
  }

  sleep(0.5);

  // DELETE product
  if (createdProductId) {
    res = http.del(`${BASE_URL}/api/produk/${createdProductId}`);
    check(res, {
      'DELETE product status is 200': (r) => r.status === 200,
      'DELETE product returns message': (r) => r.json('message') !== undefined,
    });
  }

  sleep(0.5);

  // DELETE category
  if (createdCategoryId) {
    res = http.del(`${BASE_URL}/categories/${createdCategoryId}`);
    check(res, {
      'DELETE category status is 200': (r) => r.status === 200,
      'DELETE category returns message': (r) => r.json('message') !== undefined,
    });
  }

  sleep(1);
}

// Setup function - runs once before test
export function setup() {
  console.log('🚀 Starting smoke test...');
  console.log(`📡 Base URL: ${BASE_URL}`);
}

// Teardown function - runs once after test
export function teardown(data) {
  console.log('✅ Smoke test completed!');
}
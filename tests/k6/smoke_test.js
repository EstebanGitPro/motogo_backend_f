import http from 'k6/http';
import { check, sleep } from 'k6';

// ============================================
// Smoke Test — Sanity check (1 VU, 10s)
// ============================================
// Uso: k6 run tests/k6/smoke_test.js
// Con URL custom: k6 run --env BASE_URL=http://tu-servidor:8085 tests/k6/smoke_test.js

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8085';

export const options = {
  vus: 1,
  duration: '10s',
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
};

export default function smokeTest() {
  const res = http.get(`${BASE_URL}/health`);
  check(res, {
    'health status 200': (r) => r.status === 200,
  });
  sleep(1);
}

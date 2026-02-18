import http from 'k6/http';
import { check, sleep } from 'k6';

// Stress Test: Find system limits
export const options = {
  stages: [
    { duration: '2m', target: 50 },   // Ramp-up to 50 users
    { duration: '3m', target: 50 },   // Stay at 50 users
    { duration: '2m', target: 100 },  // Ramp-up to 100 users
    { duration: '3m', target: 100 },  // Stay at 100 users
    { duration: '2m', target: 200 },  // Ramp-up to 200 users (stress)
    { duration: '3m', target: 200 },  // Stay at 200 users
    { duration: '2m', target: 0 },    // Ramp-down
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'], // 95% requests under 1s during stress
    http_req_failed: ['rate<0.10'],    // Less than 10% errors
  },
};

const BASE_URL = 'http://host.docker.internal:8085/motogo/api/v1';

export default function stressTest() {
  const responses = http.batch([
    ['GET', `${BASE_URL}/messages`, null, { tags: { name: 'GetMessages' } }],
    ['GET', `${BASE_URL}/messages/1`, null, { tags: { name: 'GetMessageById' } }],
    ['GET', `${BASE_URL}/messages?module=users`, null, { tags: { name: 'FilterMessages' } }],
  ]);

  responses.forEach((res) => {
    check(res, {
      'status is 2xx or 404': (r) => (r.status >= 200 && r.status < 300) || r.status === 404,
    });
  });

  sleep(1);
}

export function setup() {
  console.log('Running Stress Test - Finding system limits...');
}

export function teardown(data) {
  console.log('Stress Test completed!');
}

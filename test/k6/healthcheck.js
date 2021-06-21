import { check, group } from "k6";
import http from "k6/http";
import { URL } from 'https://jslib.k6.io/url/1.0.0/index.js';

export default function () {
    group('Healthcheck endpoint', () => {
        const baseURL = __ENV.EVENT_URL || 'http://localhost:8080';
        const url = new URL(`${baseURL}/healthcheck`);
        group(`GET ${url}`, () => {
            const res = http.get(url.toString());
            check(res, {
                'status is 200': (r) => r.status === 200,
                'response is OK': (r) => r.body === 'OK'
            });
        });
    });
};

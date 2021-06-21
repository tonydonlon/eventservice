import { check, group } from "k6";
import http from "k6/http";
import { URL } from 'https://jslib.k6.io/url/1.0.0/index.js';

export default () => {
    const sessionId = __ENV.SESSION_ID || '8e0b8eea-24f4-4cb6-b525-e8835f824768';
    const baseURL = __ENV.EVENT_URL || 'http://localhost:8080';
    const url = new URL(`${baseURL}/session/${sessionId}`);

    group('Get sessions endpoints', () => {
        group(`GET ${url} `, () => {
            const res = http.get(`${url}`);
            check(res, {
                'status is 200': (r) => {
                    return r.status === 200;
                },
                'has session info': (r) => {
                    const sess = r.json();
                    return sess.type === 'SESSION'
                        && !isNaN(sess.start)
                        && !isNaN(sess.end)
                },
                'has events': (r) => {
                    const events = r.json('children');
                    let eventsValidated = false;
                    const eventValidator = (evt) => {
                        return evt.type === 'EVENT'
                            && evt.hasOwnProperty('type')
                            && evt.hasOwnProperty('timestamp')
                            && evt.hasOwnProperty('name');
                    };
                    if (Array.isArray(events)) {
                        eventsValidated = events.every(eventValidator);
                    }

                    return eventsValidated;
                }
            });
        });
    });
};

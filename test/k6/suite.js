import healthcheck from './healthcheck.js';
import getsession from './getsession.js';
import websocket from './websocket.js';

export const options = {
    thresholds: {
        // add "abortOnFail: true" to exit immediately
        failedTestCases: [{ threshold: 'count==0' }],
    }
};

export default function () {
    healthcheck();
    getsession();
    websocket();
};

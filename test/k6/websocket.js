import ws from 'k6/ws';
import { check, group } from 'k6';
import uuid from './lib/uuid.js';
import { createMessageStream } from './lib/message.js';
const sessionId = __ENV.SESSION_ID || uuid.v4();

export default function () {
  group('Websocket endpoints', () => {
    const url = `ws://127.0.0.1:8080/ws?sessionId=${sessionId}`;
    const params = {};
    group(`Sending message streams to ${url}`, () => {
      const res = ws.connect(url, params, (socket) => {
        socket.on('open', () => {
          console.log('connected');

          socket.setInterval(() => {
            socket.ping();
            console.log(`Sending ping every 1sec`);
          }, 1000);

          // immediately send first set of events
          socket.send(JSON.stringify(createMessageStream(sessionId, true, false, 3)));

          socket.setTimeout(() => {
            socket.send(JSON.stringify(createMessageStream(sessionId, false, true, 2)));
            // TODO terminate session client or server side?
            //socket.close();
          }, 1000);

        });
        socket.on('message', (data) => {
          // TODO remove this...the WS implementation is unidirectional for now
          console.log('Message received: ', data);
        });
        socket.on('close', () => console.log('disconnected'));
      });

      check(res, { 'status is 101': (r) => r && r.status === 101 });
    });
  });
}

// TODO hit REST endpoint to verify those records were created
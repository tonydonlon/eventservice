import ws from 'k6/ws';
import { check, sleep } from 'k6';

export default function () {
  const url = 'ws://127.0.0.1:8080/event';
  const params = {};

  const res = ws.connect(url, params, (socket) => {
    socket.on('open', () => {
      console.log('connected');

      socket.setInterval(() => {
        socket.ping();
        console.log(`Sending ping every 1sec`);
      }, 1000);

      // immediately send first set of events
      socket.send(JSON.stringify(firstMsg));

      // wait to send second set of messages
      socket.setTimeout(() => {
        socket.send(JSON.stringify(secondMsg));
        // TODO terminate session client or server side?
        //socket.close();
      }, 3000);

    });
    socket.on('message', (data) => {
      // TODO remove this...the WS implementation is unidirectional for now
      console.log('Message received: ', data);
    });
    socket.on('close', () => console.log('disconnected'));
  });

  check(res, { 'status is 101': (r) => r && r.status === 101 });
}

const firstMsg = [
  {
    time: Date.now(),
    type: 'SESSION_START',
    session_id: '4cc700ae-4510-43f2-b939-1cd18dbf56a3',
  },
  {
    time: Date.now() + 1,
    type: 'EVENT',
    name: 'cart_load'
  },
  {
    time: Date.now() + 2,
    type: 'EVENT',
    name: 'cart_show'
  }
];

const secondMsg = [
  {
    time: Date.now(),
    type: 'EVENT',
    name: 'cart_checkout'
  },
  {
    time: Date.now(),
    type: 'SESSION_END',
    session_id: '4cc700ae-4510-43f2-b939-1cd18dbf56a3',
  },
];
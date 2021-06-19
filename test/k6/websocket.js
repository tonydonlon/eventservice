import ws from 'k6/ws';
import { check } from 'k6';

export default function () {
  const url = 'ws://127.0.0.1:8080/event';
  const params = {};

  const res = ws.connect(url, params, function (socket) {
    socket.on('open', () => {
      console.log('connected');

      socket.setInterval(function timeout() {
        socket.ping();
        const msg = Date.now();
        console.log(`Pinging every 1sec ${msg}`);
        socket.send(msg);

      }, 1000);
    });
    socket.on('message', (data) => {
      console.log('Message received: ', data);
    });
    socket.on('close', () => console.log('disconnected'));
  });

  check(res, { 'status is 101': (r) => r && r.status === 101 });
}

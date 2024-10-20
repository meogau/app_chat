import { check } from 'k6';
import ws from 'k6/ws';
import { sleep } from 'k6';

export let options = {
    stages: [
        { duration: '30s', target: 100 }, // ramp-up to 100 VUs in 30 seconds
        { duration: '1m', target: 500 },  // ramp-up to 500 VUs in 1 minute
        { duration: '2m', target: 1000 }, // ramp-up to 1000 VUs in 2 minutes
        { duration: '2m', target: 0 },    // ramp-down after test
    ],
};

export default function () {
    let url = 'ws://localhost:8888';

    let response = ws.connect(url, {}, function (socket) {
        socket.on('open', function () {
            console.log('Connected');
            socket.send('Hello from K6!');
        });

        socket.on('message', function (msg) {
            console.log(`Received message: ${msg}`);
        });

        socket.on('close', function () {
            console.log('Disconnected');
        });

        socket.on('error', function (e) {
            console.log('Error: ' + e.error());
        });

        // Keep the connection open for 30 seconds
        sleep(30);

        socket.close();
    });

    check(response, {
        'Connected successfully': (res) => res && res.status === 101,
    });
}

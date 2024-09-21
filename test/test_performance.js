import {check, sleep} from 'k6';
import ws from 'k6/ws';

export const options = {
    vus: 5000,
    duration: '1m',
};

export default function () {
    const url = 'ws://localhost:8888';
    const params = {tags: {my_tag: 'websocket'}};

    const response = ws.connect(url, params, function (socket) {
        socket.on('open', function () {
            console.log('WebSocket connection opened');
            socket.send('Hello WebSocket!');

            // Simulate messages being sent periodically
            socket.setInterval(function () {
                socket.send('ping');
            }, 1000);
        });

        socket.on('message', function (message) {
            console.log(`Received message: ${message}`);
        });

        socket.on('close', function () {
            console.log('WebSocket connection closed');
        });

        socket.on('error', function (error) {
            console.log(`WebSocket error: ${error}`);
        });

        // Keep the connection open for 10 seconds
        sleep(10);
    });

    // Check if the connection was successful
    check(response, {'WebSocket status is 101': (r) => r && r.status === 101});
}

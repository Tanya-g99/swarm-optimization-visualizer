const WebSocket = require('ws');
const { Worker } = require('worker_threads');
const http = require('http');

const server = http.createServer();
const wss = new WebSocket.Server({ noServer: true });

function createWorker(ws, workerFile, data) {
    const worker = new Worker(workerFile, { workerData: data });

    worker.on('message', (msg) => {
        ws.send(JSON.stringify(msg));
    });

    // worker.on('exit', () => {
    //     console.log(`Worker ${workerFile} finished execution`);
    // });

    ws.on('close', () => {
        worker.terminate();
        console.log(`Connection closed, worker ${workerFile} terminated`);
    });
}

server.on('upgrade', (request, socket, head) => {
    // console.log("🔄 Обновление соединения:", request.url);

    if (!['/ws/AFSA', '/ws/timer', '/ws/timer-task', '/ws/timer-step'].includes(request.url)) {
        socket.destroy();
        return;
    }

    // socket.setNoDelay(true);
    wss.handleUpgrade(request, socket, head, (ws) => {
        wss.emit('connection', ws, request);
    });
});

wss.on('connection', (ws, req) => {
    // console.log("🟢 Клиент подключился:", req.url);

    ws.on('message', (message) => {
        const data = JSON.parse(message);
        // console.log(data);

        switch (req.url) {
            case '/ws/AFSA':
                createWorker(ws, './workers/afsaWorker.js',data);
                break;
            default:
                ws.send(JSON.stringify({ error: "❌ Неизвестный маршрут" }));
                ws.close();
        }
    });

    // ws.on('close', () => {
    //     console.log('🔴 Клиент отключился:', req.url);
    // });
});

server.listen(8080, () => console.log("🚀 WebSocket сервер запущен на ws://localhost:8080"));

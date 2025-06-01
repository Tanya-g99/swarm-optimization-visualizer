const WebSocket = require('ws');

const WS_URL = 'ws://localhost:8080/ws/AFSA';
const CONCURRENT_REQUESTS = 300;
const NUM_RUNS = 5;

function generateRequestData() {
  return JSON.stringify({
    targetFunction: 'sin(x) + cos(y)',
    maxIter: 100,
    bounds: [[-10, 10], [-10, 10]],
    populationSize: 100,
    eta: 0.01,
    maxTryNum: 10,
    visual: [1, 10],
    teta: 0.5,
    seed: 1234
  });
}

function sendRequests(concurrentRequests) {
  return new Promise((resolve) => {
    const startTime = Date.now();
    let completedRequests = 0;

    const sockets = [];

    for (let i = 0; i < concurrentRequests; i++) {
      const ws = new WebSocket(WS_URL);

      ws.on('open', () => {
        ws.send(generateRequestData());
      });

      ws.on('message', () => {
        completedRequests++;
        ws.close();
        if (completedRequests === concurrentRequests) {
          const duration = Date.now() - startTime;
          resolve(duration);
        }
      });

      ws.on('error', (error) => {
        console.error('WebSocket error: ', error);
        ws.close();
      });

      sockets.push(ws);
    }
  });
}

async function runMultipleTests() {
  const durations = [];

  for (let i = 0; i < NUM_RUNS; i++) {
    console.log(`Run #${i + 1}...`);
    const duration = await sendRequests(CONCURRENT_REQUESTS);
    durations.push(duration);
    console.log(`Duration: ${duration} ms`);
  }

  const minTime = Math.min(...durations);
  console.log(`\nMinimum time over ${NUM_RUNS} runs: ${minTime} ms`);
}

runMultipleTests();

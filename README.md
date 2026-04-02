
# Temperature sensor dashboard

This repo contains the code for a dashboard application displaying the
Temperature and humidity in real time.

It is composed of several parts:

```mermaid
flowchart LR;
  Arduino sensor --> Go Worker;
  Go Worker --> Python FastAPI server;
  Python FastAPI server --> Vue.js dashboard;
```

## Workflow

### Data generation

The data is generated from an Arduino board, using a DHT11
Temperature and humidity sensor. The board generates readings every second
and writes them to the serial port.

### Reading Arduino data

The data generated from Arduino is read by a Go worker, that reads on the
same serial port, and inserts the data in a Redis stream.

### Serving Redis data

A Python FastAPI server then reads from the Redis stream, and exposes this data
via a WebSocket

### Displaying the data on a dashboard

The final part is the dashboard, built using Vue.js and chart.js. On mounting,
the dashboard connects to the FastAPI WebSocket. The displayed data is then
updated in real-time on the dashboard.

const { NodeSDK } = require('@opentelemetry/sdk-node');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-grpc');
const { OTLPMetricExporter } = require('@opentelemetry/exporter-metrics-otlp-grpc');
const { PeriodicExportingMetricReader } = require('@opentelemetry/sdk-metrics');
const { HttpInstrumentation } = require('@opentelemetry/instrumentation-http');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');

const sdk = new NodeSDK({
  traceExporter: new OTLPTraceExporter({
    url: 'http://localhost:4317', 
  }),
  metricReader: new PeriodicExportingMetricReader({
    exporter: new OTLPMetricExporter({
      url: 'http://127.0.0.1:4317',
    }),
    exportIntervalMillis: 1000, // Send data every 5 seconds
  }),
  instrumentations: [
    new HttpInstrumentation(),
    getNodeAutoInstrumentations(),
  ],
});

sdk.start();
console.log('OTel SDK: Tracing and Metrics Started');
# Event-Driven Trader

## Input
Place the historical OHLCV data in `data/input/data.csv`. Use `data.csv.example` template as a reference.


## Installation
```
make build
make up
```

## Trigger Market Data Generation
```
make run-market-data
```

## Empty Kafka topics
```
make kafka-empty-topics
```

## Events Inspection
```
make kafka-candle-closed-events
make kafka-input-completed-events
make kafka-marker-created-events
make kafka-run-completed-events
```

## Output
After run `data/result/chart.html` file will be created as output, containing candles and backtest strategy markers. Backtest metrics including (win rate, SL & TP count, ...) and strategy configs provided in html page title.


## Monitoring
```
make up-monitoring
```
Access http://localhost:9090 to execute PromQL queries, view service target health, and inspect recorded metrics or firing alerts.

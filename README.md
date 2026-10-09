# Event-Driven Trader

## Installation
```
make rebuild
```

## Empty Kafka topics
```
make kafka-empty-topics
```

## View events
```
make kafka-candle-closed-events
make kafka-input-completed-events
make kafka-marker-created-events
make kafka-run-completed-events
```

## Input and Output
Place the historical OHLCV data in `data/input/data.csv`. Use `data.csv.example` for the required column format.

After run `data/result/chart.html` file will be created as output, containing candles and backtest strategy markers. Backtest metrics including (win rate, SL & TP count, ...) and strategy configs provided in html page title.
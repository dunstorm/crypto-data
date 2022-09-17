# Crypto Data

CLI tool to download historical data from binance.

Install:

```bash
git clone github.com/dunstorm/crypto-data
cd crypto-data
go install .
```

Usage:

```bash
$ crypto-data --help
Usage: crypto-data [OPTIONS] COMMAND [ARGS]...
```

Example:

```
crypto-data download -t "BTCUSDT" -s "2020-01-01" -e "2022-09-17" -i 1h -o "btcusdt_1h.csv"
```
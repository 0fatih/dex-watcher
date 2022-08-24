# DEX Watcher

This is a tool for real-time dex transaction logging.

To use this program, you should have mongoDB on your localhost. Because, we are writing some information to the database. For example:

- Factories

  - Address

- Pairs

  - Address
  - Token0 Address
  - Token1 Address
  - Reserve0
  - Reserve1

- Tokens
  - Address
  - Name
  - Decimals

# Arguments:

```
  -initialize
    	Creates documents and writes initial data to database.
  -listen
    	Start listening to the blockchain.
  -pairs int
    	Pair amount per factory that you want to subscribe. If you want to subscribe to all pairs in a factory, enter 0. (default 5)
  -refresh
    	Fetchs latest prices and writes to database for every pair.
```

# Example Usages

For getting all pairs to the database in the factories list:

```
dex-watcher --pairs 0 --initialize
```

For getting only 5 pairs in all factories (i.e, if there are 5 factories, then you are going to scan 25 pairs):

```
dex-watcher --pairs 5 --initialize
```

To catch up reserves:

```
dex-watcher --refresh
```

If you want to listen pairs on your database (`dex-watcher/pairs`), you can use:

```
dex-watcher --listen
```

Initialize and start listening:

```
dex-watcher --initialize 100 --listen
```

# Configuration

## Factories

Currently in the `factoryList` (in `main.go`) there is only Uniswap V2 Factory's address. If you want to add/delete factories,
you only have to change this list.

## Provider

You have to specify your provider in the main.go (change `YOUR_PROVIDER_HERE` to your WEBSOCKET provider link).

# TODO

I'll (maybe) do these:

- Better configuration settings.
  - Maybe there can a be config file in `~/`.
  - Add flags to manage providers and factories.
- Calculate prices and write historically to database.
- Support other factoriy interfaces (currently only Uniswap V2 supported).

If I do these, I'll add new todos.

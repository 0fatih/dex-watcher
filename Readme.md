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

Set factories in `.env` file (if you don't have, create one) with the `FACTORY_ADDRESSES` variable.
This variable should be a string. If you want to enter more than one factories, enter with them a
`,` as a seperator. For example:

`FACTORY_ADDRESSES="0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f,0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"`

## Provider

You can set the provider in `.env` file (if you don't have, create one) with the `PROVIDER_URI` variable.

## Example `.env` File

```
PROVIDER_URI="YOUR_PROVIDER_HERE"

FACTORY_ADDRESSES="0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
```

# TODO

I'll (maybe) do these:

- Better configuration settings.
  - Add flags to manage providers and factories.
- Calculate prices and write historically to database.
- Support other factory interfaces (currently only Uniswap V2 supported).

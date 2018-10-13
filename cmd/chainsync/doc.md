# chainsync

### Prototype blockchain searchengine

Blockchain data is immutable, and most blockchain servers have only the most bare rudimentary functions for search, namely two indices one for disk block position and the second for hash to block height search.

The purpose of this server is to offload this work from the blockchain server and enable common searches of bitcoin style blockchain ledgers with very little work to implement new ones that implement the same set of JSONRPC protocols

It will also answer and pass through any RPC queries it does not implement on to the full node that is maintaining a current chain database, and pass back the responses.

Using the `codec` golang serialization library, the server will create a JSON rpc endpoint and a msgpack binary RPC endpoint. The latter is recommended as its serialization overhead is much lower.
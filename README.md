![https://gitlab.com/parallelcoin/duo/raw/master/assets/logo/logo256x256.png](https://gitlab.com/parallelcoin/duo/raw/master/assets/logo/logo256x256.png)

# Parallelcoin Proof of Work Cryptocurrency Daemon and Wallet

This is a port of the original Bitcoin-based C++ code that runs the Parallelcoin network, structurally identical but rewritten in Go. The Qt GUI wallet will be replaced with a Progressive Web App and a gRPC protocol will be added for faster RPC processing.

It will also add a Masternode system that distributes a reward randomly to any nodes serving up the blockchain and relaying data, and implements a messaging system like Bitmessage based on the wallet key cryptography. The web application will also be a server monitor, and block explorer in addition to regular wallet functions. Ecommerce needs secure communications, they should be a unified application.

The blockchain fundamentals will also be improved with the addition of merge mining to each algorithm for the coin with the same PoW and the most hashpower, in order to help get around the vulnerabilities of a small cryptocurrency to hashpower attacks. In line with this, there will also be the addition of more PoW algorithms, as many as possible, as well as reinforcing the robustness of the multi-algorithm difficulty regulation.

## Building

### Prerequisites

You need some version of BerkeleyDB C headers installed (and the library of course). This is to enable the import of `wallet.dat` files from the old parallelcoind.

Also needed is base58check:

    go get github.com/anaskhan96/base58check
    
For keeping secrets off disk and away from potential buffer exploits, we use memguard, and for that we want to check out the latest stable release tag:
    
    go get github.com/awnumar/memguard \
       && cd $GOPATH/src/github.com/awnumar/memguard \
       && git checkout v0.15.0

### Build it

    go get gitlab.com/parallelcoin/duo

### Install it

    go install gitlab.com/parallelcoin/duo

## Running

    duo

This will take care of creating a default data directory and configuration automatically.

## Progress status

See [Checklist](checklist.md)

Currently working on basic wallet I/O
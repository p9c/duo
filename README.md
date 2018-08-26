![https://gitlab.com/parallelcoin/duo/raw/master/assets/logo/logo256x256.png](https://gitlab.com/parallelcoin/duo/raw/master/assets/logo/logo256x256.png)
# Parallelcoin Daemon and Wallet

![https://goreportcard.com/report/gitlab.com/parallelcoin/duo](https://goreportcard.com/badge/gitlab.com/parallelcoin/duo)     ![https://godoc.org/gitlab.com/parallelcoin/duo?status.svg](https://godoc.org/gitlab.com/parallelcoin/duo)

This is a port of the original Bitcoin-based C++ code that runs the Parallelcoin network, structurally identical but rewritten in Go. The Qt GUI wallet will be replaced with a Progressive Web App and a gRPC protocol will be added for faster RPC processing.

It will also add a Masternode system that distributes a reward randomly to any nodes serving up the blockchain and relaying data, and implements a messaging system like Bitmessage based on the wallet key cryptography. The web application will also be a server monitor, and block explorer in addition to regular wallet functions. Ecommerce needs secure communications, they should be a unified application.

The blockchain fundamentals will also be improved with the addition of merge mining to each algorithm for the coin with the same PoW and the most hashpower, in order to help get around the vulnerabilities of a small cryptocurrency to hashpower attacks. In line with this, there will also be the addition of more PoW algorithms, as many as possible, as well as reinforcing the robustness of the multi-algorithm difficulty regulation.

## Progress

Previous work that was having serious nil pointer problems has been archived but once a robust set of data primitives are written they should be quickly refactored.
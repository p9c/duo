# Parallelcoin Blockchain Client 

[![Go Report Card](https://goreportcard.com/badge/github.com/parallelcointeam/duo?style=flat-square)](https://goreportcard.com/report/github.com/parallelcointeam/duo)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/parallelcointeam/duo)
[![Release](https://img.shields.io/github/release/parallelcointeam/duo.svg?style=flat-square)](https://github.com/parallelcointeam/duo/releases/latest)

Work in progress...

v0.0.0 is not a release but there is lots inside it, and lots more coming!

## Parallelcoin docker inside!

In the `docker/legacy` directory, if you are running linux and have docker installed, if you run `source init.sh` and `halp` you will get instructions how to get the legacy parallelcoind running in a docker (for windows/mac you will need to control it differently, possibly mac can use the aliases, on windows maybe a mingw shell). Basically, to make it run, you can do a one liner:

    .stop;.rm;.build;.run;.start

and it will be listening on 127.0.0.1:11048 for RPC and 0.0.0.0:11047 for p2p communication.
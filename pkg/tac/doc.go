// Package tac is a library for a telemetry and control interface that allows the deployment, automatic zeroconf interlinking and monitoring of a cluster of services transparently as libraries, local sockets, tcp and udp sockets and othter transports for RPC.
//
// An application creates a new cryptographic identity when it first runs and it implements a command interface that allows a user who has the secret key of the node to instruct the node and retreive telemetric data such as logs, even potentially to connect a TTY to a local shell.
//
// TaC functions like a universal pipe that connects services, and negotiates the path independently of the client application. It allows streamed and block-transfer/packet connections, proxying, relaying and routing. Transports are isolated from the calling functions and internally connect via pass by value functions, and through sockets if both are on localhost, and over tcp, udp and http.
package tac

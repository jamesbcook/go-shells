Go-Shells
========

Execute system commands on another computer


## Server Example

    ➜  server git:(master) ✗ ./server -h
       Usage of ./server:
         -port="4444": Port to listen on

## Client Example

    ➜  client git:(master) ✗ ./client -h
       Usage of ./client:
         -host="127.0.0.1": Host to connect to
         -port="4444": Port to connect to


### Basic

1. Commands are sent in the clear over the wire

### Base64

1. Commands are base64 encoded and sent over the wire

### Encrypted

1. Commands are encrypted, base64 encoded and then sent over the wire

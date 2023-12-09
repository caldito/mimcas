# Mimcas
![Apache 2.0 License](https://img.shields.io/hexpm/l/plug.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/caldito/mimcas.svg)](https://pkg.go.dev/github.com/caldito/mimcas)
[![Go Report Card](https://goreportcard.com/badge/github.com/caldito/mimcas)](https://goreportcard.com/report/github.com/caldito/mimcas)

Multithreaded In-Memory Cache Server.

Available commands:
- **set:** Sets a value for a new or existing key
- **get** Retrieves the value of a single key
- **mget** Retrieves the value of one or more keys
- **del** Removes an item from the cache
- **quit** Quit client session
- **ping** Responds "pong"
## Build
You can build the project running `make`. The binaries will then be available at `bin/mimcas-server` and `bin/mimcas-cli`.

Dependencies:
- `go` >= v1.18
- `make`

Command `make run` does the same but it will also start the server.

## Running the server
### Docker
To start the mimcas server using docker run:
```
docker run -p 20000:20000 pablogcaldito/mimcas-server:v0.1.0
```
### Standalone binary
To build and run the server run the following two commands:
```
make
./bin/mimcas-server
```
### Server flags
None of them are required. The available flags are:
- `-port`: Port to use for listening for incoming connections. By default it will be `20000`.
- `-maxmemory`: Maximum number of bytes available to use. Items will be evicted following LRU policy when that limit is crossed. By default there is no limit.

## Connecting with a client
The only client for now is the CLI one. It will available when building the source code as well.
### Docker
```
docker run --network host -it pablogcaldito/mimcas-cli:v0.1.0 mimcas-cli
```
### Standalone binary
Build and execute the resulting binary. Then connect with the CLI and you can issue commands to the server:
```
make
./bin/mimcas-cli
>> get a
NULL
>> set a 2
OK
>> get a
OK
2
>> set b 3
OK
>> mget a b
OK
2
OK
3
>> del a
OK
>> get a 
NULL
>> quit
```
### Client flags
None of them are required. The available flags are:
- `-host`: Host to use when opening a connection. By default it will be `localhost`.
- `-port`: Port to use when opening a connection. By default it will be `20000`.

## License
This project is licensed under the Apache License Version 2.0

## Contributing
Pull requests are welcomed and encouraged. For questions, feature requests and bug reports, please open an issue.

There is also a [TODO](https://github.com/caldito/mimcas/blob/main/TODO) file containing work planned to do and also [issues on GitHub](https://github.com/caldito/mimcas/issues).

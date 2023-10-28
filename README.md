# Mimcas
Multithreaded In-Memory Cache Server.

![Apache 2.0 License](https://img.shields.io/hexpm/l/plug.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/caldito/mimcas.svg)](https://pkg.go.dev/github.com/caldito/mimcas)
[![Go Report Card](https://goreportcard.com/badge/github.com/caldito/mimcas)](https://goreportcard.com/report/github.com/caldito/mimcas)

## Build
You can build the project running `make`. The binary will then be available at `bin/mimcas-server`.

Command `make run` does the same but it will also start the server.

## Usage
### Available commands
- **set:** Sets a value for a new or existing key
- **get** Retrieves the value of a single key
- **mget** Retrieves the value of one or more keys
- **del** Removes an item from the cache
- **quit** Quit client session
- **ping** Responds "pong"

### Flags
None of them are required. The available flags are:
- `-port`: Port to use for listening for incoming connections. By default it will be `20000`.
- `-maxmemory`: Maximum number of bytes available to use. Items will be evicted following LRU policy when that limit is crossed. By default there is no limit.


## Clients
The only client for now is the CLI one. It will available when building the source code as well.

First of all start the server:
```
$ make # build the program
$ ./bin/mimcas-server # start the server
```

Then connect with the CLI:
```
$ ./bin/mimcas-cli
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

## License
This project is licensed under the Apache License Version 2.0

## Contributing
Pull requests are welcomed and encouraged. For questions, feature requests and bug reports, please open an issue.

There is also a [TODO](https://github.com/caldito/mimcas/blob/main/TODO) file containing work planned to do and also [issues on GitHub](https://github.com/caldito/mimcas/issues).

# Mimcas
Multithreaded In-Memory Cache Server.

![Apache 2.0 License](https://img.shields.io/hexpm/l/plug.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/caldito/mimcas.svg)](https://pkg.go.dev/github.com/caldito/mimcas)
[![Go Report Card](https://goreportcard.com/badge/github.com/caldito/mimcas)](https://goreportcard.com/report/github.com/caldito/mimcas)

## Build
You can build the project running `make`. The binary will then be available at `bin/mimcas-server`.

With make run following command does the same but also starts the server
```
make run
```

## Usage
### Available commands
- *set:* Sets a value for a new or existing key
- *get* Retrieves the value of a single key
- *mget* Retrieves the value of one or more keys
- *del* Removes 
- *quit* Finish session
- *ping* Responds "pong"


## Clients
For now there is no client but you can connect using tools like netcat as the protocol is quite simple.

Example:
```
$ nc localhost 20000
get a
NULL

set a 2
OK

get a 
OK
2

set b 3
OK

mget a b
OK
2

OK
3

del a
OK

get a
NULL

quit
```

## License
This project is licensed under the Apache License Version 2.0

## Contributing
Pull requests are welcomed and encouraged. For questions, feature requests and bug reports, please open an issue.

There is also a [TODO](https://github.com/caldito/mimcas/blob/main/TODO) file containing work planned to do and also [issues on GitHub](https://github.com/caldito/mimcas/issues).
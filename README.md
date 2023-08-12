# mimcas
Multithreaded In-Memory Cache Server.

## Usage

### Start the server
Run the following command to start the server. It listens for connections on port `20000`.
```
make run
```

## Available commands

### set
Sets a value for a new or existing key

### get
Retrieves the value of a key

### mget
Retrieves the value of one or multiple keys

## Connect with a client and run commands
As the client as of now you can use netcat.
Example
```
$ nc localhost 20000
get a
(nil)
set a 1
OK
get a
1
set b 2
OK
mget a b
1
2
mget b a
2
1
quit
```

## License
This project is licensed under the Apache License Version 2.0

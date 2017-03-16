# lock
global lock server


## build
```bash
$ make build
```

## server run example
```bash
$ ./bin/lock -p 8888

$ ./bin/lock -s /tmp/lock.sock
```

##  example
```bash
$ telnet localhost 6800
Trying ::1...
Connected to localhost.
Escape character is '^]'.
lock key1
true
lock key1
false
lock key2
true
unlock key1
lock key1
true
```

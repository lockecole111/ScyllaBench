# ScyllaBench


## Useage
```
go build -o scyllabench
./scyllabench
  -interval int
    	sampling interval second(s) of countqps. (default 5)
  -nodes string
    	nodes' ip or domain. (default "localhost")
  -ratio int
    	ratio of write/read,and if set to 0 means read and write are independent without ratio limit. (default 1)
  -read_thread int
    	number of read thread(s). (default 10)
  -ttl int
    	TTL. (default 600)
  -write_thread int
    	number of write thread(s). (default 10)
```


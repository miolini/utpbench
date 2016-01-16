# utpbench
UTP protocol network benchmark in Go

## Usage

```bash
$ utpbench
2016/01/16 11:01:38 UTP Benchmark Tool by Artem Andreenko (miolini@gmail.com)
  -c	client mode
  -d duration
    	duration (default 10s)
  -ds duration
    	duration for stats (default 5s)
  -h string
    	host (default "127.0.0.1")
  -l int
    	length of data (default 1400)
  -p int
    	port (default 6001)
  -s	server mode
  -t int
    	threads (default 1)```
    	
## Server run

```bash
$ utpbench -s -h 0.0.0.0
```

## Client run
```bash
$ utpbench -c -h 192.168.0.1 -t 10 -d 60s -ds 2s
```

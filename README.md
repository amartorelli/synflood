# synflood
Very simple TCP packet generator to send syn packets to a remote host

## Usage
```
Usage of synflood:
  -host string
        the host you want to send the packets to
  -interval int
        interval in milliseconds to send packets (default 3000)
  -port int
        the port to send the packets to
```

## Example
`synflood -host localhost -port 8080 -interval 100`
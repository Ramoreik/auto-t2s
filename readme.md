# Auto Tun2Socks

## What ?

This program is basically a wrapper around the great `tun2socks` project, it also uses the `netlink` module.  
> https://github.com/xjasonlyu/tun2socks  
> https://github.com/vishvananda/netlink  

## Why ?

In large scope penetration tests with a large amount of subnets that are accessible from a single pivot point, this tool comes in handy.  
The idea is to automatically configure the `tun0` interface and populate the routes for `tun2socks`.  
This prevents having to run a bunch of `ip route add ...` if a lot of precise addresses have to be added to your routes.  
This way, if you have any tunnel with a supported protocol you can instantly add all your routes and pivot.  

## Usage

```
Usage of ./auto-t2s:
  -device string
        Use this device [driver://]name
  -fwmark int
        Set firewall MARK (Linux only)
  -interface string
        Interface name (default "tun0")
  -loglevel string
        Log level [debug|info|warning|error|silent] (default "info")
  -mtu int
        Set device maximum transmission unit (MTU)
  -proxy string
        Use this proxy [protocol://]host[:port]
  -restapi string
        HTTP statistic server listen address
  -route-file string
        Route File, containing one route per line to configure.
  -tcp-auto-tuning
        Enable TCP receive buffer auto-tuning
  -tcp-rcvbuf string
        Set TCP receive buffer size for netstack
  -tcp-sndbuf string
        Set TCP send buffer size for netstack
  -tun-post-up string
        Execute a command after TUN device setup
  -tun-pre-up string
        Execute a command before TUN device setup
  -udp-timeout duration
        Set timeout for each UDP session
```

routes file example:
```
10.0.2.1/24
10.4.1.1/24
```

It runs just like normal tun2socks:
```
auto-t2s -interface wlan0 -mtu 1500 -proxy socks5://localhost:8081 -device tun0 -routes-file './routes'
```

This will create the tun0 device, then bring the link up.  

Afterwards, the routes will be populated, both of the first steps use `netlink`.  
Finally, the `tun2socks` engine will be started and the forwarding of packets through the socks proxy will start.  

  
```
curl https://10.100.1.100
...
```

Now the specified routes should be accessible.  
The caveat is that since the tool uses netlink it needs higher privileges.

## Disclaimer: 

This could probably have been achieved with a `-tun-pre-up` command pointing to a script.
I wanted to play with `tun2socks`, `netlink` and `gvisor` , I feel like there is potential to learn these projects to create interesting tools.

## Installation

```bash
go get github.com/vishvananda/netlink
go get github.com/xjasonlyu/tun2socks/v2
go build .
```



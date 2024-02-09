# Auto Tun2Socks

## What ?

This program is basically a wrapper around the great tun2socks project, it also uses the netlink module.  
> https://github.com/xjasonlyu/tun2socks  
> https://github.com/vishvananda/netlink  

## Why ?

In large scope penetration tests with a large amount of subnets that are accessible from a single pivot point, this tool comes in handy.  
The idea is to automatically configure the `tun0` interface and populate the routes for `tun2socks`.  
This prevents having to run a bunch of `ip route add ...` if a lot of precise addresses have to be added to your routes.  
This way, if you have any tunnel with a supported protocol you can instantly add all your routes and pivot.  

## Usage

```
auto-t2s -interface wlan0 -mtu 1500 -proxy socks5://localhost:8081 -device tun0 -routes '10.100.1.1/24'
```
> This program needs root access to function properly.  

This will create the tun0 device, and assign it the `10.10.10.10/32` address then bring the link up.  
Afterwards, the routes will be populated, both of the first steps use `netlink`.  
Finally, the `tun2socks` engine will be started and the forwarding of packets through the socks proxy will start.  
  
```
curl https://10.100.1.100
...
```

Now the specified routes should be accessible.  

## Installation

```bash
go get github.com/vishvananda/netlink
go get github.com/xjasonlyu/tun2socks
```



package main

import (
  "flag"
  "net"
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "strings"

  "github.com/vishvananda/netlink"
  "github.com/xjasonlyu/tun2socks/v2/engine"
)

var (
  key = new(engine.Key)
  routes string
)

func init () {
  flag.IntVar(&key.Mark, "fwmark", 0, "Set firewall MARK (Linux only)")
  flag.StringVar(&key.Interface, "interface", "tun0", "Interface name") 
  flag.StringVar(&key.Proxy, "proxy", "", "Use this proxy [protocol://]host[:port]")
  flag.IntVar(&key.MTU, "mtu", 0, "Set device maximum transmission unit (MTU)")
  flag.DurationVar(&key.UDPTimeout, "udp-timeout", 0, "Set timeout for each UDP session")
  flag.StringVar(&key.Device, "device", "", "Use this device [driver://]name")
  flag.StringVar(&key.LogLevel, "loglevel", "info", "Log level [debug|info|warning|error|silent]")
  flag.StringVar(&key.RestAPI, "restapi", "", "HTTP statistic server listen address")
  flag.StringVar(&key.TCPSendBufferSize, "tcp-sndbuf", "", "Set TCP send buffer size for netstack")
  flag.StringVar(&key.TCPReceiveBufferSize, "tcp-rcvbuf", "", "Set TCP receive buffer size for netstack")
  flag.BoolVar(&key.TCPModerateReceiveBuffer, "tcp-auto-tuning", false, "Enable TCP receive buffer auto-tuning")
  flag.StringVar(&key.TUNPreUp, "tun-pre-up", "", "Execute a command before TUN device setup")
  flag.StringVar(&key.TUNPostUp, "tun-post-up", "", "Execute a command after TUN device setup")
  flag.StringVar(&routes, "routes", "", "Routes, space separated.")
  flag.Parse()
}

func startEngine () {
  fmt.Printf(("[*] Initializaing and starting tun2socks engine.\n"))
  engine.Insert(key)
  engine.Start()
  defer engine.Stop()

  configureTun()
  populateRoutes()

  fmt.Printf("[*] Waiting until told to terminate.\n")
  sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

}

func configureTun() {
  fmt.Printf("[*] Creating tunnel device: %s.\n", key.Device)
  tun0_nl,_ := netlink.LinkByName(key.Device)
  tun0_a,_ := netlink.ParseAddr("10.10.10.10/32")
  netlink.AddrAdd(tun0_nl, tun0_a)
  netlink.LinkSetUp(tun0_nl)
}

func populateRoutes() {
  fmt.Printf("[*] Populating routes.\n")
  rs := strings.Split(routes, " ") 
  for _, route := range rs {
    tun,_ := netlink.LinkByName(key.Device)
    _,ipnet,_ := net.ParseCIDR(route) 
    r := &netlink.Route{
      Scope: netlink.SCOPE_LINK,
      LinkIndex: tun.Attrs().Index,
      Dst: ipnet,
    }
    fmt.Printf("[*] Adding route: %s.\n", route)
    netlink.RouteAdd(r)
  }
}

func main () {
  fmt.Printf("[*] Using interface: %s \n", key.Interface)
  fmt.Printf("[*] Using device: %s \n", key.Device)
  fmt.Printf("[*] Using proxy: %s \n", key.Proxy)
  startEngine()
}


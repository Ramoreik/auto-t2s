package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/vishvananda/netlink"
	"github.com/xjasonlyu/tun2socks/v2/engine"
	"github.com/xjasonlyu/tun2socks/v2/log"
)

var (
  key = new(engine.Key)
  routesFile string
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
  flag.StringVar(&routesFile, "route-file", "", "Route File, containing one route per line to configure.")
  flag.Parse()
}

func readRoutesFile (routesFile string) []string {
  content,err := os.ReadFile(routesFile)
  if err != nil {
    fmt.Print(err)
  }
  routes := strings.Split(string(content), "\n")
  return routes[0:len(routes)-1]
}

func bringUpTun() {
  log.Infof("[auto-t2s] Bringing up device: %s.\n", key.Device)
  tun0_nl,_ := netlink.LinkByName(key.Device)
  netlink.LinkSetUp(tun0_nl)
}

func populateRoutes(routes []string) {
  log.Infof("[auto-t2s] Populating routes.\n")
  for _, route := range routes {
    tun,_ := netlink.LinkByName(key.Device)
    _,ipnet,err := net.ParseCIDR(route) 
    if err != nil {
      log.Warnf("Invalid CIDR provided : '%s'", route)
    }
    r := &netlink.Route{
      Scope: netlink.SCOPE_LINK,
      LinkIndex: tun.Attrs().Index,
      Dst: ipnet,
    }
    log.Infof("[auto-t2s] Adding route: %s.\n", route)
    netlink.RouteAdd(r)
  }
}

func startEngine () {
  log.Infof(("[auto-t2s] Initializaing and starting tun2socks engine.\n"))
  engine.Insert(key)
  engine.Start()
  defer engine.Stop()

  bringUpTun()
  routes := readRoutesFile(routesFile)
  populateRoutes(routes)

  log.Infof("[auto-t2s] Waiting until told to terminate.\n")
  sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

}

func main () {
  log.Infof("[auto-t2s] Using interface: %s \n", key.Interface)
  log.Infof("[auto-t2s] Using device: %s \n", key.Device)
  log.Infof("[auto-t2s] Using proxy: %s \n", key.Proxy)
  startEngine()
}


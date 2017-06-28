package misc

import (
	"log"
	"net"
)

/*
Based on:
	* URL: https://github.com/mccoyst/myip/blob/master/myip.go
	* URL: http://changsijay.com/2013/07/28/golang-get-ip-address/
*/

// GetIP returns the local network ip address.
func GetIP() string {
	addrs, e := net.InterfaceAddrs()
	if e != nil {
		log.Panic(e)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

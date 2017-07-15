package main

import (
	"net"

	"github.com/spf13/viper"
)

var whiteIPList = make([]net.IP, 0)

func initConfig() {
	a := viper.GetStringSlice("white_ip_list")
	if len(a) > 0 {
		for _, s := range a {
			ip := net.ParseIP(s)
			whiteIPList = append(whiteIPList, ip)
		}
	}
}

func isWhiteIP(ip net.IP) bool {
	if len(whiteIPList) == 0 {
		return true
	}
	for _, t := range whiteIPList {
		if t.String() == ip.String() {
			return true
		}
	}
	return false
}

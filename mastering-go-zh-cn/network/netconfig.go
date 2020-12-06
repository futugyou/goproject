package main

import (
	"fmt"
	"net"
)

func netinterfaces() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, i := range interfaces {
		fmt.Printf("interface: %v\n", i.Name)
		fmt.Printf("interface Flages: %v\n", i.Flags.String())
		fmt.Printf("interface MTU: %v\n", i.MTU)
		fmt.Printf("interface Hardware Address: %v\n", i.HardwareAddr)
		byname, err := net.InterfaceByName(i.Name)
		if err != nil {
			fmt.Println(err)
		}
		addresses, err := byname.Addrs()
		for k, v := range addresses {
			fmt.Printf("interface address #%v: %v\n", k, v.String())
		}
		fmt.Println()
	}
}

var (
	addresslist = []string{"127.0.0.1", "mtsoukalos.eu", "packtpub.com", "google.com", "www.google.com", "cnn.com"}
)

func lookip(address string) ([]string, error) {
	hosts, err := net.LookupAddr(address)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

func lookHostname(hostname string) ([]string, error) {
	ops, err := net.LookupHost(hostname)
	if err != nil {
		return nil, err
	}
	return ops, nil
}

func ndsfind() {
	for _, input := range addresslist {
		fmt.Println("input: ", input)
		ipaddress := net.ParseIP(input)
		if ipaddress == nil {
			ips, err := lookHostname(input)
			if err == nil {
				for _, singleip := range ips {
					fmt.Println("one: ", singleip)
				}
			}
		} else {
			hosts, err := lookip(input)
			if err == nil {
				for _, hostname := range hosts {
					fmt.Println("two: ", hostname)
				}
			}
		}
		fmt.Println()
	}
}

func findns() {
	for _, input := range addresslist {
		fmt.Println("input: ", input)
		nses, err := net.LookupNS(input)
		if err == nil {
			for _, ns := range nses {
				fmt.Println(ns.Host)
			}

		}
		fmt.Println()
	}
}

func findmx() {
	for _, input := range addresslist {
		fmt.Println("input: ", input)
		nses, err := net.LookupMX(input)
		if err == nil {
			for _, ns := range nses {
				fmt.Println(ns.Host)
			}

		}
		fmt.Println()
	}
}

func main() {
	findmx()
}

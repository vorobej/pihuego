package hue

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

const (
	bridgeProtocol = "http://"
)

// Bridge holds internal data for hue bridge
type Bridge struct {
	ip       string
	username string
}

// DiscoverBridge search for bridge over local network
func DiscoverBridge() Bridge {
	if useHardcodedBridge {
		bridgeIP := "http://10.0.0.57"
		fmt.Println("DiscoverBridge DEBUG: using hardcoded bridge ip: ", bridgeIP)
		return Bridge{ip: bridgeIP}
	}

	service := "239.255.255.250:1900"
	macAddress, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		fmt.Println("DiscoverBridges Error: ", err)
	}
	send, err := net.DialUDP("udp4", nil, macAddress)
	if err != nil {
		fmt.Println("DiscoverBridges Error: ", err)
	}
	defer send.Close()

	// Send SSDP Message
	ssdpDiscoveryMessage := []byte("M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: ssdp:discover\r\nMX: 10\r\nST: \"ssdp:all\"\r\n\r\n")
	_, err = send.Write(ssdpDiscoveryMessage)
	if err != nil {
		fmt.Println("DiscoverBridges Error: ", err)
	}
	fmt.Println("Searching for Philip Hue Hub (Could take up to 30 secs)...")

	// Listen for SSDP/HTTP NOTIFY over UDP
	listen, err := net.ListenMulticastUDP("udp4", nil, macAddress)
	if err != nil {
		fmt.Println("DiscoverBridges Error: ", err)
	}
	defer listen.Close()

	descriptionURL := ""
	for {
		b := make([]byte, 256)
		_, _, err := listen.ReadFromUDP(b)
		if err != nil {
			fmt.Println("DiscoverBridges Error: ", err)
		}
		payloadMessage := string(b)
		headers := strings.Split(payloadMessage, "\r\n")
		for _, header := range headers {
			datum := strings.Split(header, ": ")
			if len(datum) > 1 {
				if datum[0] == "LOCATION" {
					if strings.Contains(datum[1], "description.xml") {
						descriptionURL = datum[1]
						break
					}
				}
			}
		}
		if strings.Contains(descriptionURL, "description.xml") {
			break
		}
	}
	u, err := url.Parse(descriptionURL)
	if err != nil {
		fmt.Println("DiscoverBridges Error: ", err)
	}
	hostname := ""
	if strings.Contains(u.Host, ":") {
		h := strings.Split(u.Host, ":")
		hostname = h[0]
	} else {
		hostname = u.Host
	}
	fmt.Printf("Found Hub at %s\n", hostname)

	var result = Bridge{ip: bridgeProtocol + hostname}
	return result
}

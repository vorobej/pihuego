package main

import (
	"fmt"

	"github.com/vorobej/pihuego/hue"
)

func main() {
	var bridge = hue.DiscoverBridge()
	hue.CreateUser(&bridge)
	lights, err := hue.LightsStatus(&bridge)
	if err != nil {
		fmt.Println("Unable to get lights status: ", err)
	}

	for _, light := range lights {
		fmt.Printf("%s\n", light)
	}

	//hue.SetLightState(&bridge, nil)

	/*
		fmt.Println("starting server...")
		s, err := gossdp.NewSsdp(nil)
		if err != nil {
			log.Println("Error creating ssdp server: ", err)
			return
		}

			go s.Start()
			serverDef := gossdp.AdvertisableServer{
				ServiceType: "urn:vorobej:sparrow:light:1",
				DeviceUuid:  "hh0c2981-0029-44b7-4u04-27f187aecf78",
				Location:    GetLocalIP(),
				MaxAge:      3600,
			}
			s.AdvertiseServer(serverDef)
			time.Sleep(3000 * time.Second)
	*/
}

/*
// GetLocalIP Gets local ip of running device
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}*/

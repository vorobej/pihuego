package hue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"strconv"
	"strings"
)

const (
	bridgeProtocol = "http://"
	bridgeFilename = "bridge.dat"
)

// Bridge holds internal data for hue bridge
type Bridge struct {
	IP       string `json:"ip"`
	Username string `json:"username"`
}

// LoadBridges read config file and return list of saved bridges
func LoadBridges() (bridges []Bridge, err error) {
	data, err := ioutil.ReadFile(bridgeFilename)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &bridges); err != nil {
		return nil, err
	}

	return bridges, nil
}

// Lights return list of lights added to bridge
func (bridge *Bridge) Lights() ([]Light, error) {
	if bridge == nil {
		return nil, fmt.Errorf("bridge can't be nil")
	}
	resp := []byte(jsonStr)
	var err error

	var jsonResult map[string]interface{}
	if err = json.Unmarshal(resp, &jsonResult); err != nil {
		fmt.Printf("ERROR json: %s\n", err)
		return nil, err
	}

	lights := make([]Light, len(jsonResult))
	var index int
	for key, value := range jsonResult {
		keyInt, _ := strconv.Atoi(key)
		lightObject := value.(map[string]interface{})
		stateObject := lightObject["state"].(map[string]interface{})
		xy := stateObject["xy"].([]interface{})
		lightState := LightState{
			On:               stateObject["on"].(bool),
			Reachable:        stateObject["reachable"].(bool),
			Brightness:       (uint8)(stateObject["bri"].(float64)),
			Saturation:       (uint8)(stateObject["sat"].(float64)),
			Hue:              (uint16)(stateObject["hue"].(float64)),
			ColorTemperature: (uint16)(stateObject["ct"].(float64)),
			Alert:            stateObject["alert"].(string),
			Effect:           stateObject["effect"].(string),
			ColorMode:        stateObject["colormode"].(string),
			XY:               [2]float64{xy[0].(float64), xy[1].(float64)},
		}

		lights[index] = Light{
			id:     keyInt,
			name:   lightObject["name"].(string),
			state:  lightState,
			bridge: bridge,
		}
		index++
	}
	return lights, nil
}

// DiscoverBridge search for bridge over local network
func DiscoverBridge() Bridge {
	if useHardcodedBridge {
		bridgeIP := "http://10.0.0.57"
		fmt.Println("DiscoverBridge DEBUG: using hardcoded bridge ip: ", bridgeIP)
		return Bridge{IP: bridgeIP}
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

	var result = Bridge{IP: bridgeProtocol + hostname}
	return result
}

func (bridge Bridge) String() string {
	return fmt.Sprintf("ip<%s> username<%s>", bridge.IP, bridge.Username)
}

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/vorobej/pihuego/hue"
)

const (
	allLights = -1
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("show usage")
		return
	}

	// TODO load bridge from json file?
	bridge := hue.DiscoverBridge()

	switch os.Args[1] {
	case "off":
		fmt.Println("NOT IMPLEMENTED: lights off")
		turnLigthsOff(&bridge)
	case "on":
		fmt.Println("NOT IMPLEMENTED: lights on")
	case "list":
		listLights(&bridge)
	}
}

func turnLigthsOff(bridge *hue.Bridge) bool {
	lightID := allLights
	// check if light id is passed, if not - turn off all
	if len(os.Args) >= 3 {
		var err error
		lightID, err = strconv.Atoi(os.Args[2])
		if err != nil {
			lightID = allLights
		}
	}

	fmt.Printf("turning off light id<%d>\n", lightID)
	lights, err := hue.LightsStatus(bridge)
	if err != nil {
		fmt.Println("turnLigthsOff: Unable to get lights status: ", err)
		return false
	}

	for _, light := range lights {
		if lightID == allLights {
			light.TurnOff(bridge)
		} else if lightID == light.ID() {
			light.TurnOff(bridge)
			break
		}
	}
	return true
}

func listLights(bridge *hue.Bridge) {
	fmt.Println("List of lights:")
	lights, err := hue.LightsStatus(bridge)
	if err != nil {
		fmt.Println("Can't get lights status")
		return
	}

	for _, light := range lights {
		fmt.Printf("\t%s\n", light)
	}
}

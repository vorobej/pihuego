package main

import (
	"flag"
	"fmt"

	"github.com/vorobej/pihuego/hue"
)

const (
	allLights      = -1
	invalidID      = -1
	defaultCommand = "usage"
)

type commandHandler func(bridge *hue.Bridge, light *hue.Light) error

type commmand struct {
	handler       commandHandler
	requireBridge bool
	requireLight  bool
}

var bridgeID = flag.Int("bridge", invalidID, "bridge ID")
var lightID = flag.Int("light", invalidID, "light ID")
var cmdName = flag.String("cmd", defaultCommand, "command to execute")

func main() {
	flag.Parse()
	fmt.Printf("DEBUG: bridge<%d> light<%d> command<%s>\n", *bridgeID, *lightID, *cmdName)

	var commands = map[string]commmand{
		"usage": {
			handler: usageHandler,
		},
		"bridges": {
			handler: bridgesHandler,
		},
		"lights": {
			handler:       lightsHandler,
			requireBridge: true,
		},
		"on": {
			handler:       onHandler,
			requireBridge: true,
			requireLight:  true,
		},
		"off": {
			handler:       offHandler,
			requireBridge: true,
			requireLight:  true,
		},
	}

	var cmdBridge *hue.Bridge
	var cmdLight *hue.Light

	cmd, ok := commands[*cmdName]
	if !ok {
		fmt.Printf("ERROR: unknown command: %s\n", *cmdName)
		return
	}
	if cmd.requireBridge {
		if *bridgeID == invalidID {
			fmt.Printf("ERROR: command<%s> required bridge id\n", *cmdName)
			return
		}
		bridges, err := hue.LoadBridges()
		if err != nil {
			fmt.Printf("ERROR: unable to load bridges: %s\n", err)
			return
		} else if *bridgeID >= len(bridges) {
			// TODO print whole list?
			fmt.Printf("ERROR: bridge id is out of range\n")
			return
		}
		cmdBridge = &bridges[*bridgeID]

		if cmd.requireLight {
			if *lightID == invalidID {
				fmt.Printf("ERROR: command<%s> require light id\n", *cmdName)
				return
			}
			lights, err := cmdBridge.Lights()
			if err != nil {
				fmt.Printf("ERROR: unable to get lights from bridge<%s>\n", cmdBridge.IP)
			} else if *lightID >= len(lights) {
				// TODO print whole list?
				fmt.Printf("ERROR: light id is out of range\n")
				return
			}
			cmdLight = &lights[*lightID]
		}
	}
	// TODO check error?
	cmd.handler(cmdBridge, cmdLight)
}

func usageHandler(bridge *hue.Bridge, light *hue.Light) error {
	fmt.Printf("you have to provide command name\n")
	return nil
}

// bridgesHandler show list of saved bridges
func bridgesHandler(bridge *hue.Bridge, light *hue.Light) error {
	bridges, err := hue.LoadBridges()
	if err != nil {
		return fmt.Errorf("can't load bridges: %s", err)
	}
	if len(bridges) == 0 {
		fmt.Println("There're no saved bridges")
	} else {
		fmt.Println("#### Found bridges:")
		for ix, bridge := range bridges {
			fmt.Printf("\t[%d] %s\n", ix, bridge)
		}
	}
	return nil
}

// lightsHandler show list of lights available on bridge
func lightsHandler(bridge *hue.Bridge, light *hue.Light) error {
	fmt.Println("#### Available lights:")
	lights, err := bridge.Lights()
	if err != nil {
		return err
	}
	for ix, light := range lights {
		fmt.Printf("\t[%d] %s\n", ix, light)
	}
	return nil
}

// offHandler turn off provided light
func offHandler(bridge *hue.Bridge, light *hue.Light) error {
	return light.TurnOff()
}

// onHandler turn on provided light
func onHandler(bridge *hue.Bridge, light *hue.Light) error {
	return light.TurnOn()
}

package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/vorobej/pihuego/hue"
)

const (
	allLights      = -1
	invalidID      = -1
	defaultCommand = "usage"
)

type commandHandler func(args commandArguments) error

type commmand struct {
	handler       commandHandler
	requireBridge bool
	requireLight  bool
}

type commandArguments struct {
	bridge *hue.Bridge
	light  *hue.Light
	color  string
}

var bridgeID = flag.Int("bridge", invalidID, "bridge ID")
var lightID = flag.Int("light", invalidID, "light ID")
var cmdName = flag.String("cmd", defaultCommand, "command to execute")
var color = flag.String("color", "", "color to set light to")

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
		"discover": {
			handler: discoverHandler,
		},
		"pair": {
			handler: pairHandler,
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
		"color": {
			handler:       colorHandler,
			requireBridge: true,
			requireLight:  true,
		},
	}

	var args commandArguments
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
		args.bridge = &bridges[*bridgeID]

		if cmd.requireLight {
			if *lightID == invalidID {
				fmt.Printf("ERROR: command<%s> require light id\n", *cmdName)
				return
			}
			lights, err := args.bridge.Lights()
			if err != nil {
				fmt.Printf("ERROR: unable to get lights from bridge<%s>\n", args.bridge.IP)
			} else if *lightID >= len(lights) {
				// TODO print whole list?
				fmt.Printf("ERROR: light id is out of range\n")
				return
			}
			args.light = &lights[*lightID]
		}
		// set color
		args.color = *color
	}
	// TODO check error?
	if err := cmd.handler(args); err != nil {
		fmt.Println(err)
	}
}

func usageHandler(args commandArguments) error {
	fmt.Printf("you have to provide command name\n")
	return nil
}

// bridgesHandler show list of saved bridges
func bridgesHandler(args commandArguments) error {
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
func lightsHandler(args commandArguments) error {
	fmt.Println("#### Available lights:")
	lights, err := args.bridge.Lights()
	if err != nil {
		return err
	}
	for ix, light := range lights {
		fmt.Printf("\t[%d] %s\n", ix, light)
	}
	return nil
}

// offHandler turn off provided light
func offHandler(args commandArguments) error {
	return args.light.TurnOff()
}

// onHandler turn on provided light
func onHandler(args commandArguments) error {
	return args.light.TurnOn()
}

// discoverHandler discover and save bridge to config
func discoverHandler(args commandArguments) error {
	fmt.Println("NOT IMPLEMENTED: discoverHandler()")
	// TODO save discovered bridge to config file
	hue.DiscoverBridge()
	return nil
}

// pairHandler authorize and save username for hub
func pairHandler(args commandArguments) error {
	fmt.Println("NOT IMPLEMENTED: pairHandler()")
	return nil
}

// colorHandler set color for light
func colorHandler(args commandArguments) error {
	r, g, b, err := extractColor(*color)
	if err != nil {
		return err
	}
	fmt.Printf("set color to r<%f> g<%f> b<%f>\n", r, g, b)
	return args.light.SetColor(r, g, b)
}

// extractColor get rgb color from string. note string should starts with 0x and should have 8 chars, ie 0xDEADBEAF
func extractColor(color string) (r, g, b float64, err error) {
	if len(color) != 8 {
		err = fmt.Errorf("invalid color length")
		return
	}
	if !strings.HasPrefix(color, "0x") {
		err = fmt.Errorf("color should starts from 0x")
		return
	}
	intR, err := strconv.ParseInt(color[2:4], 16, 64)
	if err != nil {
		return
	}
	intG, err := strconv.ParseInt(color[4:6], 16, 64)
	if err != nil {
		return
	}
	intB, err := strconv.ParseInt(color[6:8], 16, 64)
	if err != nil {
		return
	}

	r = float64(intR) / 255.0
	g = float64(intG) / 255.0
	b = float64(intB) / 255.0
	return
}

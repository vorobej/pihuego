package hue

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/vorobej/pihuego/hue/request"
)

type byID []Light

// Light datastructure for light
type Light struct {
	id    int
	name  string
	state LightState

	// pointer to bridge where lights are stored TODO remove args from on/off function
	bridge *Bridge
}

// LightState current state of light
type LightState struct {
	/*
		On/Off state of the light. On=true, Off=false
	*/
	On bool `json:"on"`

	/*
		Brightness of the light. This is a scale from the minimum brightness the light is capable of, 1, to the maximum capable brightness, 254.
	*/
	Brightness uint8 `json:"bri"`

	/*
		Saturation of the light. 254 is the most saturated (colored) and 0 is the least saturated (white).
	*/
	Saturation uint8 `json:"sat"`
	/*
		Hue of the light. This is a wrapping value between 0 and 65535. Both 0 and 65535 are red, 25500 is green and 46920 is blue.
	*/
	Hue uint16 `json:"hue"`
	/*
		The Mired Color temperature of the light. 2012 connected lights are capable of 153 (6500K) to 500 (2000K).
	*/
	ColorTemperature uint16 `json:"ct"`

	/*
		The dynamic effect of the light, can either be “none” or “colorloop”.
		If set to colorloop, the light will cycle through all hues using the current brightness and saturation settings.
	*/
	Effect string `json:"effect"`

	/*
		“none” 		– The light is not performing an alert effect.
		“select” 	– The light is performing one breathe cycle.
		“lselect” 	– The light is performing breathe cycles for 15 seconds or until an "alert": "none" command is received.
	*/
	Alert string `json:"alert"`

	/*
		Indicates the color mode in which the light is working, this is the last command type it received. Values are “hs” for Hue and Saturation,
		“xy” for XY and “ct” for Color Temperature. This parameter is only present when the light supports at least one of the values.
	*/
	ColorMode string `json:"colormode"`

	/*
		Indicates if a light can be reached by the bridge.
	*/
	Reachable bool `json:"reachable"`

	/*
		The x and y coordinates of a color in CIE color space.
		The first entry is the x coordinate and the second entry is the y coordinate. Both x and y are between 0 and 1.
	*/
	XY [2]float64 `json:"xy"`
}

type SetLightStateBody struct {
	On         bool      `json:"on"`
	Brightness uint8     `json:"bri,omitempty"`
	Hue        uint16    `json:"hue,omitempty"`
	Saturation uint8     `json:"sat,omitempty"`
	XY         []float64 `json:"xy,omitempty"`
	//ct
	//alert
	//effect
	//transitiontime
	//bri_inc
	//sat_inc
	//hue_inc
	//ct_inc
	//xy_inc

}

// SetLightState set new state to selected light
func (light *Light) SetLightState() error {
	if err := verifyLight(light); err != nil {
		return nil
	}

	var body = SetLightStateBody{
		Hue:        3000,
		On:         true,
		Brightness: 200,
	}
	data, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("JSON marshaling is failing: %s", err)
	}
	url := fmt.Sprintf("%s/api/%s/lights/%d/state", light.bridge.IP, light.bridge.Username, light.id)
	data, err = request.PUT(url, data)
	if err != nil {
		return err
	}
	fmt.Println("SetLightState() response: ", string(data))
	return nil
}

// TurnOff method to turn off light
func (light *Light) TurnOff() error {
	if err := verifyLight(light); err != nil {
		return err
	}

	data, err := json.Marshal(SetLightStateBody{On: false})
	if err != nil {
		fmt.Printf("JSON marshaling is failing: %s", err)
		return err
	}

	url := fmt.Sprintf("%s/api/%s/lights/%d/state", light.bridge.IP, light.bridge.Username, light.id)
	data, err = request.PUT(url, data)
	if err != nil {
		return err
	}
	fmt.Println("TurnOff() response: ", string(data))
	return nil
}

// TurnOn restore last state of light
func (light *Light) TurnOn() error {
	if err := verifyLight(light); err != nil {
		return err
	}

	data, err := json.Marshal(SetLightStateBody{On: true})
	if err != nil {
		fmt.Printf("JSON marshaling is failing: %s", err)
		return err
	}

	url := fmt.Sprintf("%s/api/%s/lights/%d/state", light.bridge.IP, light.bridge.Username, light.id)
	data, err = request.PUT(url, data)
	if err != nil {
		return err
	}
	fmt.Println("TurnOn() response: ", string(data))
	return nil
}

// SetColor set color for light
func (light *Light) SetColor(r, g, b float64) error {
	if err := verifyLight(light); err != nil {
		return err
	}

	var red, green, blue float64

	// apply gamma correction
	if r > 0.04045 {
		red = math.Pow((r+0.055)/(1.0/0.055), 2.4)
	} else {
		red = r / 12.92
	}
	if g > 0.04045 {
		green = math.Pow((g+0.055)/(1.0+0.055), 2.4)
	} else {
		green = g / 12.92
	}
	if b > 0.04045 {
		blue = math.Pow((b+0.055)/(1.0+0.055), 2.4)
	} else {
		blue = b / 12.92
	}
	fmt.Printf("RGB<%f/%f/%f>\n", red, green, blue)
	// convert rgb to xyz
	X := red*0.664511 + green*0.154324 + blue*0.162028
	Y := red*0.283881 + green*0.668433 + blue*0.047685
	Z := red*0.000088 + green*0.072310 + blue*0.986039

	// calculate xy
	x := X / (X + Y + Z)
	y := Y / (X + Y + Z)

	// TODO check if light state is on?
	color := SetLightStateBody{On: true, XY: []float64{x, y}}
	data, err := json.Marshal(color)
	if err != nil {
		fmt.Printf("JSON marshaling is failing: %s", err)
	}
	fmt.Println(string(data))
	url := fmt.Sprintf("%s/api/%s/lights/%d/state", light.bridge.IP, light.bridge.Username, light.id)
	data, err = request.PUT(url, data)
	if err != nil {
		return err
	}
	fmt.Println("SetColor() response: ", string(data))
	return nil
}

// ID returns light id
func (light *Light) ID() int {
	return light.id
}

// Prints light info
func (light Light) String() string {
	return fmt.Sprintf("id<%d> name<%s> state<%v>", light.id, light.name, light.state)
}

// verifyLight check if light is not null and have valid bridge
func verifyLight(light *Light) error {
	if light == nil {
		return fmt.Errorf("light is nil")
	}
	if light.bridge == nil {
		return fmt.Errorf("light don't have parent bridge")
	}
	return nil
}

func (a byID) Len() int {
	return len(a)
}

func (a byID) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byID) Less(i, j int) bool {
	return a[i].id < a[j].id
}

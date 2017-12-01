package hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"

	"github.com/vorobej/pihuego/hue/request"
)

const (
	jsonStr = `{"1":{"state":{"on":true,"bri":254,"hue":34076,"sat":251,"effect":"none","xy":[0.3144,0.3301],"ct":153,"alert":"select","colormode":"xy","reachable":true},"swupdate":{"state":"noupdates","lastinstall":null},"type":"Extended color light","name":"living room","modelid":"LCT007","manufacturername":"Philips","uniqueid":"00:17:88:01:10:5a:33:78-0b","swversion":"5.50.1.19085"},"3":{"state":{"on":true,"bri":254,"hue":14956,"sat":140,"effect":"none","xy":[0.4571,0.4097],"ct":366,"alert":"select","colormode":"ct","reachable":true},"swupdate":{"state":"noupdates","lastinstall":null},"type":"Extended color light","name":"kitchen","modelid":"LCT007","manufacturername":"Philips","uniqueid":"00:17:88:01:10:55:f2:7c-0b","swversion":"5.50.1.19085"},"4":{"state":{"on":false,"bri":254,"hue":14956,"sat":140,"effect":"none","xy":[0.4571,0.4097],"ct":366,"alert":"select","colormode":"ct","reachable":true},"swupdate":{"state":"noupdates","lastinstall":null},"type":"Extended color light","name":"bedroom","modelid":"LCT007","manufacturername":"Philips","uniqueid":"00:17:88:01:10:5a:45:c6-0b","swversion":"5.50.1.19085"},"5":{"state":{"on":true,"bri":200,"hue":3000,"sat":140,"effect":"none","xy":[0.5182,0.3643],"ct":480,"alert":"select","colormode":"hs","reachable":true},"swupdate":{"state":"noupdates","lastinstall":null},"type":"Extended color light","name":"lightstrip","modelid":"LST002","manufacturername":"Philips","uniqueid":"00:17:88:01:02:af:3b:42-0b","swversion":"5.90.0.19950"}}`
)

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
	On         bool       `json:"on"`
	Brightness uint8      `json:"bri,omitempty"`
	Hue        uint16     `json:"hue,omitempty"`
	Saturation uint8      `json:"sat,omitempty"`
	XY         [2]float64 `json:"xy,omitempty"`
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
func SetLightState(bridge *Bridge, light *Light) {
	if bridge == nil /*|| light == nil */ {
		fmt.Printf("SetLightState() invalid arguments: bridge<%p> light<%p>\n", bridge, light)
		return
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
	url := fmt.Sprintf("%s/api/%s/lights/%d/state", bridge.IP, bridge.Username, light.id)
	request.PUT(url, bytes.NewReader(data))
}

// TurnOff method to turn off light
func (light *Light) TurnOff(bridge *Bridge) error {
	data, err := json.Marshal(SetLightStateBody{On: false})
	if err != nil {
		fmt.Printf("JSON marshaling is failing: %s", err)
	}
	url := fmt.Sprintf("%s/api/%s/lights/%d/state", bridge.IP, bridge.Username, light.id)
	_, err = request.PUT(url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return nil
}

// TurnOn restore last state of light
func (light *Light) TurnOn(bridge *Bridge) error {
	data, err := json.Marshal(SetLightStateBody{On: true})
	if err != nil {
		fmt.Printf("JSON marshaling is failing: %s", err)
	}
	url := fmt.Sprintf("%s/api/%s/lights/%d/state", bridge.IP, bridge.Username, light.id)
	_, err = request.PUT(url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return nil
}

// SetColor set color for light
func (light *Light) SetColor(bridge *Bridge, r, g, b float64) error {
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

	color := SetLightStateBody{On: true, XY: [2]float64{x, y}}
	data, err := json.Marshal(color)
	if err != nil {
		fmt.Printf("JSON marshaling is failing: %s", err)
	}

	fmt.Println(data)
	/*
		url := fmt.Sprintf("%s/api/%s/lights/%d/state", bridge.ip, bridge.username, light.id)
		_, err = request.PUT(url, bytes.NewReader(data))
		if err != nil {
			return err
		}*/
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

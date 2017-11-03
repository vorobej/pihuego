package hue

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/vorobej/pihuego/hue/request"
)

type LightsState struct {
	lights []Light
}

// Light datastructure for light
type Light struct {
	state LightsState
}

type LightState struct {
	on         bool
	Brightness uint8  `json:"bri"`
	Hue        uint16 `json:"hue"`
	Saturation uint8  `json:"sat"`
	//"effect":"none",
	//"xy":[
	//0.3144,
	//0.3301
	//],
	//"ct":153,
	//"alert":"select",
	//"colormode":"xy",
	//"reachable":true
}

type setLightStateBody struct {
	On         bool   `json:"on,omitempty"`
	Brightness uint8  `json:"bri,omitempty"`
	Hue        uint16 `json:"hue,omitempty"`
	Saturation uint8  `json:"sat,omitempty"`
	//XY  list
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

// LightsStatus get status of all lights
func LightsStatus(bridge *Bridge) ([]Light, error) {
	if bridge == nil {
		return nil, fmt.Errorf("bridge can't be nil")
	}
	resp, err := request.GET(bridge.ip + "/api/" + bridge.username + "/lights")
	if err != nil {
		return nil, err
	}
	fmt.Println(string(resp))
	//json.Unmarshal(resp, )

	//var result LightsState
	var result map[string]interface{}
	if err = json.Unmarshal(resp, &result); err != nil {
		fmt.Printf("ERROR json: %s\n", err)
		return nil, err
	}
	fmt.Println(result)
	return nil, nil
}

// SetLightState set new state to selected light
func SetLightState(bridge *Bridge, light *Light) {
	if bridge == nil /*|| light == nil */ {
		fmt.Printf("SetLightState() invalid arguments: bridge<%p> light<%p>\n", bridge, light)
		return
	}
	var body = setLightStateBody{
		Hue:        3000,
		On:         true,
		Brightness: 200,
	}
	data, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("JSON marshaling is failing: %s", err)
	}
	fmt.Printf("about to post request %s\n", data)
	request.PUT(bridge.ip+"/api/"+bridge.username+"/lights/5/state", bytes.NewReader(data))
}

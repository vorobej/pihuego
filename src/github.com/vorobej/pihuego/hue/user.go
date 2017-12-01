package hue

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/vorobej/pihuego/hue/request"
)

const (
	deviceName = "pihuego#device"
)

type createUserBody struct {
	DeviceType string `json:"devicetype"`
}

// CreateUser creates user for HUE bridge
func CreateUser(bridge *Bridge) (string, error) {
	if bridge == nil {
		return "", fmt.Errorf("bridge can't be nil")
	}

	var body = createUserBody{DeviceType: deviceName}
	data, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("JSON marshaling is failing: %s", err)
	}
	fmt.Printf("about to post request %s\n", data)
	request.POST(bridge.IP+"/api", bytes.NewReader(data))
	return "", nil
}

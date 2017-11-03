package hue

// Bridge holds internal data for hue bridge
type Bridge struct {
	ip       string
	username string
}

const (
	// use hardcoded bridge ip instead of real ssdp discovery
	useHardcodedBridge = true
	// use hardcoded username instead of real pairing
	useHardcodedUsername = true
)

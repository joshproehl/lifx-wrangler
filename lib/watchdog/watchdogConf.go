package watchdog

import ()

type WatchdogConf struct {
	RescanSeconds          int // How many seconds between each full re-scan of the network. (Discovers new lights)
	BulbUpdateStateMillis  int // How many milliseconds between each request to the bulb to update it's state
	BulbUpdateOtherSeconds int // How many seconds between re-fetching all non-state information from the bulb

	MQTTServer      string // The URL of the server to connect to
	MQTTTopicPrefix string // Prefix all of our channel names with this
	MQTTDeviceID    string // What should this device call itself?
}

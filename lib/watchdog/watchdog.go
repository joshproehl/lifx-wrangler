package watchdog

import (
	//"fmt"
	proto "github.com/joshproehl/go-lifx/protocol"
	jww "github.com/spf13/jwalterweatherman"
	"sync"
)

// Watchdog monitors the messages coming over the LAN and keeps information about all of the Lights it hears.
// It also acts as client, and allows interaction with lights on the LAN.
type Watchdog struct {
	connected  bool
	connection *proto.Connection
	messages   <-chan proto.Message
	errors     <-chan error
	lights     *LightCollection
	lightsLock sync.RWMutex
	conf       *WatchdogConf
	confLock   sync.RWMutex
}

// NewLifxWatchdog creates a new watchdog and starts it monitoring the LAN.
func NewLifxWatchdog(c *WatchdogConf) *Watchdog {
	w := &Watchdog{conf: c}
	if conn, err := proto.Connect(); err == nil {
		w.connection = conn
		w.connected = true
	}

	messages, errors := w.connection.Listen()

	w.messages = messages
	w.errors = errors
	w.lights = NewLightCollection(w)

	w.monitorAndUpdate()

	w.SendMessage(proto.LightGet{})

	jww.INFO.Println("Watchdog up and running.")
	return w
}

// monitorAndUpdate listens to the local network and updates our state with what it hears.
func (w *Watchdog) monitorAndUpdate() {
	go func(iw *Watchdog) {

	}(w)
}

// SendMessage sends a payload out over the network and returns
// TODO: WTF are we returning here? (Copied from bjeanes/go-lifx client)
func (w *Watchdog) SendMessage(payload proto.Payload) (data []byte, error error) {
	msg := proto.Message{}
	msg.Payload = payload

	w.connection.WriteMessage(msg)
	return data, nil
}

// GetConf returns the current configuration of the Watchdog
func (w *Watchdog) GetConf() *WatchdogConf {
	w.confLock.RLock()
	defer w.confLock.RUnlock()

	return w.conf
}

// SetConf will set a new conf object.
func (w *Watchdog) SetConf(nc *WatchdogConf) {
	w.confLock.Lock()
	defer w.confLock.Unlock()

	w.conf = nc
}

// GetLightCount returns the number of lights that we're currently tracking
func (w *Watchdog) GetLightCount() int {
	w.lightsLock.RLock()
	defer w.lightsLock.RUnlock()
	return w.lights.Count()
}

// GetForSelector takes a selector string used by the LIFX HTTP API and returns the lights found by that selector
func (w *Watchdog) GetForSelector(s string) []Light {
	// TODO: Actually search...
	return w.lights.All()
}

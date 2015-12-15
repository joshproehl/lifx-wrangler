package watchdog

import (
	//"errors"
	//"fmt"
	proto "github.com/joshproehl/go-lifx/protocol"
	jww "github.com/spf13/jwalterweatherman"
	"sync"
	"time"
)

// Watchdog monitors the messages coming over the LAN and keeps information about all of the Lights it hears.
// It also acts as client, and allows interaction with lights on the LAN.
type Watchdog struct {
	connected       bool
	connection      *proto.Connection
	messages        <-chan proto.Message
	errors          <-chan error
	LightCollection *LightCollection
	conf            *WatchdogConf
	confLock        sync.RWMutex
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
	w.LightCollection = NewLightCollection(w)

	w.monitorAndUpdate()

	w.SendMessage(proto.LightGet{})

	jww.INFO.Println("Watchdog up and running.")
	return w
}

// monitorAndUpdate listens to the local network and updates our state with what it hears.
// for anything non-trivial this loop should dispatch on a new goroutine in order to let
// the monitor keep looping and not miss any packets.
func (w *Watchdog) monitorAndUpdate() {
	go func(iw *Watchdog) {
		// TODO: What if we update the conf while the watchdog is running?
		rescanTicker := time.NewTicker(time.Duration(w.GetConf().RescanSeconds) * time.Second)

		for {
			select {
			case <-rescanTicker.C:
				iw.SendMessage(proto.LightGet{})
			case err := <-(*iw).errors:
				jww.ERROR.Println("Recieved Error:", err)
			case msg := <-(*iw).messages:
				// We've recieved a message, dispatch it!
				switch p := msg.Payload.(type) {
				case *proto.LightState:
					jww.INFO.Println("Heard updated state", p, "for ip", msg.From.String())
					go (*iw.LightCollection).updateStateForIP(p, msg.From.String())
				case *proto.DeviceStateWifiInfo:
					jww.INFO.Println("Heard Wifi Info", p, "from ip", msg.From.String())
				}
			}
		}
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
	return w.LightCollection.Count()
}

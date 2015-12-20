package watchdog

import (
	//"encoding/json"
	"fmt"
	proto "github.com/joshproehl/go-lifx/protocol"
	jww "github.com/spf13/jwalterweatherman"
	"sync"
	"time"
)

// Light represents a single physical bulb.
type Light struct {
	*Watchdog
	state    *proto.LightState
	lastseen time.Time

	sync.RWMutex
}

func NewLight(w *Watchdog) *Light {
	return &Light{
		Watchdog: w,
	}
}

// Label returns the text label of the bulb.
func (l *Light) Label() string {
	l.RLock()
	defer l.RUnlock()

	return l.state.Label.String()
}

func (l *Light) LastSeen() time.Time {
	l.RLock()
	defer l.RUnlock()

	return l.lastseen
}

// SetState takes a state recieved over the network and updates the bulb's stored state accordingly.
func (l *Light) GetState() proto.LightState {
	l.RLock()
	defer l.RUnlock()

	return *l.state
}

// SetState takes a state recieved over the network and updates the bulb's stored state accordingly.
func (l *Light) SetState(s *proto.LightState) {
	l.Lock()
	defer l.Unlock()

	l.lastseen = time.Now()
	l.state = s
}

// TurnOff sends a power-off message to the device
func (l *Light) TurnOff() {
	w := *l.Watchdog
	jww.CRITICAL.Println(fmt.Sprintf("Watchdog is %v", w))
	_, err := l.Watchdog.SendMessage(proto.DeviceSetPower{Level: 0})
	if err != nil {
		panic("FAILED TO TURN OFF")
	}
}

// TurnOn sends a power-on message to the device
func (l *Light) TurnOn() {
	l.Watchdog.SendMessage(proto.DeviceSetPower{Level: 1})
}

package watchdog

import (
	//"encoding/json"
	proto "github.com/joshproehl/go-lifx/protocol"
	//jww "github.com/spf13/jwalterweatherman"
	"sync"
	"time"
)

// Light represents a single physical bulb.
type Light struct {
	*Watchdog
	state    *proto.LightState
	lastseen time.Time

	mutex sync.RWMutex
}

func NewLight() *Light {
	return &Light{
		mutex: *new(sync.RWMutex),
	}
}

// Label returns the text label of the bulb.
func (l *Light) Label() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.state.Label.String()
}

func (l *Light) LastSeen() time.Time {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.lastseen
}

// SetState takes a state recieved over the network and updates the bulb's stored state accordingly.
func (l *Light) GetState() proto.LightState {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return *l.state
}

// SetState takes a state recieved over the network and updates the bulb's stored state accordingly.
func (l *Light) SetState(s *proto.LightState) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.lastseen = time.Now()
	l.state = s
}

// TurnOff sends a power-off message to the device
func (l Light) TurnOff() {
	l.Watchdog.SendMessage(proto.DeviceSetPower{Level: 0})
}

// TurnOn sends a power-on message to the device
func (l Light) TurnOn() {
	l.Watchdog.SendMessage(proto.DeviceSetPower{Level: 1})
}

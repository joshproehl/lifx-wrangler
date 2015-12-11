package watchdog

import (
	jww "github.com/spf13/jwalterweatherman"
	"sync"
)

// LightCollection holds a group of bulbs. The intent is that it should represent all the bulbs found on the LAN.
type LightCollection struct {
	watchdog   *Watchdog
	lights     []Light
	lightsLock sync.RWMutex
}

// NewLightCollection sets up and return a new LightCollection associated with a particular watchdog.
func NewLightCollection(w *Watchdog) *LightCollection {
	jww.DEBUG.Println("New LightCollection created for watchdog", w)
	return &LightCollection{watchdog: w}
}

// Count returns the number of Lights currently in this collection
func (lc *LightCollection) Count() int {
	lc.lightsLock.RLock()
	defer lc.lightsLock.RUnlock()
	return len(lc.lights)
}

// All gives us an immutable copy of the Lights in this collection
func (lc *LightCollection) All() []Light {
	return lc.lights
}

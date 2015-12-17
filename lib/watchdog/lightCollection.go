package watchdog

import (
	"fmt"
	proto "github.com/joshproehl/go-lifx/protocol"
	jww "github.com/spf13/jwalterweatherman"
	"sync"
)

// LightCollection holds a group of bulbs. The intent is that it should represent all the bulbs found on the LAN.
type LightCollection struct {
	watchdog   *Watchdog
	lights     []*Light
	lightsLock sync.RWMutex // This mutex should protect ALL of the various ways to access lights, as they should all be considered atomic.
	lightsByIP map[string]*Light
}

// NewLightCollection sets up and return a new LightCollection associated with a particular watchdog.
func NewLightCollection(w *Watchdog) *LightCollection {
	jww.DEBUG.Println("New LightCollection created for watchdog", w)
	return &LightCollection{watchdog: w, lightsByIP: map[string]*Light{}}
}

// Count returns the number of Lights currently in this collection
func (lc *LightCollection) Count() int {
	lc.lightsLock.RLock()
	defer lc.lightsLock.RUnlock()

	return len(lc.lights)
}

// All returns pointers to to every light *CURRENTLY* being managed. Note that it does NOT return a pointer to the
// managed array, so if a new light is added later, whatever holds this pointer won't get that update.
func (lc *LightCollection) All() []*Light {
	lc.lightsLock.RLock()
	defer lc.lightsLock.RUnlock()

	lightCopy := append([]*Light(nil), lc.lights...)

	return lightCopy
}

// updateForStateMessage takes a state message, finds the bulb in this collection, and updates it's values
func (lc *LightCollection) updateStateForIP(m *proto.LightState, ip string) {
	l := lc.GetOrCreateLightForIP(ip)
	// TODO: Only setstate and publish if deepequal of existing state and new state are not the same.
	l.SetState(m)

	lc.watchdog.mqttPublish(fmt.Sprintf("bulbs/byip/state/%s", ip), fmt.Sprintf("GotState for %s", m.Label))
}

func (lc *LightCollection) RefreshBulbStates() {

}

// LightForIP gets the pointer to the light for given IP string. If no light is found it creates a new one in the collection.
func (lc *LightCollection) GetOrCreateLightForIP(ips string) *Light {
	lc.lightsLock.RLock()
	l, found := lc.lightsByIP[ips]
	lc.lightsLock.RUnlock()
	if found {
		return l
	} else {
		lc.lightsLock.Lock() // TODO: We have a potential race condition here if something ELSE is trying to access the lock and gets inbetween these two calls!

		nl := NewLight(lc.watchdog)
		lc.lights = append(lc.lights, nl)
		lc.lightsByIP[ips] = nl
		lc.lightsLock.Unlock()
		return nl
	}
}

// GetForLabel returns immutable copies of the lights which have the given label.
// Ideally this should only return a single light, but we want to account for the possibility of duplicate labels, confusing though that may make things.
func (lc *LightCollection) GetForLabel(lbl string) []*Light {
	ret := make([]*Light, 0)

	lc.lightsLock.RLock()
	defer lc.lightsLock.RUnlock()

	for _, l := range lc.lights {
		if l.Label() == lbl {
			ret = append(ret, l)
		}
	}

	return ret
}

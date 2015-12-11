package lifxapiv1

import (
	"github.com/go-zoo/bone"
	wd "github.com/joshproehl/lifx-wrangler/lib/watchdog"
	jww "github.com/spf13/jwalterweatherman"
	"net/http"
	//"strings"
)

type v1 struct {
	watchdog *wd.Watchdog
}

func NewLifxCloudV1APICloneHandlers(w *wd.Watchdog) *v1 {
	return &v1{watchdog: w}
}

// Handle a request tho the root reesource. No results found here. TODO: Clone whatever the HTTP API does for this
func (v *v1) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{\"error\":\"Please visit one of the API URL's\"}"))
}

/*
  {
    "id": "d3b2f2d97452",
    "uuid": "8fa5f072-af97-44ed-ae54-e70fd7bd9d20",
    "label": "Left Lamp",
    "connected": true,
    "power": "on",
    "color": {
      "hue": 250.0,
      "saturation": 0.5,
      "kelvin": 3500
    },
    "brightness": 0.5,
    "group": {
      "id": "1c8de82b81f445e7cfaafae49b259c71",
      "name": "Lounge"
    },
    "location": {
      "id": "1d6fe8ef0fde4c6d77b0012dc736662c",
      "name": "Home"
    },
    "last_seen": "2015-03-02T08:53:02.867+00:00",
    "seconds_since_seen": 0.002869418,
    "product": {
      "name": "Original 1000",
      "company": "LIFX",
      "identifier": "lifx_original_1000",
      "capabilities": {
        "has_color": true,
        "has_variable_color_temp": true
      }
    }
  }
*/

// LightsHandler handles the /lights/:selector route. It returns all lights for the given selector
func (v *v1) LightsHandler(w http.ResponseWriter, r *http.Request) {
	selector := bone.GetValue(r, "selector")
	res := v.watchdog.GetForSelector(selector)
	jww.INFO.Println("LightsHandler() got", len(res), "lights for selector", selector)
	w.Write([]byte("{\"error\":\"Can't write currently...\"}"))
}

// LightsStateHandler handles the /lights/:selector/state route.
func (v *v1) LightsStateHandler(w http.ResponseWriter, r *http.Request) {

}

// LightsToggleHandler handles the /lights/:selector/toggle route. It toggles the power state of the select lights
func (v *v1) LightsToggleHandler(w http.ResponseWriter, r *http.Request) {

}

// LightsBreatheHandler handles the /lights/:selector/effects/breathe" route.
func (v *v1) LightsBreatheHandler(w http.ResponseWriter, r *http.Request) {

}

// LightsPulseHandler handles the /lights/:selector/effects/pulse" route.
func (v *v1) LightsPulseHandler(w http.ResponseWriter, r *http.Request) {

}

// LightsCycleHandler handles the /lights/:selector/cycle" route.
func (v *v1) LightsCycleHandler(w http.ResponseWriter, r *http.Request) {

}

// LightsMultiStateHandler handles the /lights/states route.
func (v *v1) LightsMultiStateHandler(w http.ResponseWriter, r *http.Request) {

}

// ScenesHandler handles the /scenes route.
func (v *v1) ScenesHandler(w http.ResponseWriter, r *http.Request) {

}

// SceneActivateHandler handles the /scenes/scene_id::scene_uuid/activate route.
func (v *v1) SceneActivateHandler(w http.ResponseWriter, r *http.Request) {

}

// ColorHandler handles the /color route
func (v *v1) ColorHandler(w http.ResponseWriter, r *http.Request) {

}

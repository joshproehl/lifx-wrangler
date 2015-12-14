package lifxapiv1

import (
	"bytes"
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

// LightsHandler handles the /lights/:selector route. It returns all lights for the given selector
func (v *v1) LightsHandler(w http.ResponseWriter, r *http.Request) {
	selector := bone.GetValue(r, "selector")
	res := v.watchdog.GetForSelector(selector)
	jww.INFO.Println("LightsHandler() got", len(res), "lights for selector", selector)

	buf := new(bytes.Buffer)
	buf.WriteString("[")
	resCount := len(res)
	for i, l := range res {
		buf.WriteString(lightsHandlerJSON(l))
		if i < resCount-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("]")
	w.Write(buf.Bytes())
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

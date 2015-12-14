package lifxapiv1

import (
	"bytes"
	//"time"
	//proto "github.com/joshproehl/go-lifx/lib/protocol"
	wd "github.com/joshproehl/lifx-wrangler/lib/watchdog"
)

const (
	LIFXTimeFormat = "2006-01-02T15:04:05.999+07:00" // Define the ISO 8601 format used in the LIFX API.
)

/*
	lightsHandlerJSON takes a watchdog.Light object, and marshals it into the specific JSON required to mimic the lifx v1 api

	Example JSON object:
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
func lightsHandlerJSON(l wd.Light) string {
	s := l.GetState()
	buf := new(bytes.Buffer)

	buf.WriteString("{")
	buf.WriteString("\"id\": \"\",")
	buf.WriteString("\"uuid\": \"\",")
	buf.WriteString("\"label\": \"" + s.Label.String() + "\",")
	buf.WriteString("\"connected\": \"\",")
	buf.WriteString("\"power\": \"\",")
	buf.WriteString("\"color\": {")
	buf.WriteString("\"hue\":\"" + s.Color.Hue.String() + "\",")
	buf.WriteString("\"saturation\":\"" + s.Color.Saturation.String() + "\",")
	buf.WriteString("\"kelvin\":\"" + s.Color.Kelvin.String() + "\"")
	buf.WriteString("},")
	buf.WriteString("\"brightness\":\"" + "" + "\",")
	buf.WriteString("\"group\": {")
	buf.WriteString("\"id\":\"" + "" + "\",")
	buf.WriteString("\"name\":\"" + "" + "\"")
	buf.WriteString("},")
	buf.WriteString("\"location\": {")
	buf.WriteString("\"id\":\"" + "" + "\",")
	buf.WriteString("\"name\":\"" + "" + "\"")
	buf.WriteString("},")
	buf.WriteString("\"last_seen\":\"" + l.LastSeen().Format(LIFXTimeFormat) + "\",")
	buf.WriteString("\"seconds_since_seen\":1,") // time.Now() - l.LastSeen()
	buf.WriteString("\"product\": {")
	buf.WriteString("\"name\":\"" + "" + "\",")
	buf.WriteString("\"company\":\"" + "" + "\",")
	buf.WriteString("\"identifier\":\"" + "" + "\",")
	buf.WriteString("\"capabilities\": {")
	buf.WriteString("}")
	buf.WriteString("}")
	buf.WriteString("}")

	return buf.String()
}

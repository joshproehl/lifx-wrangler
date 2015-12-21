package lifxapiv1

import (
	"bytes"
	//"time"
	"github.com/pdf/golifx/common"
	jww "github.com/spf13/jwalterweatherman"
	"strconv"
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
func lightsHandlerJSON(l common.Light) string {
	buf := new(bytes.Buffer)

	lbl, err := l.GetLabel()
	if err != nil {
		jww.ERROR.Println("lightsHandlerJSON error getting label:", err)
	}

	buf.WriteString("{")
	buf.WriteString("\"id\": \"" + strconv.Itoa(int(l.ID())) + "\",")
	buf.WriteString("\"uuid\": \"\",")
	buf.WriteString("\"label\": \"" + lbl + "\",")
	buf.WriteString("\"connected\": \"\",")
	buf.WriteString("\"power\": \"\",")
	buf.WriteString("\"color\": {")
	buf.WriteString("\"hue\":\"" + strconv.Itoa(int(l.CachedColor().Hue)) + "\",")
	buf.WriteString("\"saturation\":\"" + strconv.Itoa(int(l.CachedColor().Saturation)) + "\",")
	buf.WriteString("\"kelvin\":\"" + strconv.Itoa(int(l.CachedColor().Kelvin)) + "\"")
	buf.WriteString("},")
	buf.WriteString("\"brightness\":\"" + strconv.Itoa(int(l.CachedColor().Brightness)) + "\",")
	buf.WriteString("\"group\": {")
	buf.WriteString("\"id\":\"" + "" + "\",")
	buf.WriteString("\"name\":\"" + "" + "\"")
	buf.WriteString("},")
	buf.WriteString("\"location\": {")
	buf.WriteString("\"id\":\"" + "" + "\",")
	buf.WriteString("\"name\":\"" + "" + "\"")
	buf.WriteString("},")
	buf.WriteString("\"last_seen\":\"" + "" + "\",")
	//buf.WriteString("\"last_seen\":\"" + l.LastSeen().Format(LIFXTimeFormat) + "\",")
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

/*
	lightsToggleHandlerJSON takes a watchdog.Light object, and marshals it into the specific JSON required to mimic the lifx v1 api toggle route

	Example JSON object:
	{
		"results": [
			{
				"id": "d3b2f2d97452",
				"label": "Left Lamp",
				"status": "ok"
			}
		]
	}

	Note that we're actually creating one of the objects from inside the results array here.
*/
func lightsToggleHandlerJSON(l common.Light) string {
	lbl, err := l.GetLabel()
	if err != nil {
		jww.ERROR.Println("lightsToggleHandlerJSON error getting label:", err)
	}

	buf := new(bytes.Buffer)

	buf.WriteString("{")
	buf.WriteString("\"id\": \"\",")
	buf.WriteString("\"label\": \"" + lbl + "\",")
	buf.WriteString("\"status\": \"\"")
	buf.WriteString("}")

	return buf.String()
}

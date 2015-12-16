package main

import (
	"flag"
	"fmt"
	"github.com/go-zoo/bone"
	"github.com/joshproehl/lifx-wrangler/lib/lifxapiv1"
	wd "github.com/joshproehl/lifx-wrangler/lib/watchdog"
	jww "github.com/spf13/jwalterweatherman"
	"net/http"
	"time"
)

var watchdog *wd.Watchdog

func main() {
	// Set up command line flags
	flgVerbose := flag.Bool("verbose", false, "Output additional debugging information to both STDOUT and the log file")
	flgPortNum := flag.Int("port", 7688, "The port to run the HTTP server on.") // 7688 = "LX"
	flgUpdateMillis := flag.Int("updateMillis", 1000, "How many milliseconds between each re-check of the bulb. Should not be lower than 50. (20 requests per second)")
	flgMQTTServer := flag.String("mqttServer", "", "The MQTT URL (tcp://localhost:1833) to use. If not set MQTT will be disabled.")
	flgMQTTTopicPrefix := flag.String("mqttTopicPrefix", "lifx-wrangler/", "All our mqtt topics get this prefix")
	flgMQTTDeviceID := flag.String("mqttDeviceId", "lifx-wrangler", "The MQTT Device ID to use")
	flag.Parse()

	// Note at this point only WARN or above is actually logged to file, and ERROR or above to console.
	jww.SetLogFile("lifx-wrangler.log")

	// Set extra logging if the command line flag was set
	if *flgVerbose {
		jww.SetLogThreshold(jww.LevelDebug)
		jww.SetStdoutThreshold(jww.LevelInfo)
		jww.INFO.Println("Verbose debug level set.")
	} else {
		// Set custom default logging verbosity.
		jww.SetLogThreshold(jww.LevelWarn)
		jww.SetStdoutThreshold(jww.LevelError)
	}

	jww.INFO.Println("Starting run at", time.Now().Format("2006-01-02 15:04:05"))

	// Sanity checking inputs
	if *flgUpdateMillis < 50 {
		*flgUpdateMillis = 50
		jww.CRITICAL.Println("UpdateMillis was set to low, resetting to 50.")
	}
	if (*flgMQTTTopicPrefix)[len(*flgMQTTTopicPrefix)-1:] != "/" {
		// TODO: Append trailing / to topic prefix
	}

	// Build a configuration and get a new watchdog for it.
	watchdog = wd.NewLifxWatchdog(&wd.WatchdogConf{
		RescanSeconds:          30,
		BulbUpdateStateMillis:  *flgUpdateMillis,
		BulbUpdateOtherSeconds: 10, // TODO: Tune this.
		MQTTServer:             *flgMQTTServer,
		MQTTTopicPrefix:        *flgMQTTTopicPrefix,
		MQTTDeviceID:           *flgMQTTDeviceID,
	})

	// Set up the HTTP router, followed by all the routes
	router := bone.New()

	/*
		// Redirect static resources, and then handle the static resources (/gui/) routes with the static asset file
		router.Handle("/", http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			http.Redirect(response, request, "/gui/", 302)
		}))
		router.Get("/gui/", http.StripPrefix("/gui/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: ""})))
	*/
	router.GetFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Lifx-Wrangler!"))
	})

	// Define the API (JSON) routes
	v1 := lifxapiv1.NewLifxCloudV1APICloneHandlers(watchdog)
	router.GetFunc("/api/lfxc/v1", v1.RootHandler)

	// Define the /v1/ API routes from the lifx public API to mimic
	router.GetFunc("/api/lfxc/v1/lights/:selector", v1.LightsHandler)
	router.PutFunc("/api/lfxc/v1/lights/:selector/state", v1.LightsStateHandler)
	router.PostFunc("/api/lfxc/v1/lights/:selector/toggle", v1.LightsToggleHandler)
	router.PostFunc("/api/lfxc/v1/lights/:selector/effects/breathe", v1.LightsBreatheHandler)
	router.PostFunc("/api/lfxc/v1/lights/:selector/effects/pulse", v1.LightsPulseHandler)
	router.PostFunc("/api/lfxc/v1/lights/:selector/cycle", v1.LightsCycleHandler)
	router.PutFunc("/api/lfxc/v1/lights/states", v1.LightsMultiStateHandler)
	router.GetFunc("/api/lfxc/v1/scenes", v1.ScenesHandler)
	router.PutFunc("/api/lfxc/v1/scenes/scene_id::scene_uuid/activate", v1.SceneActivateHandler)
	router.GetFunc("/api/lfxc/v1/color", v1.ColorHandler)

	// Start the HTTP server
	fmt.Println("Starting API server on port", *flgPortNum, ". Press Ctrl-C to quit.")
	http.ListenAndServe(fmt.Sprintf(":%d", *flgPortNum), router)
}

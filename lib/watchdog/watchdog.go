package watchdog

import (
	//"errors"
	"fmt"
	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	proto "github.com/joshproehl/go-lifx/protocol"
	jww "github.com/spf13/jwalterweatherman"
	"strings"
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
	mqttClient      *mqtt.Client
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

	if c.MQTTServer != "" {
		opts := mqtt.NewClientOptions().AddBroker(c.MQTTServer).SetClientID(c.MQTTDeviceID).SetCleanSession(true)
		// TODO: Subscribe to correct topics and handle errors
		opts.OnConnect = func(c *mqtt.Client) {
			if token := c.Subscribe("topic", 1, mqttMessageReceived); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
		}

		w.mqttClient = mqtt.NewClient(opts)

		if token := w.mqttClient.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	// Start the watchdog listening for packets from bulbs!
	w.monitorAndUpdate()

	w.SendMessage(proto.LightGet{})

	jww.INFO.Println("Watchdog up and running.")
	return w
}

// Shutdown stops the watchdog and closes all resources.
func (w *Watchdog) Shutdown() {
	w.mqttClient.Disconnect(250)
}

func mqttMessageReceived(client *mqtt.Client, msg mqtt.Message) {
	topics := strings.Split(msg.Topic(), "/")
	msgFrom := topics[len(topics)-1]
	fmt.Print(msgFrom + ": " + string(msg.Payload()))
}

func (w *Watchdog) mqttPublish(topic string, message string) error {
	if w.mqttClient == nil {
		return fmt.Errorf("Watchdog has no MQTT client")
	}

	fullTopic := fmt.Sprintf("%s%s", w.conf.MQTTTopicPrefix, topic) // TODO: This isn't how we should access conf...
	if token := w.mqttClient.Publish(fullTopic, 1, false, message); token.Wait() && token.Error() != nil {
		fmt.Println("Failed to send message")
	}
	return nil
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
func (w *Watchdog) SendMessage(payload proto.Payload) ([]byte, error) {
	msg := proto.Message{}
	msg.Payload = payload
	//data := make([]byte, 0)

	err := w.connection.WriteMessage(msg)
	return nil, err
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

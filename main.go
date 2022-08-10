package main

import (
	"fmt"
	"time"
	"os"
	"encoding/json"
	"github.com/VictoriaMetrics/metrics"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/syoder89/tasmota-monitor/vmclient"
)

type TasmotaMsg struct {
	Time string
	ENERGY struct {
		TotalStartTime string
		Power float64
		ApparentPower float64
		ReactivePower float64
		Factor float64
		Voltage float64
		Current float64
	}
}

var tmsg TasmotaMsg
var sensor
// tcp://mosquitto
var broker = "tcp://mosquitto:1883"
// http://172.20.1.4:8428/api/v1/import/prometheus
var vmPushURL = "http://victoria-metrics-victoria-metrics-single-server:8428/api/v1/import/prometheus"

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	json.Unmarshal([]byte(msg.Payload()), &tmsg)
	fmt.Println(tmsg)
	vmclient.Push(vmPushURL, 20*time.Second, `sensor="`+sensor+`"`, false)
}

func main() {
	if val, ok := os.LookupEnv("SENSOR"); ok {
		sensor = val
	}
	else {
		panic("No sensor name provided!")
	}

	if val, ok := os.LookupEnv("BROKER"); ok {
		broker = val
	}
	if val, ok := os.LookupEnv("VM_PUSH_URL"); ok {
		vmPushURL = val
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("monitor-"+sensor)
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	metrics.NewGauge(`power`, func() float64 { return tmsg.ENERGY.Power })
	metrics.NewGauge(`apparent_power`, func() float64 { return tmsg.ENERGY.ApparentPower })
	metrics.NewGauge(`reactive_power`, func() float64 { return tmsg.ENERGY.ReactivePower })
	metrics.NewGauge(`power_factor`, func() float64 { return tmsg.ENERGY.Factor })
	metrics.NewGauge(`voltage`, func() float64 { return tmsg.ENERGY.Voltage })
	metrics.NewGauge(`current`, func() float64 { return tmsg.ENERGY.Current })

	sub(client)

	for true {
		time.Sleep(time.Second)
	}
}

// Received message: {"Time":"2022-08-07T02:39:55","ENERGY":{"TotalStartTime":"2022-08-02T20:37:49","Total":0.006,"Yesterday":0.000,"Today":0.000,"Period": 0,"Power": 0,"ApparentPower": 0,"ReactivePower": 0,"Factor":0.00,"Voltage":121,"Current":0.000}} from topic: tele/taylor_water/SENSOR

func sub(client mqtt.Client) {
//	topic := "#"
	topic := "tele/"+sensor+"/SENSOR"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}

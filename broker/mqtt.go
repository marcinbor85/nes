package broker

import (
	"time"
	"math/rand"
	"github.com/eclipse/paho.mqtt.golang"

	"github.com/marcinbor85/nes/protocol"
)

type FrameHandler func(client *Client, frame *protocol.Frame)

type Client struct {
	BrokerAddress	string
	Handler			mqtt.MessageHandler
	MqttClient		mqtt.Client
	OnFrame			FrameHandler
}

func GenerateClientID(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func messageTopic(recipient string) string {
	topic := "nes/" + recipient + "/message"
	return topic
}

func wait(token mqtt.Token) error {
	token.Wait()
	return token.Error()
}

func (client *Client) Connect() error {
	rand.Seed(time.Now().UnixNano())

	clientID := GenerateClientID(8)

	opts := mqtt.NewClientOptions().AddBroker(client.BrokerAddress).SetClientID(clientID)

	onMessage := func(c mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload())
		frame, err := protocol.FrameFromString(payload)
		if err != nil {
			return
		}

		if client.OnFrame != nil {
			client.OnFrame(client, frame)
		}
	}

	opts.SetDefaultPublishHandler(onMessage)

	mqttClient := mqtt.NewClient(opts)
	token := mqttClient.Connect()

	err := wait(token)
	if err != nil {
		return err
	}

	client.MqttClient = mqttClient
	return nil
}

func (client *Client) Recv(recipient string) error {
	topic := messageTopic(recipient)
	token := client.MqttClient.Subscribe(topic, 2, nil);

	err := wait(token)
	return err
}

func (client *Client) Send(frame *protocol.Frame, recipient string) error {
	message := frame.String()
	topic := messageTopic(recipient)
	token := client.MqttClient.Publish(topic, 2, false, message)

	err := wait(token)
	return err
}

func (client *Client) Disconnect() {
	client.MqttClient.Disconnect(250)
}

/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: mqtt client
 *
 */

package module

import (
	"fmt"
	"time"

	"clc.hmu/app/public/log"
	"clc.hmu/app/public/log/elog"
	"clc.hmu/app/public/sys"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

// MQTTClient mqtt client
type MQTTClient struct {
	Client MQTT.Client
}

// NewMQTTClient new mqtt client
func NewMQTTClient(handler MQTT.MessageHandler, willtopic, willpayload, conntopic, connpayload string) MQTTClient {
	// 每次启动都使用一个新的client进行连接
	// config.Configuration.MQTT.ClientID
	// "The ClientID (optional)"
	id := uuid.New().String()

	cfg := sys.GetBusManagerCfg()

	// broker := config.Configuration.MQTT.Broker          // "The broker URI. ex: tcp://10.10.1.1:1883"
	broker := "tcp://" + cfg.MQTT.Host + ":" + cfg.MQTT.Port
	user := cfg.MQTT.User              // "The User (optional)"
	password := cfg.MQTT.Password      // "The password (optional)"
	cleansess := cfg.MQTT.CleanSession // "Set Clean Session (default false)"
	store := cfg.MQTT.Store            // "The Store Directory (default use memory store)"

	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(id)
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetCleanSession(cleansess)
	opts.SetStore(MQTT.NewFileStore(store))
	opts.SetAutoReconnect(true)
	opts.SetDefaultPublishHandler(handler)
	opts.SetOnConnectHandler(func(client MQTT.Client) {
		token := client.Publish(conntopic, cfg.MQTT.Qos, true, connpayload)
		token.Wait()

		// set led interval
		LEDSetInterval = 1000
		SetNetworkStatus(Online, true)
		elog.LOG.Infof("connect mqtt server success, publish status result:{%v}", token.Error())
	})
	opts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		LEDSetInterval = 50
		SetNetworkStatus(Offline, true)
		elog.LOG.Info("disconnect mqtt server")
	})
	opts.SetWill(willtopic, willpayload, cfg.MQTT.Qos, true)

	return MQTTClient{Client: MQTT.NewClient(opts)}
}

// ConnectServer connect server
func (c *MQTTClient) ConnectServer() error {
	if c.Client == nil {
		return fmt.Errorf("mqtt client unavailable")
	}

	token := c.Client.Connect()
	token.Wait()

	return token.Error()
}

// ReconnectServer reconnect server
func (c *MQTTClient) ReconnectServer() error {
	retry := time.NewTicker(5 * time.Second)
RetryLoop:
	for {
		select {
		case <-retry.C:
			if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
				//handle error
				log.Info("retry connect mqtt fail")
			} else {
				log.Info("retry connect mqtt success")
				retry.Stop()
				break RetryLoop
			}
		}
	}

	return nil
}

// DisconnectServer disconnect
func (c *MQTTClient) DisconnectServer() error {
	if c.Client == nil {
		return fmt.Errorf("mqtt client unavailable")
	}

	c.Client.Disconnect(250)

	return nil
}

// PublishSampleValues public messages
func (c *MQTTClient) PublishSampleValues(topic, payload string) error {
	if c.Client == nil {
		return fmt.Errorf("mqtt client unavailable")
	}

	cfg := sys.GetBusManagerCfg()

	// log.Debug("PublishSamleValues:", topic, payload)
	token := c.Client.Publish(topic, cfg.MQTT.Qos, true, payload)
	token.Wait()

	return token.Error()
}

// Subscribe subscribe
func (c *MQTTClient) Subscribe(topic string) error {
	if c.Client == nil {
		return fmt.Errorf("mqtt client unavailable")
	}

	log.Debug("Subscribe:", topic)

	token := c.Client.Subscribe(topic, sys.GetBusManagerCfg().MQTT.Qos, nil)
	token.Wait()

	return token.Error()
}

/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2019/05/15
 * Despcription: frpc ini
 *
 */

package web

import (
	"fmt"

	ini "gopkg.in/ini.v1"
)

const (
	sectionCommon = "common"
	sectionSSH    = "ssh"

	keyServerAddress = "server_addr"
	keyServerPort    = "server_port"

	keyType       = "type" // default "tcp"
	keyLocalIP    = "local_ip"
	keyLocalPort  = "local_port"
	keyRemotePort = "remote_port"

	typeTCP = "tcp"
)

// IniParser ini
type IniParser struct {
	reader *ini.File // config reader
}

// Load load
func (ip *IniParser) Load(filename string) error {
	conf, err := ini.Load(filename)
	if err != nil {
		ip.reader = nil
		return err
	}

	ip.reader = conf
	return nil
}

// Value get value
func (ip *IniParser) Value(section string, key string) string {
	if ip.reader == nil {
		return ""
	}

	s := ip.reader.Section(section)
	if s == nil {
		return ""
	}

	return s.Key(key).String()
}

// SetValue set value
func (ip *IniParser) SetValue(section, key, value string) {
	if ip.reader == nil {
		return
	}

	ip.reader.Section(section).Key(key).SetValue(value)
}

// DeleteSection delete section
func (ip *IniParser) DeleteSection(section string) {
	if ip.reader == nil {
		return
	}

	ip.reader.DeleteSection(section)
}

// SaveTo save to file
func (ip *IniParser) SaveTo(filename string) error {
	if ip.reader == nil {
		return fmt.Errorf("reader do not init")
	}

	return ip.reader.SaveTo(filename)
}

// Agency agency config
type Agency struct {
	Name       string `json:"name"`
	LocalIP    string `json:"localIP"`
	LocalPort  string `json:"localPort"`
	RemotePort string `json:"remotePort"`
}

// FrpcMapConfig frpc map config
type FrpcMapConfig struct {
	ServerIP   string   `json:"serverIP"`
	ServerPort string   `json:"serverPort"`
	Agencys    []Agency `json:"agencys"`
}

// Read read config
func (ip *IniParser) Read() FrpcMapConfig {
	var cfg FrpcMapConfig

	if ip.reader == nil {
		return cfg
	}

	sections := ip.reader.Sections()
	ls := len(sections)
	for i := 1; i < ls; i++ {
		secname := sections[i].Name()
		lk := len(sections[i].Keys())

		if secname == sectionCommon {
			for j := 0; j < lk; j++ {
				k := sections[i].Keys()[j].Name()
				v := sections[i].Keys()[j].Value()

				switch k {
				case keyServerAddress:
					cfg.ServerIP = v
				case keyServerPort:
					cfg.ServerPort = v
				}
			}
		} else {
			var agc Agency
			agc.Name = secname

			for j := 0; j < lk; j++ {
				k := sections[i].Keys()[j].Name()
				v := sections[i].Keys()[j].Value()

				switch k {
				case keyLocalIP:
					agc.LocalIP = v
				case keyLocalPort:
					agc.LocalPort = v
				case keyRemotePort:
					agc.RemotePort = v
				}
			}

			cfg.Agencys = append(cfg.Agencys, agc)
		}
	}

	return cfg
}

// SaveFrpConfig save config
func (ip *IniParser) SaveFrpConfig(cfg FrpcMapConfig, filename string) error {
	if ip.reader == nil {
		return fmt.Errorf("reader do not init")
	}

	sections := ip.reader.Sections()
	ls := len(sections)
	for i := 1; i < ls; i++ {
		ip.reader.DeleteSection(sections[i].Name())
	}

	ip.reader.Section(sectionCommon).Key(keyServerAddress).SetValue(cfg.ServerIP)
	ip.reader.Section(sectionCommon).Key(keyServerPort).SetValue(cfg.ServerPort)

	for _, c := range cfg.Agencys {
		ip.reader.Section(c.Name).Key(keyType).SetValue(typeTCP)
		ip.reader.Section(c.Name).Key(keyLocalIP).SetValue(c.LocalIP)
		ip.reader.Section(c.Name).Key(keyLocalPort).SetValue(c.LocalPort)
		ip.reader.Section(c.Name).Key(keyRemotePort).SetValue(c.RemotePort)
	}

	// default ssh, do not change
	muid := getUUID()
	if len(muid) == 13 {
		remoteport := muid[6:8] + muid[10:13]
		ip.reader.Section(sectionSSH).Key(keyType).SetValue(typeTCP)
		ip.reader.Section(sectionSSH).Key(keyLocalIP).SetValue("127.0.0.1")
		ip.reader.Section(sectionSSH).Key(keyLocalPort).SetValue("22")
		ip.reader.Section(sectionSSH).Key(keyRemotePort).SetValue(remoteport)
	}

	return ip.SaveTo(filename)
}

// AddFrpcCommonSection add common section
func (ip *IniParser) AddFrpcCommonSection(host, port, filename string) error {
	if ip.reader == nil {
		return fmt.Errorf("reader do not init")
	}

	ip.reader.Section(sectionCommon).Key(keyServerAddress).SetValue(host)
	ip.reader.Section(sectionCommon).Key(keyServerPort).SetValue(port)

	return ip.SaveTo(filename)
}

// AddFrpcAgencySection add agency section
func (ip *IniParser) AddFrpcAgencySection(section, host, localport, remoteport, filename string) error {
	if ip.reader == nil {
		return fmt.Errorf("reader do not init")
	}

	ip.reader.Section(section).Key(keyType).SetValue(typeTCP)
	ip.reader.Section(section).Key(keyLocalIP).SetValue(host)
	ip.reader.Section(section).Key(keyLocalPort).SetValue(localport)
	ip.reader.Section(section).Key(keyRemotePort).SetValue(remoteport)

	return ip.SaveTo(filename)
}

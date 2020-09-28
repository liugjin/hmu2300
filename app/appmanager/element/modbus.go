/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/04
 * Despcription: channel mapping
 *
 */

package element

// Channel channel
type Channel struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	DataType string      `json:"datatype"`
	Value    interface{} `json:"value"`
}

// Setting setting
type Setting struct {
	Throttle  int     `json:"throttle"`
	Optimized bool    `json:"optimized"`
	StepDiff  int     `json:"stepDiff"`
	MaxDiff   int     `json:"maxDiff"`
	COV       float64 `json:"cov"`
}

// ChannelMapping channel mapping
type ChannelMapping struct {
	// modbus
	Code     interface{} `json:"code,omitempty"`
	Address  int32       `json:"address"`
	Quantity int32       `json:"quantity,omitempty"`
	Format   interface{} `json:"format,omitempty"`

	// pmbus
	CID1         byte   `json:"cid1,omitempty"`
	CID2         byte   `json:"cid2,omitempty"`
	COMMAND      uint16 `json:"command"`
	Offset       int    `json:"offset"`
	Length       int    `json:"length,omitempty"`
	CommandGroup byte   `json:"commandgroup,omitempty"`
	CommandType  byte   `json:"commandtype,omitempty"`

	// oid
	OID string `json:"oid,omitempty"`

	// entrance guard
	Sequence string `json:"seqno,omitempty"`
	Group    int    `json:"group,omitempty"`
	Type     string `json:"type,omitempty"`

	Expression string `json:"expression"`
	ChannelID  string `json:"channel"`

	// channel cov
	COV float64 `json:"cov"`
}

// Mapping mapping
type Mapping struct {
	Protocol        string           `json:"protocol"`
	Type            string           `json:"type"`
	Setting         Setting          `json:"setting"`
	ChannelMappings []ChannelMapping `json:"mapping"`
}

// Element modbus element
type Element struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
	Channels    []Channel `json:"channels"`
	Mappings    []Mapping `json:"mappings"`
}

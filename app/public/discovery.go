/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: discovery difinition
 *
 */

package public

// Discovery discovery
type Discovery struct {
	// ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Vendor      string `json:"vendor"`
	Model       string `json:"model"`
	Station     string `json:"station"`
	StationName string `json:"stationName"`
	Project     string `json:"project"`
	User        string `json:"user"`
}

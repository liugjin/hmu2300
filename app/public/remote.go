package public

// RemoteInfo remote info
type RemoteInfo struct {
	Operation string `json:"operation"`
	Host      string `json:"host"`
	Port      string `json:"port"`
}

// RemotePayload remote payload
type RemotePayload struct {
	MonitoringUnit string `json:"monitoringUnit"`
	SampleUnit     string `json:"sampleUnit"`
	Channel        string `json:"channel"`
	Parameters     struct {
		Value string `json:"value"`
	} `json:"parameters"`
	Phase     string `json:"phase"`
	Timeout   int    `json:"timeout"`
	Operator  string `json:"operator"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Result    string `json:"result"`
}

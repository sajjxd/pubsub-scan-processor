package types

import "time"

const (
	_ = iota // Ignore 0
	V1
	V2
)

type Scan struct {
	Ip          string      `json:"ip"`
	Port        uint32      `json:"port"`
	Service     string      `json:"service"`
	Timestamp   int64       `json:"timestamp"`
	DataVersion int         `json:"data_version"`
	Data        interface{} `json:"data"`
}

type V1Data struct {
	ResponseBytesUtf8 []byte `json:"response_bytes_utf8"`
}

type V2Data struct {
	ResponseStr string `json:"response_str"`
}

type ScanRecord struct {
	Ip          string
	Port        uint32
	Service     string
	Response    string
	LastScanned time.Time
}

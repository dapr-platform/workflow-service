package entity

import "github.com/dapr-platform/common"

type DeviceMirror struct {
	State                 DeviceMirrorState     `json:"state"`
	Metadata              DeviceMirrorMetadata  `json:"metadata"`
	Timestamp             int64                 `json:"timestamp"`
	TimestampStr          string                `json:"timestamp_str"`
	LatestRawMsg          string                `json:"latest_raw_msg"`
	Alerts                []map[string]any      `json:"alerts"`
	RecentEvents          []DeviceMirrorEvent   `json:"recent_events"`
	CurrentKpiDatas       []CurrentKpiData      `json:"current_kpi_datas"`
	Recent_InvokeServices []DeviceMirrorService `json:"recent_invoke_services"`
	Version               int64                 `json:"version"`
}
type DeviceMirrorEvent struct {
	Timestamp    int64  `json:"timestamp"`
	TimestampStr string `json:"timestamp_str"`
	Identifier   string `json:"identifier"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Desc         string `json:"desc"`
	Method       string `json:"method"`
	Required     bool   `json:"required"`
	OutputData   any    `json:"output_data"`
}
type DeviceMirrorService struct {
	Timestamp    int64  `json:"timestamp"`
	TimestampStr string `json:"timestamp_str"`
	Identifier   string `json:"identifier"`
	Name         string `json:"name"`
	CallType     string `json:"call_type"`
	Method       string `json:"method"`
	Desc         string `json:"desc"`
	Required     bool   `json:"required"`
	InputData    any    `json:"input_data"`
	OutputData   any    `json:"output_data"`
}
type CurrentKpiData struct {
	Name  string           `json:"name"`
	Id    string           `json:"id"`
	Ts    common.LocalTime `json:"ts"`
	Value any              `json:"value"`
	Unit  string           `json:"unit"`
}

type DeviceMirrorState struct {
	Reported map[string]interface{} `json:"reported"`
	Desired  map[string]interface{} `json:"desired"`
}

type DeviceMirrorMetadata struct {
	Reported map[string]interface{} `json:"reported"`
	Desired  map[string]interface{} `json:"desired"`
}

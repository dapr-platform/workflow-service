package flow_dsl

import (
	"encoding"
	"encoding/json"
	"github.com/spf13/cast"
)

var _ encoding.BinaryMarshaler = new(SecondInterval)
var _ encoding.BinaryUnmarshaler = new(SecondInterval)

type SecondInterval struct {
	Second int64 `json:"second" redis:"second"`
}

func (m *SecondInterval) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m *SecondInterval) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)

}

func GetPayloadId(payload any) string {
	m := cast.ToStringMap(payload)
	if m["id"] == nil {
		return ""
	}
	return m["id"].(string)
}

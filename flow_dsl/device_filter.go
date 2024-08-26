package flow_dsl

import (
	"github.com/spf13/cast"
)

type DeviceFilter struct {
	Type      string   `json:"type"`
	IDs       []string `json:"ids"`
	Tags      []string `json:"tags"`
	ProductId string   `json:"product_id"`
}

func (f *DeviceFilter) Filter(device map[string]any) bool {
	switch f.Type {
	case "device":
		if f.IDs != nil {
			for _, id := range f.IDs {
				if cast.ToString(device["id"]) == id {
					return true
				}
			}
		}
	case "product":
		if f.ProductId != "" {
			return cast.ToString(device["product_id"]) == f.ProductId
		}
	case "tag":
		if f.Tags != nil {
			for _, tag := range f.Tags {
				deviceTag := cast.ToStringSlice(device["tags"])
				if deviceTag != nil {
					for _, t := range deviceTag {
						if t == tag {
							return true
						}
					}

				}
			}
		}
	}
	return false
}

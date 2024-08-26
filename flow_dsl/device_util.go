package flow_dsl

import (
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func IsIdEntity(entity any) bool {
	if entity == nil {
		return false
	}
	switch entity.(type) {
	case map[string]any:
		return entity.(map[string]any)["id"] != nil
	}
	return false
}
func GetIdEntityId(entity any) (id string) {
	if entity == nil {
		return
	}
	switch entity.(type) {
	case map[string]any:
		return cast.ToString(entity.(map[string]any)["id"])
	}
	return
}

func CheckDeviceDataMatchProperty(property map[string]any, deviceWithTags map[string]any) (match bool, err error) {
	key := cast.ToString(property["key"])
	val, exist := deviceWithTags[key]
	if !exist {
		err = errors.New(key + " not exist")
		return
	}
	op := cast.ToString(property["op"])
	if op == "" {
		op = "eq"
	}
	value := property["value"]
	value2 := property["value2"]
	switch op {
	case "gt":
		if cast.ToFloat64(val) > cast.ToFloat64(value) {
			match = true
			return
		}
	case "ge":
		if cast.ToFloat64(val) >= cast.ToFloat64(value) {
			match = true
			return
		}
	case "lt":
		if cast.ToFloat64(val) < cast.ToFloat64(value) {
			match = true
			return
		}
	case "le":
		if cast.ToFloat64(val) <= cast.ToFloat64(value) {
			match = true
			return
		}
	case "eq":
		if cast.ToString(val) == cast.ToString(value) {
			match = true
			return
		}
	case "noeq":
		if cast.ToString(val) != cast.ToString(value) {
			match = true
			return
		}
	case "range1":
		if cast.ToFloat64(value) <= cast.ToFloat64(val) && cast.ToFloat64(val) <= cast.ToFloat64(value2) {
			match = true
			return
		}
	case "range2":
		if cast.ToFloat64(value) <= cast.ToFloat64(val) && cast.ToFloat64(val) < cast.ToFloat64(value2) {
			match = true
			return
		}
	case "range3":
		if cast.ToFloat64(value) < cast.ToFloat64(val) && cast.ToFloat64(val) <= cast.ToFloat64(value2) {
			match = true
			return
		}
	case "range4":
		if cast.ToFloat64(value) < cast.ToFloat64(val) && cast.ToFloat64(val) < cast.ToFloat64(value2) {
			match = true
			return
		}

	}
	return
}

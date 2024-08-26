package flow_dsl

import "testing"

func TestCheckDeviceDataMatchProperty(t *testing.T) {
	property := make(map[string]interface{})
	//[dataType:product key:NewWindSwitch meta_type:properties name:查询当前状态 product:w4Ee6dKj9Vu1hsXLwIt-T productName:户内新风智控 type:DeviceStateQuery value:0]
	property["key"] = "NewWindSwitch"
	property["dataType"] = "product"
	property["meta_type"] = "properties"
	property["name"] = "查询定时"
	property["product"] = "w4Ee6dKj9Vu1hsXLwIt-T"
	property["productName"] = "户内新风景名"
	property["type"] = "DeviceStateQuery"
	property["value"] = 0
	device := make(map[string]interface{})
	//map[NewWindSwitch:0 identifier:HWIU-01-01-01 tags:[网关:DDC-1-B1-01 区域:1号楼低区西侧 户号:101 楼号:1号楼]]
	device["NewWindSwitch"] = 0
	device["identifier"] = "HWIU-01-01-01"
	device["tags"] = "爱奇良名"
	match, err := CheckDeviceDataMatchProperty(property, device)
	if err != nil {
		t.Error(err)
		return
	}
	if !match {
		t.Error("device data not match property")
	}
	t.Log(match)
}

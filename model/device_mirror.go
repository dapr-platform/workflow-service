package model

import (
	"database/sql"
	"github.com/dapr-platform/common"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = common.LocalTime{}
)

/*
DB Table Details
-------------------------------------


Table: o_device_mirror
[ 0] id                                             VARCHAR(32)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 1] json_data                                      TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] updated_time                                   TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "nDefBLoGkeMvBXcDjIXwDliYl",    "json_data": "ZvlNGabQGchEAWyiArXbjYejt",    "updated_time": 26}



*/

var (
	Device_mirror_FIELD_NAME_id = "id"

	Device_mirror_FIELD_NAME_json_data = "json_data"

	Device_mirror_FIELD_NAME_updated_time = "updated_time"
)

// Device_mirror struct is a row record of the o_device_mirror table in the  database
type Device_mirror struct {
	ID          string           `json:"id"`           //唯一标识
	JSONData    string           `json:"json_data"`    //当前数据json
	UpdatedTime common.LocalTime `json:"updated_time"` //更新时间

}

var Device_mirrorTableInfo = &TableInfo{
	Name: "o_device_mirror",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
			Comment:            `唯一标识`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(32)",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       32,
			GoFieldName:        "ID",
			GoFieldType:        "string",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "json_data",
			Comment:            `当前数据json`,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "JSONData",
			GoFieldType:        "string",
			JSONFieldName:      "json_data",
			ProtobufFieldName:  "json_data",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "updated_time",
			Comment:            `更新时间`,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TIMESTAMP",
			DatabaseTypePretty: "TIMESTAMP",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMP",
			ColumnLength:       -1,
			GoFieldName:        "UpdatedTime",
			GoFieldType:        "common.LocalTime",
			JSONFieldName:      "updated_time",
			ProtobufFieldName:  "updated_time",
			ProtobufType:       "uint64",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *Device_mirror) TableName() string {
	return "o_device_mirror"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *Device_mirror) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *Device_mirror) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *Device_mirror) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *Device_mirror) TableInfo() *TableInfo {
	return Device_mirrorTableInfo
}

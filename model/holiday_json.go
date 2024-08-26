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


Table: o_holiday_json
[ 0] id                                             VARCHAR(32)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 1] json_data                                      TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "egqIsIxcudCNCdLdZuDUZbiKP",    "json_data": "PaAAgFefqVUoViLaVveGTcAJf"}



*/

var (
	Holiday_json_FIELD_NAME_id = "id"

	Holiday_json_FIELD_NAME_json_data = "json_data"
)

// Holiday_json struct is a row record of the o_holiday_json table in the  database
type Holiday_json struct {
	ID       string `json:"id"`        //id(年份，例如2023）
	JSONData string `json:"json_data"` //json数据

}

var Holiday_jsonTableInfo = &TableInfo{
	Name: "o_holiday_json",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
			Comment:            `id(年份，例如2023）`,
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
			Comment:            `json数据`,
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
	},
}

// TableName sets the insert table name for this struct type
func (h *Holiday_json) TableName() string {
	return "o_holiday_json"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (h *Holiday_json) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (h *Holiday_json) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (h *Holiday_json) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (h *Holiday_json) TableInfo() *TableInfo {
	return Holiday_jsonTableInfo
}

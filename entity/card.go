package entity

import "encoding/json"

type Cell struct {
	Shape string         `json:"shape"`
	Id    string         `json:"id"`
	Label string         `json:"label"`
	Data  map[string]any `json:"data"`
}

func (c *Cell) FromMap(m map[string]any) (err error) {
	buf, _ := json.Marshal(m)
	err = json.Unmarshal(buf, c)
	return
}

type Edge struct {
	Shape  string   `json:"shape"`
	Id     string   `json:"id"`
	Source EdgePort `json:"source"`
	Target EdgePort `json:"target"`
	ZIndex int      `json:"zIndex"`
}

func (c *Edge) FromMap(m map[string]any) (err error) {
	buf, _ := json.Marshal(m)
	err = json.Unmarshal(buf, c)
	return
}

type EdgePort struct {
	Cell string `json:"cell"`
	Port string `json:"port"`
}

type Node struct {
	Id         string         `json:"id"`
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties"`
	Ports      []Port         `json:"ports"`
}

type Port struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Business string `json:"business"`
}

package entity

type Holiday struct {
	Schema string       `json:"$schema"`
	Id     string       `json:"$id"`
	Year   int          `json:"year"`
	Papers []string     `json:"papers"`
	Days   []HolidayDay `json:"days"`
}
type HolidayDay struct {
	Name     string `json:"name"`
	Date     string `json:"date"` //YYYY-MM-DD
	IsOffDay bool   `json:"isOffDay"`
}

package holiday

import (
	"context"
	"encoding/json"
	"github.com/dapr-platform/common"
	"time"
	"workflow-service/entity"
	"workflow-service/model"
)

var holidayMap = make(map[string]map[string]entity.HolidayDay)

func init() {
	time.Sleep(time.Second * 10)
	//go refreshData()
}

func refreshData() {
	for {
		allData, err := common.DbQuery[model.Holiday_json](context.Background(), common.GetDaprClient(), model.Holiday_jsonTableInfo.Name, "")
		if err != nil {
			time.Sleep(time.Second * 5)
			common.Logger.Error("query holiday data error", err)
			continue
		}
		if len(allData) == 0 {
			time.Sleep(time.Second * 5)
			common.Logger.Error("query holiday data not found")
			continue
		}
		for _, data := range allData {
			var holiday entity.Holiday
			err1 := json.Unmarshal([]byte(data.JSONData), &holiday)
			if err1 != nil {
				common.Logger.Error("query holiday data error", err1)
				return
			}
			dayMap := make(map[string]entity.HolidayDay)
			for _, day := range holiday.Days {
				dayMap[day.Date] = day
			}
			holidayMap[holiday.Id] = dayMap
		}
		return
	}

}

func IsHoliday(now time.Time) bool {

	dayStr := now.Format("2006-01-02")
	yearStr := now.Format("2006")
	if _, ok := holidayMap[yearStr]; ok {
		if day, ok := holidayMap[yearStr][dayStr]; ok {
			return day.IsOffDay
		}
	}
	return false
}
func IsWorkday(now time.Time) bool {

	week := now.Weekday()

	if week == time.Saturday || week == time.Sunday {
		if !IsHoliday(now) { //倒休
			return true
		}
		return false
	} else {
		return !IsHoliday(now)
	}

}

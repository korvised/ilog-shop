package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func Debug(obj any) {
	raw, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(raw))
}

func LocalTime() time.Time {
	loc, _ := time.LoadLocation("asia/Bangkok")
	return time.Now().In(loc)
}

func ConvertStringToTime(t string) time.Time {
	layout := "2006-01-02T15:04:05.999 -0700 MST"
	result, err := time.Parse(layout, t)
	if err != nil {
		log.Printf("Error when convert string to time: %v", err)
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	return result.In(loc)
}

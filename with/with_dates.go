package with

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/validation"
	"strings"
)

// DatesFields is a struct that contains dates for indexing
type DatesFields struct {
	Dates   []string `json:"dates,omitempty" dalgo:"dates,omitempty" firestore:"dates,omitempty"`
	DateMin string   `json:"dateMin,omitempty" dalgo:"dateMin,omitempty" firestore:"dateMin,omitempty"`
	DateMax string   `json:"dateMax,omitempty" dalgo:"dateMax,omitempty" firestore:"dateMax,omitempty"`
}

func (v *DatesFields) UpdatesWhenDatesChanged() []dal.Update {
	updates := []dal.Update{
		{Field: "dates", Value: v.Dates},
		{Field: "dateMin", Value: v.DateMin},
		{Field: "dateMax", Value: v.DateMax},
	}
	if len(v.Dates) == 0 {
		updates[0].Value = dal.DeleteField
	}
	return updates
}

func (v *DatesFields) AddDate(date string) {
	v.Dates = append(v.Dates, date)
	if v.DateMax == "" || date > v.DateMax {
		v.DateMax = date
	}
	if v.DateMin == "" || date < v.DateMin {
		v.DateMin = date
	}
}

// Validate returns error if not valid
func (v *DatesFields) Validate() error {
	v.DateMin = ""
	v.DateMax = ""
	for i, date := range v.Dates {
		if strings.TrimSpace(date) == "" {
			return validation.NewErrRecordIsMissingRequiredField(fmt.Sprintf("dates[%v]", i))
		}
		if _, err := ValidateDateString(date); err != nil {
			return validation.NewErrBadRecordFieldValue("date", err.Error())
		}

		for j, date2 := range v.Dates {
			if j != i && date2 == date {
				return validation.NewErrBadRecordFieldValue("dates", "duplicate value: "+date)
			}
		}
		if date < v.DateMin || v.DateMin == "" {
			v.DateMin = date
		}
		if date > v.DateMax {
			v.DateMax = date
		}
	}
	return nil
}

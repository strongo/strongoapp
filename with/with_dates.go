package with

import (
	"fmt"
	"github.com/dal-go/dalgo/update"
	"github.com/strongo/validation"
	"slices"
	"strings"
)

// DatesFields is a struct that contains dates for indexing
type DatesFields struct {
	Dates   []string `json:"dates,omitempty" dalgo:"dates,omitempty" firestore:"dates,omitempty"`
	DateMin string   `json:"dateMin,omitempty" dalgo:"dateMin,omitempty" firestore:"dateMin,omitempty"`
	DateMax string   `json:"dateMax,omitempty" dalgo:"dateMax,omitempty" firestore:"dateMax,omitempty"`
}

func (v *DatesFields) UpdatesWhenDatesChanged() []update.Update {
	updates := []update.Update{
		update.ByFieldName("dates", v.Dates),
		update.ByFieldName("dateMin", v.DateMin),
		update.ByFieldName("dateMax", v.DateMax),
	}
	if len(v.Dates) == 0 {
		updates[0] = update.ByFieldName("dates", update.DeleteField)
	}
	return updates
}

func (v *DatesFields) AddDate(date string) (updates []update.Update) {
	if slices.Contains(v.Dates, date) {
		return
	}
	v.Dates = append(v.Dates, date)
	updates = []update.Update{
		update.ByFieldName("dates", v.Dates),
	}
	if v.DateMax == "" || date > v.DateMax {
		v.DateMax = date
		updates = append(updates, update.ByFieldName("dateMax", v.DateMax))
	}
	if v.DateMin == "" || date < v.DateMin {
		v.DateMin = date
		updates = append(updates, update.ByFieldName("dateMin", v.DateMin))
	}
	return
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

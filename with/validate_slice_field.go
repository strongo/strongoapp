package with

import (
	"fmt"
	"github.com/strongo/validation"
	"slices"
	"strings"
)

func ValidateSetSliceField(field string, v []string, isRecordID bool) error {
	count := len(v)
	for i, s := range v {
		if s2 := strings.TrimSpace(s); s2 == "" {
			return validation.NewErrRecordIsMissingRequiredField(fmt.Sprintf("%v[%v]", field, i))
		} else if s2 != s {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("%v[%v]", field, i), "starts or ends with spaces")
		}
		if i < count {
			if slices.Contains(v[i+1:], s) {
				return validation.NewErrBadRecordFieldValue(field,
					fmt.Sprintf("duplicate value at indexes %d & %d: %s", i, slices.Index(v[i+1:], s), s))
			}
		}
		if isRecordID {
			if err := ValidateRecordID(s); err != nil {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("%s[%v]", field, i), err.Error())
			}
		}
	}
	return nil
}

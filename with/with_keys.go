package with

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/validation"
	"strings"
)

type KeysField struct {
	Keys []string `json:"keys,omitempty" dalgo:"keys,omitempty"  firestore:"keys,omitempty"`
}

func (v KeysField) Validate() error {
	for i, k := range v.Keys {
		if s := strings.TrimSpace(k); s == "" {
			return validation.NewErrRecordIsMissingRequiredField(fmt.Sprintf("keys[%v]", i))
		} else if s != k {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("keys[%v]", i), fmt.Sprintf("should be trimmed, got: [%v]", k))
		}
		for j, k2 := range v.Keys {
			if j != i && k == k2 {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("keys[%v]", i), fmt.Sprintf("duplicate value: [%v]", k))
			}
		}
	}
	return nil
}

func (v KeysField) UpdatesWhenKeysChanged() []dal.Update {
	if len(v.Keys) == 0 {
		return []dal.Update{
			{Field: "keys", Value: dal.DeleteField},
		}
	}
	return []dal.Update{
		{Field: "keys", Value: v.Keys},
	}
}

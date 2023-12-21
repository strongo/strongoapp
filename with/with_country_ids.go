package with

import (
	"fmt"
)

// CountryIDsField defines a record with a Country IDs
type CountryIDsField struct {
	CountryIDs []string `json:"countryIDs,omitempty" dalgo:"countryIDs,omitempty" firestore:"countryIDs,omitempty"`
}

func (v CountryIDsField) Validate() error {
	for i, countryID := range v.CountryIDs {
		if err := ValidateRequiredCountryID(fmt.Sprintf("countryIDs[%v]", i), countryID); err != nil {
			return err
		}
	}
	return nil
}

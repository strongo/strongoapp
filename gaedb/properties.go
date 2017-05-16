package gaedb

import (
	"google.golang.org/appengine/datastore"
	"time"
	"fmt"
	"github.com/pkg/errors"
)

type IsOkToRemove func(p datastore.Property) bool // TODO: Open source + article?

func IsObsolete(_ datastore.Property) bool {
	return true
}

func IsFalse(p datastore.Property) bool {
	return !p.Value.(bool)
}

func IsZeroInt(p datastore.Property) bool {
	return p.Value.(int64) == 0
}

func IsZeroFloat(p datastore.Property) bool {
	return p.Value.(float64) == 0
}

func IsZeroTime(p datastore.Property) bool {
	return p.Value.(time.Time).IsZero()
}

func IsEmptyString(p datastore.Property) bool {
	return p.Value.(string) == ""
}

func IsEmptyByteArray(p datastore.Property) bool {
	v := p.Value.([]uint8)
	return  v == nil || len(v) == 0
}

func IsEmptyStringOrSpecificValue(v string) func(p datastore.Property) bool {
	return func(p datastore.Property) bool {
		s := p.Value.(string)
		return s == "" || s == v
	}
}

// Removes properties in place and returns filtered slice
func CleanProperties(properties []datastore.Property, filters map[string]IsOkToRemove) (filtered []datastore.Property, err error) {
	var (
		i int
		p datastore.Property
	)

	if properties == nil {
		return properties, errors.New("properties == nil")
	}

	if filters == nil {
		return properties, errors.New("filters == nil")
	}

	if len(filters) == 0 {
		return properties, errors.New("len(filters) == 0")
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Failed to process property %v=%T(%v), recovered: %v", p.Name, p.Value, p.Value, r))
		}
		filtered = properties[:i]
	}()

	for _, p = range properties {
		if filter, ok := filters[p.Name]; !ok || !filter(p) {
			properties[i] = p
			i += 1
		}
	}
	return
}

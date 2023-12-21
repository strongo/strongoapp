package person

import (
	"errors"
	"github.com/alexsergivan/transliterator"
	"github.com/strongo/random"
	"slices"
	"strings"
)

// CleanForID removes non-ASCII characters and converts to a lower case
func CleanForID(s string) string {
	var result strings.Builder
	result.Grow(len(s))
	for i := 0; i < len(s); i++ {
		b := s[i]
		if 'a' <= b && b <= 'z' ||
			'0' <= b && b <= '9' {
			result.WriteByte(b)
		} else if 'A' <= b && b <= 'Z' {
			result.WriteByte(b + 'a' - 'A') // to lower case
		}
	}
	return result.String()
}

// GenerateIDFromNameOrRandom generates ContactID from name or random
func GenerateIDFromNameOrRandom(name *NameFields, existingIDs []string) (id string, err error) {
	trans := transliterator.NewTransliterator(nil)
	//
	if nick := CleanForID(trans.Transliterate(name.NickName, "")); nick != "" {
		if !slices.Contains(existingIDs, id) {
			return nick, nil
		}
	}

	first := CleanForID(trans.Transliterate(name.FirstName, ""))
	last := CleanForID(trans.Transliterate(name.LastName, ""))
	middle := CleanForID(trans.Transliterate(name.MiddleName, ""))

	if first == "" || last == "" || middle == "" {
		if first == "" && last == "" && middle == "" {
			if name.FullName != "" {
				if names := strings.Split(name.FullName, " "); len(names) > 0 {
					for _, n := range names {
						n = CleanForID(trans.Transliterate(n, ""))
						if len(n) > 0 {
							id += n[0:1]
						}
					}
					if len(id) > 0 && !slices.Contains(existingIDs, id) {
						return id, nil
					}
					if len(names) == 2 {
						first = names[0]
						last = names[1]
					}
				}
			}
		}

		if first != "" && last != "" {
			// Try to use 1st chars of first & last names
			if id = first[0:1] + last[0:1]; !slices.Contains(existingIDs, id) {
				return id, nil
			}
		}
		if first != "" && middle != "" && last != "" {
			// Try to user 1st chars of all first, middle, last names
			if id = first[0:1] + middle[0:1] + last[0:1]; !slices.Contains(existingIDs, id) {
				return id, nil
			}

			// Try to user 1st chars of all first, last names, middle names
			if id = first[0:1] + last[0:1] + middle[0:1]; !slices.Contains(existingIDs, id) {
				return id, nil
			}
		}
		// Try to use 1st char of first name
		if first != "" {

			// Try to use 1st char of first name
			if id = first[0:1]; !slices.Contains(existingIDs, id) {
				return id, nil
			}

			// Try to use 1st and last char of first name
			if id = first[0:1] + first[len(first)-1:]; slices.Contains(existingIDs, id) {
				return id, nil
			}

			// Try to use the whole first name
			if id = first; !slices.Contains(existingIDs, id) {
				return id, nil
			}
		}
		if first != "" && last != "" {
			// Try to use full first name and 1st char of last name
			if id = first + last[0:1]; slices.Contains(existingIDs, id) {
				return id, nil
			}
		}
		if last != "" {
			if id = last; slices.Contains(existingIDs, id) {
				return id, nil
			}
		}
	}

	return NewUniqueRandomID(existingIDs, 3)
}

// NewUniqueRandomID generate new team ContactID
func NewUniqueRandomID(existingIDs []string, idLength int) (id string, err error) {
	randomIDAttempts := 0
OUTER:
	for {
		id = random.ID(idLength)
		for _, v := range existingIDs {
			if v == id {
				randomIDAttempts++
				if randomIDAttempts > 100 {
					return "", errors.New("too many attempts to generate random member ContactID")
				}
				continue OUTER
			}
		}
		break
	}
	return
}

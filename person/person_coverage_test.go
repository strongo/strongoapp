package person

import (
	"strings"
	"testing"
)

func TestNameFields_SetNames_And_GetName(t *testing.T) {
	nf := &NameFields{}

	names := []Name{
		{Field: Username, Value: "user1"},
		{Field: FullName, Value: "Full Name"},
		{Field: FirstName, Value: "First"},
		{Field: MiddleName, Value: "Middle"},
		{Field: LastName, Value: "Last"},
		{Field: NickName, Value: "Nick"},
		{Field: ScreenName, Value: "Screen"},
	}

	if err := nf.SetNames(names...); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if nf.GetName(Username) != "user1" {
		t.Fatalf("expected user1, got %s", nf.GetName(Username))
	}
	if nf.GetName(FullName) != "Full Name" {
		t.Fatalf("expected Full Name, got %s", nf.GetName(FullName))
	}
	if nf.GetName(FirstName) != "First" {
		t.Fatalf("expected First, got %s", nf.GetName(FirstName))
	}
	if nf.GetName(MiddleName) != "Middle" {
		t.Fatalf("expected Middle, got %s", nf.GetName(MiddleName))
	}
	if nf.GetName(LastName) != "Last" {
		t.Fatalf("expected Last, got %s", nf.GetName(LastName))
	}
	if nf.GetName(NickName) != "Nick" {
		t.Fatalf("expected Nick, got %s", nf.GetName(NickName))
	}
	if nf.GetName(ScreenName) != "Screen" {
		t.Fatalf("expected Screen, got %s", nf.GetName(ScreenName))
	}
	if nf.GetName(NameField(999)) != "" {
		t.Fatalf("expected empty string for unsupported NameField")
	}

	// Unsupported field error in SetNames
	if err := nf.SetNames(Name{Field: NameField(999), Value: "bad"}); err == nil {
		t.Fatal("expected error on unsupported NameField")
	}

	// String()
	str := nf.String()
	if !strings.Contains(str, `UserName="user1"`) {
		t.Fatalf("unexpected string output: %s", str)
	}
}

func TestNameFields_GetFullName(t *testing.T) {
	tests := []struct {
		name string
		nf   NameFields
		want string
	}{
		{
			name: "explicit full name",
			nf:   NameFields{FullName: "Full Name Custom", FirstName: "First", LastName: "Last"},
			want: "Full Name Custom",
		},
		{
			name: "first nick last",
			nf:   NameFields{FirstName: "First", NickName: "Nick", LastName: "Last"},
			want: "First (Nick) Last",
		},
		{
			name: "first last",
			nf:   NameFields{FirstName: "First", LastName: "Last"},
			want: "First Last",
		},
		{
			name: "first nick",
			nf:   NameFields{FirstName: "First", NickName: "Nick"},
			want: "First (Nick)",
		},
		{
			name: "last nick",
			nf:   NameFields{LastName: "Last", NickName: "Nick"},
			want: "Last (Nick)",
		},
		{
			name: "first only",
			nf:   NameFields{FirstName: "First"},
			want: "First",
		},
		{
			name: "last only",
			nf:   NameFields{LastName: "Last"},
			want: "Last",
		},
		{
			name: "nick only",
			nf:   NameFields{NickName: "Nick"},
			want: "Nick",
		},
		{
			name: "username only",
			nf:   NameFields{UserName: "user123"},
			want: "user123",
		},
		{
			name: "screen name only",
			nf:   NameFields{ScreenName: "screen123"},
			want: "screen123",
		},
		{
			name: "empty",
			nf:   NameFields{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.nf.GetFullName(); got != tt.want {
				t.Fatalf("GetFullName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNameFields_Equal(t *testing.T) {
	var n1 *NameFields
	var n2 *NameFields
	if !n1.Equal(n2) {
		t.Fatal("expected nil == nil")
	}

	n1Val := NameFields{FirstName: "A", MiddleName: "B", LastName: "C", UserName: "D", NickName: "E", FullName: "F"}
	n1 = &n1Val
	if n1.Equal(nil) || n2.Equal(n1) {
		t.Fatal("expected nil != non-nil")
	}

	n2Val := n1Val
	n2 = &n2Val
	if !n1.Equal(n2) {
		t.Fatal("expected equal structs")
	}

	// Difference
	n2Val.NickName = "Diff"
	if n1.Equal(n2) {
		t.Fatal("expected false on difference")
	}
}

func TestNameFields_Validate_And_IsEmpty(t *testing.T) {
	var nilNF *NameFields
	if !nilNF.IsEmpty() {
		t.Fatal("expected nil to be empty")
	}
	if err := nilNF.Validate(); err != nil {
		t.Fatalf("expected nil to be valid, got %v", err)
	}

	emptyNF := &NameFields{}
	if !emptyNF.IsEmpty() {
		t.Fatal("expected empty struct to be empty")
	}
	if err := emptyNF.Validate(); err == nil {
		t.Fatal("expected error on empty struct ValidateAtLeast1Name")
	}

	// Non-empty
	validNF := &NameFields{FirstName: "John"}
	if validNF.IsEmpty() {
		t.Fatal("expected non-empty")
	}
	if err := validNF.Validate(); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}

	// Space validations
	badFirst := &NameFields{FirstName: " John"}
	if err := badFirst.Validate(); err == nil {
		t.Fatal("expected error for space in FirstName")
	}

	badLast := &NameFields{LastName: "Doe "}
	if err := badLast.Validate(); err == nil {
		t.Fatal("expected error for space in LastName")
	}

	badMiddle := &NameFields{MiddleName: " M "}
	if err := badMiddle.Validate(); err == nil {
		t.Fatal("expected error for space in MiddleName")
	}

	badFull := &NameFields{FullName: " John Doe "}
	if err := badFull.Validate(); err == nil {
		t.Fatal("expected error for space in FullName")
	}
}

func TestGenerateIDFromNameOrRandom_Coverage(t *testing.T) {
	// 1. name is nil -> defaults to empty &NameFields{}, generates random ID
	id1, err := GenerateIDFromNameOrRandom(nil, nil)
	if err != nil || len(id1) != 3 {
		t.Fatalf("expected 3-char random id, got %s, err: %v", id1, err)
	}

	// 2. NickName provided
	// Case 2a: existingIDs does not contain "" -> returns nick
	id2, err := GenerateIDFromNameOrRandom(&NameFields{NickName: "TheRock"}, []string{"other"})
	if err != nil || id2 != "therock" {
		t.Fatalf("expected therock, got %s, err: %v", id2, err)
	}

	// Case 2b: existingIDs contains "therock" -> nick branch does not return, continues
	// With FullName "John Doe", 2 words:
	id2b, err := GenerateIDFromNameOrRandom(&NameFields{NickName: "TheRock", FullName: "John Doe"}, []string{"therock"})
	if err != nil || id2b != "jd" {
		t.Fatalf("expected jd, got %s, err: %v", id2b, err)
	}

	// 3. FullName with >2 words: "John Middle Doe"
	// id += n[0:1] -> "jmd"
	id3, err := GenerateIDFromNameOrRandom(&NameFields{FullName: "John Middle Doe"}, nil)
	if err != nil || id3 != "jmd" {
		t.Fatalf("expected jmd, got %s, err: %v", id3, err)
	}

	// 4. FullName where id exists in existingIDs, but len(names) == 2 -> sets first and last
	// existingIDs contains "jd" -> first="john", last="doe", first[0:1]+last[0:1] is "jd" which exists,
	// then falls to first != "" -> first[0:1] = "j", not in existingIDs -> returns "j"
	id4, err := GenerateIDFromNameOrRandom(&NameFields{FullName: "John Doe"}, []string{"jd"})
	if err != nil || id4 != "j" {
		t.Fatalf("expected j, got %s, err: %v", id4, err)
	}

	// 5. first != "" && middle != "" && last != ""
	// 5a. first[0:1] + last[0:1] is "jl". If "jl" is in existingIDs, tries first[0:1] + middle[0:1] + last[0:1] = "jml"
	id5a, err := GenerateIDFromNameOrRandom(&NameFields{FirstName: "John", MiddleName: "Michael", LastName: "Locke"}, []string{"jl"})
	if err != nil || id5a != "jml" {
		t.Fatalf("expected jml, got %s, err: %v", id5a, err)
	}
	// 5b. "jl" and "jml" in existingIDs -> tries first[0:1] + last[0:1] + middle[0:1] = "jlm"
	id5b, err := GenerateIDFromNameOrRandom(&NameFields{FirstName: "John", MiddleName: "Michael", LastName: "Locke"}, []string{"jl", "jml"})
	if err != nil || id5b != "jlm" {
		t.Fatalf("expected jlm, got %s, err: %v", id5b, err)
	}

	// 6. first != "":
	// 6a. first[0:1] = "j" in existingIDs -> returns first[0:1] + first[len-1:] = "jn"
	id6a, err := GenerateIDFromNameOrRandom(&NameFields{FirstName: "John"}, []string{"j"})
	if err != nil || id6a != "jn" {
		t.Fatalf("expected jn, got %s, err: %v", id6a, err)
	}

	// 6b. "j" and "jn" in existingIDs -> returns "john"
	id6b, err := GenerateIDFromNameOrRandom(&NameFields{FirstName: "John"}, []string{"j", "jn"})
	if err != nil || id6b != "john" {
		t.Fatalf("expected john, got %s, err: %v", id6b, err)
	}

	// 7. first != "" && last != "" where "jd", "j", "jn", "john" in existingIDs
	// returns "johnd"
	id7, err := GenerateIDFromNameOrRandom(&NameFields{FirstName: "John", LastName: "Doe"}, []string{"jd", "j", "jn", "john"})
	if err != nil || id7 != "johnd" {
		t.Fatalf("expected johnd, got %s, err: %v", id7, err)
	}

	// 8. last != ""
	id8, err := GenerateIDFromNameOrRandom(&NameFields{LastName: "Doe"}, nil)
	if err != nil || id8 != "doe" {
		t.Fatalf("expected doe, got %s, err: %v", id8, err)
	}

	// 9. All options exhausted for last != "" -> falls through to NewUniqueRandomID
	id9, err := GenerateIDFromNameOrRandom(&NameFields{LastName: "Doe"}, []string{"doe"})
	if err != nil || len(id9) != 3 {
		t.Fatalf("expected 3-char random id, got %s, err: %v", id9, err)
	}
}

func TestNewUniqueRandomID_Error(t *testing.T) {
	// When existingIDs contains "" and idLength is 0, random.ID(0) produces ""
	// which collides every time, hitting the >100 retry limit and returning error.
	_, err := NewUniqueRandomID([]string{""}, 0)
	if err == nil || !strings.Contains(err.Error(), "too many attempts") {
		t.Fatalf("expected too many attempts error, got %v", err)
	}
}

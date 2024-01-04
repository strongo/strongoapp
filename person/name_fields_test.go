package person

import "testing"

func TestDeductNamesFromFullNames(t *testing.T) {
	type args struct {
		fullName string
	}
	tests := []struct {
		name          string
		args          args
		wantFirstName string
		wantLastName  string
	}{
		{
			name:          "empty",
			args:          args{fullName: ""},
			wantFirstName: "",
			wantLastName:  "",
		},
		{
			name:          "John Smith",
			args:          args{fullName: "John Smith"},
			wantFirstName: "John",
			wantLastName:  "Smith",
		},
		{
			name:          "John  Smith",
			args:          args{fullName: "John  Smith"},
			wantFirstName: "John",
			wantLastName:  "Smith",
		},
		{
			name:          " John Smith",
			args:          args{fullName: " John Smith"},
			wantFirstName: "John",
			wantLastName:  "Smith",
		},
		{
			name:          "John Smith ",
			args:          args{fullName: "John Smith "},
			wantFirstName: "John",
			wantLastName:  "Smith",
		},
		{
			name:          "John Jr Smith",
			args:          args{fullName: "John Jr Smith"},
			wantFirstName: "",
			wantLastName:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFirstName, gotLastName := DeductNamesFromFullNames(tt.args.fullName)
			if gotFirstName != tt.wantFirstName {
				t.Errorf("DeductNamesFromFullNames() gotFirstName = %v, want %v", gotFirstName, tt.wantFirstName)
			}
			if gotLastName != tt.wantLastName {
				t.Errorf("DeductNamesFromFullNames() gotLastName = %v, want %v", gotLastName, tt.wantLastName)
			}
		})
	}
}

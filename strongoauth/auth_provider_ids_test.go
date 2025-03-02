package strongoauth

import (
	"testing"
)

func TestValidateAuthProviderID(t *testing.T) {
	SetKnownAuthProviderIDs([]string{"email", "telegram", "firebase"})
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{

		{
			name: "valid input", // Test case with expected valid input
			args: args{
				v: "firebase",
			},
			wantErr: "",
		},
		{
			name: "empty input", // Test case with empty input
			args: args{
				v: "",
			},
			wantErr: "empty string",
		},
		{
			name: "unknown", // Test case with invalid characters
			args: args{
				v: "someProvider",
			},
			wantErr: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAuthProviderID(tt.args.v)
			if tt.wantErr == "" && err != nil || tt.wantErr != "" && err == nil {
				t.Errorf("ValidateAuthProviderID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

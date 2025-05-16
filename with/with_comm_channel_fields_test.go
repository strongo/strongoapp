package with

import (
	"github.com/stretchr/testify/assert"
	"github.com/strongo/validation"
	"strings"
	"testing"
)

func TestCommChannelFields_Validate(t *testing.T) {
	type fields struct {
		EmailsField EmailsField
		PhonesField PhonesField
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{name: "empty", fields: fields{}, wantErr: assert.NoError},
		{
			name: "invalid_emails",
			fields: fields{
				EmailsField: EmailsField{
					Emails: map[string]*CommunicationChannelProps{
						"": {},
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				if !validation.IsBadFieldValueError(err) {
					return assert.Fail(t, "expected bad field value error", i...)
				}
				if !strings.Contains(err.Error(), EmailsFieldName) {
					return assert.Fail(t, "expected error to contain 'phones'", i...)
				}
				return true
			},
		},
		{
			name: "invalid_phones",
			fields: fields{
				PhonesField: PhonesField{
					Phones: map[string]*CommunicationChannelProps{
						"": {},
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				if !validation.IsBadFieldValueError(err) {
					return assert.Fail(t, "expected bad field value error", i...)
				}
				if !strings.Contains(err.Error(), PhonesFieldName) {
					return assert.Fail(t, "expected error to contain 'phones'", i...)
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &CommChannelFields{
				EmailsField: tt.fields.EmailsField,
				PhonesField: tt.fields.PhonesField,
			}
			tt.wantErr(t, v.Validate(), "Validate()")
		})
	}
}

func TestCommChannelFields_GetCommChannels(t *testing.T) {
	type fields struct {
		EmailsField EmailsField
		PhonesField PhonesField
	}
	type args struct {
		t CommChannelType
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]*CommunicationChannelProps
	}{
		{
			name: "get_emails",
			args: args{
				t: CommChannelTypeEmail,
			},
			fields: fields{
				EmailsField: EmailsField{
					Emails: map[string]*CommunicationChannelProps{
						"email@example.com": {},
					},
				},
				PhonesField: PhonesField{
					Phones: map[string]*CommunicationChannelProps{
						"+1234567890": {},
					},
				},
			},
		},
		{
			name: "get_phone",
			args: args{
				t: CommChannelTypePhone,
			},
			fields: fields{
				EmailsField: EmailsField{
					Emails: map[string]*CommunicationChannelProps{
						"email@example.com": {},
					},
				},
				PhonesField: PhonesField{
					Phones: map[string]*CommunicationChannelProps{
						"+1234567890": {},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &CommChannelFields{
				EmailsField: tt.fields.EmailsField,
				PhonesField: tt.fields.PhonesField,
			}
			var expected map[string]*CommunicationChannelProps
			switch tt.args.t {
			case CommChannelTypeEmail:
				expected = tt.fields.EmailsField.Emails
			case CommChannelTypePhone:
				expected = tt.fields.PhonesField.Phones
			}
			assert.Equalf(t, expected, v.GetCommChannels(tt.args.t), "GetCommChannels(%v)", tt.args.t)
		})
	}
}

package with

import (
	"fmt"
	"github.com/strongo/validation"
	"net/mail"
)

const EmailsFieldName = "emails"

type EmailsField struct {
	Emails map[string]*CommunicationChannelProps `json:"emails,omitempty" firestore:"emails,omitempty"`
}

func (v *EmailsField) Validate() error {
	if err := validateCommunicationChannelsField(EmailsFieldName, v.Emails, func(k string, _ *CommunicationChannelProps) error {
		if _, err := mail.ParseAddress(k); err != nil {
			return validation.NewErrBadRecordFieldValue(EmailsFieldName+fmt.Sprintf("[%s]", k), "invalid email address: "+err.Error())
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

package with

type CommChannelType = string

const (
	CommChannelTypeEmail CommChannelType = "email"
	CommChannelTypePhone CommChannelType = "phone"
)

type CommChannelFields struct {
	EmailsField
	PhonesField
}

func (v *CommChannelFields) Validate() error {
	if err := v.EmailsField.Validate(); err != nil {
		return err
	}
	if err := v.PhonesField.Validate(); err != nil {
		return err
	}
	return nil
}

func (v *CommChannelFields) GetCommChannels(t CommChannelType) (channels map[string]*CommunicationChannelProps, channelsFieldName string) {
	switch t {
	case CommChannelTypeEmail:
		if v.Emails == nil {
			v.Emails = make(map[string]*CommunicationChannelProps)
		}
		return v.Emails, EmailsFieldName
	case CommChannelTypePhone:
		if v.Phones == nil {
			v.Phones = make(map[string]*CommunicationChannelProps)
		}
		return v.Phones, PhonesFieldName
	default:
		return nil, ""
	}
}

package with

const PhonesFieldName = "phones"

type PhonesField struct {
	Phones map[string]*CommunicationChannelProps `json:"phones,omitempty" firestore:"phones,omitempty"`
}

func (v PhonesField) Validate() error {
	return validateCommunicationChannelsField(PhonesFieldName, v.Phones, nil)
}

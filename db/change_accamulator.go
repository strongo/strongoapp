package db

type Changes struct {
	entityHolders []EntityHolder

}

func (changes Changes) IsChanged(entityHolder EntityHolder) bool {
	for i := range changes.entityHolders {
		if changes.entityHolders[i] == entityHolder {
			return true
		}
	}
	return false
}

func (changes *Changes) FlagAsChanged(entityHolder EntityHolder) {
	if entityHolder == nil {
		panic("entityHolder == nil")
	}
	for _, eh := range changes.entityHolders {
		if eh == entityHolder {
			return
		}
	}
	changes.entityHolders = append(changes.entityHolders, entityHolder)
}

func (changes Changes) EntityHolders() (entityHolders []EntityHolder) {
	entityHolders = make([]EntityHolder, len(changes.entityHolders))
	copy(entityHolders, changes.entityHolders)
	return
}

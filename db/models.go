package db

type NoStrID struct {}

func (_ NoStrID) StrID() string {
	return ""
}

func (_ *NoStrID) SetStrID(id string) {
	panic("String ID is not supported, use SetIntID()")
}

type NoIntID struct {}

func (_ NoIntID) IntID() int64 {
	return 0
}

func (_ *NoIntID) SetIntID(id int64) {
	panic("Integer ID is not supported, use SetStrID()")
}

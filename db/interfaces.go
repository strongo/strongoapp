package db

import "golang.org/x/net/context"

type EntityHolder interface {
	Kind() string
	IntID() int64
	StrID() string
	Entity() interface{}
	SetEntity(interface{})
	SetIntID(id int64)
	SetStrID(id string)
}

type MultiUpdater interface {
	UpdateMulti(c context.Context, entityHolders []EntityHolder) error
}

type MultiGetter interface {
	GetMulti(c context.Context, entityHolders []EntityHolder) error
}

type RunOptions map[string]interface{}

type TransactionCoordinator interface {
	RunInTransaction(c context.Context, f func(c context.Context) error, options RunOptions) error
}

type Database interface {
	TransactionCoordinator
	MultiGetter
	MultiUpdater
}

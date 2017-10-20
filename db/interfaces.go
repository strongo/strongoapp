package db

import (
	"golang.org/x/net/context"
	"time"
	"math/rand"
)

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

type Getter interface {
	Get(c context.Context, entityHolder EntityHolder) error
}

type Updater interface {
	Update(c context.Context, entityHolder EntityHolder) error
}

type RunOptions map[string]interface{}

type TransactionCoordinator interface {
	RunInTransaction(c context.Context, f func(c context.Context) error, options RunOptions) (err error)
	IsInTransaction(c context.Context) bool
	NonTransactionalContext(tc context.Context) context.Context
}

type Database interface {
	TransactionCoordinator
	Getter
	Updater
	MultiGetter
	MultiUpdater
}

var (
	CrossGroupTransaction     = RunOptions{"XG": true}
	SingleGroupTransaction    = RunOptions{}
)

const idChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // Removed 1, I and 0, O as can be messed with l/1 and 0.

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomStringID(n uint8) string {
	b := make([]byte, n)
	lettersCount := len(idChars)
	for i := range b {
		b[i] = idChars[random.Intn(lettersCount)]
	}
	return string(b)
}


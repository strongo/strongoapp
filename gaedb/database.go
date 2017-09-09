package gaedb

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"github.com/strongo/app/db"
)

type database struct {
}

func NewDatabase() db.Database {
	return database{}
}

func (_ database) UpdateMulti(c context.Context, entityHolders []db.EntityHolder) (err error) { // TODO: Rename to PutMulti?
	keys := make([]*datastore.Key, len(entityHolders))
	vals := make([]interface{}, len(entityHolders))
	insertedIndexes := make([]int, 0, len(entityHolders))
	for i, entityHolder := range entityHolders {
		intID := entityHolder.IntID()
		strID := entityHolder.StrID()
		if intID == 0 && strID == "" {
			keys[i] = NewIncompleteKey(c, entityHolder.Kind(), nil)
			insertedIndexes = append(insertedIndexes, i)
		} else if intID != 0 && strID != "" {
			return errors.New(fmt.Sprintf("Entity #%d has both IDs: %v(intID=%d, strID=%v)", i, entityHolder.Kind(), intID, strID))
		} else {
			keys[i] = NewKey(c, entityHolder.Kind(), strID, intID, nil)
		}
		vals[i] = entityHolder.Entity()
	}
	if keys, err = PutMulti(c, keys, vals); err != nil {
		return err
	}
	for _, i := range insertedIndexes {
		entityHolders[i].SetIntID(keys[i].IntID())
	}
	return nil
}

func (_ database) GetMulti(c context.Context, entityHolders []db.EntityHolder) error {
	keys := make([]*datastore.Key, len(entityHolders))
	vals := make([]interface{}, len(entityHolders))
	for i, entityHolder := range entityHolders {
		intID := entityHolder.IntID()
		strID := entityHolder.StrID()
		if intID != 0 && strID != "" {
			return errors.New("intID != 0 && strID is NOT empty string")
		} else if intID == 0 && strID == "" {
			return errors.New("intID == 0 && strID is empty string")
		}
		keys[i] = NewKey(c, entityHolder.Kind(), strID, intID, nil)
		vals[i] = entityHolder.Entity()
	}
	return GetMulti(c, keys, vals)
}

var xgTransaction = &datastore.TransactionOptions{XG: true}

func (_ database) RunInTransaction(c context.Context, f func(c context.Context) error, options db.RunOptions) error {
	var to *datastore.TransactionOptions
	if xg, ok := options["XG"]; ok && xg.(bool) == true {
		to = xgTransaction
	}
	return RunInTransaction(c, f, to)
}

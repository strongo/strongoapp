package gaedb

import (
	"bytes"
	"fmt"
	"github.com/qedus/nds"
	"github.com/strongo/app/log"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"os"
	"github.com/strongo/app/db"
	"strconv"
	"strings"
	"github.com/pkg/errors"
)

var (
	LoggingEnabled = true
	mockDB           MockDB
	NewIncompleteKey = datastore.NewIncompleteKey
	NewKey           = datastore.NewKey

	RunInTransaction = func (c context.Context, f func(tc context.Context) error, opts *datastore.TransactionOptions) error {
		if LoggingEnabled {
			log.Debugf(c, "Starting transaction...")
		}
		if err := nds.RunInTransaction(c, f, opts); err != nil {
			if LoggingEnabled {
				if strings.Contains(err.Error(),"nested transactions are not supported") {
					panic(err)
				}
				log.Debugf(c, errors.WithMessage(err, "transaction failed").Error())
			}
			return err
		}
		if LoggingEnabled {
			log.Debugf(c, "Transaction successful")
		}
		return nil
	}

	Put = func(c context.Context, key *datastore.Key, val interface{}) (*datastore.Key, error) {
		if LoggingEnabled {
			log.Debugf(c, "nds.Put(%v, %T=%+v)", key2str(key), val, val)
			if entity, ok := val.(datastore.PropertyLoadSaver); ok {
				if props, err := entity.Save(); err != nil {
					return nil, errors.WithMessage(err, "failed to Save() to properties")
				} else {
					log.Debugf(c, "properties: %v", props)
				}

			}
		}
		return nds.Put(c, key, val)
	}

	PutMulti = func(c context.Context, keys []*datastore.Key, vals interface{}) ([]*datastore.Key, error) {
		if LoggingEnabled {
			logKeys(c, "nds.PutMulti", keys)
		}
		return nds.PutMulti(c, keys, vals)
	}

	Get = func(c context.Context, key *datastore.Key, val interface{}) error {
		if LoggingEnabled {
			log.Debugf(c, "nds.Get(%v)", key2str(key))
		}
		if key.IntID() == 0 && key.StringID() == "" {
			panic("key.IntID() == 0 && key.StringID() is empty string")
		}
		return nds.Get(c, key, val)
	}

	GetMulti = func(c context.Context, keys []*datastore.Key, vals interface{}) error {
		if LoggingEnabled {
			logKeys(c, "nds.GetMulti", keys)
		}
		return nds.GetMulti(c, keys, vals)
	}

	Delete = func(c context.Context, key *datastore.Key) error {
		if LoggingEnabled {
			log.Debugf(c, "gaedb.Delete(%v)", key2str(key))
		}
		return nds.Delete(c, key)
	}

	DeleteMulti = func(c context.Context, keys []*datastore.Key) error {
		if LoggingEnabled {
			logKeys(c, "nds.DeleteMulti", keys)
		}
		return nds.DeleteMulti(c, keys)
	}
)

func key2str(key *datastore.Key) string {
	kind := key.Kind()
	if intID := key.IntID(); intID != 0 {
		return kind + ":int=" + strconv.FormatInt(intID, 10)
	} else if strID := key.StringID(); strID != "" {
		return kind + ":str=" + strID
	} else {
		return kind + ":new"
	}
}


func logKeys(c context.Context, f string, keys []*datastore.Key) {
	var buffer bytes.Buffer
	buffer.WriteString(f + "(\n")
	for _, key := range keys {
		buffer.WriteString(fmt.Sprintf("\t%v\n", key2str(key)))
	}
	buffer.WriteString(")")
	log.Debugf(c, buffer.String())
}

func SetupNdsMock() {
	if err := os.Setenv("GAE_LONG_APP_ID", "debtstracker"); err != nil {
		panic(err)
	}
	if err := os.Setenv("GAE_PARTITION", "DEVTEST"); err != nil {
		panic(err)
	}
	mockDB = MockDB{EntitiesByKind: make(EntitiesStorage)}

	Get = func(c context.Context, key *datastore.Key, val interface{}) error {
		if c == nil {
			panic("c == nil")
		}
		if key == nil {
			panic("key == nil")
		}
		log.Debugf(c, "gaedb.Get(key=%v:%v)", key.Kind(), key.IntID())
		kind := key.Kind()

		if entitiesByKey, ok := mockDB.EntitiesByKind[kind]; !ok {
			return datastore.ErrNoSuchEntity
		} else {
			mockKey := NewMockKeyFromDatastoreKey(key)
			if p, ok := entitiesByKey[mockKey]; !ok {
				return datastore.ErrNoSuchEntity
			} else {
				if e, ok := val.(datastore.PropertyLoadSaver); ok {
					return e.Load(p)
				} else {
					return datastore.LoadStruct(e, p)
				}
			}
		}
	}

	Put = func(c context.Context, key *datastore.Key, val interface{}) (*datastore.Key, error) {
		if c == nil {
			panic("c == nil")
		}
		kind := key.Kind()
		entitiesByKey, ok := mockDB.EntitiesByKind[kind]
		if !ok {
			entitiesByKey = make(map[MockKey][]datastore.Property)
			mockDB.EntitiesByKind[kind] = entitiesByKey
		}
		mockKey := NewMockKeyFromDatastoreKey(key)
		if key.StringID() == "" {
			for k, _ := range entitiesByKey {
				if k.Kind == key.Kind() && k.IntID > mockKey.IntID {
					mockKey.IntID = k.IntID + 1
				}
			}
		}

		var p []datastore.Property
		var err error
		if e, ok := val.(datastore.PropertyLoadSaver); ok {
			if p, err = e.Save(); err != nil {
				return key, err
			}
		} else {
			if p, err = datastore.SaveStruct(val); err != nil {
				return key, err
			}
		}
		entitiesByKey[mockKey] = p
		return NewKey(c, mockKey.Kind, mockKey.StringID, mockKey.IntID, nil), nil
	}

	PutMulti = func(c context.Context, keys []*datastore.Key, vals interface{}) ([]*datastore.Key, error) {
		entityHolders := vals.([]db.EntityHolder)
		var err error
		var errs []error
		for i, key := range keys {
			if key, err = Put(c, key, entityHolders[i]); err != nil {
				errs = append(errs, err)
			}
			keys[i] = key
		}
		if len(errs) > 0 {
			return keys, appengine.MultiError(errs)
		}
		return keys, nil
	}
}

type MockKey struct {
	Kind     string
	IntID    int64
	StringID string
}

type EntitiesStorage map[string]map[MockKey][]datastore.Property

type MockDB struct {
	EntitiesByKind EntitiesStorage
}

func NewMockKeyFromDatastoreKey(key *datastore.Key) MockKey {
	return MockKey{Kind: key.Kind(), IntID: key.IntID(), StringID: key.StringID()}
}

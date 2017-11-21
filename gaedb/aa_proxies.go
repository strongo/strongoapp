package gaedb

import (
	"bytes"
	"fmt"
	//"github.com/strongo/nds"
	"github.com/strongo/app/log"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"os"
	"github.com/strongo/app/db"
	"strconv"
	"strings"
	"github.com/pkg/errors"
	"github.com/strongo/nds"
)

var (
	LoggingEnabled = true // TODO: move to Context.WithValue()
	mockDB           MockDB
	NewIncompleteKey = datastore.NewIncompleteKey
	NewKey           = datastore.NewKey

	//dbRunInTransaction = datastore.RunInTransaction
	//dbGet = datastore.Get
	//dbGetMulti = datastore.GetMulti
	//dbPut = datastore.Put
	//dbPutMulti = datastore.PutMulti
	//dbDelete = datastore.Delete
	//dbDeleteMulti = datastore.DeleteMulti

	dbRunInTransaction = nds.RunInTransaction
	dbGet = nds.Get
	dbGetMulti = datastore.GetMulti
	dbPut = nds.Put
	dbPutMulti = nds.PutMulti
	dbDelete = nds.Delete
	dbDeleteMulti = nds.DeleteMulti

	RunInTransaction = func (c context.Context, f func(tc context.Context) error, opts *datastore.TransactionOptions) error {
		if LoggingEnabled {
			if opts == nil {
				log.Debugf(c, "gaedb.RunInTransaction(): starting transaction, opts=nil...")
			} else {
				log.Debugf(c, "gaedb.RunInTransaction(): starting transaction, opts=%+v...", *opts)
			}
		}
		attempt := 0
		fWrapped := func(c context.Context) (err error) {
			attempt += 1
			log.Debugf(c, "tx attempt #%d", attempt)
			if err = f(c); err != nil {
				m := fmt.Sprintf("tx attempt #%d failed: ", attempt)
				if err == datastore.ErrConcurrentTransaction {
					log.Warningf(c, m + err.Error())
				} else {
					log.Errorf(c, m + err.Error())
				}
			}
			return
		}
		if err := dbRunInTransaction(c, fWrapped, opts); err != nil {
			if LoggingEnabled {
				if strings.Contains(err.Error(),"nested transactions are not supported") {
					panic(err)
				}
				log.Errorf(c, errors.WithMessage(err, "transaction failed").Error())
			}
			return err
		}
		if LoggingEnabled {
			log.Debugf(c, "Transaction successful")
		}
		return nil
	}

	Put = func(c context.Context, key *datastore.Key, val interface{}) (*datastore.Key, error) {
		if val == nil {
			panic("val == nil")
		}
		var err error
		isPartialKey := key.Incomplete()
		if LoggingEnabled {
			buf := new(bytes.Buffer)
			fmt.Fprintf(buf, "dbPut(%v) => properties:", key2str(key))
			if props, err := datastore.SaveStruct(val); err != nil {
				return nil, errors.WithMessage(err, fmt.Sprintf("failed to SaveStruct(%v) to properties", val))
			} else {
				var prevPropName string
				for _, prop := range props {
					if prop.Name == prevPropName {
						fmt.Fprintf(buf, ", %v", prop.Value)
					} else {
						fmt.Fprintf(buf, "\n\t%v: %v", prop.Name, prop.Value)
					}
					prevPropName = prop.Name
				}
			}
			log.Debugf(c, buf.String())
		}
		if key, err = dbPut(c, key, val); err != nil {
			return key, errors.WithMessage(err, "failed to put to db " + key2str(key))
		} else if LoggingEnabled && isPartialKey {
			log.Debugf(c, "dbPut() inserted new record with key: " + key2str(key))
		}
		return key, err
	}

	PutMulti = func(c context.Context, keys []*datastore.Key, vals interface{}) ([]*datastore.Key, error) {
		if LoggingEnabled {
			//buf := new(bytes.Buffer)
			//buf.WriteString("dbPutMulti(")
			//for i, key := range keys {
			// Need to use reflection...
			//}
			//buf.WriteString(")")
			logKeys(c, "dbPutMulti", keys)
		}
		return dbPutMulti(c, keys, vals)
	}

	Get = func(c context.Context, key *datastore.Key, val interface{}) error {
		if LoggingEnabled {
			log.Debugf(c, "dbGet(%v)", key2str(key))
		}
		if key.IntID() == 0 && key.StringID() == "" {
			panic("key.IntID() == 0 && key.StringID() is empty string")
		}
		return dbGet(c, key, val)
	}

	GetMulti = func(c context.Context, keys []*datastore.Key, vals interface{}) error {
		if LoggingEnabled {
			logKeys(c, "dbGetMulti", keys)
		}
		return dbGetMulti(c, keys, vals)
	}

	Delete = func(c context.Context, key *datastore.Key) error {
		log.Warningf(c, "gaedb.Delete(%v)", key2str(key))
		return dbDelete(c, key)
	}

	DeleteMulti = func(c context.Context, keys []*datastore.Key) error {
		log.Warningf(c, "Deleting %v entities", len(keys))
		logKeys(c, "dbDeleteMulti", keys)
		return dbDeleteMulti(c, keys)
	}
)

func key2str(key *datastore.Key) string {
	if key == nil {
		return "nil"
	}
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
	prevKey := "-"
	for _, key := range keys {
		ks := key2str(key)
		if ks == prevKey {
			log.Errorf(c, "Duplicate keys: " + ks)
		}
		buffer.WriteString(fmt.Sprintf("\t%v\n", ks))
		prevKey = ks
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

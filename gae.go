package strongo

import (
	"time"
	"encoding/binary"
	"google.golang.org/appengine"
	"encoding/base64"
	"golang.org/x/net/context"
)

func SignIdWithExpiry(c context.Context, id int64, expires time.Time) string {
	toSign := make([]byte, 16)
	binary.LittleEndian.PutUint64(toSign[:7], uint64(id))
	binary.LittleEndian.PutUint64(toSign[8:], uint64(expires.Unix()))
	_, signature, err := appengine.SignBytes(c, toSign)
	if err != nil {
		panic(err.Error())
	}
	return base64.URLEncoding.EncodeToString(signature)
}



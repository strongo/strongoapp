package delaying

import (
	"context"
	"fmt"
	"github.com/strongo/log"
)

func VoidWithLog(key string, i any) Function {
	doNothing := func() {}
	return NewFunction(key, doNothing,
		func(c context.Context, params Params, args ...interface{}) error {
			log.Debugf(c, fmt.Sprintf("%s.EnqueueWork(%+v): %+v", key, args, params))
			return nil
		},
		func(c context.Context, params Params, args ...[]interface{}) error {
			log.Debugf(c, fmt.Sprintf("%s.EnqueueWorkMulti(%+v): %+v", key, args, params))
			return nil
		},
	)
}

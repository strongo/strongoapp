package gae

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/taskqueue"
	"net/url"
	"github.com/strongo/app/log"
	"github.com/strongo/app/gaedb"
)

//TODO: Document why whe need this
var CallDelayFunc = func(c context.Context, queueName, subPath string, f *delay.Function, args ...interface{}) error {
	if task, err := CreateDelayTask(queueName, subPath, f, args...); err != nil {
		return err
	} else {
		_, err = AddTaskToQueue(c, task, queueName)
		return err
	}
}

//TODO: Document why whe need this
func CreateDelayTask(queueName, subPath string, f *delay.Function, args ...interface{}) (*taskqueue.Task, error) {
	if queueName == "" {
		return nil, errors.New("queueName is empty")
	}
	if queueName == "default" {
		return nil, errors.New("queueName is 'default'")
	}
	if subPath == "" {
		return nil, errors.New("subPath is empty")
	}
	if task, err := f.Task(args...); err != nil {
		return task, err
	} else {
		task.Path += fmt.Sprintf("?task=%v&queue=%v", url.QueryEscape(subPath), url.QueryEscape(queueName))
		return task, nil
	}
}

//TODO: Document why whe need this
var AddTaskToQueue = func(c context.Context, task *taskqueue.Task, queueName string) (*taskqueue.Task, error) {
	if queueName == "" {
		return nil, errors.New("queueName is empty")
	}
	if queueName == "default" {
		return nil, errors.New("queueName is 'default'")
	}
	isInTransaction := gaedb.NewDatabase().IsInTransaction(c)
	log.Debugf(c, "Adding task to queue '%v', tx=%v): %+v", queueName, isInTransaction, task)
	return taskqueue.Add(c, task, queueName)
}

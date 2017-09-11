package db

import (
	"fmt"
	"github.com/pkg/errors"
)

type ErrDuplicateUser struct { // TODO: Should it be moved out of this package?
	SearchCriteria   string
	DuplicateUserIDs []int64
}

func (err ErrDuplicateUser) Error() string {
	return fmt.Sprintf("Multiple users by given search criteria[%v]: %v", err.SearchCriteria, err.DuplicateUserIDs)
}

var (
	ErrRecordNotFound     = errors.New("Record not found")
)

func IsNotFound(err error) bool {
	_, ok := err.(ErrNotFoundByID)
	return ok || errors.Cause(err) == ErrRecordNotFound
}

type IntIdentifier interface {
	IntID() int64
}

type StrIdentifier interface {
	StrID() string
}

type ErrNotFoundByID interface {
	error
	IntIdentifier
	StrIdentifier
}

type errNotFoundByID struct {
	kind  string
	intID int64
	strID string
	cause error
}

func (e errNotFoundByID) IntID() int64 {
	return e.intID
}

func (e errNotFoundByID) StrID() string {
	return e.strID
}

func (e errNotFoundByID) ID() interface{} {
	if e.intID != 0 && e.strID == "" {
		return e.intID
	} else if e.strID != "" && e.intID == 0 {
		return e.strID
	} else if e.intID == 0 && e.strID == "" {
		return nil
	} else {
		panic("e.intID != 0 && e.strID is not empty string")
	}
}

func (e errNotFoundByID) Cause() error {
	return e.cause
}

func (e errNotFoundByID) Error() string {
	return fmt.Sprintf("'%v' not found by id=%v: %v", e.kind, e.ID(), e.cause)
}

func NewErrNotFoundByIntID(kind string, id int64, cause error) error {
	return errNotFoundByID{kind: kind, intID: id, cause: errNotFoundCause(cause)}
}

func NewErrNotFoundByStrID(kind string, id string, cause error) error {
	return errNotFoundByID{kind: kind, strID: id, cause: errNotFoundCause(cause)}
}

func errNotFoundCause(cause error) error {
	if cause == nil || cause == ErrRecordNotFound {
		return ErrRecordNotFound
	} else {
		return errors.WithMessage(ErrRecordNotFound, cause.Error())
	}
}

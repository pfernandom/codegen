package errors

import (
	"fmt"
	"log"
)

type mustError[T any] struct {
	V   T
	Err error
}

func (me mustError[T]) Error() string {
	return me.Err.Error()
}

func Must(err error) mustError[any] {
	return mustError[any]{Err: err}
}

func MustNow(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func MustReturn[T any](v T, err error) mustError[T] {
	return mustError[T]{V: v, Err: err}
}

func (me mustError[T]) OrFail() T {
	if me.Err != nil {
		log.Fatal(me)
	}
	return me.V
}

func (me mustError[T]) OrFailWith(message string, args ...any) T {
	if me.Err != nil {
		all_args := append(args, me.Err)
		log.Fatal(fmt.Errorf(message, all_args...))
	}
	return me.V
}

func MustOrMessage(err error, message string, args ...any) {
	if err != nil {
		log.Fatal(fmt.Errorf(message, args...))
	}
}

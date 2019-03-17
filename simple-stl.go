package simple_stl

import (
	"reflect"
	"sync"
)

type Reader interface {
	Read () Reader
}
type Writer interface {
	Write (Writer)
}

type Executer interface {

}


type Sdata struct {
	Model reflect.Type
	Data []interface{}
	Rwm sync.RWMutex
	Err error
}

func NewSdata(model interface{})Sdata{
	return Sdata{
		Model:reflect.TypeOf(model),
		Data: make([]interface{},0),
	}
}


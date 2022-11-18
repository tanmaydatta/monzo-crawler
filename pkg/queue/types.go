package queue

import "errors"

var EOF = errors.New("EOF")

type IQueue interface {
	EnQueue(Element) error
	DeQueue() (Element, error)
	Close() error
	GetName() string
}

type Element interface {
	GetBaseURL() string
	GetData() interface{}
	GetType() string
	SetData(interface{}) error
}

type IReader interface {
	Read() (Element, error)
	Close() error
}

type IWriter interface {
	Write(Element) error
	Close() error
}

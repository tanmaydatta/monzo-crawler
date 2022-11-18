package queue

import (
	"errors"
	"fmt"
)

type queueElement struct {
	baseURL string
	typ     string
}

type queueElementToFetch struct {
	queueElement
	data *FetchElementData
}

type queueElementFetched struct {
	queueElement
	data *FetchedElementData
}

type FetchedElementData struct {
	Urls    []string
	Depth   int
	BaseUrl string
	CurUrl  string
	Path    string
	Robots  string
}

type FetchElementData struct {
	Robots  string // TODO:
	Path    string
	BaseUrl string
	CurUrl  string
	Depth   int
}

func (q *queueElement) GetBaseURL() string {
	return q.baseURL
}

func (q *queueElementToFetch) GetData() interface{} {
	return q.data
}

func (q *queueElementFetched) GetData() interface{} {
	return q.data
}

func (q *queueElement) GetType() string {
	return q.typ
}

func (q *queueElementToFetch) SetData(i interface{}) error {
	if data, ok := i.(*FetchElementData); ok {
		q.data = data
	}
	return errors.New(fmt.Sprintf("data should be of type *ElementData but got %v\n", i))
}

func (q *queueElementFetched) SetData(i interface{}) error {
	if data, ok := i.(*FetchedElementData); ok {
		q.data = data
	}
	return errors.New(fmt.Sprintf("data should be of type []string but got %v\n", i))
}

func NewFetchQueueElement(e *FetchElementData, baseURL string, typ string) Element {
	return &queueElementToFetch{
		queueElement: queueElement{baseURL: baseURL, typ: typ}, data: e,
	}
}

func NewFetchedQueueElement(data *FetchedElementData, baseURL string, typ string) Element {
	return &queueElementFetched{
		queueElement: queueElement{baseURL: baseURL, typ: typ}, data: data,
	}
}

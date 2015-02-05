package goloc

import ()

type LinkedList struct {
	Size  int
	First *LinkedElement
	Last  *LinkedElement
}

func NewLinkedList() *LinkedList {
	list := new(LinkedList)
	return list
}

type LinkedElement struct {
	Element interface{}
	Next    *LinkedElement
}

func (list *LinkedList) AddFirst(element interface{}) *LinkedElement {
	nl := new(LinkedElement)
	nl.Element = element
	nl.Next = list.First
	list.First = nl
	list.Size++
	return nl
}

func (list *LinkedList) AddLast(element interface{}) *LinkedElement {
	nl := new(LinkedElement)
	nl.Element = element
	if list.Last != nil {
		list.Last.Next = nl
	}
	list.Last = nl
	if list.First == nil {
		list.First = nl
	}
	list.Size++
	return nl
}

func (list *LinkedList) ToArray() []interface{} {
	array := make([]interface{}, list.Size)
	var elem *LinkedElement
	var i int
	for elem = list.First; elem != nil; elem = elem.Next {
		array[i] = elem.Element
		i++
	}
	return array
}

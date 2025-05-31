package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	firstItem *ListItem
	lastItem  *ListItem
	length    int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Back() *ListItem {
	return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	l.pushFront(item)
	return item
}

func (l *list) pushFront(item *ListItem) {
	if l.lastItem == nil {
		l.lastItem = item
	}

	if l.firstItem != nil {
		l.firstItem.Prev = item
		item.Next = l.firstItem
	}

	l.firstItem = item
	l.length++
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
	}

	if l.firstItem == nil {
		l.firstItem = item
	}

	if l.lastItem != nil {
		l.lastItem.Next = item
		item.Prev = l.lastItem
	}

	l.lastItem = item
	l.length++

	return l.lastItem
}

func (l *list) Remove(item *ListItem) {
	l.remove(item)
}

func (l *list) remove(item *ListItem) {
	if item.Next != nil {
		item.Next.Prev = item.Prev
	}

	if item.Prev != nil {
		item.Prev.Next = item.Next
	}

	if item == l.firstItem {
		l.firstItem = item.Next
	} else if item == l.lastItem {
		l.lastItem = l.lastItem.Prev
	}

	item.Next = nil
	item.Prev = nil

	l.length--
}

func (l *list) MoveToFront(item *ListItem) {
	l.remove(item)
	l.pushFront(item)
}

func NewList() List {
	return new(list)
}

package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	Clear()
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFirst(v interface{}) *ListItem {
	l.len++
	newEl := ListItem{v, nil, l.front}
	l.front = &newEl
	l.back = &newEl
	return &newEl
}

func (l *list) PushFront(v interface{}) *ListItem {
	if l.Len() == 0 {
		return l.PushFirst(v)
	}
	l.len++
	newEl := ListItem{v, l.front, nil}
	l.front.Prev = &newEl
	l.front = &newEl
	return &newEl
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.Len() == 0 {
		return l.PushFirst(v)
	}
	l.len++
	newEl := ListItem{v, nil, l.back}
	l.back.Next = &newEl
	l.back = &newEl
	return &newEl
}

func (l *list) Remove(i *ListItem) {
	switch {
	case i == l.Back():
		l.back = i.Prev
		i.Prev.Next = nil
	case i == l.Front():
		l.front = i.Next
		i.Next.Prev = nil
	default:
		tmp := i.Prev.Next
		i.Prev.Next = i.Next
		i.Next.Prev = tmp
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.Front() {
		return
	}
	l.Remove(i)
	l.len++
	i.Next = l.Front()
	i.Next.Prev = i
	i.Prev = nil
	l.front = i
}

func (l *list) Clear() {
	l.len = 0
	l.back = nil
	l.front = nil
}

func NewList() List {
	return new(list)
}

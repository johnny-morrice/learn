package collections

type List[T any] struct {
	Len      int
	zeroNode *ListNode[T]
	lastNode *ListNode[T]
}

func (l List[T]) Copy() List[T] {
	if l.Len == 0 {
		return List[T]{}
	}
	fromNode := l.zeroNode
	toNodeRoot := &ListNode[T]{}
	toNode := toNodeRoot
	toNode.Value = fromNode.Value
	for fromNode.Next != nil {
		fromNode = fromNode.Next
		nextToNode := &ListNode[T]{}
		nextToNode.Value = fromNode.Value
		toNode.Next = nextToNode
		toNode = nextToNode
	}
	return List[T]{
		zeroNode: toNodeRoot,
		lastNode: toNode,
		Len:      l.Len,
	}
}

func (l List[T]) Slice() []T {
	sl := []T{}
	node := l.zeroNode
	for i := 0; i < l.Len; i++ {
		sl = append(sl, node.Value)
		node = node.Next
	}
	return sl
}

func (l List[T]) Append(val T) List[T] {
	node := &ListNode[T]{
		Value: val,
	}
	newList := l.Copy()
	if newList.Len == 0 {
		return List[T]{
			zeroNode: node,
			lastNode: node,
			Len:      1,
		}
	}
	newList.lastNode.Next = node
	newList.lastNode = node
	newList.Len++
	return newList
}

type ListNode[T any] struct {
	Value T
	Next  *ListNode[T]
}

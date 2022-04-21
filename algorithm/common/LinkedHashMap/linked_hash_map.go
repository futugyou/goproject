package LinkedHashMap

type Node struct {
	Key   int
	Value int
	Prev  *Node
	Next  *Node
}

func NewNode(k, v int) *Node {
	n := Node{}
	n.Key = k
	n.Value = v
	return &n
}

type DoubleList struct {
	head *Node
	tail *Node
	size int
}

func NewDoubleList() *DoubleList {
	d := DoubleList{}
	d.head = NewNode(0, 0)
	d.tail = NewNode(0, 0)
	d.head.Next = d.tail
	d.tail.Prev = d.head
	d.size = 0
	return &d
}

func (d *DoubleList) AddLast(x *Node) {
	x.Prev = d.tail.Prev
	x.Next = d.tail
	d.tail.Prev.Next = x
	d.tail.Prev = x
	d.size++
}

func (d *DoubleList) Remove(x *Node) {
	x.Prev.Next = x.Next
	x.Next.Prev = x.Prev
	d.size--
}

func (d *DoubleList) RemoveFirst() *Node {
	if d.head.Next == d.tail {
		return nil
	}
	first := d.head.Next
	d.Remove(first)
	return first
}

func (d *DoubleList) Size() int {
	return d.size
}

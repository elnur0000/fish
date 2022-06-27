package datastruct

import "fmt"

type Node struct {
	Val    any
	Parent *LinkedList
	Prev   *Node
	Next   *Node
}

type LinkedList struct {
	Head   *Node
	Tail   *Node
	Length int
}

func (l *LinkedList) Insert(newNode *Node) {
	l.Length++

	if l.Head == nil {
		l.Head = newNode
		l.Tail = newNode
		return
	}

	l.Tail.Next = newNode
	newNode.Prev = l.Tail
	l.Tail = newNode

}

func (l *LinkedList) Remove(node *Node) {
	if node.Parent != l {
		return
	}

	if node.Prev == nil {
		l.Head = node.Next
	} else {
		node.Prev.Next = node.Next
	}

	if node.Next == nil {
		l.Tail = node.Prev
	} else {
		node.Next.Prev = node.Prev
	}

	node.Prev = nil
	node.Next = nil

	l.Length--
}

func (l *LinkedList) Print() {
	for node := l.Head; node != nil; node = node.Next {
		fmt.Println(node)
	}
}

func (l *LinkedList) Len() int {
	length := 0
	for node := l.Head; node != nil; node = node.Next {
		length++
	}
	return length
}

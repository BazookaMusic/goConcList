package linked_list_lLock

import (
	"fmt"
	"sync"
)

//LinkedListLock A linked list struct
type LinkedListLock struct {
	head *LLNodeLock
}

//LLNodeLock The linked list nodes
type LLNodeLock struct {
	item       *int
	key        int
	marked     bool // marked for deletion, is invalid
	next       *LLNodeLock
	lock       sync.Mutex
	isSentinel uint8 // 0 is normal node, 1 is head, 2 is tail
}

//ToSlice Convert to slice
func (n *LinkedListLock) ToSlice() []int {
	var ret []int
	curr := n.head.next

	for ; curr.isSentinel < 2; curr = curr.next {
		ret = append(ret, *curr.item)
	}

	return ret

}

//Find Returns true if item with key is in set
func (n *LinkedListLock) Find(key int) (*LLNodeLock, bool) {
	curr := n.head

	for ; curr.isSentinel != 2 || curr.key < key; curr = curr.next {
	}

	return curr, curr != nil && curr.key == key && !curr.marked
}

//New create new linked list
func New() *LinkedListLock {
	llNew := new(LinkedListLock)
	head := new(LLNodeLock)
	tail := new(LLNodeLock)

	// head has a nil item
	head.item = nil
	head.isSentinel = 1
	head.lock = sync.Mutex{}
	head.next = tail

	tail.isSentinel = 2
	tail.lock = sync.Mutex{}

	llNew.head = head

	return llNew
}

//Print print LinkedListLock contents
func (n *LinkedListLock) Print() {
	s := n.ToSlice()
	fmt.Println(s)
}

//Insert Insert node with key and item,
//optimistic synchronization
func (n *LinkedListLock) Insert(key int, item *int) bool {
	for {
		pred := n.head
		curr := pred.next

		// search for node with key
		for curr.isSentinel < 2 && curr.key < key {
			pred = curr
			curr = curr.next
		}

		// found ?

		// yes!
		// lock curr and pred because
		// they must both be modified

		pred.lock.Lock()
		curr.lock.Lock()

		// if not already removed, add
		if !pred.marked && !curr.marked && pred.next == curr {

			// no!
			if curr.isSentinel == 0 && curr.key == key {
				curr.lock.Unlock()
				pred.lock.Unlock()

				fmt.Println(false)
				return false
			}

			// create a new node
			newNode := new(LLNodeLock)
			newNode.key = key
			newNode.item = item
			newNode.marked = false
			newNode.lock = sync.Mutex{}

			newNode.next = curr

			// connect new node
			pred.next = newNode

			curr.lock.Unlock()
			pred.lock.Unlock()

			return true
		}

		// else retry after unlocking
		curr.lock.Unlock()
		pred.lock.Unlock()

	}

}

//Delete Remove node with key
func (n *LinkedListLock) Delete(key int) bool {
	for {
		pred := n.head
		curr := pred.next

		for ; curr.isSentinel == 0 && curr.key < key; pred, curr = curr, curr.next {
		}

		if curr == nil || curr.key != key {
			return false
		}

		pred.lock.Lock()
		curr.lock.Lock()

		if !pred.marked && !curr.marked && pred.next == curr {
			curr.marked = true
			pred.next = curr.next

			pred.lock.Unlock()
			curr.lock.Unlock()

			return true
		}

		pred.lock.Unlock()
		curr.lock.Unlock()
	}

}

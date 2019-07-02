package ll

import (
	"fmt"
	"sync"
)

//ll A linked list struct
type ll struct {
	head *ll_node
}

//ll_node The linked list nodes
type ll_node struct {
	item   *int
	key    int
	marked bool // marked for deletion, is invalid
	next   *ll_node
	lock   sync.Mutex
	isHead bool
}

//ToSlice Convert to slice
func (n *ll) ToSlice() []*int {
	var ret []*int
	curr := n.head.next

	for ; curr != nil; curr = curr.next {
		ret = append(ret, curr.item)
	}

	return ret

}

//Find Returns true if item with key is in set
func (n *ll) Find(key int) (*ll_node, bool) {
	curr := n.head

	for ; curr != nil && (curr.isHead || curr.key < key); curr = curr.next {
	}

	return curr, curr != nil && curr.key == key && !curr.marked
}

//New create new linked list
func New() *ll {
	ll_n := new(ll)
	head := new(ll_node)

	// head has a nil item
	head.item = nil
	head.isHead = true
	ll_n.head = head
	return ll_n
}

//Print print ll contents
func (n *ll) Print() {
	s := n.ToSlice()
	fmt.Println(s)
}

//Insert Insert node with key and item,
//optimistic synchronization
func (n *ll) Insert(key int, item *int) bool {
	for {
		pred := n.head
		curr := pred.next

		// search for node with key
		for ; curr != nil && (curr.isHead || curr.key < key); pred, curr = curr, curr.next {
		}

		// found ?

		// no!
		if curr != nil && curr.key == key {
			return false
		}

		// yes!
		// lock curr and pred because
		// they must both be modified

		pred.lock.Lock()
		if curr != nil {
			curr.lock.Lock()
		}

		// if not already removed, add
		if !pred.marked && (curr == nil || !curr.marked) && pred.next == curr {
			new_node := new(ll_node)
			new_node.key = key
			new_node.item = item
			new_node.marked = false
			new_node.lock = sync.Mutex{}
			pred.next = new_node
			new_node.next = curr

			pred.lock.Unlock()
			if curr != nil {
				curr.lock.Unlock()
			}

			return true
		}

		// else retry after unlocking

		pred.lock.Unlock()
		if curr != nil {
			curr.lock.Unlock()
		}
	}

}

//Delete Remove node with key
func (n *ll) Delete(key int) bool {
	for {
		pred := n.head
		curr := pred.next

		for ; curr != nil && (curr.isHead || curr.key < key); pred, curr = curr, curr.next {
		}

		if curr == nil || curr.key != key {
			return false
		}

		pred.lock.Lock()
		curr.lock.Lock()

		if !pred.marked && !curr.marked && pred.next == curr {
			pred.next = curr.next

			pred.lock.Unlock()
			curr.lock.Unlock()

			return true
		}

		pred.lock.Unlock()
		curr.lock.Unlock()
	}

}

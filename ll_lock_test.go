package linked_list_lock

import (
	"sync"
	"testing"
)

var iterations = 5000

func insertAndWait(index int, item *int, a *LinkedListLock, wg *sync.WaitGroup) {
	defer wg.Done()
	a.Insert(index, item)
}

func delAndWait(index int, a *LinkedListLock, wg *sync.WaitGroup) {
	defer wg.Done()
	a.Delete(index)
}
func TestInsertDelete(t *testing.T) {

	a := New()
	var wg sync.WaitGroup

	wg.Add(iterations)

	for index := iterations; index > 0; index-- {
		i := new(int)
		*i = index
		go insertAndWait(index, i, a, &wg)

	}
	wg.Wait()
	wg.Add(iterations / 2)

	for index := iterations / 2; index > 0; index-- {
		go delAndWait(index, a, &wg)

	}

	wg.Wait()
	a.Print()

	counter := iterations/2 + 1
	for _, item := range a.ToSlice() {
		if item < iterations/2 {
			t.Error("Node which should be deleted found")
		}

		if item != counter {
			t.Errorf("Inserted item %d not found", item)
		}
		counter++

	}

}

package ll

import (
	"sync"
	"testing"
)

var iterations = 500

func insertAndWait(index int, item *int, a *ll, wg *sync.WaitGroup) {
	defer wg.Done()
	a.Insert(index, item)
}

func delAndWait(index int, a *ll, wg *sync.WaitGroup) {
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

}

package ll

import (
	"fmt"
	"testing"
)

func Verify(a *ll) {

}
func TestInsert(t *testing.T) {

	a := New()
	defer a.Print()

	c := 13
	b := 14
	fmt.Println(a.Insert(12, &c))

	fmt.Println(a.Insert(13, &b))

	for index := 14; index > 0; index-- {
		i := new(int)
		*i = 13
		go a.Insert(index, i)

	}

	for index := 5; index > 0; index-- {
		go a.Delete(index)

	}

}

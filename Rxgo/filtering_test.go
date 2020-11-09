package rxgo_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pmlpml/rxgo"
)

func TestDebounce(t *testing.T) {
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Map(func(x int) int {
		if x != 0 {
			time.Sleep(1 * time.Millisecond)
		}
		return x
	}).Debounce(2 * time.Millisecond).Subscribe(func(x int) {
		if x != 9 {
			fmt.Printf("error Debounce with %d\n", x)
			os.Exit(-1)
		}
		fmt.Printf("Debunce %d\n", x)
	})
}

func TestDistinct(t *testing.T) {
	var all = map[int]bool{}
	rxgo.Just(0, 1, 2, 3, 6, 5, 6, 2, 3, 9).Distinct().Subscribe(func(x int) {
		if _, ok := all[x]; !ok {
			all[x] = true
			return
		}
		fmt.Printf("error Distinct with %d\n", x)
		os.Exit(-1)
	})
}

func TestElementAt(t *testing.T) {
	rxgo.Just(0, 1, 2, 3, 6, 5, 6, 2, 3, 9).ElementAt(1).Subscribe(func(x int) {
		if x != 1 {
			fmt.Printf("error ElementAt with %d\n", x)
			os.Exit(-1)
		}
	})
}

func TestFirst(t *testing.T) {
	rxgo.Just(0, 1, 2, 3, 6, 5, 6, 2, 3, 9).First().Subscribe(func(x int) {
		if x != 0 {
			fmt.Printf("error First with %d\n", x)
			os.Exit(-1)
		}
	})
}

func TestIgnoreElements(t *testing.T) {
	rxgo.Just(0, 1, 2, 3, 6, 5, 6, 2, 3, 9).IgnoreElements().Subscribe(func(x int) {
		fmt.Printf("error IgnoreElements with %d\n", x)
		os.Exit(-1)
	})
}

/*func TestLast(t *testing.T) {
	rxgo.Just(0, 1, 2, 3, 6, 5, 6, 2, 3, 9).Last().Subscribe(func(x int) {
		if x != 9 {
			fmt.Printf("error Last with %d\n", x)
			os.Exit(-1)
		}
	})
}*/
func TesteLast() {
	rxgo.Just(33, 1, 0, 215, 4, 6).Last().Subscribe(func(x int) {
		fmt.Print(x)
	})
	fmt.Println()
	//Output:6
}

func TestSample(t *testing.T) {
	var samplearr = []int{2, 4, 6, 8, 9}
	var count = 0
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Map(func(x int) int {
		if x != 0 {
			time.Sleep(500 * time.Millisecond)
		}
		return x
	}).Sample(1 * time.Second).Subscribe(func(x int) {
		if x != samplearr[count] {
			fmt.Printf("error Sample with %d\n", x)
			os.Exit(-1)
		}
		count++
	})
}

func TestSkip(t *testing.T) {
	var skiparr = []int{4, 5, 6, 7, 8, 9}
	var count = 0
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Skip(4).Subscribe(func(x int) {
		if x != skiparr[count] {
			fmt.Printf("error Skip with %d\n", x)
			os.Exit(-1)
		}
		count++
	})
}

func TestSkiplast(t *testing.T) {
	var skiparr = []int{0, 1, 2}
	var count = 0
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Skiplast(7).Subscribe(func(x int) {
		if x != skiparr[count] {
			fmt.Printf("error Skiplast with %d\n", x)
			os.Exit(-1)
		}
		count++
	})
}

func TestTake(t *testing.T) {
	var takearr = []int{0, 1, 2}
	var count = 0
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Take(3).Subscribe(func(x int) {
		if x != takearr[count] {
			fmt.Printf("error Take with %d\n", x)
			os.Exit(-1)
		}
		count++
	})
}

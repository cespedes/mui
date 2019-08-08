package frontend

import (
	"fmt"
	"sort"
	"os"
)

type FrontEnd interface {
	Name() string
	Priority() int
	Available() bool
	Question() int
	Input() string
}

var frontendList = make([]FrontEnd, 0)

func RegisterFrontEnd(new_frontend FrontEnd) {
	frontendList = append(frontendList, new_frontend)
	sort.Slice(frontendList, func(i, j int) bool {
		return frontendList[i].Priority() > frontendList[j].Priority()
	})
}

func chooseFrontend() FrontEnd {
	for _, f := range frontendList {
		if f.Available() {
			fmt.Printf("DEBUG: Found frontend with priority %d\n", f.Priority())
			return f
		}
	}
	return nil
}

func Question(args []string) {
	f := chooseFrontend()
	if f == nil {
		fmt.Println("DEBUG: frontend.Question(): No frontends are available")
		os.Exit(1)
	}
	os.Exit(f.Question())
}

func Input(args []string) {
	f := chooseFrontend()
	if f == nil {
		fmt.Println("DEBUG: frontend.Question(): No frontends are available")
		os.Exit(1)
	}
	s := f.Input()
	fmt.Println(s)
	os.Exit(0)
}

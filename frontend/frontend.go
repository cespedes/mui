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
}

var frontendList = make([]FrontEnd, 0)

func RegisterFrontEnd(new_frontend FrontEnd) {
	frontendList = append(frontendList, new_frontend)
	sort.Slice(frontendList, func(i, j int) bool {
		return frontendList[i].Priority() > frontendList[j].Priority()
	})
}

func Question(args []string) {
	for _, f := range frontendList {
		if f.Available() {
			fmt.Printf("DEBUG: Found frontend with priority %d\n", f.Priority())
			os.Exit(f.Question())
			return
		}
	}
	fmt.Println("DEBUG: frontend.Question(): No frontends are available")
	os.Exit(1)
}

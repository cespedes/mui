package frontend

import (
	"fmt"
	"sort"
	"os"
)

type FrontEnd interface {
	Priority() int
	Available() bool
	Question()
}

var frontendList = make([]FrontEnd, 0)

func RegisterFrontEnd(new_frontend FrontEnd) {
	frontendList = append(frontendList, new_frontend)
	sort.Slice(frontendList, func(i, j int) bool {
		return frontendList[i].Priority() > frontendList[j].Priority()
	})
}

func Question() {
	for _, f := range frontendList {
		if f.Available() {
			f.Question()
			return
		}
	}
	fmt.Println("frontend.Question(): No frontends are available")
	os.Exit(1)
}

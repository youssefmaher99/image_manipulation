package data

import "fmt"

type InMemory []string

var InMemoryArchives InMemory

func (m *InMemory) FileExist(fileName string) bool {
	for i := 0; i < len(*m); i++ {
		if (*m)[i] == fileName {
			return true
		}
	}
	return false
}

func (m *InMemory) Remove(fileName string) {
	fmt.Printf("Before : %v\n", *m)
	for i := 0; i < len(*m); i++ {
		if (*m)[i] == fileName {
			*m = append((*m)[:i], (*m)[i+1:]...)
			break
		}
	}
	fmt.Printf("After : %v\n", *m)

}

func (m *InMemory) Add(fileName string) {
	(*m) = append((*m), fileName)
}

package data

type InMemory []string

var InMemoryArchives InMemory
var InMemoryUUID InMemory

func (m *InMemory) ItemExist(item string) bool {
	for i := 0; i < len(*m); i++ {
		if (*m)[i] == item {
			return true
		}
	}
	return false
}

func (m *InMemory) Remove(item string) {
	for i := 0; i < len(*m); i++ {
		if (*m)[i] == item {
			*m = append((*m)[:i], (*m)[i+1:]...)
			break
		}
	}

}

func (m *InMemory) Add(item string) {
	(*m) = append((*m), item)
}

package data

type InMemory map[string]struct{}

var InMemoryArchives InMemory = make(InMemory)
var InMemoryUUID InMemory = make(InMemory)

func (m InMemory) ItemExist(item string) bool {
	if _, ok := m[item]; ok {
		return true
	}
	return false
}

func (m InMemory) Remove(item string) {
	delete(m, item)
}

func (m InMemory) Add(item string, s struct{}) {
	m[item] = s
}

package util

type HashSet map[interface{}]struct{}

func NewHashSet() HashSet {
	return make(HashSet)
}

func (set HashSet) Add(item interface{}) {
	set[item] = struct{}{}
}

func (set HashSet) Remove(item interface{}) {
	delete(set, item)
}

func (set HashSet) Contains(item interface{}) bool {
	_, found := set[item]
	return found
}

func (set HashSet) Size() int {
	return len(set)
}

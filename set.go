package main

type Set map[string]struct{}

func NewSet() Set {
	return make(map[string]struct{})
}

func Init(values []string) Set {
	set := NewSet()
	for _, value := range values {
		set.Add(value)
	}
	return set
}

func (s *Set) Add(value string) {
	(*s)[value] = struct{}{}
}

func (s *Set) Remove(value string) {
	delete((*s), value)
}

func (s *Set) Contains(value string) bool {
	_, contains := (*s)[value]
	return contains
}

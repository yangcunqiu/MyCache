package cache

type String string

func (s String) Len() int {
	return len(s)
}

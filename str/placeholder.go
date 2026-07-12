package str

import (
	"fmt"
	"strings"
)

type Entry struct {
	Key   string
	Value string
}
type Placeholder []Entry

func (self *Placeholder) Append(key string, value any) {
	*self = append(*self, Entry{key, fmt.Sprintf("%v", value)})
}

func (self *Placeholder) Remove(key string) {
	for i, e := range *self {
		if e.Key == key {
			*self = append((*self)[:i], (*self)[i+1:]...)
			return
		}
	}
}

func (self *Placeholder) Rename(oldKey string, newKey string) {
	for _, entry := range *self {
		if entry.Key == oldKey {
			entry.Key = newKey
			return
		}
	}
}

// InteporlatePlaceholder Reemplaza los placeholders en el texto
func (ph Placeholder) Replace(text string) string {
	for _, entry := range ph {
		text = strings.ReplaceAll(text, ":"+entry.Key, entry.Value)
	}
	return text
}

// forma de uso:
// ph := NewPlaceholder("key1", "value1", "key2", "value2")
func NewPlaceholder(s ...string) Placeholder {
	ph := Placeholder{}
	for i := 0; i < len(s); i += 2 {
		ph = append(ph, Entry{s[i], s[i+1]})
	}
	return ph
}

// forma de uso:
// ph := NewPlaceholders("key1:value1", "key2:value2")
func NewPlaceholders(s ...string) Placeholder {
	ph := Placeholder{}
	for _, v := range s {
		phe := strings.Split(v, ":")
		ph = append(ph, Entry{phe[0], phe[1]})
	}
	return ph
}

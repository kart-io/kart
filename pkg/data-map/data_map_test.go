package data_map

import (
	"fmt"
	"testing"
)

func Test_IntNewDataMapOption(t *testing.T) {
	opt := NewDataMapOption(map[int]string{
		1:  "12312",
		11: "123",
	})

	data := opt.Keys()
	fmt.Println(data)

	fmt.Println(opt.Option(1))
	fmt.Println(opt.Options())
}

func TestName(t *testing.T) {
	data := MapKeys(map[int]int{
		1: 1,
	})
	fmt.Println(data)
}

func MapKeys[Key comparable, Val any](m map[Key]Val) []Key {
	s := make([]Key, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

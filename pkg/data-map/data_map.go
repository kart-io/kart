package data_map

import "sort"

type DataType interface {
	int | int8 | int16 | int32 | int64 | string | float32 | float64
}

type Slice[T DataType] []T

type DataMap[KEY DataType, VALUE DataType] map[KEY]VALUE

type DataMapOption[KEYS DataType, VALUE DataType] struct {
	keys    Slice[KEYS]
	options DataMap[KEYS, VALUE]
}

func NewDataMapOption[KEYS, VALUE DataType](options DataMap[KEYS, VALUE]) *DataMapOption[KEYS, VALUE] {
	keys := make(Slice[KEYS], 0, len(options))
	for k := range options {
		keys = append(keys, k)
	}
	return &DataMapOption[KEYS, VALUE]{
		keys:    keys,
		options: options,
	}
}

func (o *DataMapOption[KEYS, VALUE]) Keys() []KEYS {
	sort.Slice(o.keys, func(i, j int) bool {
		return o.keys[j] > o.keys[i]
	})
	return o.keys
}

func (o *DataMapOption[KEYS, VALUE]) Option(key KEYS) VALUE {
	return o.options[key]
}

func (o *DataMapOption[KEYS, VALUE]) Options() DataMap[KEYS, VALUE] {
	return o.options
}

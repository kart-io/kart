package data_map

type DataType any

type DataMapOption struct {
	keys    []any
	options map[any]any
}

func NewDataMapOption(options map[any]any) *DataMapOption {
	keys := make([]any, 0, len(options))
	for k := range options {
		keys = append(keys, k)
	}
	return &DataMapOption{
		options: options,
		keys:    keys,
	}
}

func (o *DataMapOption) Keys() []any {
	return o.keys
}

func (o *DataMapOption) Option(key any) any {
	return o.options[key]
}

func (o *DataMapOption) Options() map[any]any {
	return o.options
}

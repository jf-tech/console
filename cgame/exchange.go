package cgame

type Exchange struct {
	BoolData    map[interface{}]bool
	IntData     map[interface{}]int
	StringData  map[interface{}]string
	GenericData map[interface{}]interface{}
}

func newExchange() *Exchange {
	return &Exchange{
		BoolData:    map[interface{}]bool{},
		IntData:     map[interface{}]int{},
		StringData:  map[interface{}]string{},
		GenericData: map[interface{}]interface{}{},
	}
}

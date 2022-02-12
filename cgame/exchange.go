package cgame

type Exchange struct {
	BData map[interface{}]bool
}

func newExchange() *Exchange {
	return &Exchange{
		BData: map[interface{}]bool{},
	}
}

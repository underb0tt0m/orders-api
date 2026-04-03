package order

type Order struct {
	Name   string
	Count  int
	Status string
}

func NewOrder(name string, count int, status string) *Order {
	return &Order{
		name,
		count,
		status,
	}
}

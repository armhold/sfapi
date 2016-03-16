package sfapi

import (
	"fmt"
)

// for POSTing
type OrderRequest struct {
	Account   string
	Venue     string
	Stock     string
	Price     int64
	Qty       int64
	Direction string
	OrderType string
}

func (or *OrderRequest) String() string {
	return fmt.Sprintf("OrderRequest: %s %d @ %s", or.Direction, or.Qty, Dollars(or.Price))
}

package rpcexample

import (
	"errors"
)

/*type Args struct {
	A, B int
}*/

type Quotient struct {
	Quo, Rem int
}

type Args struct{
	Vect []int
}
type Sum int


// Arith service for RPC
type Arith int

//Result of RPC call
type Result int

// Every method that we want to export must satisfy the conditions:
// (1) the method has two arguments, both exported (or builtin) types
// (2) the method's second argument is a pointer
// (3) the method has return type error

// Arith service has Multiply which takes numbers A, B
// as arguments and returns error or stores product in reply
/*func (t *Arith) Multiply(args Args, Result *int) error {
	*Result = args.A * args.B
	return nil
}*/

// Arith service has  Divide which takes numbers A, B
// as arguments and returns error or stores quotient and remainder in reply
/*func (t *Arith) Divide(args Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("Divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}*/

func (t *Arith) Sum(args Args, sum *int) error {

	if len(args.Vect) == 0 {
		return errors.New("Array vuoto!")
	}


	*sum = 0
	for i := 0; i < len(args.Vect); i++ {
		*sum += args.Vect[i]
	}
	return nil

}

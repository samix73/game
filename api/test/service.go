package test

type Service interface {
	Hello(args *Args, reply *Reply) error
}

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

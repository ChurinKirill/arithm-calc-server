package main

var (
	InvalidExpression = 1
	InnerExpression   = 2
	Ok                = 3
)

type CustomPanic struct {
	Type int
	Text error
}

package calculation

import (
	"fmt"
)

type iNode interface {
	calculate() (float64, error)
}

type node struct {
	left      iNode
	right     iNode
	operation rune
}

type nodeConst struct {
	constVal float64
}

type nodeBrackets struct {
	innerExp iNode
}

func (nc nodeConst) calculate() (float64, error) {
	return nc.constVal, nil
}

func (n nodeBrackets) calculate() (float64, error) {
	return n.innerExp.calculate()
}

func (n node) calculate() (float64, error) {
	l, errl := n.left.calculate()
	r, errr := n.right.calculate()
	if errl != nil {
		return 0, errl
	} else if errr != nil {
		return 0, errr
	}
	switch n.operation {
	case '+':
		return l + r, nil
	case '-':
		return l - r, nil
	case '*':
		return l * r, nil
	}
	if r == 0 {
		return 0, fmt.Errorf("Division by zero")
	}
	return l / r, nil
}

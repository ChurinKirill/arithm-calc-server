package calculation

import (
	logger "arithm-calc-server/logger"
	"fmt"
)

func createNode(tokens []iToken) (iNode, CustomPanic) {
	if len(tokens) == 0 || len(tokens) == 1 && tokens[0].getType() != numT {
		return nil, CustomPanic{
			Type: InvalidExpression,
			Text: fmt.Errorf("createNode: incorrect input format"),
		}
	} else if len(tokens) == 1 { // single number
		return nodeConst{constVal: tokens[0].getValue().(float64)}, CustomPanic{Type: Ok, Text: nil}
	} else if tokens[0].getValue() == '(' && tokens[len(tokens)-1].getValue() == ')' { // the bracket-value, e.g."(a+b*c/d)" => nodeBrackets
		inner, err := createNode(tokens[1 : len(tokens)-1])
		if err.Type != Ok {
			return nil, err
		}
		return nodeBrackets{innerExp: inner}, CustomPanic{Type: Ok, Text: nil}
	}
	i := len(tokens) - 1
	for i >= 0 { // low priority operations first (will be calculated last)
		if tokens[i].getValue() == ')' { // multi-value that has brackets (ignoring all in brackets - this will calculate first)
			i--
			var s int
			ok := false
			for i >= 0 {
				if tokens[i].getValue() == ')' {
					s++
				} else if tokens[i].getValue() == '(' && s != 0 {
					s--
				} else if tokens[i].getValue() == '(' && s == 0 {
					if i > 0 {
						i--
					}
					ok = true
					break
				}
				i--
			}
			if !ok {
				return nil, CustomPanic{
					Type: InvalidExpression,
					Text: fmt.Errorf("createNode: incorrect input format"),
				}
			} else if i == 0 {
				continue
			}
		}
		if tokens[i].getType() == opT && (tokens[i].getValue() == '+' || tokens[i].getValue() == '-') {
			leftN, errl := createNode(tokens[:i])
			rightN, errr := createNode(tokens[i+1:])
			var er CustomPanic
			er.Type = Ok
			if errl.Type != Ok {
				er = errl
			} else if errr.Type != Ok {
				er = errr
			}
			if er.Type != Ok {
				return nil, er
			}
			return node{left: leftN, right: rightN, operation: tokens[i].getValue().(rune)}, CustomPanic{
				Type: Ok,
				Text: nil,
			}
		}
		i--
	}
	// if low priority operations does not exist
	i = len(tokens) - 1
	for i >= 0 { // high priority operations (will be calculated first)
		if tokens[i].getValue() == ')' { // multi-value that has brackets (ignoring all in brackets - this will calculate firts)
			i--
			var s int
			ok := false
			for i >= 0 {
				if tokens[i].getValue() == ')' {
					s++
				} else if tokens[i].getValue() == '(' && s != 0 {
					s--
				} else if tokens[i].getValue() == '(' && s == 0 {
					i--
					ok = true
					break
				}
				i--
			}
			if !ok {
				return nil, CustomPanic{
					Type: InvalidExpression,
					Text: fmt.Errorf("createNode: incorrect input format"),
				}
			}
		}
		if tokens[i].getType() == opT && (tokens[i].getValue() == '*' || tokens[i].getValue() == '/') {
			leftN, errl := createNode(tokens[:i])
			rightN, errr := createNode(tokens[i+1:])
			//log.Printf("main_calc.go, errl: type=%d, text=%s", errl.Type, errl.Text)
			//log.Printf("main_calc.go, errr: type=%d, text=%s", errr.Type, errr.Text)
			er := CustomPanic{Type: Ok, Text: nil}
			if errl.Type != Ok {
				er = errl
			} else if errr.Type != Ok {
				er = errr
			}
			if er.Type != Ok {
				return nil, er
			}
			return node{left: leftN, right: rightN, operation: tokens[i].getValue().(rune)}, CustomPanic{Type: Ok, Text: nil}
		}
		i--
	}
	return nodeConst{constVal: 0}, CustomPanic{Type: Ok, Text: nil} // dummy return
}

func Calc(expression string) (float64, CustomPanic) {
	logger.Log(fmt.Sprintf("main_calc.go: given expression = %s", expression))
	tokens, errt := Tokenize(expression)
	if errt.Type != Ok {
		return 0, errt
	}
	rootNode, errn := createNode(tokens)
	//log.Printf("main_calc.go, errn: type=%d, text=%s", errn.Type, errn.Text)
	if errn.Type != Ok {
		return 0, errn
	}
	res, err := rootNode.calculate()
	if err != nil {
		return 0, CustomPanic{Type: InnerExpression, Text: err}
	}
	return res, CustomPanic{Type: Ok, Text: nil}
}

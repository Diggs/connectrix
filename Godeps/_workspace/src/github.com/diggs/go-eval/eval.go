package goeval

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

// EvalBool evaluates a boolean expression as a string
// It returns the result of the expression and any parsing errors encountered
func EvalBool(expression string) (bool, error) {
	result, tk, err := Eval(expression)
	if err != nil {
		return false, err
	}
	if tk != token.STRING {
		return false, fmt.Errorf("Expected expression to evaluate to type STRING but got '%v'.", tk)
	}
	return strconv.ParseBool(result)
}

// EvalArithmetic evaluates an arithmetic expression as a string
// It returns the result of the expression and any parsing errors encountered
func EvalArithmetic(expression string) (int, error) {
	result, tk, err := Eval(expression)
	if err != nil {
		return 0, err
	}
	if tk != token.INT {
		return 0, fmt.Errorf("Expected expression to evaluate to type INT but got '%v'.", tk)
	}
	i, err := strconv.ParseInt(result, 0, 0)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// Eval evaluates boolean or arithmetic expressions as a string
// It returns the result of the expression, the type of the result, and any parsing errors encountered
func Eval(expression string) (string, token.Token, error) {
	expr, err := parser.ParseExpr(expression)
	if err != nil {
		return "", token.IDENT, err
	}
	result, tk, err := evaluateExpr(expr)
	if err != nil {
		return result, tk, err
	}

	return result, tk, err
}

func evaluateExpr(expr ast.Expr) (string, token.Token, error) {

	switch exprType := expr.(type) {
	case *ast.BinaryExpr:
		op := expr.(*ast.BinaryExpr).Op

		lVal, lType, lerr := evaluateExpr(expr.(*ast.BinaryExpr).X)
		if lerr != nil {
			return lVal, lType, lerr
		}

		rVal, rType, rerr := evaluateExpr(expr.(*ast.BinaryExpr).Y)
		if rerr != nil {
			return rVal, rType, rerr
		}

		if lType != rType {
			return "", token.IDENT, fmt.Errorf("Operands must be of same type - x:%v  y:%v", lType, rType)
		}

		switch lType {
		case token.INT:
			left, _ := strconv.Atoi(lVal)
			right, _ := strconv.Atoi(rVal)
			switch op {
			case token.EQL:
				return strconv.FormatBool(left == right), token.STRING, nil
			case token.NEQ:
				return strconv.FormatBool(left != right), token.STRING, nil
			case token.GEQ:
				return strconv.FormatBool(left >= right), token.STRING, nil
			case token.LEQ:
				return strconv.FormatBool(left <= right), token.STRING, nil
			case token.GTR:
				return strconv.FormatBool(left > right), token.STRING, nil
			case token.LSS:
				return strconv.FormatBool(left < right), token.STRING, nil
			case token.ADD:
				return strconv.FormatInt(int64(left+right), 10), token.INT, nil
			case token.SUB:
				return strconv.FormatInt(int64(left-right), 10), token.INT, nil
			case token.MUL:
				return strconv.FormatInt(int64(left*right), 10), token.INT, nil
			case token.QUO:
				return strconv.FormatInt(int64(left/right), 10), token.INT, nil
			default:
				return "", token.IDENT, fmt.Errorf("Unsupported operator '%v' for int", op)
			}

		case token.STRING:
			left := lVal
			right := rVal

			switch op {
			case token.LAND:
				left, _ := strconv.ParseBool(left)
				right, _ := strconv.ParseBool(right)
				return strconv.FormatBool(left && right), token.STRING, nil
			case token.LOR:
				left, _ := strconv.ParseBool(left)
				right, _ := strconv.ParseBool(right)
				return strconv.FormatBool(left || right), token.STRING, nil
			case token.EQL:
				return strconv.FormatBool(left == right), token.STRING, nil
			case token.NEQ:
				return strconv.FormatBool(left != right), token.STRING, nil
			case token.GEQ:
				return strconv.FormatBool(left >= right), token.STRING, nil
			case token.LEQ:
				return strconv.FormatBool(left <= right), token.STRING, nil
			case token.GTR:
				return strconv.FormatBool(left > right), token.STRING, nil
			case token.LSS:
				return strconv.FormatBool(left < right), token.STRING, nil
			default:
				return "", token.IDENT, fmt.Errorf("Unsupported operator '%v' for string", op)
			}
		}

	case *ast.BasicLit:
		return expr.(*ast.BasicLit).Value, expr.(*ast.BasicLit).Kind, nil

	case *ast.UnaryExpr:
		lVal, lType, lerr := evaluateExpr(expr.(*ast.UnaryExpr).X)
		if lerr != nil {
			return lVal, lType, lerr
		}
		if lType != token.INT {
			return "", token.IDENT, errors.New("Unary operations only supported for ints")
		}
		if expr.(*ast.UnaryExpr).Op == token.SUB {
			lVal = fmt.Sprintf("-%s", lVal)
		}
		return lVal, lType, lerr

	case *ast.ParenExpr:
		return evaluateExpr(expr.(*ast.ParenExpr).X)

	default:
		return "", token.IDENT, fmt.Errorf("Unsupported expr type: %v", exprType)
	}

	return "", token.IDENT, errors.New("You shouldn't be here.")
}

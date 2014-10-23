package rules

import (
	"errors"
	"fmt"
	"github.com/diggs/glog"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func Evaluate(expression string) (string, token.Token, error) {
	expr, err := parser.ParseExpr(expression)
	if err != nil {
		return "", token.IDENT, err
	}

	// uncomment to print AST for debugging purposes
	// fset := token.NewFileSet()
	// ast.Print(fset, expr)

	result, tk, err := evaluateExpr(expr)
	if err != nil {
		return result, tk, err
	}

	glog.Debugf("Evaluated expression '%s' with result '%s'", expression, result)

	return result, tk, err
}

func evaluateExpr(expr ast.Expr) (string, token.Token, error) {

	/*
		Limitations:
			floats not supported, they will be rounded to whole numbers
	*/

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
			return "", token.IDENT, errors.New(fmt.Sprintf("Operands must be of same type - x:%v  y:%v", lType, rType))
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
				return "", token.IDENT, errors.New(fmt.Sprintf("Unsupported operator '%v' for int", op))
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
				return "", token.IDENT, errors.New(fmt.Sprintf("Unsupported operator '%v' for string", op))
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
		return "", token.IDENT, errors.New(fmt.Sprintf("Unsupported expr type: %v", exprType))
	}

	return "", token.IDENT, errors.New("You shouldn't be here.")
}

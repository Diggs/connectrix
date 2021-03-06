package goeval

import (
	"github.com/stretchr/testify/assert"
	"go/token"
	"testing"
)

func TestEvalBool(t *testing.T) {
	result, err := EvalBool("1 > 2 && 2 < 3")
	assert.Nil(t, err)	
	assert.Equal(t, false, result)

	result, err = EvalBool("1 + 2")
	assert.NotNil(t, err)	
	assert.Equal(t, false, result)
}

func TestEvalArithmetic(t *testing.T) {
	result, err := EvalArithmetic("1 + 2")
	assert.Nil(t, err)	
	assert.Equal(t, 3, result)

	result, err = EvalArithmetic("1 > 2")
	assert.NotNil(t, err)	
	assert.Equal(t, 0, result)
}

func TestParentheses(t *testing.T) {

	result, tk, err := Eval("(1 < 2 || 2 > 3) && (2 > 1)")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	result, tk, err = Eval("(1 + (2 * 10)) + 4")
	assert.Nil(t, err)
	assert.Equal(t, token.INT, tk)
	assert.Equal(t, "25", result)
}

func TestEnforcesOperandsOfSameType(t *testing.T) {

	// mix boolean expression with arithmetic
	result, tk, err := Eval("1 > 2 && 2 + 1")
	assert.NotNil(t, err)
	assert.Equal(t, token.IDENT, tk)
	assert.Equal(t, result, "")
}

func TestLogicalAnd(t *testing.T) {

	// both true
	result, tk, err := Eval("1 < 2 && 2 < 3")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	// left false, right true
	result, tk, err = Eval("1 > 2 && 2 < 3")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	// left true, right false
	result, tk, err = Eval("1 < 2 && 2 > 3")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	// both false
	result, tk, err = Eval("1 > 2 && 2 > 3")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)
}

func TestLogicalOr(t *testing.T) {

	// both true
	result, tk, err := Eval("1 < 2 || 2 < 3")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	// left false, right true
	result, tk, err = Eval("1 > 2 || 2 < 3")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	// left true, right false
	result, tk, err = Eval("1 < 2 || 2 > 3")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	// both false
	result, tk, err = Eval("1 > 2 || 2 > 3")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)
}

func TestEQ(t *testing.T) {

	result, tk, err := Eval("2 == 1")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	result, tk, err = Eval("2 == 2")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	result, tk, err = Eval(`"FOO" == "BAR"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	result, tk, err = Eval(`"FOO" == "FOO"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)
}

func TestNEQ(t *testing.T) {

	result, tk, err := Eval("2 != 1")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	result, tk, err = Eval("2 != 2")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	result, tk, err = Eval(`"FOO" != "BAR"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	result, tk, err = Eval(`"FOO" != "FOO"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)
}

func TestGEQ(t *testing.T) {

	result, tk, err := Eval("2 >= 3")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	result, tk, err = Eval("2 >= 2")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	result, tk, err = Eval("2 >= 1")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	result, tk, err = Eval(`"FOO" >= "BAR"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	result, tk, err = Eval(`"FOO" >= "FOO"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	result, tk, err = Eval(`"FOO" >= "GOO"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)
}

func TestLEQ(t *testing.T) {

	result, tk, err := Eval("2 <= 1")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	result, tk, err = Eval("2 <= 2")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	result, tk, err = Eval("2 <= 3")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	result, tk, err = Eval(`"FOO" <= "BAR"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	result, tk, err = Eval(`"FOO" <= "FOO"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	result, tk, err = Eval(`"FOO" <= "GOO"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)
}

func TestGTR(t *testing.T) {

	result, tk, err := Eval("2 > 1")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	result, tk, err = Eval("2 > 2")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	result, tk, err = Eval("2 > 3")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	result, tk, err = Eval(`"FOO" > "BAR"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	result, tk, err = Eval(`"FOO" > "FOO"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	result, tk, err = Eval(`"FOO" > "GOO"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)
}

func TestLSS(t *testing.T) {

	result, tk, err := Eval("2 < 1")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	result, tk, err = Eval("2 < 2")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	result, tk, err = Eval("2 < 3")
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)

	result, tk, err = Eval(`"FOO" < "BAR"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	result, tk, err = Eval(`"FOO" < "FOO"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "false", result)

	result, tk, err = Eval(`"FOO" < "GOO"`)
	assert.Nil(t, err)
	assert.Equal(t, token.STRING, tk)
	assert.Equal(t, "true", result)
}

func TestADD(t *testing.T) {

	result, tk, err := Eval("2 + 1")
	assert.Nil(t, err)
	assert.Equal(t, token.INT, tk)
	assert.Equal(t, "3", result)

	result, tk, err = Eval("2 + -1")
	assert.Nil(t, err)
	assert.Equal(t, token.INT, tk)
	assert.Equal(t, "1", result)

	result, tk, err = Eval(`"FOO" + "BAR"`)
	assert.NotNil(t, err)
	assert.Equal(t, token.IDENT, tk)
	assert.Equal(t, "", result)
}

func TestSUB(t *testing.T) {

	result, tk, err := Eval("2 - 1")
	assert.Nil(t, err)
	assert.Equal(t, token.INT, tk)
	assert.Equal(t, "1", result)

	result, tk, err = Eval("2 - -1")
	assert.Nil(t, err)
	assert.Equal(t, token.INT, tk)
	assert.Equal(t, "3", result)

	result, tk, err = Eval(`"FOO" - "BAR"`)
	assert.NotNil(t, err)
	assert.Equal(t, token.IDENT, tk)
	assert.Equal(t, "", result)
}

func TestMUL(t *testing.T) {

	result, tk, err := Eval("2 * 2")
	assert.Nil(t, err)
	assert.Equal(t, token.INT, tk)
	assert.Equal(t, "4", result)

	result, tk, err = Eval("2 * -2")
	assert.Nil(t, err)
	assert.Equal(t, token.INT, tk)
	assert.Equal(t, "-4", result)

	result, tk, err = Eval(`"FOO" * "BAR"`)
	assert.NotNil(t, err)
	assert.Equal(t, token.IDENT, tk)
	assert.Equal(t, "", result)
}

func TestQUO(t *testing.T) {

	result, tk, err := Eval("4 / 2")
	assert.Nil(t, err)
	assert.Equal(t, token.INT, tk)
	assert.Equal(t, "2", result)

	result, tk, err = Eval(`"FOO" / "BAR"`)
	assert.NotNil(t, err)
	assert.Equal(t, token.IDENT, tk)
	assert.Equal(t, "", result)
}

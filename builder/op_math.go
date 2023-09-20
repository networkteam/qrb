package builder

// Math operators

const (
	opPlus   Operator = "+"
	opMinus  Operator = "-"
	opDivide Operator = "/"
	opMult   Operator = "*"
	opMod    Operator = "%"
	opPow    Operator = "^"
)

// Plus builds the + operator (addition) for numeric types.
func (b ExpBase) Plus(rgt Exp) ExpBase {
	return b.Op(opPlus, rgt)
}

// Minus builds the - operator (subtraction) for numeric types.
func (b ExpBase) Minus(rgt Exp) ExpBase {
	return b.Op(opMinus, rgt)
}

// Neg builds the - unary operator (negation) for numeric types.
func Neg(exp Exp) ExpBase {
	return ExpBase{
		Exp: unaryExp{
			prefix:     "-",
			exp:        exp,
			precedence: 4, // see opPrecedence for unary minus
		},
	}
}

// Mult builds the * operator (multiplication) for numeric types.
func (b ExpBase) Mult(rgt Exp) ExpBase {
	return b.Op(opMult, rgt)
}

// Divide builds the / operator (division) for numeric types.
// Do not confuse with the div function that returns the integer part of the division.
func (b ExpBase) Divide(rgt Exp) ExpBase {
	return b.Op(opDivide, rgt)
}

// Mod builds the % operator (remainder) for numeric types.
func (b ExpBase) Mod(rgt Exp) ExpBase {
	return b.Op(opMod, rgt)
}

// Pow builds the ^ operator (exponentiation) for numeric types.
func (b ExpBase) Pow(rgt Exp) ExpBase {
	return b.Op(opPow, rgt)
}

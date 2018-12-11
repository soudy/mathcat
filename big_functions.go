package mathcat

import "math/big"

// RationalToInteger converts a rational number to an integer
func RationalToInteger(n *big.Rat) *big.Int {
	return new(big.Int).Div(n.Num(), n.Denom())
}

// Factorial calculates the factorial of rational number n
func Factorial(n *big.Rat) *big.Rat {
	integer := RationalToInteger(n)
	fact := new(big.Int).MulRange(1, integer.Int64())
	return new(big.Rat).SetInt(fact)
}

// Gcd calculates the greatest common divisor of the numbers x and y
func Gcd(x, y *big.Rat) *big.Rat {
	xInt := RationalToInteger(x)
	yInt := RationalToInteger(y)
	gcd := new(big.Int).GCD(nil, nil, xInt, yInt)

	return new(big.Rat).SetInt(gcd)
}

// Max gives the maximum of two rational numbers
func Max(a, b *big.Rat) *big.Rat {
	if a.Cmp(b) == 1 {
		return a
	}

	return b
}

// Min gives the minimum of two rational numbers
func Min(a, b *big.Rat) *big.Rat {
	if a.Cmp(b) == -1 {
		return a
	}

	return b
}

// Floor returns the floor of a rational number
func Floor(n *big.Rat) *big.Rat {
	return new(big.Rat).SetInt(RationalToInteger(n))
}

// Ceil returns the ceil of a rational number
func Ceil(n *big.Rat) *big.Rat {
	floor := Floor(n.Neg(n))
	return new(big.Rat).Neg(floor)
}

// Mod returns x % y
func Mod(x, y *big.Rat) *big.Rat {
	res := new(big.Rat)
	quo := Floor(res.Quo(x, y))
	return res.Sub(x, res.Mul(y, quo))
}

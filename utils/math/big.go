package math

import "math/big"

func Mul(x, y *big.Int) *big.Int {
	return big.NewInt(0).Mul(x, y)
}

func Div(x, y *big.Int) *big.Int {
	return big.NewInt(0).Div(x, y)
}

func Add(x, y *big.Int) *big.Int {
	return big.NewInt(0).Add(x, y)
}

func Sub(x, y *big.Int) *big.Int {
	return big.NewInt(0).Sub(x, y)
}

func Neg(x *big.Int) *big.Int {
	return big.NewInt(0).Neg(x)
}

func ToBigInt(s string) *big.Int {
	res := big.NewInt(0)
	res.SetString(s, 10)
	return res
}

func Max(a, b *big.Int) *big.Int {
	if a.Cmp(b) == 1 {
		return a
	} else {
		return b
	}
}

func IsZero(x *big.Int) bool {
	if x.Cmp(big.NewInt(0)) == 0 {
		return true
	} else {
		return false
	}
}

func IsEqual(x, y *big.Int) bool {
	if x.Cmp(y) == 0 {
		return true
	} else {
		return false
	}
}

func IsGreaterThan(x, y *big.Int) bool {
	if x.Cmp(y) == 1 {
		return true
	} else {
		return false
	}
}

func IsSmallerThan(x, y *big.Int) bool {
	if x.Cmp(y) == -1 {
		return true
	} else {
		return false
	}
}

func IsEqualOrGreaterThan(x, y *big.Int) bool {
	return (IsEqual(x, y) || IsGreaterThan(x, y))
}

func IsEqualOrSmallerThan(x, y *big.Int) bool {
	return (IsEqual(x, y) || IsSmallerThan(x, y))
}

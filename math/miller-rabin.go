package math

import "math/big"

func MillerRabinIsPrime(p *big.Int, k int) bool {
	return p.ProbablyPrime(k)
}
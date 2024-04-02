package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"pgen/math"
	"runtime"
	"sync"
	"time"
)

var bits int
var thrs int
var chks int
var debug bool

var wg sync.WaitGroup

func init() {
	flag.IntVar(&bits, "bits", 2048, "Set the n of bits of the generated prime")
	flag.IntVar(&thrs, "thrs", runtime.NumCPU(),
		"Set how many threads beign used to generate the prime number",
	)
	flag.IntVar(&chks, "checks", 10,
		"Set the n of checks that are used to confirm if is a prime number",
	)
	flag.BoolVar(&debug, "debug", false, "Use debug mode")
	
	flag.Parse()
}

func main() {
	var fp *big.Int
	ctx, cancel := context.WithCancel(context.Background())

	start := time.Now()
	for t := 0; t < thrs; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

		primeLoop:
			for {
				select {
				case <-ctx.Done():
					break primeLoop
				default:
					p := big.NewInt(0)

					if bits == 2 {
						p.SetInt64(rand.Int63n(2) + 2)
						fp = p
						break primeLoop
					}

					p.SetBit(p, 0, 1)
					p.SetBit(p, bits-1, 1)
					for i := 1; i < bits-1; i++ {
						p.SetBit(p, i, uint(rand.Intn(2)))
					}

				fCheck:
					for {
						for _, v := range math.HC_PRIMES {
							pb5 := big.NewInt(v)
							mod := big.NewInt(0).Mod(p, pb5)
							if mod.Int64() == 0 && p.Cmp(pb5) == 1 {
								p.Add(p, big.NewInt(2))
								continue fCheck
							}
						}
						break fCheck
					}

					if !math.MillerRabinIsPrime(p, chks) {
						continue primeLoop
					}

					fp = p
					cancel()
					break primeLoop
				}
			}
		}()
	}
	wg.Wait()
	end := time.Now()
	ellapsed := end.Sub(start)

	fmt.Println(fp.String())
	if debug {
		fmt.Printf("fp.BitLen(): %v\n", fp.BitLen())
		fmt.Printf("ellapsed.String(): %v\n", ellapsed.String())
	}
}

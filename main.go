package main

import (
	"fmt"
	"log"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"

	"github.com/rs/zerolog"
)

// Circuit: Prove that Min ≤ Age ≤ Max
type Circuit struct {
	// Private input: the user's age
	Age frontend.Variable `gnark:"age"`

	// Public inputs: range bounds
	Min frontend.Variable `gnark:",public"`
	Max frontend.Variable `gnark:",public"`
}

// rangeNonNeg constrains v >= 0 by forcing v to be representable
// as a small non-negative integer using 'bits' bits.
func rangeNonNeg(api frontend.API, v frontend.Variable, bits int) {
	bin := api.ToBinary(v, bits) // constrain 0 ≤ v < 2^bits
	for _, b := range bin {
		api.AssertIsBoolean(b)
	}
	// Reconstruct v from bits and assert equality
	reconstructed := frontend.Variable(0)
	for i, b := range bin {
		reconstructed = api.Add(reconstructed, api.Mul(b, 1<<i))
	}
	api.AssertIsEqual(v, reconstructed)
}

// Define: enforce Min ≤ Age ≤ Max
func (c *Circuit) Define(api frontend.API) error {
	const bits = 16 // plenty for realistic ages

	lower := api.Sub(c.Age, c.Min) // Age - Min ≥ 0  ⇒ Age ≥ Min
	rangeNonNeg(api, lower, bits)

	upper := api.Sub(c.Max, c.Age) // Max - Age ≥ 0  ⇒ Age ≤ Max
	rangeNonNeg(api, upper, bits)

	return nil
}

func main() {
	// Disable gnark debug logs
	zerolog.SetGlobalLevel(zerolog.Disabled)

	// -----------------------------
	// Ask user for inputs
	// -----------------------------
	var age, min, max int
	fmt.Print("Enter Age (private): ")
	_, err := fmt.Scan(&age)
	if err != nil {
		log.Fatalf("failed to read Age: %v", err)
	}

	fmt.Print("Enter Min bound (public): ")
	_, err = fmt.Scan(&min)
	if err != nil {
		log.Fatalf("failed to read Min: %v", err)
	}

	fmt.Print("Enter Max bound (public): ")
	_, err = fmt.Scan(&max)
	if err != nil {
		log.Fatalf("failed to read Max: %v", err)
	}

	// -----------------------------
	// 1) Compile circuit
	// -----------------------------
	var circuit Circuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		log.Fatalf("compile error: %v", err)
	}

	// -----------------------------
	// 2) Trusted setup (Groth16)
	// -----------------------------
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		log.Fatalf("setup error: %v", err)
	}

	// -----------------------------
	// 3) Assign inputs (witness)
	// -----------------------------
	assignment := Circuit{
		Age: age, // private
		Min: min, // public
		Max: max, // public
	}

	fmt.Println("\n=== Inputs ===")
	fmt.Printf("Private:  Age = %v\n", age)
	fmt.Printf("Public:   Min = %v\n", min)
	fmt.Printf("Public:   Max = %v\n", max)
	fmt.Println("Proving statement: Min ≤ Age ≤ Max ?")

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		log.Fatalf("witness error: %v", err)
	}
	publicWitness, err := witness.Public()
	if err != nil {
		log.Fatalf("public witness error: %v", err)
	}

	// -----------------------------
	// 4) Prove
	// -----------------------------
	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		fmt.Println("Prove: ❌ FAILED (witness does not satisfy constraints)")
		log.Fatalf("Reason: %v\n", err)
	}

	// -----------------------------
	// 5) Verify
	// -----------------------------
	if err := groth16.Verify(proof, vk, publicWitness); err != nil {
		fmt.Println("Verification: ❌ FAILED")
		fmt.Printf("Reason: %v\n", err)
		return
	}
	fmt.Println("Verification: ✅ SUCCESS (Min ≤ Age ≤ Max proven zero-knowledge)")
}
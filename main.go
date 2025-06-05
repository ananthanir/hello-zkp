package main

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// Circuit defines a simple circuit: prove that Age > 10
type Circuit struct {
	Age frontend.Variable `gnark:"age"`
}

// Define declares the circuit constraints
func (circuit *Circuit) Define(api frontend.API) error {
	adjustedAge := api.Sub(circuit.Age, 11)

	// Decompose adjustedAge into 8 bits
	bits := api.ToBinary(adjustedAge, 8)
	for _, bit := range bits {
		api.AssertIsBoolean(bit)
	}
	// Reconstruct adjustedAge from bits
	var reconstructed frontend.Variable = 0
	for i, bit := range bits {
		reconstructed = api.Add(reconstructed, api.Mul(bit, 1<<i))
	}
	api.AssertIsEqual(adjustedAge, reconstructed)

	return nil
}

func main() {
	// Compile circuit into R1CS
	var circuit Circuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// groth16 Setup
	pk, vk, _ := groth16.Setup(ccs)

	// Witness assignment (Prover's secret)
	assignment := Circuit{Age: 25}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	// groth16: Prove & Verify
	proof, _ := groth16.Prove(ccs, pk, witness)
	groth16.Verify(proof, vk, publicWitness)
}

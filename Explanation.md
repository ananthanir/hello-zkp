# Zero-Knowledge Proof Demo with gnark (Go)

This example shows how to use **gnark** with the **Groth16** zk-SNARK algorithm to prove a simple statement:  
👉 *"My age is between Min and Max, without revealing the exact age."*

---

## Code Overview

We’ll go through the Go code step by step.

---

### 1. Define the Circuit

```go
// Circuit: Prove that Min ≤ Age ≤ Max
type Circuit struct {
    Age frontend.Variable `gnark:"age"`        // Private input: secret age
    Min frontend.Variable `gnark:",public"`    // Public input: lower bound
    Max frontend.Variable `gnark:",public"`    // Public input: upper bound
}

// Constraints: Age must be between Min and Max
func (c *Circuit) Define(api frontend.API) error {
    const bits = 16
    lower := api.Sub(c.Age, c.Min) // Age - Min ≥ 0
    rangeNonNeg(api, lower, bits)

    upper := api.Sub(c.Max, c.Age) // Max - Age ≥ 0
    rangeNonNeg(api, upper, bits)
    return nil
}
```

- The circuit enforces two constraints: `Age ≥ Min` and `Age ≤ Max`.
- Range checks (`rangeNonNeg`) ensure no field wrap-around.

---

### 2. Compile the Circuit

```go
var circuit Circuit
ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
```

- Compiles the circuit into **R1CS** form (Rank-1 Constraint System).
- This is the math-friendly representation used by zk-SNARKs.

---

### 3. Trusted Setup (Groth16)

```go
pk, vk, _ := groth16.Setup(ccs)
```

- **Groth16 algorithm** generates:  
  - **Proving Key (pk)** → used by Prover to create proofs  
  - **Verification Key (vk)** → used by Verifier to check proofs  
- Requires randomness. Leftover randomness is called **toxic waste** (must be destroyed).

---

### 4. Inputs (Witness)

```go
assignment := Circuit{
    Age: 25,  // private
    Min: 18,  // public
    Max: 30,  // public
}

witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
publicWitness, _ := witness.Public()
```

- **Private input:** `Age = 25`  
- **Public inputs:** `Min = 18`, `Max = 30`  
- Full witness contains both; verifier only sees the **publicWitness**.

---

### 5. Prove

```go
proof, _ := groth16.Prove(ccs, pk, witness)
```

- **Groth16.Prove** generates a zk-SNARK proof.
- Prover uses: compiled circuit, proving key, full witness.

---

### 6. Verify

```go
groth16.Verify(proof, vk, publicWitness)
```

- **Groth16.Verify** checks the proof.
- Verifier uses: proof, verification key, public inputs.  
- ✅ If valid → convinced age is in range.  
- ❌ If invalid → reject.

---

## Workflow Summary

1. **Circuit Creation** → define rules  
2. **Compile** → circuit → R1CS  
3. **Setup (Groth16)** → generate pk, vk  
4. **Inputs** → witness (private + public)  
5. **Prove** → prover generates proof  
6. **Verify** → verifier checks proof with vk + public inputs  

---
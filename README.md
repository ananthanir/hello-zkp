# Proving Age Over 10 with Zero Knowledge Proofs (Groth16 + gnark)

This project is a minimal and clean example showing how Zero Knowledge Proofs can be used to **prove statements** like "**I am over 10 years old**" without **revealing your actual age**.

## 🛠️ We use
 * **Groth16** zkSNARK proving system
 * **[gnark](https://github.com/consensys/gnark)** — a modern zk framework in Go
 * **BN254** elliptic curve (BLS12-381 and others are also supported)
 * **[Go](https://go.dev/doc/install)** as the programming language

## 🎯 What This Code Does
 The code defines a cryptographic circuit that:
 * Takes an input **Age** (private / secret input).
 * Proves that **Age > 10** without leaking Age.
 * Uses binary decomposition to enforce non-negativity (Age - 11 >= 0).
Uses the efficient Groth16 proving system to generate and verify the proof.

## ✅ Proof guarantees
 * **Completeness**: If the prover knows an Age > 10, they can generate a valid proof.
 * **Soundness**: If the prover does not have Age > 10, they cannot create a valid proof.
 * **Zero Knowledge**: The verifier learns nothing about the actual Age.

## 📦 How to Run
Clone the repository.
```
git clone https://github.com/ananthanir/hello-zkp.git
```
Run the code.
```
cd hello-zkp
go run .
```

## 🛡️ About Groth16 and Trusted Setup
 * Groth16 is a popular zkSNARK proving system with extremely small proof sizes (~200 bytes) and fast verification.
 * Trusted Setup: Groth16 requires a one-time setup for each circuit to generate proving and verification keys. If the setup is compromised, the security guarantees are broken.
 * In practice, trusted setups are generated through multi-party ceremonies to ensure security.
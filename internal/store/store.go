package store

// This is expected to be a well known name.
// Repo doesn't have to be this name, but if it is we can automatically find it.
const StandardCrevProofRepoName = "crev-proofs"

type ProofStore struct {
	Dir string
}

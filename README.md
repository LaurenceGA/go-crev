<p align="center">
    <img src="https://github.com/LaurenceGA/go-crev/workflows/Master%20checks/badge.svg" />
</p>

# go-crev
A cryptographically verifiable code review system for go packages.


[Crev](https://github.com/crev-dev/crev/) is a language and ecosystem agnostic, distributed code review system.

go-crev is an implementation of Crev as a command line tool integrated with go dependency management. This tool helps gophers evaluate the quality and trustworthiness of their package dependencies.

## Concepts

### Proof
A cryptographically signed document. Either a `Review` or `Trust`.

#### Review
Describes the results of a code review.

#### Trust
Expresses trust for another `CrevID`.

### Proof repository
A collection of `proofs`.
Typically a git repository.

### CrevID
Every user/entity has their own identity.
This is tied to every proofs they create.

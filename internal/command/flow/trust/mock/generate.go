package mock

//go:generate mockgen -package=mock -destination=mock.go github.com/LaurenceGA/go-crev/internal/command/flow/trust ConfigReader,Github,Prompter,KeyLoader,StoreWriter
//go:generate mockgen -package=mock -destination=mock_signer.go golang.org/x/crypto/ssh Signer

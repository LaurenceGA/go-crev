package mock

//go:generate mockgen -package=mock -destination=mock.go github.com/LaurenceGA/go-crev/internal/command/flow/trust ConfigReader,Github

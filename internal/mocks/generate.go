package mocks

//go:generate mockgen -package=mocks -destination=mock_store.go github.com/LaurenceGA/go-crev/internal/store GitCloner,FileDirs

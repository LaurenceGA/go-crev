package mocks

//go:generate mockgen -package=mocks -destination=mock_store.go github.com/LaurenceGA/go-crev/internal/store GitCloner,FileDirs
//go:generate mockgen -package=mocks -destination=mock_mod.go github.com/LaurenceGA/go-crev/mod ModulesWrapper
//go:generate mockgen -package=mocks -destination=mock_id_flow.go github.com/LaurenceGA/go-crev/internal/command/flow ConfigManipulator,Github

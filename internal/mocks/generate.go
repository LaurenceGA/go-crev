package mocks

//go:generate mockgen -package=mocks -destination=mock_fetcher.go github.com/LaurenceGA/go-crev/internal/store/fetcher GitCloner
//go:generate mockgen -package=mocks -destination=mock_mod.go github.com/LaurenceGA/go-crev/mod ModulesWrapper
//go:generate mockgen -package=mocks -destination=mock_id_flow.go github.com/LaurenceGA/go-crev/internal/command/flow ConfigManipulator,Github,RepoFetcher
//go:generate mockgen -package=mocks -destination=mock_files.go github.com/LaurenceGA/go-crev/internal/files AppDirs

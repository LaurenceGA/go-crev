package store

func New() *Fetcher {
	return &Fetcher{}
}

type Fetcher struct {
}

func (f *Fetcher) Fetch(url string) error {
	return nil
}
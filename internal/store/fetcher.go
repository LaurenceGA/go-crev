package store

func New() *Fetcher {
	return &Fetcher{}
}

type Fetcher struct {
}

// Fetch will download a store from a URL to the cache.
func (f *Fetcher) Fetch(fetchURL string) error {
	return nil
}

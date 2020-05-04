package id

type CrevIdentity string

const (
	CrevID CrevIdentity = "CrevID"
	PGP    CrevIdentity = "pgp"
	Github CrevIdentity = "Github"
)

type ID struct {
	ID    string       `yaml:"id,omitempty"`
	Type  CrevIdentity `yaml:"id-type,omitempty"`
	URL   string       `yaml:"url,omitempty"`
	Alias string       `yaml:"alias:omitempty"`
}

package id

import "net/url"

type CrevIdentity string

const (
	CrevID CrevIdentity = "CrevID"
	PGP    CrevIdentity = "pgp"
	Github CrevIdentity = "Github"
)

type ID struct {
	ID   string       `json:"id"`
	Type CrevIdentity `json:"id-type"`
	URL  *url.URL     `json:"url"`
}

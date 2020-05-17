package review

import "context"

func NewCreator() *Creator {
	return &Creator{}
}

type Creator struct {
}

type CreatorOptions struct {
	IdentityFile string
}

func (c *Creator) CreateReview(ctx context.Context, packageName string, options CreatorOptions) error {
	return nil
}

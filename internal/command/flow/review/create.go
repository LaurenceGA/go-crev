package review

func NewCreator() *Creator {
	return &Creator{}
}

type Creator struct {
}

func (c *Creator) CreateReview() error {
	return nil
}

package mixin

type MixinService interface {
	Transfer() error
	SendMessage() error
	NFTS() error
}

type Mixin struct {
}

func (m *Mixin) Transfer() error {
	return nil
}

func (m *Mixin) SendMessage() error {
	return nil
}

func (m *Mixin) NFTS() error {
	return nil
}

func NewMixinService() MixinService {
	return &Mixin{}
}

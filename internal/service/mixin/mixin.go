package mixin

import (
	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/lixvyang/betxin.one/config"
)

type MixinService interface {
	Transfer() error
	SendMessage() error
}

type MixinCli struct {
	client *mixin.Client
}

func New(conf config.MixinKeyStore) *MixinCli {
	client, err := mixin.NewFromKeystore(&mixin.Keystore{
		PinToken:          conf.PinToken,
		Scope:             conf.Scope,
		SessionID:         conf.SessionID,
		ServerPublicKey:   conf.ServerPublicKey,
		ClientID:          conf.ClientID,
		PrivateKey:        conf.PrivateKey,
		AppID:             conf.AppID,
		SessionPrivateKey: conf.SessionPrivateKey,
	})
	if err != nil {
		panic(err)
	}
	return &MixinCli{
		client: client,
	}
}

func (m *MixinCli) Transfer() error {
	// m.client.sa

	return nil
}

func (m *MixinCli) SendMessage() error {
	return nil
}

func (m *MixinCli) NFTS() error {
	return nil
}

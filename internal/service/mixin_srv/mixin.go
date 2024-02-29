package mixin_srv

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/fox-one/mixin-sdk-go/v2/mixinnet"
	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

const (
	ASSET_CNB = "965e5c6e-434c-3fa9-b780-c50f43cd955c"

	MAX_UTXO_NUM = 255
)

type MixinService interface {
	Transfer() error
	SendMessage() error
}

type MixinCli struct {
	client *mixin.Client
	db     database.Database

	spendKey mixinnet.Key

	conf *config.MixinConfig
}

func New(conf *config.MixinConfig, db database.Database) *MixinCli {
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

	user, err := client.UserMe(context.Background())
	if err != nil {
		panic(err)
	}

	if conf.SpendKey == "" {
		conf.SpendKey = os.Getenv("SPEND_KEY")
	}

	fmt.Println(conf.SpendKey)

	spendKey, err := mixinnet.ParseKeyWithPub(conf.SpendKey, user.SpendPublicKey)
	if err != nil {
		panic(err)
	}

	return &MixinCli{
		client,
		db,
		spendKey,
		conf,
	}
}

type TransferRequest struct {
	RequestId    string
	MemberAmount []struct {
		Member []string
		Amount decimal.Decimal
	}
	Memo string
}

func (m *MixinCli) Transfer(ctx context.Context, req *TransferRequest) (*mixin.SafeTransactionRequest, error) {
	totalAmount := decimal.NewFromInt(0)
	lo.ForEach(req.MemberAmount, func(item struct {
		Member []string
		Amount decimal.Decimal
	}, index int) {
		totalAmount = totalAmount.Add(item.Amount)
	})

	// 1. 将utxos聚合
	utxos, _ := m.client.SafeListUtxos(ctx, mixin.SafeListUtxoOption{
		Asset:     ASSET_CNB,
		State:     mixin.SafeUtxoStateUnspent,
		Threshold: 1,
		Limit:     500,
	})

	// 1: select utxos
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Amount.LessThanOrEqual(utxos[j].Amount)
	})

	var useAmount decimal.Decimal
	var useUtxos []*mixin.SafeUtxo
	for _, utxo := range utxos {
		if useAmount.GreaterThanOrEqual(totalAmount) {
			break
		}
		if len(useUtxos) > MAX_UTXO_NUM {
			useUtxos = useUtxos[1:]
		}
		useUtxos = append(useUtxos, utxo)
	}

	if useAmount.LessThan(totalAmount) {
		return nil, fmt.Errorf("not enough utxos")
	}

	// 2: build transaction
	b := mixin.NewSafeTransactionBuilder(useUtxos)
	b.Memo = req.Memo

	txOutout := make([]*mixin.TransactionOutput, len(req.MemberAmount))
	for i := 0; i < len(req.MemberAmount); i++ {
		txOutout[i] = &mixin.TransactionOutput{
			Address: mixin.RequireNewMixAddress(req.MemberAmount[i].Member, byte(len(req.MemberAmount[i].Member))),
		}
	}

	tx, err := m.client.MakeTransaction(ctx, b, txOutout)
	if err != nil {
		return nil, err
	}

	raw, err := tx.Dump()
	if err != nil {
		return nil, err
	}

	// 3. create transaction
	request, err := m.client.SafeCreateTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID:      req.RequestId,
		RawTransaction: raw,
	})
	if err != nil {
		return nil, err
	}
	// 4. sign transaction
	err = mixin.SafeSignTransaction(
		tx,
		m.spendKey,
		request.Views,
		0,
	)
	if err != nil {
		return nil, err
	}
	signedRaw, err := tx.Dump()
	if err != nil {
		return nil, err
	}

	// 5. submit transaction
	_, err = m.client.SafeSubmitTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID:      req.RequestId,
		RawTransaction: signedRaw,
	})
	if err != nil {
		return nil, err
	}

	// 6. read transaction
	req1, err := m.client.SafeReadTransactionRequest(ctx, req.RequestId)
	if err != nil {
		return nil, err
	}
	return req1, nil
}

// 定期聚合utxos
func (m *MixinCli) ArrgegateUtxos(ctx context.Context, requestId string) error {
	var req *mixin.SafeTransactionRequest

	for {
		utxos, _ := m.client.SafeListUtxos(ctx, mixin.SafeListUtxoOption{
			Asset:     ASSET_CNB,
			State:     mixin.SafeUtxoStateUnspent,
			Threshold: 1,
		})

		if len(utxos) <= MAX_UTXO_NUM {
			break
		}

		// 将utxos分割成 255 个大小的数组
		utxoSlice := make([]*mixin.SafeUtxo, 0)
		utxoSliceAmount := decimal.NewFromInt(0)
		for i := 0; i <= MAX_UTXO_NUM; i++ {
			utxoSlice = append(utxoSlice, utxos[i])
			utxoSliceAmount = utxoSliceAmount.Add(utxos[i].Amount)
		}

		// 2: build transaction
		b := mixin.NewSafeTransactionBuilder(utxoSlice)
		b.Memo = "arrgegate utxos"

		tx, err := m.client.MakeTransaction(ctx, b, []*mixin.TransactionOutput{
			{
				Address: mixin.RequireNewMixAddress([]string{m.client.ClientID}, 1),
				Amount:  utxoSliceAmount,
			},
		})
		if err != nil {
			return err
		}

		raw, err := tx.Dump()
		if err != nil {
			return err
		}
		// 3. create transaction
		request, err := m.client.SafeCreateTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
			RequestID:      requestId,
			RawTransaction: raw,
		})
		if err != nil {
			return err
		}

		// 4. sign transaction
		err = mixin.SafeSignTransaction(
			tx,
			m.spendKey,
			request.Views,
			0,
		)
		if err != nil {
			return err
		}

		signedRaw, err := tx.Dump()
		if err != nil {
			return err
		}

		// 5. submit transaction
		_, err = m.client.SafeSubmitTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
			RequestID:      requestId,
			RawTransaction: signedRaw,
		})
		if err != nil {
			return err
		}

		// 6. read transaction
		req, err = m.client.SafeReadTransactionRequest(ctx, requestId)
		if err != nil {
			return err
		}

		// 写入聚合请求
		if req != nil {
			// m.db.CreateBonuse()
		}

	}
	return nil
}

func (m *MixinCli) SendMessage() error {
	return nil
}

func (m *MixinCli) NFTS() error {
	return nil
}

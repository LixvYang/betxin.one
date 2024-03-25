package mixin_srv

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/fox-one/mixin-sdk-go/v2/mixinnet"
	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils/convert"
	"github.com/rs/zerolog/log"

	"github.com/shopspring/decimal"
)

type MixinSrv struct {
	User     *mixin.User
	Client   *mixin.Client
	SpendKey mixinnet.Key

	transferMutex *sync.Mutex
	appConf       *config.AppConfig
}

func New(appConf *config.AppConfig) *MixinSrv {
	conf := appConf.MixinConfig
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

	spendKey, err := mixinnet.ParseKeyWithPub(conf.SpendKey, user.SpendPublicKey)
	if err != nil {
		panic(err)
	}

	mixinSrv := &MixinSrv{
		user,
		client,
		spendKey,
		&sync.Mutex{},
		appConf,
	}

	return mixinSrv
}

func (m *MixinSrv) NFTS() error {
	return nil
}

// 主动聚合utxos 至 utxo 数量不超过 255 个
func (m *MixinSrv) SyncArrgegateUtxos(ctx context.Context) error {
	m.transferMutex.Lock()
	defer m.transferMutex.Unlock()

	for {
		requestId := convert.NewUUID()
		utxos, err := m.Client.SafeListUtxos(ctx, mixin.SafeListUtxoOption{
			Asset:     consts.ASSET_CNB,
			State:     mixin.SafeUtxoStateUnspent,
			Threshold: 1,
		})

		if err != nil {
			log.Error().Err(err).Msg("list utxos failed")
			continue
		}

		if len(utxos) <= consts.MAX_UTXO_NUM {
			// 主动聚合完成
			break
		}

		// 将utxos分割成 255 个大小的数组
		utxoSlice := make([]*mixin.SafeUtxo, 0)
		utxoSliceAmount := decimal.NewFromInt(0)
		for i := 0; i < len(utxos); i++ {
			utxoSlice = append(utxoSlice, utxos[i])
			utxoSliceAmount = utxoSliceAmount.Add(utxos[i].Amount)
		}

		// 2: build transaction
		b := mixin.NewSafeTransactionBuilder(utxoSlice)
		b.Memo = "arrgegate utxos"

		tx, err := m.Client.MakeTransaction(ctx, b, []*mixin.TransactionOutput{
			{
				Address: mixin.RequireNewMixAddress([]string{m.Client.ClientID}, 1),
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
		request, err := m.Client.SafeCreateTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
			RequestID:      requestId,
			RawTransaction: raw,
		})
		if err != nil {
			return err
		}

		// 4. sign transaction
		err = mixin.SafeSignTransaction(
			tx,
			m.SpendKey,
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
		_, err = m.Client.SafeSubmitTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
			RequestID:      requestId,
			RawTransaction: signedRaw,
		})
		if err != nil {
			return err
		}

		// 重试读取交易状态
		const defaultMaxRetryTimes = 3
		retryTimes := 0
		for {
			if retryTimes >= defaultMaxRetryTimes {
				break
			}

			retryTimes++
			time.Sleep(time.Second * time.Duration(retryTimes))
			_, err := m.Client.SafeReadTransactionRequest(ctx, requestId)
			if err != nil {
				continue
			} else {
				break
			}
		}

		// 等待 250ms
		time.Sleep(time.Second >> 2)
	}
	return nil
}

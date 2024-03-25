package mixin_srv

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/lixvyang/betxin.one/internal/consts"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type TransferOneRequest struct {
	RequestId string
	Member    string
	Amount    decimal.Decimal
	Memo      string
}

type MemberAmount struct {
	Member []string
	Amount decimal.Decimal
}

type TransferManyRequest struct {
	RequestId    string
	MemberAmount []MemberAmount
	Memo         string
}

func (m *MixinSrv) TransferManyWithRetry(ctx context.Context, req *TransferManyRequest) error {
	var err error
	for i := 0; i < defaultMaxMixinRetry; i++ {
		if _, err = m.TransferMany(ctx, req); err != nil {
			log.Error().Err(err).Msg("send transfer many failed, retrying...")
			time.Sleep(time.Second << i)
			continue
		} else {
			return nil
		}
	}
	return err
}

func (m *MixinSrv) TransferOneWithRetry(ctx context.Context, req *TransferOneRequest) error {
	var err error
	for i := 0; i < defaultMaxMixinRetry; i++ {
		if _, err = m.TransferOne(ctx, req); err != nil {
			log.Error().Err(err).Msg("send transfer one failed, retrying...")
			time.Sleep(time.Second << i)
			continue
		} else {
			return nil
		}
	}
	return err
}

func (m *MixinSrv) TransferOne(ctx context.Context, req *TransferOneRequest) (*mixin.SafeTransactionRequest, error) {
	err := m.SyncArrgegateUtxos(ctx)
	if err != nil {
		return nil, err
	}

	m.transferMutex.Lock()
	defer m.transferMutex.Unlock()

	// 1. 将utxos聚合
	utxos, _ := m.Client.SafeListUtxos(ctx, mixin.SafeListUtxoOption{
		Asset:     consts.ASSET_CNB,
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
		if useAmount.GreaterThanOrEqual(req.Amount) {
			break
		}
		useAmount = useAmount.Add(utxo.Amount)
		useUtxos = append(useUtxos, utxo)
		if len(useUtxos) > consts.MAX_UTXO_NUM {
			useUtxos = useUtxos[1:]
		}
	}

	if useAmount.LessThanOrEqual(req.Amount) {
		return nil, fmt.Errorf("not enough utxos")
	}

	// 2: build transaction
	b := mixin.NewSafeTransactionBuilder(useUtxos)
	b.Memo = req.Memo

	txOutout := &mixin.TransactionOutput{
		Address: mixin.RequireNewMixAddress([]string{req.Member}, 1),
		Amount:  req.Amount,
	}

	tx, err := m.Client.MakeTransaction(ctx, b, []*mixin.TransactionOutput{txOutout})
	if err != nil {
		return nil, err
	}

	raw, err := tx.Dump()
	if err != nil {
		return nil, err
	}

	// 3. create transaction
	request, err := m.Client.SafeCreateTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID:      req.RequestId,
		RawTransaction: raw,
	})
	if err != nil {
		return nil, err
	}
	// 4. sign transaction
	err = mixin.SafeSignTransaction(
		tx,
		m.SpendKey,
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
	_, err = m.Client.SafeSubmitTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID:      req.RequestId,
		RawTransaction: signedRaw,
	})
	if err != nil {
		return nil, err
	}

	// 6. read transaction
	req1, err := m.Client.SafeReadTransactionRequest(ctx, req.RequestId)
	if err != nil {
		return nil, err
	}
	return req1, nil
}

func (m *MixinSrv) TransferMany(ctx context.Context, req *TransferManyRequest) (*mixin.SafeTransactionRequest, error) {
	err := m.SyncArrgegateUtxos(ctx)
	if err != nil {
		return nil, err
	}

	m.transferMutex.Lock()
	defer m.transferMutex.Unlock()

	totalAmount := decimal.NewFromInt(0)
	
	lo.ForEach(req.MemberAmount, func(item MemberAmount, index int) {
		totalAmount = totalAmount.Add(item.Amount)
	})

	// 1. 将utxos聚合
	utxos, _ := m.Client.SafeListUtxos(ctx, mixin.SafeListUtxoOption{
		Asset:     consts.ASSET_CNB,
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
		useAmount = useAmount.Add(utxo.Amount)
		useUtxos = append(useUtxos, utxo)
		if len(useUtxos) > consts.MAX_UTXO_NUM {
			useUtxos = useUtxos[1:]
		}
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
			Amount:  req.MemberAmount[i].Amount,
		}
	}

	tx, err := m.Client.MakeTransaction(ctx, b, txOutout)
	if err != nil {
		return nil, err
	}

	raw, err := tx.Dump()
	if err != nil {
		return nil, err
	}

	// 3. create transaction
	request, err := m.Client.SafeCreateTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID:      req.RequestId,
		RawTransaction: raw,
	})
	if err != nil {
		return nil, err
	}
	// 4. sign transaction
	err = mixin.SafeSignTransaction(
		tx,
		m.SpendKey,
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
	_, err = m.Client.SafeSubmitTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID:      req.RequestId,
		RawTransaction: signedRaw,
	})
	if err != nil {
		return nil, err
	}

	// 6. read transaction
	req1, err := m.Client.SafeReadTransactionRequest(ctx, req.RequestId)
	if err != nil {
		return nil, err
	}
	return req1, nil
}

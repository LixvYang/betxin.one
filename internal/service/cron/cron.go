package cron

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/service/mixin_srv"
	"github.com/lixvyang/betxin.one/internal/utils/convert"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

type CronService struct {
	mixinSrv *mixin_srv.MixinSrv
	storage  *mongo.MongoService
}

func NewCronService(mixinSrv *mixin_srv.MixinSrv, storage *mongo.MongoService) *CronService {
	return &CronService{
		mixinSrv: mixinSrv,
		storage:  storage,
	}
}

func (s *CronService) Run() {

}

// 定期查询入账记录，并根据入账记录更新用户的余额
func (s *CronService) CronQueryDeposit() {
	ctx := context.Background()
	now := time.Now()

	ticker := time.NewTicker(time.Second * 3)

	// 1. 拿到所有的入账记录
	// 2. 遍历入账记录，更新用户余额
	// 3. 保存用户余额
	for range ticker.C {
		snapshots, err := s.mixinSrv.Client.ReadSafeSnapshots(ctx, consts.ASSET_CNB, now.Add(-time.Hour>>1), "DESC", 500)
		if err != nil {
			log.Error().Err(err).Msg("read safe snapshots failed")
			continue
		}
		count, err := s.storage.GetSnapshotCount(ctx)
		if err != nil {
			continue
		}

		if count == 0 {
			for _, snapshot := range snapshots {
				if snapshot.Amount.IsPositive() {
					err = s.handleDeposit(snapshot)
					if err != nil {
						log.Error().Err(err).Msg("handle deposit failed")
						// todo: 记录错误日志 发送mixin消息
						continue
					}
				}
			}
		} else {
			lastestSnapshot, err := s.storage.GetLastestSnapshot(ctx)
			if err != nil {
				continue
			}

			for _, snapshot := range snapshots {
				if snapshot.RequestID == lastestSnapshot.RequestID {
					break
				}

				// 聚合 utxo 和 memo 为空 忽略
				if snapshot.Memo == "arrgegate utxos" || snapshot.Memo == "" {
					continue
				}

				// TODO
				if snapshot.Amount.IsPositive() {
					err = s.handleDeposit(snapshot)
					if err != nil {
						log.Error().Err(err).Msg("handle deposit failed")
						// 回退金额
						// 发送mixin消息
						continue
					} else {
						// 发送话题已购买记录
						// 发送mixin消息
					}
				}
			}
		}

	}
}

/*
memo格式 :
{
	"t": "uuid",    话题id
	"a": 0 | 1, YES / NO
}
*/

type MemoAction struct {
	Tid    string `json:"t"`
	Action bool   `json:"a"`
}

func (s *CronService) handleDeposit(snapshot *mixin.SafeSnapshot) error {
	ctx := context.Background()
	now := time.Now()
	_ = s.storage.InsertSnapshot(ctx, &schema.Snapshot{
		SnapshotID: snapshot.SnapshotID,
		RequestID:  snapshot.RequestID,
		UserID:     snapshot.UserID,
		AssetID:    consts.ASSET_CNB,
		Memo:       snapshot.Memo,
		CreatedAt:  snapshot.CreatedAt,
	})

	var memoAction MemoAction
	// 1. 解析memo
	err := json.Unmarshal([]byte(snapshot.Memo), &memoAction)
	if err != nil {
		return err
	}

	// 2. 处理memo
	_, err = convert.VaildUUID(memoAction.Tid)
	if err != nil {
		return err
	}

	topic, err := s.storage.GetTopicByTid(ctx, memoAction.Tid)
	if err != nil {
		return err
	}

	if topic.EndTime.Before(now) || topic.IsStop {
		// TODO 退款
		return errors.New("topic already stop")
	}

	return s.storage.HandleMixinTopicDepositAction(ctx, &schema.TopicBuyAction{
		Tid:       memoAction.Tid,
		Action:    memoAction.Action,
		Uid:       snapshot.UserID,
		Amount:    snapshot.Amount,
		RequestID: convert.NewUUID(),
	})
}

// 定期聚合utxos
func (m *CronService) ArrgegateUtxos(ctx context.Context, requestId string) error {
	var req *mixin.SafeTransactionRequest

	ticker := time.NewTicker(time.Second * 10)
	for range ticker.C {
		utxos, err := m.mixinSrv.Client.SafeListUtxos(ctx, mixin.SafeListUtxoOption{
			Asset:     consts.ASSET_CNB,
			State:     mixin.SafeUtxoStateUnspent,
			Threshold: 1,
		})

		if err != nil {
			log.Error().Err(err).Msg("list utxos failed")
			continue
		}

		if len(utxos) <= consts.MAX_UTXO_NUM {
			break
		}

		// 将utxos分割成 255 个大小的数组
		utxoSlice := make([]*mixin.SafeUtxo, 0)
		utxoSliceAmount := decimal.NewFromInt(0)
		for i := 0; i <= consts.MAX_UTXO_NUM; i++ {
			utxoSlice = append(utxoSlice, utxos[i])
			utxoSliceAmount = utxoSliceAmount.Add(utxos[i].Amount)
		}

		// 2: build transaction
		b := mixin.NewSafeTransactionBuilder(utxoSlice)
		b.Memo = "arrgegate utxos"

		tx, err := m.mixinSrv.Client.MakeTransaction(ctx, b, []*mixin.TransactionOutput{
			{
				Address: mixin.RequireNewMixAddress([]string{m.mixinSrv.Client.ClientID}, 1),
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
		request, err := m.mixinSrv.Client.SafeCreateTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
			RequestID:      requestId,
			RawTransaction: raw,
		})
		if err != nil {
			return err
		}

		// 4. sign transaction
		err = mixin.SafeSignTransaction(
			tx,
			m.mixinSrv.SpendKey,
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
		_, err = m.mixinSrv.Client.SafeSubmitTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
			RequestID:      requestId,
			RawTransaction: signedRaw,
		})
		if err != nil {
			return err
		}

		// 6. read transaction
		req, err = m.mixinSrv.Client.SafeReadTransactionRequest(ctx, requestId)
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

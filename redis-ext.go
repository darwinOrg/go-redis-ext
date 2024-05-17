package redisdkext

import (
	daogext "github.com/darwinOrg/daog-ext"
	dgcoll "github.com/darwinOrg/go-common/collection"
	dgctx "github.com/darwinOrg/go-common/context"
	dgcron "github.com/darwinOrg/go-cron"
	dglock "github.com/darwinOrg/go-dlock"
	dglogger "github.com/darwinOrg/go-logger"
	redisdk "github.com/darwinOrg/go-redis"
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
	"time"
)

const defaultMessagePageSize = 1000

func XAddValuesWithExpiration(ctx *dgctx.DgContext, stream string, values any, expiration time.Duration) (string, error) {
	messageId, err := redisdk.XAddValues(stream, values)
	if err != nil {
		return messageId, err
	}
	defer func() {
		if err != nil {
			_, de := redisdk.XDel(stream, messageId)
			if de != nil {
				dglogger.Errorf(ctx, "redisdk.XDel | streamKey: %s | error: %v", stream, de)
			}
		}
	}()

	now := time.Now()
	err = daogext.Write(ctx, func(tc *daog.TransContext) error {
		message := &RedisMessage{
			StreamKey: stream,
			MessageId: messageId,
			ExpiredAt: ttypes.NormalDatetime(now.Add(expiration)),
			CreatedBy: ctx.UserId,
			CreatedAt: ttypes.NormalDatetime(now),
		}
		_, err = RedisMessageDao.Insert(tc, message)
		if err != nil {
			dglogger.Errorf(ctx, "RedisMessageDao.Insert error: %v", err)
			return err
		}

		return nil
	})

	return messageId, err
}

func StartDeleteExpiredMessageTask(db daog.Datasource) {
	cron := dgcron.NewAndStart(dglock.NewDbLocker(db))
	cron.AddJobWithLock("定时删除过期redis消息", "0 0/5 * * * ?", 5000, func(ctx *dgctx.DgContext) {
		ctx.NotLogSQL = true
		needContinue := true

		for needContinue {
			needContinue, _ = daogext.WriteWithResult(ctx, func(tc *daog.TransContext) (bool, error) {
				matcher := daog.NewMatcher().Lt(RedisMessageFields.ExpiredAt, time.Now())
				count, err := RedisMessageDao.Count(tc, matcher)
				if err != nil {
					dglogger.Errorf(ctx, "RedisMessageDao.Count error: %v", err)
					return false, err
				}
				if count == 0 {
					return false, nil
				}

				messages, err := RedisMessageDao.QueryPageListMatcher(tc, matcher, daog.NewPager(defaultMessagePageSize, 1))
				if err != nil {
					dglogger.Errorf(ctx, "RedisMessageDao.QueryListMatcherWithViewColumns error: %v", err)
					return false, err
				}
				if len(messages) == 0 {
					return false, nil
				}

				streamKey2MessageIdsMap := dgcoll.Extract2KeyListMap(messages, func(message *RedisMessage) string {
					return message.StreamKey
				}, func(message *RedisMessage) string {
					return message.MessageId
				})
				for streamKey, messageIds := range streamKey2MessageIdsMap {
					_, err = redisdk.XDel(streamKey, messageIds...)
					if err != nil {
						dglogger.Errorf(ctx, "redisdk.XDel | streamKey: %s | error: %v", streamKey, err)
						return false, err
					}
				}

				_, err = RedisMessageDao.DeleteByIds(tc, dgcoll.MapToList(messages, func(message *RedisMessage) int64 {
					return message.Id
				}))
				if err != nil {
					dglogger.Errorf(ctx, "RedisMessageDao.DeleteByIds error: %v", err)
					return false, err
				}

				return count > defaultMessagePageSize, nil
			})
		}
	})
}

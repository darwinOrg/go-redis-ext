package redisdkext

import (
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
)

var RedisMessageFields = struct {
	Id        string
	StreamKey string
	MessageId string
	ExpiredAt string
	CreatedBy string
	CreatedAt string
}{
	"id",
	"stream_key",
	"message_id",
	"expired_at",
	"created_by",
	"created_at",
}

var RedisMessageMeta = &daog.TableMeta[RedisMessage]{
	Table: "redis_message",
	Columns: []string{
		"id",
		"stream_key",
		"message_id",
		"expired_at",
		"created_by",
		"created_at",
	},
	AutoColumn: "id",
	LookupFieldFunc: func(columnName string, ins *RedisMessage, point bool) any {
		if "id" == columnName {
			if point {
				return &ins.Id
			}
			return ins.Id
		}
		if "stream_key" == columnName {
			if point {
				return &ins.StreamKey
			}
			return ins.StreamKey
		}
		if "message_id" == columnName {
			if point {
				return &ins.MessageId
			}
			return ins.MessageId
		}
		if "expired_at" == columnName {
			if point {
				return &ins.ExpiredAt
			}
			return ins.ExpiredAt
		}
		if "created_by" == columnName {
			if point {
				return &ins.CreatedBy
			}
			return ins.CreatedBy
		}
		if "created_at" == columnName {
			if point {
				return &ins.CreatedAt
			}
			return ins.CreatedAt
		}

		return nil
	},
}

var RedisMessageDao daog.QuickDao[RedisMessage] = &struct {
	daog.QuickDao[RedisMessage]
}{
	daog.NewBaseQuickDao(RedisMessageMeta),
}

type RedisMessage struct {
	Id        int64
	StreamKey string
	MessageId string
	ExpiredAt ttypes.NormalDatetime
	CreatedBy int64
	CreatedAt ttypes.NormalDatetime
}

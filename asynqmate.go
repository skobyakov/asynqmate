package main

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type TaskState string

type Task struct {
	ID  string
	Msg string
}

const (
	queuesKey = "asynq:queues"

	TaskStatePending TaskState = "pending"
)

var searchTasksByStateCmd = redis.NewScript(`
local res = {}
local tasks_ids = redis.call("LRANGE", "asynq:{asynqmate}:" ..ARGV[1], 0, -1)
local hits = 0
for _, tid in pairs(tasks_ids) do
	local msg = redis.call("HGET", "asynq:{asynqmate}:t:" ..tid , "msg")
	local index = string.find(msg, ARGV[2])

	if index ~= nil then
		hits = hits + 1
		table.insert(res, tid)
		table.insert(res, msg)
	end

	if hits == 10 then return res end
end
return res
`)

type AsynqMate struct {
	rc redis.UniversalClient
}

func NewAsynqMate(addr string) *AsynqMate {
	rc := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{addr},
	})

	return &AsynqMate{
		rc: rc,
	}
}

func (am *AsynqMate) ListAllQueuesNames(ctx context.Context) ([]string, error) {
	return am.rc.SMembers(ctx, queuesKey).Result()
}

func (am *AsynqMate) SearchForTasks(ctx context.Context, state TaskState, msg string) ([]*Task, error) {
	res, err := searchTasksByStateCmd.Run(ctx, am.rc, []string{}, []interface{}{string(state), msg}).Result()
	if err != nil {
		return nil, err
	}

	items, ok := res.([]interface{})
	if !ok {
		return nil, errors.New("can't cast Redis response")
	}

	tasks := make([]*Task, 0, len(items))
	for i, it := range items {
		t, ok := it.(string)
		if !ok {
			return nil, errors.New("can't cast Redis response")
		}

		if i%2 == 0 {
			tasks = append(tasks, &Task{ID: t})
		} else {
			tasks[len(tasks)-1].Msg = t
		}
	}

	return tasks, nil
}

func (am *AsynqMate) Close() {
	am.rc.Close()
}

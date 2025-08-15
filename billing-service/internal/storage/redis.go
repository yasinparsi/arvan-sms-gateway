package storage

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func (r *RedisClient) Ping(context context.Context) {
	panic("unimplemented")
}

func (r *RedisClient) Client() {
	panic("unimplemented")
}

func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisClient{client: rdb}
}

var luaCheckAmount = redis.NewScript(`
local current = redis.call("GET", KEYS[1])
if not current then
  return {-1, 0}
end
local amount = tonumber(current)
local required = tonumber(ARGV[1])
if amount >= required then
  local new = redis.call("DECRBY", KEYS[1], required)
  return {1, new}
else
  return {0, amount}
end
`)

// CheckAmountAtomically returns (allowed, currentAmount, error)
func (r *RedisClient) CheckAmountAtomically(ctx context.Context, userID string, requiredAmount int64) (bool, int64, error) {
	key := "user:amount:" + userID
	result, err := luaCheckAmount.Run(ctx, r.client, []string{key}, requiredAmount).Result()
	if err != nil {
		return false, 0, err
	}

	arr, ok := result.([]interface{})
	if !ok || len(arr) != 2 {
		return false, 0, nil
	}

	status, ok1 := arr[0].(int64)
	amount, ok2 := arr[1].(int64)
	if !ok1 || !ok2 {
		return false, 0, nil
	}

	switch status {
	case -1:
		return false, 0, nil
	case 1:
		return true, amount, nil
	case 0:
		return false, amount, nil
	}
	return false, amount, nil
}

func (r *RedisClient) ChargeUser(ctx context.Context, userID string, chargeAmount int64) error {
	key := "user:amount:" + userID
	return r.client.IncrBy(ctx, key, chargeAmount).Err()
}

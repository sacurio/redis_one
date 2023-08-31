package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

const (
	noExists = 0
	exists   = 1
)

type rdbHandler struct {
	ctx    context.Context
	client *redis.Client
}

func NewRDBHandler(ctx context.Context, addr string) rdbHandler {
	return rdbHandler{
		ctx: ctx,
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func main() {
	rdb := NewRDBHandler(context.Background(), "localhost:6379")
	defer rdb.client.Close()

	status := rdb.client.Ping(rdb.ctx)
	fmt.Printf("%+v\n", status)

	rdb.add("name", "Sandy")
	rdb.add("lastname", "Acurio")
	rdb.add("age", "41")

	rdb.delete("age")

	rdb.print("*")
}

func (h rdbHandler) add(k, v string) {
	result, err := h.client.Exists(h.ctx, k).Result()
	if err != nil {
		fmt.Printf("error on try to find %s key\n", k)
	}

	if result == noExists {
		err := h.client.Set(h.ctx, k, v, 0).Err()
		if err != nil {
			fmt.Printf("error: %s", err.Error())
		}
		fmt.Printf("key: %s saved with value: %s successfully\n", k, v)
		return
	}
	fmt.Printf("key: %s already exists\n", k)
}

func (h rdbHandler) print(pattern string) {
	keys, err := h.client.Keys(h.ctx, pattern).Result()
	fmt.Println()
	if err != nil {
		fmt.Printf("error retrieving keys, %s\n", err.Error())
	}

	for _, k := range keys {
		v, err := h.client.Get(h.ctx, k).Result()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Printf("%s:%+v\n", k, v)
	}
}

func (h rdbHandler) delete(k string) {
	if err := h.client.Del(h.ctx, k).Err(); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%s key deleted\n", k)
}

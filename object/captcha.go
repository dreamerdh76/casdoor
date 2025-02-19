// Copyright 2021 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/casdoor/casdoor/conf"
	"github.com/dchest/captcha"
	"github.com/redis/go-redis/v9"
)

var (
	redisStore  *RedisStore
	RedisClient *redis.Client
	ctx         = context.Background()
)

type RedisStore struct {
	Expire time.Duration
}

func (rs *RedisStore) Set(id string, digits []byte) {
	RedisClient.Set(ctx, id, string(digits), rs.Expire)
}

func (rs *RedisStore) Get(id string, clear bool) (digits []byte) {
	bs, err := RedisClient.Get(ctx, id).Bytes()
	if err != nil {
		log.Println(err)
	}

	return bs
}

func GetCaptcha() (string, []byte, error) {
	if err := NewRedisClient(); err == nil {
		redisStore = &RedisStore{Expire: 2 * time.Minute}
		captcha.SetCustomStore(redisStore)
	}

	id := captcha.NewLen(5)

	var buffer bytes.Buffer

	err := captcha.WriteImage(&buffer, id, 200, 80)
	if err != nil {
		return "", nil, err
	}

	return id, buffer.Bytes(), nil
}

func VerifyCaptcha(id string, digits string) bool {
	res := captcha.VerifyString(id, digits)

	return res
}

func NewRedisClient() error {
	if conf.GetConfigString("redisEndpoint") == "" {
		return errors.New("error: redis not configurated")
	}
	readTimeout := 2 * time.Second
	writeTimeout := 2 * time.Second
	dns := strings.Split(conf.GetConfigString("redisEndpoint"), ",")

	if len(dns) == 3 {
		dbName, _ := strconv.Atoi(dns[1])
		RedisClient = redis.NewClient(&redis.Options{
			Network:      "tcp",
			Addr:         dns[0],
			Password:     dns[2],
			DB:           dbName,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		})
	} else {
		RedisClient = redis.NewClient(&redis.Options{
			Network:      "tcp",
			Addr:         conf.GetConfigString("redisEndpoint"),
			Password:     "", // no password set
			DB:           0,  // use default DB
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		})
	}

	_, err := RedisClient.Ping(ctx).Result()
	return err
}

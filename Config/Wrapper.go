package Config

import (
	"fmt"
	redis "gopkg.in/redis.v2"
	"strconv"
)

/*
	Wrapper
*/
func NewWrapper(f File) (w Wrapper, err error) {
	var r *redis.Client
	if r, err = f.Connect(); err != nil {
		return
	}

	return Wrapper{Redis: r}, nil
}

type Wrapper struct {
	Redis *redis.Client
}

func (self Wrapper) FetchIds(key string, start int64, end int64) (ids []int64, err error) {
	req := self.Redis.LRange(key, start, end)
	if err = req.Err(); err != nil {
		return
	}

	// optionally halting
	length := len(req.Val())
	if length == 0 {
		return ids, nil
	}

	// converting them
	ids = make([]int64, length)
	var i int
	for k, v := range req.Val() {
		i, err = strconv.Atoi(v)
		if err != nil {
			return
		}
		ids[k] = int64(i)
	}

	return ids, nil
}

func (self Wrapper) IncrAll(key string, count int) (ids []int64, err error) {
	// misc
	var cmds []redis.Cmder
	pipe := self.Redis.Pipeline()

	// running the pipeline
	for i := 0; i < count; i++ {
		pipe.Incr(key)
	}
	cmds, err = pipe.Exec()
	if err != nil {
		return
	}

	// gathering for ids and checking for errors
	ids = make([]int64, len(cmds))
	for i, cmd := range cmds {
		if err = cmd.Err(); err != nil {
			return
		}

		ids[i] = cmd.(*redis.IntCmd).Val()
	}

	return ids, nil
}

func (self Wrapper) RPushAll(key string, values []string) (err error) {
	// misc
	var cmds []redis.Cmder
	pipe := self.Redis.Pipeline()

	// running the pipeline
	for _, v := range values {
		pipe.RPush(key, v)
	}
	cmds, err = pipe.Exec()
	if err != nil {
		return
	}

	// checking for errors
	for _, cmd := range cmds {
		if err = cmd.Err(); err != nil {
			return
		}
	}

	return nil
}

func (self Wrapper) SAddAll(key string, values []string) (err error) {
	// misc
	var cmds []redis.Cmder
	pipe := self.Redis.Pipeline()

	// running the pipeline
	for _, v := range values {
		pipe.SAdd(key, v)
	}
	cmds, err = pipe.Exec()
	if err != nil {
		return
	}

	// checking for errors
	for _, cmd := range cmds {
		if err = cmd.Err(); err != nil {
			return
		}
	}

	return nil
}

func (self Wrapper) SMembers(key string) (members []string, err error) {
	cmd := self.Redis.SMembers(key)
	if err = cmd.Err(); err != nil {
		return
	}
	return cmd.Val(), nil
}

func (self Wrapper) SIsMember(key string, value string) (isMember bool, err error) {
	values := []string{value}
	var isMembers []bool
	if isMembers, err = self.SIsMemberAll(key, values); err != nil {
		return
	}

	return isMembers[0], nil
}

func (self Wrapper) SIsMemberAll(key string, values []string) (isMembers []bool, err error) {
	var cmds []redis.Cmder
	pipe := self.Redis.Pipeline()

	// running the pipeline
	for _, v := range values {
		pipe.SIsMember(key, v)
	}
	if cmds, err = pipe.Exec(); err != nil {
		return
	}

	// checking for errors
	isMembers = make([]bool, len(cmds))
	for i, cmd := range cmds {
		if err = cmd.Err(); err != nil {
			return
		}

		isMembers[i] = cmd.(*redis.BoolCmd).Val()
	}

	return isMembers, nil
}

func (self Wrapper) SetAll(values map[string]string) (err error) {
	var cmds []redis.Cmder
	pipe := self.Redis.Pipeline()

	// running the pipeline
	for k, v := range values {
		pipe.Set(k, v)
	}
	if cmds, err = pipe.Exec(); err != nil {
		return
	}

	// checking for errors
	for _, cmd := range cmds {
		if err = cmd.Err(); err != nil {
			return
		}
	}

	return nil
}

func (self Wrapper) Get(k string) (v string, err error) {
	cmd := self.Redis.Get(k)
	if err = cmd.Err(); err != nil && err != redis.Nil {
		return
	}

	return cmd.Val(), nil
}

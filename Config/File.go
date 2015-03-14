package Config

import (
	"encoding/json"
	"errors"
	redis "gopkg.in/redis.v2"
	"io/ioutil"
	"path/filepath"
)

/*
	File
*/
func NewFile(configPath string) (configFile File, err error) {
	if len(configPath) == 0 {
		err = errors.New("Config path was blank!")
		return
	}

	var fullConfigPath string
	if fullConfigPath, err = filepath.Abs(configPath); err != nil {
		return
	}

	var b []byte
	if b, err = ioutil.ReadFile(fullConfigPath); err != nil {
		return
	}

	if err = json.Unmarshal(b, &configFile); err != nil {
		return
	}

	return
}

type File struct {
	Host     string
	Password string
	Db       int64
}

func (self File) Connect() (r *redis.Client, err error) {
	r = redis.NewTCPClient(&redis.Options{
		Addr:     self.Host,
		Password: self.Password,
		DB:       self.Db,
	})

	ping := r.Ping()
	if err = ping.Err(); err != nil {
		return
	}

	return
}

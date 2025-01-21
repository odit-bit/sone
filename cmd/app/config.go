package app

import (
	"errors"
	"fmt"
)

type Config struct {
	Http struct {
		Port int
	}

	Rtmp struct {
		Port int
	}

	Rpc struct {
		Port int
	}

	Logging struct {
		// print all logging into os.stdout
		Debug bool

		// // verbosity of logging , default "verbose"
		// Level string
	}

	// media storage to store segment of video,
	// it could be in memory or any other object storage implementation (minio, s3)
	MediaStorage int

	Minio struct {
		// dsn string connection or address
		Address   string
		AccessKey string // username
		SecretKey string // password
	}

	Filesystem struct {
		// to use memory leave it blank or "memory"
		Path string
	}

	// sql configuration backend
	SQL struct {
		// dsn string connection, for memory leave it blank or "memory"
		DSN string
	}

	// kv store, currently still use in-memory go map
	KV struct {
		DSN string
	}
}

func (c *Config) Validate() error {
	var err error

	if c.Http.Port == 0 {
		err = errors.Join(err, fmt.Errorf("http port must specified"))
	}
	if c.Rtmp.Port == 0 {
		err = errors.Join(err, fmt.Errorf("rtmp port must specified"))
	}
	if c.Rpc.Port == 0 {
		err = errors.Join(err, fmt.Errorf("rpc port must specified"))
	}

	return err
}

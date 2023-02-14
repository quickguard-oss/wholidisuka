package cache

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type cache struct {
	dir    string
	silent bool
}

func New(dir string, silent bool) (cache, error) {
	c := cache{
		dir:    dir,
		silent: silent,
	}

	if err := c.mkCacheDir(); err != nil {
		return cache{}, err
	}

	return c, nil
}

func (c cache) Set(key string, data []byte) error {
	path := c.path(key)

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	if _, err = file.Write(data); err != nil {
		return err
	}

	return nil
}

func (c cache) Get(key string, expire time.Duration) ([]byte, error) {
	path := c.path(key)

	stat, err := os.Stat(path)

	if err != nil {
		c.printf("Cache file '%s' does not exist.\n", path)

		return nil, nil
	}

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	if expire < time.Since(stat.ModTime()) {
		c.printf("'%s' has expired.\n", path)

		return nil, nil
	}

	data, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c cache) mkCacheDir() error {
	return os.MkdirAll(c.dir, os.ModePerm)
}

func (c cache) path(key string) string {
	return filepath.Join(c.dir, key)
}

func (c cache) printf(format string, a ...interface{}) {
	if !c.silent {
		fmt.Printf(format, a...)
	}
}

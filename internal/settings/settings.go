package settings

import (
	"fmt"
	"os"
	"path"

	"github.com/yohamta/dagman/internal/utils"
)

var ErrConfigNotFound = fmt.Errorf("config not found")

var cache map[string]string = nil

const (
	CONFIG__DATA_DIR = "DAGMAN__DATA"
	CONFIG__LOGS_DIR = "DAGMAN__LOGS"
)

func MustGet(name string) string {
	val, err := Get(name)
	if err != nil {
		panic(fmt.Errorf("failed to get %s : %w", name, err))
	}
	return val
}

func init() {
	load()
}

func Get(name string) (string, error) {
	if val, ok := cache[name]; ok {
		return val, nil
	}
	return "", ErrConfigNotFound
}

func load() {
	dir := utils.MustGetUserHomeDir()

	cache = map[string]string{}
	cache[CONFIG__DATA_DIR] = config(
		CONFIG__DATA_DIR,
		path.Join(dir, "/.dagman/data"))
	cache[CONFIG__LOGS_DIR] = config(
		CONFIG__LOGS_DIR,
		path.Join(dir, "/.dagman/logs"))
}

func InitTest(dir string) {
	os.Setenv("HOME", dir)
	load()
}

func config(env, def string) string {
	val := os.ExpandEnv(fmt.Sprintf("${%s}", env))
	if val == "" {
		return def
	}
	return val
}

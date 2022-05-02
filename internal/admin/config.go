package admin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/yohamta/dagu/internal/utils"
)

var tickerMatcher *regexp.Regexp

func init() {
	tickerMatcher = regexp.MustCompile("`[^`]+`")
}

type Config struct {
	Host               string
	Port               string
	Env                []string
	DAGs               string
	Command            string
	WorkDir            string
	IsBasicAuth        bool
	BasicAuthUsername  string
	BasicAuthPassword  string
	LogEncodingCharset string
}

func (c *Config) Init() {
	if c.Env == nil {
		c.Env = []string{}
	}
}

func (c *Config) setup() error {
	if c.Command == "" {
		c.Command = "dagu"
	}
	if c.DAGs == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		c.DAGs = wd
	}
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}
	if c.Port == "" {
		c.Port = "8000"
	}
	if len(c.Env) == 0 {
		env := utils.DefaultEnv()
		env, err := loadVariables(env)
		if err != nil {
			return err
		}
		c.Env = buildConfigEnv(env)
	}
	return nil
}

func buildFromDefinition(def *configDefinition) (c *Config, err error) {
	c = &Config{}
	c.Init()

	env, err := loadVariables(def.Env)
	if err != nil {
		return nil, err
	}
	c.Env = buildConfigEnv(env)

	c.Host, err = parseVariable(def.Host)
	if err != nil {
		return nil, err
	}
	c.Port = strconv.Itoa(def.Port)

	jd, err := parseVariable(def.Dags)
	if err != nil {
		return nil, err
	}
	if !filepath.IsAbs(jd) {
		return nil, fmt.Errorf("DAGs directory should be absolute path. was %s", jd)
	}
	c.DAGs, err = filepath.Abs(jd)
	if err != nil {
		return nil, err
	}
	c.Command, err = parseVariable(def.Command)
	if err != nil {
		return nil, err
	}
	c.WorkDir, err = parseVariable(def.WorkDir)
	if err != nil {
		return nil, err
	}
	if c.WorkDir == "" {
		c.WorkDir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	c.IsBasicAuth = def.IsBasicAuth
	c.BasicAuthUsername, err = parseVariable(def.BasicAuthUsername)
	if err != nil {
		return nil, err
	}
	c.BasicAuthPassword, err = parseVariable(def.BasicAuthPassword)
	if err != nil {
		return nil, err
	}
	c.LogEncodingCharset, err = parseVariable(def.LogEncodingCharset)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func buildConfigEnv(vars map[string]string) []string {
	ret := []string{}
	for k, v := range vars {
		ret = append(ret, fmt.Sprintf("%s=%s", k, v))
	}
	return ret
}

func loadVariables(strVariables map[string]string) (map[string]string, error) {
	vars := map[string]string{}
	for k, v := range strVariables {
		parsed, err := parseVariable(v)
		if err != nil {
			return nil, err
		}
		vars[k] = parsed
		err = os.Setenv(k, parsed)
		if err != nil {
			return nil, err
		}
	}
	return vars, nil
}

func parseVariable(value string) (string, error) {
	val, err := parseCommand(os.ExpandEnv(value))
	if err != nil {
		return "", err
	}
	return val, nil
}

func parseCommand(value string) (string, error) {
	matches := tickerMatcher.FindAllString(strings.TrimSpace(value), -1)
	if matches == nil {
		return value, nil
	}
	ret := value
	for i := 0; i < len(matches); i++ {
		command := matches[i]
		out, err := exec.Command(strings.ReplaceAll(command, "`", "")).Output()
		if err != nil {
			return "", err
		}
		ret = strings.ReplaceAll(ret, command, strings.TrimSpace(string(out[:])))

	}
	return ret, nil
}

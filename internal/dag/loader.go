package dag

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/yohamta/dagu/internal/pipeline"
	"go.starlark.net/starlark"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"github.com/yohamta/dagu/internal/utils"

	"gopkg.in/yaml.v2"
)

var ErrDAGNotFound = errors.New("DAG was not found")

// Loader is a config loader.
type Loader struct {
	BaseConfig string
}

// Load loads config from file.

func (cl *Loader) Load(f, params string) (*DAG, error) {
	return cl.loadDAG(f,
		&BuildDAGOptions{
			parameters: params,
		},
	)
}

// LoadwIithoutEval loads config from file without evaluating env variables.
func (cl *Loader) LoadWithoutEval(f string) (*DAG, error) {
	return cl.loadDAG(f,
		&BuildDAGOptions{
			parameters: "",
			headOnly:   false,
			noEval:     true,
			noSetenv:   true,
		},
	)
}

// LoadHeadOnly loads config from file and returns only the headline data.
func (cl *Loader) LoadHeadOnly(f string) (*DAG, error) {
	return cl.loadDAG(f,
		&BuildDAGOptions{
			parameters: "",
			headOnly:   true,
			noEval:     true,
			noSetenv:   true,
		},
	)
}

// LoadData loads config from given data.
func (cl *Loader) LoadData(data []byte) (*DAG, error) {
	raw, err := cl.unmarshalData(data)
	if err != nil {
		return nil, err
	}
	def, err := cl.decode(raw)
	if err != nil {
		return nil, err
	}
	b := &builder{
		BuildDAGOptions: BuildDAGOptions{
			headOnly: false,
			noEval:   true,
			noSetenv: true,
		},
	}
	return b.buildFromDefinition(def, nil)
}

func (cl *Loader) loadBaseConfig(file string, opts *BuildDAGOptions) (*DAG, error) {
	if !utils.FileExists(file) {
		return nil, nil
	}

	def, err := cl.loadFromYAML(file)
	if err != nil {
		return nil, err
	}
	buildOpts := *opts
	buildOpts.headOnly = false
	buildOpts.defaultEnv = utils.DefaultEnv()
	b := &builder{
		BuildDAGOptions: buildOpts,
	}
	return b.buildFromDefinition(def, nil)
}

func (cl *Loader) loadDAG(f string, opts *BuildDAGOptions) (*DAG, error) {
	if f == "" {
		return nil, fmt.Errorf("config file was not specified")
	}
	ext := filepath.Ext(f)
	if ext == "" {
		f = fmt.Sprintf("%s.yaml", f)
	}
	file, err := filepath.Abs(f)
	if err != nil {
		return nil, err
	}

	var dst *DAG = nil

	if !opts.headOnly && cl.BaseConfig != "" {
		dst, err = cl.loadBaseConfig(cl.BaseConfig, opts)
		if err != nil {
			return nil, err
		}
	}

	if dst == nil {
		dst = &DAG{}
		dst.Init()
	}

	dst.Name = strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

	var def *configDefinition
	if ext == ".star" {
		def, err = cl.loadFromStarLark(file)
	} else {
		def, err = cl.loadFromYAML(file)
	}
	if err != nil {
		return nil, err
	}

	b := builder{BuildDAGOptions: *opts}
	c, err := b.buildFromDefinition(def, dst)

	if err != nil {
		return nil, err
	}

	err = cl.merge(dst, c)
	if err != nil {
		return nil, err
	}

	dst.Location = file

	if !opts.noSetenv {
		dst.setup()
	}

	return dst, nil
}

func (cl *Loader) loadFromYAML(file string) (*configDefinition, error) {
	raw, err := cl.load(file)
	if err != nil {
		return nil, err
	}

	def, err := cl.decode(raw)
	if err != nil {
		return nil, err
	}
	return def, nil
}

func (cl *Loader) loadFromStarLark(file string) (*configDefinition, error) {
	graph, err := pipeline.NewPipeline(file)
	if err != nil {
		return nil, err
	}
	c := &configDefinition{}
	c.Name = path.Base(file)
	nodes := graph["nodes"].(*starlark.Dict)
	edges := graph["edges"].(*starlark.List)
	dependencyMap := pipeline.EdgesToDependencyMap(edges)
	for _, name := range nodes.Keys() {
		value, found, err := nodes.Get(name)
		if err != nil {
			return nil, err
		}
		if !found {
			return nil, fmt.Errorf("node %v not found", name)
		}
		node := value.(*starlark.Dict)
		task, found, err := node.Get(starlark.String("task"))
		if task == nil || !found {
			return nil, fmt.Errorf("task of node %v missing", name)
		}
		if err != nil {
			return nil, err
		}
		taskDict := task.(*starlark.Dict)
		nameStr := name.(starlark.String).GoString()
		executor := lookupDict(taskDict, "executor")
		command := lookupDict(taskDict, "command")
		script := lookupDict(taskDict, "script")
		step := &stepDef{
			Name:     nameStr,
			Depends:  dependencyMap[nameStr],
			Executor: executor,
			Command:  command,
			Script:   script,
		}
		c.Steps = append(c.Steps, step)
	}
	return c, nil
}

func lookupDict(dict *starlark.Dict, name string) string {
	value, found, err := dict.Get(starlark.String(name))
	if !found {
		fmt.Printf("missing %v in dict", name)
	}
	if err != nil {
		fmt.Printf("lookup error: %v", err)
	}
	return value.(starlark.String).GoString()
}

type mergeTranformer struct {
}

var _ mergo.Transformers = (*mergeTranformer)(nil)

func (mt *mergeTranformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(MailOn{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				dst.Set(src)
			}
			return nil
		}
	}
	return nil
}

func (cl *Loader) merge(dst, src *DAG) error {
	err := mergo.Merge(dst, src, mergo.WithOverride,
		mergo.WithTransformers(&mergeTranformer{}))
	return err
}

func (cl *Loader) load(file string) (config map[string]interface{}, err error) {
	return cl.readFile(file)
}

func (cl *Loader) readFile(file string) (config map[string]interface{}, err error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", file, err)
	}
	return cl.unmarshalData(data)
}

func (cl *Loader) unmarshalData(data []byte) (map[string]interface{}, error) {
	var cm map[string]interface{}
	err := yaml.NewDecoder(bytes.NewReader(data)).Decode(&cm)
	return cm, err
}

func (cl *Loader) decode(cm map[string]interface{}) (*configDefinition, error) {
	c := &configDefinition{}
	md, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: true,
		Result:      c,
		TagName:     "",
	})
	err := md.Decode(cm)
	return c, err
}

package envvar

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

var (
	// TagKey -
	TagKey = "envvar"
	// ErrInvalidType -
	ErrInvalidType = errors.New("struct pointer required")
)

// Inject injects environment variable to corresponding fields.
func Inject(s interface{}) error {
	return InjectEnvs(s, os.Environ())
}

// InjectEnvs injects given environments variable to corresponding fields.
func InjectEnvs(sp interface{}, envs []string) error {
	v := reflect.ValueOf(sp)

	// type check
	if v.Kind() != reflect.Ptr {
		return ErrInvalidType
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrInvalidType
	}

	m := toMap(envs)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		vf := v.Field(i)
		if !vf.CanSet() {
			continue
		}
		tag := readTag(t.Field(i))
		if tag.skip() {
			continue
		}

		// 環境変数から探す
		envValue, found := m.lookup(tag.envKey)
		if tag.required && !found {
			return fmt.Errorf("%s is required, but not found", tag.envKey)
		}

		// 現状はstringのみ
		if vf.Kind() != reflect.String {
			continue
		}

		if !found {
			vf.SetString(tag.defaultValue)
			continue
		}

		vf.SetString(envValue)
	}

	return nil
}

func readTag(sf reflect.StructField) tag {
	raw := sf.Tag.Get(TagKey)
	specs := strings.Split(raw, ",")

	tag := tag{}
	for i, spec := range specs {
		if i == 0 {
			tag.envKey = spec
			continue
		}
		spec = strings.TrimSpace(spec)
		if spec == "required" {
			tag.required = true
			continue
		}
		if strings.HasPrefix(spec, "default") {
			kv := strings.SplitN(spec, "=", 2)
			tag.defaultValue = kv[1]
			continue
		}
	}
	return tag
}

type tag struct {
	envKey       string
	defaultValue string
	required     bool
}

func (t tag) skip() bool {
	return t.envKey == ""
}

type envMap map[string]string

func toMap(envs []string) envMap {
	m := make(envMap)
	for _, env := range envs {
		kv := strings.SplitN(env, "=", 2)
		m[kv[0]] = kv[1]
	}
	return m
}

func (e envMap) lookup(key string) (string, bool) {
	v, found := e[key]
	return v, found
}

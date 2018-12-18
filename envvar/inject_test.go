package envvar_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ymgyt/appkit/envvar"
)

func TestInjectEnvs(t *testing.T) {

	type sample struct {
		Name string `json:"name,omitempty" envvar:"ENVVAR_NAME" datastore:"name,noindex"`
		Must string `envvar:"ENVVAR_MUST, required"`
		Path string `envvar:"ENVVAR_PATH,default=/ymgyt/home"`
	}

	tests := []struct {
		desc string
		envs []string
		want sample
	}{
		{
			desc: "simple",
			envs: []string{
				"ENVVAR_NAME=gopher",
				"ENVVAR_MUST=ok",
			},
			want: sample{
				Name: "gopher",
				Must: "ok",
				Path: "/ymgyt/home",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s := &sample{}
			err := envvar.InjectEnvs(s, tc.envs)
			if err != nil {
				t.Fatal()
			}
			if diff := cmp.Diff(s, &tc.want); diff != "" {
				t.Errorf("(-got +want)\n%s", diff)
			}
		})
	}
}

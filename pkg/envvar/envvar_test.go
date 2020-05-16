package envvar_test

import (
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/jpdejavite/rtg-go-toolkit/pkg/envvar"
)

func TestLoadAllMissing(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("LoadAll() should have a recover")
		}
	}()

	os.Setenv("MY_ENVVAR_1", "fks9qe")
	envvar.LoadAll([]string{"MY_ENVVAR_1", "MY_ENVVAR_2"})
}

func TestLoadAllOk(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			t.Errorf("LoadAll() should not have a recover %v", r)
		}
	}()

	os.Setenv("MY_ENVVAR_4", "12310ro")
	os.Setenv("MY_ENVVAR_5", "lkvjsdiwr014i")
	envvar.LoadAll([]string{"MY_ENVVAR_4", "MY_ENVVAR_5"})

	if diff := deep.Equal(envvar.GetEnvVar("MY_ENVVAR_4"), "12310ro"); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(envvar.GetEnvVar("MY_ENVVAR_5"), "lkvjsdiwr014i"); diff != nil {
		t.Error(diff)
	}
}

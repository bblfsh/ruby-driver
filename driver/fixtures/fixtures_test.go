package fixtures

import (
	"path/filepath"
	"testing"

	"github.com/bblfsh/ruby-driver/driver/normalizer"
	"gopkg.in/bblfsh/sdk.v2/sdk/driver"
	"gopkg.in/bblfsh/sdk.v2/sdk/driver/fixtures"
)

const projectRoot = "../../"

var Suite = &fixtures.Suite{
	Lang: "ruby",
	Ext:  ".rb",
	Path: filepath.Join(projectRoot, fixtures.Dir),
	NewDriver: func() driver.BaseDriver {
		return driver.NewExecDriverAt(filepath.Join(projectRoot, "native/exe/native"))
	},
	Transforms: driver.Transforms{
		Native: normalizer.Native,
		Code:   normalizer.Code,
	},
	BenchName: "class_complete",
}

func TestRubyDriver(t *testing.T) {
	Suite.RunTests(t)
}

func BenchmarkRubyDriver(b *testing.B) {
	Suite.RunBenchmarks(b)
}

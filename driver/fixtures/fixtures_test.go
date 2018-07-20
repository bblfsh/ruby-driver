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
		return driver.NewExecDriverAt(filepath.Join(projectRoot, "build/bin/native"))
	},
	Transforms: normalizer.Transforms,
	BenchName: "class_complete",
	Semantic: fixtures.SemanticConfig{
		BlacklistTypes: []string{
			"def",
			"str",
			"splay",
			"lvar",
			"ivar",
			"gvar",
			"cvar",
			"Sym",
			"Const",
			"flip_1",
			"flip_2",
			"arg",
			"kwarg",
			"optarg",
			"kwoptarg",
			"restarg",
			"kwrestarg",
			"comment",
		},
	},
	Docker: fixtures.DockerConfig{
		Image: "ruby:2.4",
	},
}

func TestRubyDriver(t *testing.T) {
	Suite.RunTests(t)
}

func BenchmarkRubyDriver(b *testing.B) {
	Suite.RunBenchmarks(b)
}

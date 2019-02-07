package fixtures

import (
	"path/filepath"
	"testing"

	"github.com/bblfsh/ruby-driver/driver/normalizer"
	"gopkg.in/bblfsh/sdk.v2/driver"
	"gopkg.in/bblfsh/sdk.v2/driver/fixtures"
	"gopkg.in/bblfsh/sdk.v2/driver/native"
)

const projectRoot = "../../"

var Suite = &fixtures.Suite{
	Lang: "ruby",
	Ext:  ".rb",
	Path: filepath.Join(projectRoot, fixtures.Dir),
	NewDriver: func() driver.Native {
		return native.NewDriverAt(filepath.Join(projectRoot, "build/bin/native"), native.UTF8)
	},
	Transforms: normalizer.Transforms,
	BenchName:  "class_complete",
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
			"blockarg",
			"true",
			"false",
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

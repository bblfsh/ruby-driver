package fixtures

import (
	"path/filepath"
	"testing"

	"github.com/bblfsh/ruby-driver/driver/normalizer"
	"github.com/bblfsh/sdk/v3/driver"
	"github.com/bblfsh/sdk/v3/driver/fixtures"
	"github.com/bblfsh/sdk/v3/driver/native"
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
			"Const",
			"Sym",
			"arg",
			"begin",
			"blockarg",
			"comment",
			"cvar",
			"def",
			"false",
			"flip_1",
			"flip_2",
			"gvar",
			"ivar",
			"kwarg",
			"kwoptarg",
			"kwrestarg",
			"lvar",
			"optarg",
			"restarg",
			"splay",
			"str",
			"true",
		},
	},
}

func TestRubyDriver(t *testing.T) {
	Suite.RunTests(t)
}

func BenchmarkRubyDriver(b *testing.B) {
	Suite.RunBenchmarks(b)
}

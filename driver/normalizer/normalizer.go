package normalizer

import (
	"github.com/bblfsh/sdk/uast"
	"github.com/bblfsh/sdk/uast/ann"
)

var NativeToNoder = &uast.BaseToNoder{}
var AnnotationRules *ann.Rule = nil

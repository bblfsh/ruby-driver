package normalizer

import (
	//"errors"

	"github.com/bblfsh/ruby-driver/driver/normalizer/rubyast"

	"gopkg.in/bblfsh/sdk.v1/uast"
	. "gopkg.in/bblfsh/sdk.v1/uast/ann"
	"gopkg.in/bblfsh/sdk.v1/uast/transformer"
	"gopkg.in/bblfsh/sdk.v1/uast/transformer/annotatter"
	"gopkg.in/bblfsh/sdk.v1/uast/transformer/positioner"
)

// Transformers is the of list `transformer.Transfomer` to apply to a UAST, to
// learn more about the Transformers and the available ones take a look to:
// https://godoc.org/gopkg.in/bblfsh/sdk.v1/uast/transformers
var Transformers = []transformer.Tranformer{
	annotatter.NewAnnotatter(AnnotationRules),
	positioner.NewFillOffsetFromLineCol(),
}

// Nodes doc:
// https://github.com/whitequark/parser/blob/master/doc/AST_FORMAT.md

// AnnotationRules describes how a UAST should be annotated with `uast.Role`.
//
// https://godoc.org/gopkg.in/bblfsh/sdk.v1/uast/ann
//var AnnotationRules = On(Any).Self(
var AnnotationRules = On(HasInternalRole("module")).Roles(uast.Module, uast.File).Descendants(
	On(rubyast.Begin).Roles(uast.Block),
	//On(rubyast.Module).Roles(uast.File, uast.Module).Descendants(
		//On(rubyast.Begin).Roles(uast.Block),
	//),
)

// Identifiers:
// lvasgn.target/token
// ivasgn.target/token
// ivar.token
// target.send.selector
// send.selector / send.token (note, sometimes is "foo=")
// self => needs synth token "self"
// csend.token

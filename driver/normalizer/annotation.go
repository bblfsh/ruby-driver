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

var isSomeOperator = Or(HasToken("+"), HasToken("-"), HasToken("*"), HasToken("/"),
	HasToken("%"), HasToken("**"), HasToken("=="), HasToken("!="), HasToken("!"),
	HasToken("<=>"), HasToken("==="), HasToken("eql?"), HasToken("equal?"),
	HasToken("<="), HasToken(">="), rubyast.And, rubyast.Or,
)

// Nodes doc:
// https://github.com/whitequark/parser/blob/master/doc/AST_FORMAT.md

// AnnotationRules describes how a UAST should be annotated with `uast.Role`.
//
// https://godoc.org/gopkg.in/bblfsh/sdk.v1/uast/ann
//var AnnotationRules = On(Any).Self(
var AnnotationRules = On(Any).Roles(uast.Module, uast.File).Descendants(
	On(Or(rubyast.Begin, rubyast.Block)).Roles(uast.Block).Children(
		On(rubyast.Body).Roles(uast.Body),
	),
	// *Asgn with two children = binary and value have the "Right" role but with a single children = multiple assignment target :-/
	On(rubyast.LVAsgn).Roles(uast.Assignment, uast.Binary, uast.Identifier, uast.Left).Children(
		On(HasInternalRole("value")).Roles(uast.Right),
	),
	// is also member
	On(rubyast.IVAsgn).Roles(uast.Assignment, uast.Binary, uast.Identifier, uast.Left, uast.Incomplete).Children(
		On(HasInternalRole("value")).Roles(uast.Right),
	),
	On(rubyast.GVAsgn).Roles(uast.Assignment, uast.Binary, uast.Identifier, uast.Left, uast.Visibility, uast.World).Children(
		On(HasInternalRole("value")).Roles(uast.Right),
	),
	// constant
	On(rubyast.CAsgn).Roles(uast.Assignment, uast.Binary, uast.Identifier, uast.Left).Children(
		On(HasInternalRole("value")).Roles(uast.Right),
	),
	// class
	On(rubyast.CVAsgn).Roles(uast.Assignment, uast.Binary, uast.Identifier, uast.Left).Children(
		On(HasInternalRole("value")).Roles(uast.Binary, uast.Right),
	),
	// instance (member)
	On(rubyast.IVAsgn).Roles(uast.Assignment, uast.Binary, uast.Identifier, uast.Left).Children(
		On(HasInternalRole("value")).Roles(uast.Binary, uast.Right),
	),
	// Multiple assignment; second element (whatever it is) must have the "Right" role
	On(rubyast.MAsgn).Roles(uast.Assignment, uast.Incomplete).Children(
		On(HasInternalRole("values")).Roles(uast.Binary, uast.Right),
	),
	On(rubyast.MultipleLeftSide).Roles(uast.Left, uast.Incomplete),

	// Types
	On(rubyast.Module).Roles(uast.Module, uast.Identifier),
	On(rubyast.Int).Roles(uast.Literal, uast.Number, uast.Primitive),
	On(rubyast.Str).Roles(uast.Literal, uast.String, uast.Primitive),
	On(rubyast.Pair).Roles(uast.Literal, uast.Tuple, uast.Primitive),
	On(rubyast.Array).Roles(uast.Literal, uast.List, uast.Primitive),
	On(rubyast.Hash).Roles(uast.Literal, uast.Map, uast.Primitive),
	On(rubyast.KwSplat).Roles(uast.Incomplete),

	// splat (*a)
	On(rubyast.Splat).Roles(uast.Identifier, uast.Incomplete),

	// local var ::var
	On(rubyast.LVar).Roles(uast.Identifier),
	// instance var  @var
	On(rubyast.IVar).Roles(uast.Identifier, uast.Visibility, uast.Instance),
	// global var $var
	On(rubyast.GVar).Roles(uast.Identifier, uast.Visibility, uast.World),
	// class var @@var
	On(rubyast.CVar).Roles(uast.Identifier, uast.Visibility, uast.Type),

	// Augmented assignment (op-asgn)
	On(rubyast.OpAsgn).Roles(uast.Operator, uast.Binary, uast.Assignment).Self(
		On(HasProperty("operator", "+")).Roles(uast.Arithmetic, uast.Add),
		On(HasProperty("operator", "-")).Roles(uast.Arithmetic, uast.Substract),
		On(HasProperty("operator", "*")).Roles(uast.Arithmetic, uast.Multiply),
		On(HasProperty("operator", "/")).Roles(uast.Arithmetic, uast.Divide),
		On(HasProperty("operator", "%")).Roles(uast.Arithmetic, uast.Modulo),
		// Pow
		On(HasProperty("operator", "**")).Roles(uast.Arithmetic, uast.Incomplete),
		On(HasProperty("operator", "&")).Roles(uast.Bitwise, uast.And),
		On(HasProperty("operator", "|")).Roles(uast.Bitwise, uast.Or),
		On(HasProperty("operator", "^")).Roles(uast.Bitwise, uast.Xor),
		// Complement
		On(HasProperty("operator", "~")).Roles(uast.Bitwise, uast.Incomplete),
		On(HasProperty("operator", "<<")).Roles(uast.Bitwise, uast.LeftShift),
		On(HasProperty("operator", ">>")).Roles(uast.Bitwise, uast.RightShift),
	).Children(
		On(HasInternalRole("target")).Roles(uast.Binary, uast.Left),
		On(HasInternalRole("value")).Roles(uast.Binary, uast.Right),
	),

	// Or Assign (a ||= b), And Assign (a &&= b)
	On(rubyast.AndAsgn).Roles(uast.Operator, uast.And, uast.Bitwise).Children(
		On(HasInternalRole("target")).Roles(uast.Binary, uast.Left),
		On(HasInternalRole("value")).Roles(uast.Binary, uast.Right),
	),
	On(rubyast.OrAsgn).Roles(uast.Operator, uast.Or, uast.Bitwise).Children(
		On(HasInternalRole("target")).Roles(uast.Binary, uast.Left),
		On(HasInternalRole("value")).Roles(uast.Binary, uast.Right),
	),

	// a.b.c.d would generate the tree d=->c->b->a where "a", "b" and "c" will be
	// Qualified+Identifier and "d" will be just Identifier.

	// send is used for qualified identifiers (foo.bar), method calls (puts "foo")
	// and a lot of other things...
	On(rubyast.Send).Self(
		On(And(HasInternalRole("base"),
			Not(isSomeOperator), Not(HasToken("continue")),
			Not(HasInternalRole("condition")))).Roles(uast.Qualified, uast.Identifier),

		On(HasChild(HasInternalRole("base"))).Roles(uast.Identifier),

		On(And(Or(rubyast.BodyRole, HasInternalRole("module")), Not(HasToken("continue")),
			Not(isSomeOperator))).Roles(uast.Call, uast.Identifier).Children(
			On(rubyast.Values).Roles(uast.Argument, uast.Identifier),
		),

		On(HasInternalRole("blockdata")).Self(
			On(HasToken("each")).Roles(uast.For, uast.Iterator),
			On(HasToken("lambda")).Roles(uast.Declaration, uast.Function, uast.Anonymous),
		),

		On(isSomeOperator).Roles(uast.Operator).Self(
			On(HasToken("+")).Roles(uast.Arithmetic, uast.Add),
			On(HasToken("-")).Roles(uast.Arithmetic, uast.Substract),
			On(HasToken("*")).Roles(uast.Arithmetic, uast.Multiply),
			On(HasToken("/")).Roles(uast.Arithmetic, uast.Divide),
			On(HasToken("%")).Roles(uast.Arithmetic, uast.Modulo),
			// Pow
			On(HasToken("**")).Roles(uast.Arithmetic, uast.Incomplete),
			On(HasToken("&")).Roles(uast.Bitwise, uast.And),
			On(HasToken("|")).Roles(uast.Bitwise, uast.Or),
			On(HasToken("^")).Roles(uast.Bitwise, uast.Xor),
			// Complemen
			On(HasToken("~")).Roles(uast.Bitwise, uast.Incomplete),
			On(HasToken("<<")).Roles(uast.Bitwise, uast.LeftShift),
			On(HasToken(">>")).Roles(uast.Bitwise, uast.RightShift),
			On(HasToken("==")).Roles(uast.Relational, uast.Equal),
			On(HasToken(">=")).Roles(uast.Relational, uast.GreaterThanOrEqual),
			On(HasToken("<=")).Roles(uast.Relational, uast.LessThanOrEqual),
			On(HasToken("!=")).Roles(uast.Relational, uast.Equal, uast.Not),
			On(HasToken("!")).Roles(uast.Relational, uast.Not),
			// Incomplete: check type (1 !eql? 1.0) but not being the same object like equal?
			On(HasToken("eql?")).Roles(uast.Relational, uast.Identical, uast.Incomplete),
			On(HasToken("equal?")).Roles(uast.Relational, uast.Identical, uast.Identical),
			// Combined comparison operator
			On(HasToken("<==>")).Roles(uast.Relational, uast.Incomplete),
		),

		On(HasToken("continue")).Roles(uast.Statement, uast.Continue),
	),

	// FIXME: needs Range role or similar
	On(Or(rubyast.IFlipFlop, rubyast.EFlipFlop)).Roles(uast.Incomplete, uast.List).Children(
		On(Any).Roles(uast.Identifier, uast.Incomplete),
	),
	On(rubyast.ERange).Roles(uast.Tuple, uast.Incomplete),
	On(rubyast.IRange).Roles(uast.Tuple, uast.Incomplete),
	On(rubyast.RegExp).Roles(uast.Expression, uast.Regexp),
	On(rubyast.RegExpBackRef).Roles(uast.Regexp, uast.Incomplete),
	On(rubyast.RegExpRef).Roles(uast.Regexp, uast.Incomplete),
	On(rubyast.RegOpt).Roles(uast.Regexp, uast.Incomplete),
	On(rubyast.Options).Roles(uast.Regexp, uast.Incomplete),
	On(rubyast.Symbol).Roles(uast.Identifier),
	On(rubyast.Sym).Roles(uast.Identifier),
	On(rubyast.Const).Roles(uast.Identifier, uast.Incomplete).Children(
		On(rubyast.CBase).Roles(uast.Identifier, uast.Qualified, uast.Incomplete),
	),
	// Interpolated symbols on strings
	On(rubyast.DSym).Roles(uast.String, uast.Incomplete),
	On(rubyast.Self).Roles(uast.This, uast.Left),

	On(HasInternalRole("condition")).Roles(uast.Condition),
	On(rubyast.If).Roles(uast.If, uast.Statement).Children(
		On(HasInternalRole("body")).Roles(uast.If, uast.Then),
		On(HasInternalRole("condition")).Roles(uast.If),
		On(HasInternalRole("else")).Roles(uast.If, uast.Else),
	),

	On(rubyast.Class).Roles(uast.Type, uast.Declaration, uast.Identifier).Children(
		On(And(rubyast.Block, HasInternalRole("body"))).Roles(uast.Body),
	),
	// Singleton class
	On(rubyast.SClass).Roles(uast.Type, uast.Declaration, uast.Incomplete),

	// Arguments grouping node, needs uast.Group or similar
	On(rubyast.Args).Roles(uast.Argument, uast.Incomplete).Children(
		On(rubyast.Arg).Roles(uast.Argument, uast.Name),
		On(rubyast.KwArg).Roles(uast.Argument, uast.Name, uast.Map),
		On(rubyast.KwOptArg).Roles(uast.Argument, uast.Name, uast.Incomplete),
		On(rubyast.KwRestArg).Roles(uast.Argument, uast.Incomplete).Self(
			On(Not(HasToken(""))).Roles(uast.Identifier),
		),
		On(rubyast.OptArg).Roles(uast.Argument, uast.Name).Children(
			On(Any).Roles(uast.Argument, uast.Default),
		),
	),
	On(rubyast.Alias).Roles(uast.Alias),
	On(rubyast.Def).Roles(uast.Function, uast.Declaration, uast.Identifier).Children(),
	// Singleton method
	On(rubyast.Defs).Roles(uast.Function, uast.Declaration, uast.Identifier, uast.Incomplete).Children(),
	On(rubyast.NilClass).Roles(uast.Type, uast.Null),
	On(Or(rubyast.Until, rubyast.UntilPost)).Roles(uast.Incomplete), // Complete annotations below
	On(Or(rubyast.Until, rubyast.UntilPost, rubyast.While, rubyast.WhilePost)).Roles(uast.While).Children(
		On(HasInternalRole("body")).Roles(uast.While, uast.Body),
		On(HasInternalRole("condition")).Roles(uast.While),
	),

	On(rubyast.For).Roles(uast.Statement, uast.For).Children(
		On(HasInternalRole("body")).Roles(uast.For, uast.Body),
		On(HasInternalRole("iterated")).Roles(uast.For, uast.Update),
		On(HasInternalRole("iterators")).Roles(uast.For, uast.Iterator),
	),

	On(rubyast.True).Roles(uast.Boolean, uast.Literal),
	On(rubyast.False).Roles(uast.Boolean, uast.Literal),
	On(rubyast.And).Roles(uast.Binary, uast.Expression, uast.Operator, uast.Boolean, uast.And),
	On(rubyast.Or).Roles(uast.Binary, uast.Expression, uast.Operator, uast.Boolean, uast.Or),
	On(HasInternalRole("left")).Roles(uast.Expression, uast.Left),
	On(HasInternalRole("right")).Roles(uast.Expression, uast.Right),
	On(HasToken("raise")).Roles(uast.Throw),

	// Exceptions
	On(rubyast.KwBegin).Roles(uast.Block).Self(
		On(Or(HasChild(rubyast.Rescue), HasChild(rubyast.Ensure))).Roles(uast.Try).Children(
			On(rubyast.Rescue).Roles(uast.Try, uast.Body).Children(
				On(rubyast.ResBody).Roles(uast.Catch).Children(
					On(rubyast.Retry).Roles(uast.Statement, uast.Call),
				),
			),
			On(rubyast.Ensure).Roles(uast.Finally, uast.Body),
		),
	),

	On(rubyast.Case).Roles(uast.Switch).Children(
		On(rubyast.When).Roles(uast.Case),
	),
	On(rubyast.Break).Roles(uast.Break),
	On(rubyast.Undef).Roles(uast.Statement, uast.Incomplete),
)

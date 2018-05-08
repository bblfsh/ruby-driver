package normalizer

import (
	"gopkg.in/bblfsh/sdk.v2/uast"
	"gopkg.in/bblfsh/sdk.v2/uast/role"
	. "gopkg.in/bblfsh/sdk.v2/uast/transformer"
	"gopkg.in/bblfsh/sdk.v2/uast/transformer/positioner"
)

var Native = Transformers([][]Transformer{
	{
		ResponseMetadata{
			TopLevelIsRootNode: true,
		},
	},
	{Mappings(Annotations...)},
	{RolesDedup()},
}...)

var Code = []CodeTransformer{
	positioner.NewFillOffsetFromLineCol(),
}

// FIXME: move to the SDK and remove from here and the python driver
func annotateTypeToken(typ, token string, roles ...role.Role) Mapping {
	return AnnotateType(typ,
		FieldRoles{
			uast.KeyToken: {Add: true, Op: String(token)},
		}, roles...)
}

// FIXME: move to the SDK and remove from here and the python driver
func mapInternalProperty(key string, roles ...role.Role) Mapping {
	return Map(key,
		Part("other", Obj{
			key: ObjectRoles(key),
		}),
		Part("other", Obj{
			key: ObjectRoles(key, roles...),
		}),
	)
}

// Nodes doc:
// https://github.com/whitequark/parser/blob/master/doc/AST_FORMAT.md

//var isSomeOperator = Or(HasToken("+"), HasToken("-"), HasToken("*"), HasToken("/"),
	//HasToken("%"), HasToken("**"), HasToken("=="), HasToken("!="), HasToken("!"),
	//HasToken("<=>"), HasToken("==="), HasToken("eql?"), HasToken("equal?"),
	//HasToken("<="), HasToken(">="), rubyast.And, rubyast.Or,
//)

var Annotations = []Mapping{
	ObjectToNode{
		LineKey:   "pos_line_start",
		ColumnKey: "pos_col_start",
	}.Mapping(),
	ObjectToNode{
		EndLineKey:   "pos_line_end",
		EndColumnKey: "pos_col_end",
	}.Mapping(),

	AnnotateType("file", nil, role.File),
	// XXX token
	AnnotateType("body", nil, role.Body),
	mapInternalProperty("body", role.Body),
	// XXX check that these really work
	mapInternalProperty("left", role.Left),
	mapInternalProperty("right", role.Right),
	mapInternalProperty("condition", role.Expression, role.Condition),
	mapInternalProperty("target", role.Binary, role.Left),
	mapInternalProperty("value", role.Binary, role.Right),
	mapInternalProperty("_1", role.Tuple, role.Value),
	mapInternalProperty("_2", role.Tuple, role.Value),

	// Types
	// XXX tokens
	AnnotateType("module", nil, role.Statement, role.Module, role.Identifier),
	AnnotateType("block", nil, role.Block),
	AnnotateType("int", nil, role.Expression, role.Literal, role.Number, role.Primitive),
	AnnotateType("str", nil, role.Expression, role.Literal, role.String, role.Primitive),
	AnnotateType("pair", nil, role.Expression, role.Literal, role.Tuple, role.Primitive),
	AnnotateType("array", nil, role.Expression, role.Literal, role.List, role.Primitive),
	AnnotateType("hash", nil, role.Expression, role.Literal, role.Map, role.Primitive),
	AnnotateType("class", nil, role.Statement, role.Type, role.Declaration, role.Identifier),

	// splats (*a)
	AnnotateType("kwsplat", nil, role.Expression, role.Incomplete),
	AnnotateType("splat", nil, role.Expression, role.Identifier, role.Incomplete),

	// Vars
	// local
	AnnotateType("lvar", nil, role.Expression, role.Identifier),
	// instance
	AnnotateType("ivar", nil, role.Expression, role.Identifier, role.Visibility, role.Instance),
	// global
	AnnotateType("gvar", nil, role.Expression, role.Identifier, role.Visibility, role.World),
	// class
	AnnotateType("cvar", nil, role.Expression, role.Identifier, role.Visibility, role.Type),

	// Singleton class
	AnnotateType("sclass", nil, role.Expression, role.Type, role.Declaration, role.Incomplete),

	AnnotateType("alias", nil, role.Statement, role.Alias),
	AnnotateType("def", nil, role.Statement, role.Function, role.Declaration, role.Identifier),
	// Singleton method
	AnnotateType("defs", nil, role.Statement, role.Function, role.Declaration, role.Identifier, role.Incomplete),
	AnnotateType("NilClass", nil, role.Statement, role.Type, role.Null),
	AnnotateType("break", nil, role.Statement, role.Break),
	AnnotateType("undef", nil, role.Statement, role.Incomplete),
	AnnotateType("case", nil, role.Statement, role.Switch),
	AnnotateType("when", nil, role.Expression, role.Case),

	// Exceptions
	AnnotateType("kwbegin", nil, role.Expression, role.Block),
	AnnotateType("rescue", nil, role.Expression, role.Try, role.Body),
	AnnotateType("resbody", nil, role.Expression, role.Catch),
	AnnotateType("retry", nil, role.Expression, role.Statement, role.Call, role.Incomplete),
	AnnotateType("ensure", nil, role.Expression, role.Finally),

	// Arguments
	// grouping node, need grouping role
	AnnotateType("args", nil, role.Expression, role.Argument, role.Incomplete),
	AnnotateType("kwarg", nil, role.Expression, role.Argument, role.Name, role.Map),
	AnnotateType("kwoptarg", nil, role.Expression, role.Argument, role.Name, role.Incomplete),
	AnnotateType("kwrestarg", nil, role.Expression, role.Argument, role.Identifier, role.Incomplete),
	AnnotateType("optarg", nil, role.Expression, role.Argument, role.Name),

	// Assigns
	// *Asgn with two children = binary and value have the "Right" role but with a single children = multiple assignment target :-/
	AnnotateType("lvasgn", nil, role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	// is also a member
	AnnotateType("ivasgn", nil, role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	AnnotateType("gvasgn", nil, role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left,
	// constant assign
	AnnotateType("casgn", nil, role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	// class assign
	AnnotateType("cvasgn", nil, role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	// instance member
	AnnotateType("ivasgn", nil, role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	// multiple
	AnnotateType("masgn", nil, role.Expression, role.Assignment, role.Incomplete),
	// Or Assign (a ||= b), And Assign (a &&= b)
	AnnotateType("and_asgn", nil, role.Expression, role.Operator, role.And, role.Bitwise),
	AnnotateType("or_asgn", nil, role.Expression, role.Operator, role.Or, role.Bitwise),

	// Misc
	// multiple left side
	AnnotateType("mlhs", nil, role.Left, role.Incomplete),
	AnnotateType("erange", nil, role.Expression, role.Tuple, role.Incomplete),
	AnnotateType("irange", nil, role.Expression, role.Tuple, role.Incomplete),
	AnnotateType("regexp", nil, role.Expression, role.Regexp),
	// regexp back reference
	AnnotateType("back_ref", nil, role.Expression, role.Regexp, role.Incomplete),
	// regexp reference
	AnnotateType("nth_ref", nil, role.Expression, role.Regexp, role.Incomplete),
	// regexp option/s
	AnnotateType("regopt", nil, role.Expression, role.Regexp, role.Incomplete),
	AnnotateType("options", nil, role.Expression, role.Regexp, role.Incomplete),

	AnnotateType("Symbol", nil, role.Expression, role.Identifier),
	AnnotateType("sym", nil, role.Expression, role.Identifier),
	// Interpolated symbols on strings
	AnnotateType("dsym", nil, role.Expression, role.String, role.Incomplete),
	AnnotateType("self", nil, role.Expression, role.This, role.Left),
	annotateTypeToken("true", "true", role.Expression, role.Boolean, role.Literal),
	annotateTypeToken("false", "false", role.Expression, role.Boolean, role.Literal),
	annotateTypeToken("and", "and", role.Expression, role.Binary, role.Operator, role.Boolean, role.And),
	annotateTypeToken("or", "or", role.Expression, role.Binary, role.Operator, role.Boolean, role.Or),
	annotateTypeToken("raise", "raise", role.Statement, role.Throw),

	AnnotateType("const", nil, role.Expression, role.Identifier, role.Incomplete),
	AnnotateType("cbase", nil, role.Expression, role.Identifier, role.Qualified, role.Incomplete),

/*
	// Augmented assignment (op-asgn)
	On(rubyast.OpAsgn).Roles(uast.Expression, uast.Operator, uast.Binary, uast.Assignment).Self(
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
	)

	// a.b.c.d would generate the tree d=->c->b->a where "a", "b" and "c" will be
	// Qualified+Identifier and "d" will be just Identifier.

	// send is used for qualified identifiers (foo.bar), method calls (puts "foo")
	// and a lot of other things...
	On(rubyast.Send).Self(
		On(And(HasInternalRole("base"),
			Not(isSomeOperator), Not(HasToken("continue")),
			Not(HasInternalRole("condition")))).Roles(uast.Expression, uast.Qualified, uast.Identifier),

		On(HasChild(HasInternalRole("base"))).Roles(uast.Expression, uast.Identifier),

		On(And(Or(rubyast.BodyRole, HasInternalRole("module")), Not(HasToken("continue")),
			Not(isSomeOperator))).Roles(uast.Expression, uast.Call, uast.Identifier).Children(
			On(rubyast.Values).Roles(uast.Expression, uast.Argument, uast.Identifier),
		),

		On(HasInternalRole("blockdata")).Self(
			On(HasToken("each")).Roles(uast.Statement, uast.For, uast.Iterator),
			On(HasToken("lambda")).Roles(uast.Expression, uast.Declaration, uast.Function, uast.Anonymous),
		),

		On(isSomeOperator).Roles(uast.Expression, uast.Operator).Self(
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
	On(Or(rubyast.IFlipFlop, rubyast.EFlipFlop)).Roles(uast.Expression, uast.Incomplete, uast.List).Children(
		On(Any).Roles(uast.Identifier, uast.Incomplete),
	),

	On(rubyast.If).Roles(uast.Statement, uast.If).Children(
		On(HasInternalRole("body")).Roles(uast.Expression, uast.If, uast.Then),
		On(HasInternalRole("condition")).Roles(uast.Expression, uast.If),
		On(HasInternalRole("else")).Roles(uast.Expression, uast.If, uast.Else),
	),

	// Singleton method
	On(Or(rubyast.Until, rubyast.UntilPost)).Roles(uast.Incomplete), // Complete annotations below
	On(Or(rubyast.Until, rubyast.UntilPost, rubyast.While, rubyast.WhilePost)).Roles(uast.Statement, uast.While).Children(
		On(HasInternalRole("body")).Roles(uast.Expression, uast.While, uast.Body),
		On(HasInternalRole("condition")).Roles(uast.Expression, uast.While),
	),

	On(rubyast.For).Roles(uast.Statement, uast.For).Children(
		On(HasInternalRole("body")).Roles(uast.Expression, uast.For, uast.Body),
		On(HasInternalRole("iterated")).Roles(uast.Expression, uast.For, uast.Update),
		On(HasInternalRole("iterators")).Roles(uast.Expression, uast.For, uast.Iterator),
	),
)
*/
}

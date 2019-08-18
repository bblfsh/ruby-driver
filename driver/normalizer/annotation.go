package normalizer

import (
	"github.com/bblfsh/sdk/v3/uast"
	"github.com/bblfsh/sdk/v3/uast/role"
	. "github.com/bblfsh/sdk/v3/uast/transformer"
)

var Native = Transformers([][]Transformer{
	{Mappings(Annotations...)},
	{RolesDedup()},
}...)

// FIXME: move to the SDK and remove from here and the python driver
func annotateTypeToken(typ, token string, roles ...role.Role) Mapping {
	return AnnotateType(typ,
		FieldRoles{
			uast.KeyToken: {Add: true, Op: String(token)},
		}, roles...)
}

func annotateTypeTokenField(typ, tokenfield string, roles ...role.Role) Mapping {
	return AnnotateType(typ, FieldRoles{
		tokenfield: {Rename: uast.KeyToken},
	}, roles...)
}

// FIXME: move to the SDK and remove from here and the python driver
func mapInternalProperty(key string, roles ...role.Role) Mapping {
	return Map(
		Part("other", Obj{key: ObjectRoles(key)}),
		Part("other", Obj{key: ObjectRoles(key, roles...)}),
	)
}

// Nodes doc:
// https://github.com/whitequark/parser/blob/master/doc/AST_FORMAT.md

var operatorRoles = StringToRolesMap(map[string][]role.Role{
	"+": {role.Arithmetic, role.Add},
	"-": {role.Arithmetic, role.Substract},
	"*": {role.Arithmetic, role.Multiply},
	"/": {role.Arithmetic, role.Divide},
	"%": {role.Arithmetic, role.Modulo},
	// pow
	"**": {role.Arithmetic, role.Incomplete},
	"&":  {role.Bitwise, role.And},
	"|":  {role.Bitwise, role.Or},
	"^":  {role.Bitwise, role.Xor},
	// Complement
	"~":  {role.Bitwise, role.Incomplete},
	"<<": {role.Bitwise, role.LeftShift},
	">>": {role.Bitwise, role.RightShift},
	"==": {role.Equal, role.Relational},
	"<=": {role.LessThanOrEqual, role.Relational},
	">=": {role.GreaterThanOrEqual, role.Relational},
	"!=": {role.Equal, role.Not, role.Relational},
	"!":  {role.Not, role.Relational},
	// Incomplete: check type (1 !eql? 1.0) but not being the same object like equal?
	"eql?":   {role.Identical, role.Relational},
	"equal?": {role.Identical, role.Relational},
	// rocket ship operator
	"===":  {role.Identical, role.Relational},
	"<==>": {role.Identical, role.Incomplete},
})

var Annotations = []Mapping{
	AnnotateType("file", nil, role.File),
	AnnotateType("begin", nil, role.Block),
	AnnotateType("body", nil, role.Body),
	mapInternalProperty("body", role.Body),
	mapInternalProperty("left", role.Left),
	mapInternalProperty("right", role.Right),
	mapInternalProperty("condition", role.Expression, role.Condition),
	mapInternalProperty("target", role.Binary, role.Left),
	mapInternalProperty("value", role.Binary, role.Right),
	mapInternalProperty("_1", role.Tuple, role.Value),
	mapInternalProperty("_2", role.Tuple, role.Value),

	// Types
	AnnotateType("module", nil, role.Statement, role.Module),
	AnnotateType("comment", nil, role.Noop, role.Comment),
	AnnotateType("module", nil, role.Statement, role.Module, role.Identifier),
	AnnotateType("block", nil, role.Block),
	AnnotateType("int", nil, role.Expression, role.Literal, role.Number, role.Primitive),
	AnnotateType("NilNode", nil, role.Null),
	AnnotateType("nil", nil, role.Null),
	AnnotateType("return", nil, role.Statement, role.Return),
	AnnotateType("float", nil, role.Expression, role.Literal, role.Number, role.Primitive),
	AnnotateType("complex", nil, role.Expression, role.Literal, role.Number, role.Primitive, role.Incomplete),
	AnnotateType("rational", nil, role.Expression, role.Literal, role.Number, role.Primitive, role.Incomplete),
	AnnotateType("str", nil, role.Expression, role.Literal, role.String, role.Primitive),
	AnnotateType("xstr", nil, role.Expression, role.Literal, role.String, role.Block, role.Incomplete),
	AnnotateType("dstr", nil, role.Expression, role.String, role.Block, role.Incomplete),
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
	AnnotateType("when", nil, role.Expression, role.Switch, role.Case),
	AnnotateType("super", nil, role.Expression, role.Call, role.Base),
	AnnotateType("zsuper", nil, role.Expression, role.Call, role.Base),
	AnnotateType("yield", nil, role.Return, role.Incomplete),

	// Exceptions
	AnnotateType("kwbegin", nil, role.Expression, role.Block),
	AnnotateType("rescue", nil, role.Expression, role.Try, role.Body),
	AnnotateType("resbody", nil, role.Expression, role.Catch),
	AnnotateType("retry", nil, role.Expression, role.Statement, role.Call, role.Incomplete),
	AnnotateType("ensure", nil, role.Expression, role.Finally),

	// Arguments
	// grouping node for function definition (not for calls which just use send.values), need grouping role
	AnnotateType("args", nil, role.Expression, role.Argument, role.Incomplete),
	AnnotateType("arg", nil, role.Expression, role.Argument, role.Name, role.Identifier),
	AnnotateType("blockarg", nil, role.Expression, role.Argument, role.Name, role.Identifier),
	AnnotateType("kwarg", nil, role.Expression, role.Argument, role.Name, role.Map),
	AnnotateType("optarg", nil, role.Expression, role.Argument, role.Name, role.Default),
	AnnotateType("kwoptarg", nil, role.Expression, role.Argument, role.Name, role.Incomplete),
	AnnotateType("restarg", nil, role.Expression, role.Argument, role.Identifier, role.List),
	AnnotateType("kwrestarg", nil, role.Expression, role.Argument, role.Identifier, role.Incomplete),

	// Assigns
	// constant assign
	annotateTypeTokenField("casgn", "selector", role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	// multiple
	AnnotateType("masgn", nil, role.Expression, role.Assignment, role.Incomplete),
	// *Asgn with two children = binary and value have the "Right" role but with a single children = multiple assignment target :-/
	annotateTypeTokenField("lvasgn", "target", role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	annotateTypeTokenField("gvasgn", "target", role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	AnnotateType("gvasgn", nil, role.Expression, role.Assignment, role.Binary),
	// class assign
	annotateTypeTokenField("cvasgn", "target", role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	// instance member
	annotateTypeTokenField("ivasgn", "target", role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
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

	AnnotateType("values", nil, role.Expression, role.Argument, role.Identifier),

	// For
	AnnotateType("for", ObjRoles{
		"body":      {role.Expression, role.For, role.Body},
		"iterated":  {role.Expression, role.For, role.Update},
		"iterators": {role.Expression, role.For, role.Iterator},
	}, role.Statement, role.For),

	// While/Until
	AnnotateType("while", nil, role.Statement, role.While),
	AnnotateType("while_post", nil, role.Statement, role.While),
	AnnotateType("until", nil, role.Statement, role.While),
	AnnotateType("until_post", nil, role.Statement, role.While),

	// If
	AnnotateType("if", ObjRoles{
		"body": {role.Expression, role.Then},
		"else": {role.Expression, role.Else},
	}, role.Statement, role.If),

	AnnotateTypeCustom("op_asgn", MapObj(
		Obj{
			"operator": Var("op"),
		},
		Fields{
			{Name: "operator", Op: Operator("op", operatorRoles, role.Binary)},
		}),
		LookupArrOpVar("op", operatorRoles),
		role.Expression, role.Binary, role.Assignment, role.Operator),

	AnnotateType("iflipflop", nil, role.Expression, role.List, role.Incomplete),
	AnnotateType("flip_1", nil, role.Identifier, role.Value, role.Incomplete),
	AnnotateType("flip_2", nil, role.Identifier, role.Value, role.Incomplete),

	// The many faces of Ruby's "send" start here ===>
	AnnotateType("send_statement", MapObj(
		Obj{
			"selector": String("continue"),
		},
		Obj{
			uast.KeyToken: String("continue"),
		}), role.Statement, role.Continue),

	AnnotateType("send_statement", MapObj(
		Obj{
			"selector": String("lambda"),
		},
		Obj{
			uast.KeyToken: String("lambda"),
		}), role.Expression, role.Declaration, role.Function, role.Anonymous),

	AnnotateType("send_require", MapObj(
		Obj{
			"selector": String("require"),
		},
		Obj{
			uast.KeyToken: String("require"),
		}), role.Expression, role.Import),

	AnnotateType("send_statement", MapObj(
		Obj{
			"selector": String("each"),
		},
		Obj{
			uast.KeyToken: String("each"),
		}), role.Statement, role.For, role.Iterator),

	AnnotateType("send_statement", MapObj(
		Obj{
			"selector": String("public"),
		},
		Obj{
			uast.KeyToken: String("public"),
		}), role.Statement, role.Visibility, role.World),

	AnnotateType("send_statement", MapObj(
		Obj{
			"selector": String("protected"),
		},
		Obj{
			uast.KeyToken: String("protected"),
		}), role.Statement, role.Visibility, role.Subtype),

	AnnotateType("send_statement", MapObj(
		Obj{
			"selector": String("private"),
		},
		Obj{
			uast.KeyToken: String("private"),
		}), role.Statement, role.Visibility, role.Instance),

	// Operator expression "send"
	AnnotateTypeCustom("send_operator", MapObj(
		Obj{
			"base":     ObjectRoles("bs"),
			"values":   EachObjectRoles("values"),
			"selector": Var("selector"),
		},
		Obj{
			"base":        ObjectRoles("bs", role.Left),
			"values":      EachObjectRoles("values", role.Right),
			uast.KeyToken: Var("selector"),
		}),
		LookupArrOpVar("selector", operatorRoles),
		role.Expression, role.Binary, role.Operator),

	// Same without values (unary)
	AnnotateTypeCustom("send_operator", MapObj(
		Obj{
			"selector": Var("selector"),
		},
		Obj{
			uast.KeyToken: Var("selector"),
		}),
		LookupArrOpVar("selector", operatorRoles),
		role.Expression, role.Unary, role.Operator),

	// Assignment "send" (self.foo = 1)
	AnnotateType("send_assign", MapObj(
		Obj{
			"base":        Var("base"),
			"values":      EachObjectRoles("values"),
			"selector":    Var("sel"),
			uast.KeyToken: Var("tk"),
		},
		Obj{
			"base":        Var("base"),
			"values":      EachObjectRoles("values", role.Assignment, role.Right),
			"selector":    Var("sel"),
			uast.KeyToken: Var("tk"),
		}), role.Expression, role.Assignment, role.Left),

	// Qualified identifier "send" (other than the parent of the last one that will
	// match the rule above)
	AnnotateType("send_qualified", nil, role.Expression, role.Qualified, role.Identifier),
	AnnotateType("send_identifier", nil, role.Expression, role.Identifier),

	// Function call "send" with arguments
	AnnotateType("send_call", MapObj(
		Obj{
			"base":     Var("base"),
			"selector": Var("selector"),
			"values":   EachObjectRoles("values"),
		},
		Obj{
			"base":        Var("base"),
			"values":      EachObjectRoles("values", role.Function, role.Call, role.Argument),
			uast.KeyToken: Var("selector"),
		}), role.Expression, role.Function, role.Call),

	// Function call "send" without arguments
	AnnotateType("send_call", MapObj(
		Obj{
			"base":     Var("base"),
			"selector": Var("selector"),
		},
		Obj{
			"base":        Var("base"),
			uast.KeyToken: Var("selector"),
		}), role.Expression, role.Function, role.Call),

	AnnotateType("send_array", nil, role.Expression, role.List),
}

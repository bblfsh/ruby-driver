package normalizer

import (
	"strings"

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

func annotateTypeTokenField(typ, tokenfield string, roles ...role.Role) Mapping {
	return AnnotateType(typ, FieldRoles{
		tokenfield: {Rename: uast.KeyToken},
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

type opSendAssign struct {
	op        Op
}

func (op opSendAssign) Check(st *State, n uast.Node) (bool, error) {
	s, ok := n.(uast.String)
	if !ok {
		return false, nil
	}

	gostr := string(s)
	if !strings.HasSuffix(gostr, "=") {
		return false, nil
	}

	return op.op.Check(st, uast.String(gostr[:len(gostr)-1]))
}

func (op opSendAssign) Construct(st *State, n uast.Node) (uast.Node, error) {
	n, err := op.op.Construct(st, n)
	if err != nil {
		return nil, err
	}

	v, ok := n.(uast.String)
	if !ok {
		return nil, ErrExpectedValue.New(n)
	}

	gostr := string(v) + "="
	return uast.String(gostr), nil
}

type opSendOperator struct {
	op	Op
}

func (op opSendOperator) Check(st *State, n uast.Node) (bool, error) {
	s, ok := n.(uast.String)
	if !ok {
		return false, nil
	}

	if _, ok := operatorRoles[s]; !ok {
		return false, nil
	}

	return op.op.Check(st, n)
}

func (op opSendOperator) Construct(st *State, n uast.Node) (uast.Node, error) {
	return op.op.Construct(st, n)
}


// Nodes doc:
// https://github.com/whitequark/parser/blob/master/doc/AST_FORMAT.md

var	operatorRoles = StringToRolesMap(map[string][]role.Role{
	"+":   {role.Arithmetic, role.Add},
	"-":   {role.Arithmetic, role.Substract},
	"*":   {role.Arithmetic, role.Multiply},
	"/":   {role.Arithmetic, role.Divide},
	"%":   {role.Arithmetic, role.Modulo},
	// pow
	"**":  {role.Arithmetic, role.Incomplete},
	"&":   {role.Bitwise, role.And},
	"|":   {role.Bitwise, role.Or},
	"^":   {role.Bitwise, role.Xor},
	// Complement
	"~":   {role.Bitwise, role.Incomplete},
	"<<":  {role.Bitwise, role.LeftShift},
	">>":  {role.Bitwise, role.RightShift},
	"==":  {role.Equal, role.Relational},
	"<=":  {role.LessThanOrEqual, role.Relational},
	">=":  {role.GreaterThanOrEqual, role.Relational},
	"!=":  {role.Equal, role.Not, role.Relational},
	"!":   {role.Not, role.Relational},
	// Incomplete: check type (1 !eql? 1.0) but not being the same object like equal?
	"eql?":   {role.Identical, role.Relational},
	"equal?":   {role.Identical, role.Relational},
	// rocket ship operator
	"<==>":   {role.Identical, role.Incomplete},
})


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
	annotateTypeTokenField("module", "name", role.Statement, role.Module, role.Identifier),
	AnnotateType("block", nil, role.Block),
	annotateTypeTokenField("int", "token", role.Expression, role.Literal, role.Number, role.Primitive),
	annotateTypeTokenField("float", "token", role.Expression, role.Literal, role.Number, role.Primitive),
	annotateTypeTokenField("complex", "token", role.Expression, role.Literal, role.Number, role.Primitive, role.Incomplete),
	annotateTypeTokenField("rational", "token", role.Expression, role.Literal, role.Number, role.Primitive, role.Incomplete),
	annotateTypeTokenField("str", "token", role.Expression, role.Literal, role.String, role.Primitive),
	AnnotateType("pair", nil, role.Expression, role.Literal, role.Tuple, role.Primitive),
	AnnotateType("array", nil, role.Expression, role.Literal, role.List, role.Primitive),
	AnnotateType("hash", nil, role.Expression, role.Literal, role.Map, role.Primitive),
	annotateTypeTokenField("class", "name", role.Statement, role.Type, role.Declaration, role.Identifier),

	// splats (*a)
	AnnotateType("kwsplat", nil, role.Expression, role.Incomplete),
	AnnotateType("splat", nil, role.Expression, role.Identifier, role.Incomplete),

	// Vars
	// local
	annotateTypeTokenField("lvar", "token", role.Expression, role.Identifier),
	// instance
	annotateTypeTokenField("ivar", "token", role.Expression, role.Identifier, role.Visibility, role.Instance),
	// global
	annotateTypeTokenField("gvar", "token", role.Expression, role.Identifier, role.Visibility, role.World),
	// class
	annotateTypeTokenField("cvar", "token", role.Expression, role.Identifier, role.Visibility, role.Type),

	// Singleton class
	AnnotateType("sclass", nil, role.Expression, role.Type, role.Declaration, role.Incomplete),

	AnnotateType("alias", nil, role.Statement, role.Alias),
	annotateTypeTokenField("def", "name", role.Statement, role.Function, role.Declaration, role.Identifier),
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
	// grouping node for function definition (not for calls which just use send.values), need grouping role
	AnnotateType("args", nil, role.Expression, role.Argument, role.Incomplete),
	annotateTypeTokenField("arg", "token", role.Expression, role.Argument, role.Name, role.Identifier),
	annotateTypeTokenField("kwarg", "token", role.Expression, role.Argument, role.Name, role.Map),
	annotateTypeTokenField("kwoptarg", "token", role.Expression, role.Argument, role.Name, role.Incomplete),
	annotateTypeTokenField("restarg", "name", role.Expression, role.Argument, role.Identifier, role.List),
	annotateTypeTokenField("kwrestarg", "name", role.Expression, role.Argument, role.Identifier, role.Incomplete),

	// Assigns
	// constant assign
	annotateTypeTokenField("casgn", "selector", role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	// multiple
	AnnotateType("masgn", nil, role.Expression, role.Assignment, role.Incomplete),
	// *Asgn with two children = binary and value have the "Right" role but with a single children = multiple assignment target :-/
	annotateTypeTokenField("lvasgn", "target", role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	// is also a member
	annotateTypeTokenField("ivasgn", "target", role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	annotateTypeTokenField("gvasgn", "target", role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
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

	annotateTypeTokenField("Symbol", "token", role.Expression, role.Identifier),
	annotateTypeTokenField("sym", "token", role.Expression, role.Identifier),
	// Interpolated symbols on strings
	AnnotateType("dsym", nil, role.Expression, role.String, role.Incomplete),
	AnnotateType("self", nil, role.Expression, role.This, role.Left),
	annotateTypeToken("true", "true", role.Expression, role.Boolean, role.Literal),
	annotateTypeToken("false", "false", role.Expression, role.Boolean, role.Literal),
	annotateTypeToken("and", "and", role.Expression, role.Binary, role.Operator, role.Boolean, role.And),
	annotateTypeToken("or", "or", role.Expression, role.Binary, role.Operator, role.Boolean, role.Or),
	annotateTypeToken("raise", "raise", role.Statement, role.Throw),

	annotateTypeTokenField("const", "token", role.Expression, role.Identifier, role.Incomplete),
	AnnotateType("cbase", nil, role.Expression, role.Identifier, role.Qualified, role.Incomplete),

	AnnotateType("values", nil, role.Expression, role.Argument, role.Identifier),

	// For
	AnnotateType("for", ObjRoles{
		"body": {role.Expression, role.For, role.Body},
		"iterated": {role.Expression, role.For, role.Update},
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

	MapASTCustom("op_asgn",
		Obj{
			"operator": Var("op"),
		}, Fields{
			{Name: "operator", Op: Operator("op", operatorRoles, role.Binary)},
		},
		LookupArrOpVar("op", operatorRoles),
		role.Expression, role.Binary, role.Assignment, role.Operator),

	AnnotateType("iflipflop", ObjRoles{
		"_1": {role.Identifier, role.Incomplete},
		"_2": {role.Identical, role.Incomplete},
	}, role.Expression, role.List, role.Incomplete),

	// The many faces of Ruby's "send" start here ===>

	MapAST("send", Obj{
		"selector": String("continue"),
	}, Obj{
		"selector": String("continue"),
	}, role.Statement, role.Continue),

	MapAST("send", Obj{
		"selector": String("lambda"),
	}, Obj{
		"selector": String("lambda"),
	}, role.Expression, role.Declaration, role.Function, role.Anonymous),

	MapAST("send", Obj{
		"selector": String("require"),
	}, Obj{
		"selector": String("require"),
	}, role.Expression, role.Import),

	MapAST("send", Obj{
		"selector": String("each"),
	}, Obj{
		"selector": String("each"),
	}, role.Statement, role.For, role.Iterator),

	// Operator expression "send"
	MapASTCustom("send", Obj{
		"base": Check(Not(Is(nil)), ObjectRoles("base")),
		"values": Check(Not(Is(nil)), EachObjectRoles("values")),
		"selector": opSendOperator{op: Var("selector")},
	}, Obj{
		"base": ObjectRoles("base", role.Left),
		"values": EachObjectRoles("values", role.Right),
		uast.KeyToken: Var("selector"),
	},
	LookupArrOpVar("selector", operatorRoles),
	role.Expression, role.Binary, role.Operator),

	// Assignment "send"
	MapAST("send", Obj{
		"base": Var("base"),
		"values": EachObjectRoles("values"),
		"selector": opSendAssign{op: Var("selector")},
	}, Obj{
		"base": Var("base"),
		"values": EachObjectRoles("values", role.Assignment, role.Right),
		uast.KeyToken: Var("selector"),
	}, role.Expression, role.Assignment, role.Left),

	// Parent of the last element of the qualified identifier (annotates
	// it and the child which is a normal identifier)
	MapAST("send", Obj{
		"base": Obj{
			uast.KeyType: String("send"),
			uast.KeyStart: Var("childstart"),
			uast.KeyEnd: Var("childend"),
			"base": Check(Is(nil), Var("childbase")),
			"selector": Var("childselector"),
		},
		"selector": Var("selector"),
	}, Obj{
		"base": Obj{
			uast.KeyType: String("send"),
			uast.KeyStart: Var("childstart"),
			uast.KeyEnd: Var("childend"),
			"base": Var("childbase"),
			"selector": Var("childselector"),
			uast.KeyRoles: Roles(role.Identifier),
			"__notcall": Bool(true),
		},
		uast.KeyToken: Var("selector"),
	}, role.Expression, role.Qualified, role.Identical),

	// Qualified identifier "send" (other than the parent of the last one that will
	// match the rule above)
	MapAST("send", Obj{
		"base":     Check(Not(Is(nil)), Var("base")),
		"selector": Var("selector"),
	}, Obj{
		"base": Var("base"),
		uast.KeyToken: Var("selector"),
	}, role.Expression, role.Qualified, role.Identifier),

	// Function call "send" without arguments
	MapAST("send", Obj{
		"base":     Var("base"),
		"selector": Var("selector"),
		"__notcall": Check(Is(nil), Var("__notcall")),
	}, Obj{
		"base": Var("base"),
		uast.KeyToken: Var("selector"),
	}, role.Expression, role.Function, role.Call),

	// Function call "send" with arguments
	MapAST("send", Obj{
		"base":     Var("base"),
		"selector": Var("selector"),
		"values": EachObjectRoles("values"),
	}, Obj{
		"base": Var("base"),
		"values": EachObjectRoles("values", role.Function, role.Call, role.Argument),
		uast.KeyToken: Var("selector"),
	}, role.Expression, role.Function, role.Call),

}

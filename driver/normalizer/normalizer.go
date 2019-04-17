package normalizer

import (
	"github.com/bblfsh/sdk/v3/uast"
	"github.com/bblfsh/sdk/v3/uast/role"
	. "github.com/bblfsh/sdk/v3/uast/transformer"
	"github.com/bblfsh/sdk/v3/uast/transformer/positioner"
)

var Preprocess = Transformers([][]Transformer{
	{Mappings(Preprocessors...)},
}...)

var Normalize = Transformers([][]Transformer{
	{Mappings(Normalizers...)},
}...)

var Preprocessors = []Mapping{}

var PreprocessCode = []CodeTransformer{
	positioner.FromLineCol(),
}

func mapIdentifier(key string) Mapping {
	return MapSemantic(key, uast.Identifier{}, MapObj(
		Obj{uast.KeyToken: Var("val")},
		Obj{"Name": Var("val")},
	))
}

func tokenIsIdentifier(typ, tokenKey string, roles ...role.Role) Mapping {
	return AnnotateType(typ, MapObj(
		Fields{
			{Name: tokenKey, Op: Var("name")},
		},
		Fields{
			{Name: tokenKey, Op: UASTType(uast.Identifier{}, Obj{
				"Name": Var("name"),
			})},
		}),
		roles...)
}

func identifierWithPos(nameVar string) ObjectOp {
	return UASTType(uast.Identifier{}, Obj{
		uast.KeyPos: UASTType(uast.Positions{}, Obj{
			uast.KeyStart: Var(uast.KeyStart),
			uast.KeyEnd:   Var(uast.KeyEnd),
		}),
		"Name": Var(nameVar),
	})
}

var Normalizers []Mapping = []Mapping{
	MapSemantic("str", uast.String{}, MapObj(
		Obj{uast.KeyToken: Var("val")},
		Obj{
			"Value":  Var("val"),
			"Format": String(""),
		},
	)),

	MapSemantic("true", uast.Bool{}, MapObj(Obj{}, Obj{"Value": Bool(true)})),
	MapSemantic("false", uast.Bool{}, MapObj(Obj{}, Obj{"Value": Bool(false)})),

	mapIdentifier("splay"),
	mapIdentifier("lvar"),
	mapIdentifier("ivar"),
	mapIdentifier("gvar"),
	mapIdentifier("cvar"),
	mapIdentifier("Symbol"),
	mapIdentifier("Sym"),
	mapIdentifier("Const"),

	tokenIsIdentifier("casgn", "selector", role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	tokenIsIdentifier("lvasgn", "target", role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	tokenIsIdentifier("ivasgn", "target", role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	tokenIsIdentifier("gvasgn", "target", role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	tokenIsIdentifier("cvasgn", "target", role.Expression, role.Assignment, role.Binary, role.Identifier, role.Left),
	tokenIsIdentifier("send_assign", uast.KeyToken, role.Expression, role.Assignment, role.Left),
	tokenIsIdentifier("module", uast.KeyToken, role.Statement, role.Module),
	tokenIsIdentifier("sym", uast.KeyToken, role.Expression, role.Identifier),
	tokenIsIdentifier("const", uast.KeyToken, role.Expression, role.Identifier, role.Incomplete),
	tokenIsIdentifier("send_call", "selector", role.Expression, role.Function, role.Call),

	// iflipflop / eflipflop, have selector but not names
	MapSemantic("flip_1", uast.Identifier{}, MapObj(
		Obj{
			uast.KeyToken: Var("ident"),
			// could not find any case where this wasn't nil, but this will
			// detect if we find any
			"base": Is(nil),
			// same as ident, but with a slightly different format,
			// so we just drop it
			"qnames": Any(),
		},
		Obj{"Name": Var("ident")},
	)),
	MapSemantic("flip_2", uast.Identifier{}, MapObj(
		Obj{
			uast.KeyToken: Var("ident"),
			// could not find any case where this wasn't nil, but this will
			// detect if we find any
			"base": Is(nil),
			// same as ident, but with a slightly different format,
			// so we just drop it
			"qnames": Any(),
		},
		Obj{"Name": Var("ident")},
	)),

	MapSemantic("comment", uast.Comment{}, MapObj(
		Obj{
			uast.KeyToken: CommentText([2]string{"#", ""}, "comm"),
			// TODO(juanjux): map these two booleans in some way
			"documentation": Any(),
			"inline":        Any(),
		},
		CommentNode(false, "comm", nil),
	)),

	AnnotateType("send_require", MapObj(
		Obj{
			"base":   Var("path"),
			"values": Each("vals", Var("name")),
		},
		Obj{
			"Path": Var("path"),
			"Names": Each("vals", UASTType(uast.RuntimeImport{}, Obj{
				"Path": Var("name"),
			})),
		},
	), role.Expression, role.Import),

	MapSemantic("arg", uast.Argument{}, MapObj(
		Obj{uast.KeyToken: Var("name")},
		Obj{"Name": identifierWithPos("name")},
	)),

	MapSemantic("kwarg", uast.Argument{}, MapObj(
		Obj{uast.KeyToken: Var("name")},
		Obj{"Name": identifierWithPos("name")},
	)),

	MapSemantic("blockarg", uast.Argument{}, MapObj(
		Obj{uast.KeyToken: Var("name")},
		Obj{"Name": identifierWithPos("name")},
	)),

	MapSemantic("optarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
			"default":     Var("init"),
		},
		Obj{
			"Name": identifierWithPos("name"),
			"Init": Var("init"),
		},
	)),

	MapSemantic("kwoptarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
			"default":     Var("init"),
		},
		Obj{
			"Name": identifierWithPos("name"),
			"Init": Var("init"),
		},
	)),

	MapSemantic("restarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
		},
		Obj{
			"Name":     identifierWithPos("name"),
			"Variadic": Bool(true),
		},
	)),

	MapSemantic("restarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
			"default":     Var("init"),
		},
		Obj{
			"Name":     identifierWithPos("name"),
			"Init":     Var("init"),
			"Variadic": Bool(true),
		},
	)),

	MapSemantic("kwrestarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
			"default":     Var("init"),
		},
		Obj{
			"Name":        identifierWithPos("name"),
			"Init":        Var("init"),
			"MapVariadic": Bool(true),
		},
	)),

	MapSemantic("kwrestarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
		},
		Obj{
			"Name":        identifierWithPos("name"),
			"MapVariadic": Bool(true),
		},
	)),

	AnnotateType("class", MapObj(
		Fields{
			{Name: uast.KeyToken, Op: Var("name")},
		},
		Fields{
			{Name: uast.KeyToken, Op: UASTType(uast.Identifier{}, Obj{
				"Name": Var("name"),
			})},
		}),
		role.Statement, role.Type, role.Declaration, role.Identifier),

	AnnotateType("defs", MapObj(
		Fields{
			{Name: uast.KeyToken, Op: Var("name")},
		},
		Fields{
			{Name: uast.KeyToken, Op: UASTType(uast.Identifier{}, Obj{
				"Name": Var("name"),
			})},
		}),
		role.Statement, role.Function, role.Declaration, role.Identifier, role.Incomplete),

	MapSemantic("begin", uast.Block{}, MapObj(
		Fields{
			{Name: "body", Op: Var("body")},
			// TODO(juanjux): dont drop these once we've a better solution
			{Name: "comments", Optional: "opt_comments", Op: Any()},
		},
		Obj{
			"Statements": Var("body"),
		},
	)),

	MapSemantic("def", uast.FunctionGroup{}, MapObj(
		Fields{
			{Name: "body", Op: Var("body")},
			{Name: uast.KeyToken, Op: Var("name")},
			{Name: "args", Op: Cases("case_args",
				Obj{
					uast.KeyType: String("args"),
					uast.KeyPos:  Var("_pos"),
					"children":   Var("args"),
				},
				Check(
					Not(Has{"children": Var("args")}),
					Var("_nochildren"),
				),
			)},
			{Name: "comments", Optional: "comments_opt", Op: Var("comments")},
		},
		Obj{
			"Nodes": Arr(
				Fields{
					{Name: "comments", Optional: "comments_opt", Op: Var("comments")},
				},
				UASTType(uast.Alias{}, Obj{
					"Name": UASTType(uast.Identifier{}, Obj{
						"Name": Var("name"),
					}),
					"Node": UASTType(uast.Function{}, Obj{
						"Type": UASTType(uast.FunctionType{},
							CasesObj("case_args",
								Obj{},
								Objs{
									{"Arguments": Var("args")},
									{"Arguments": Arr()},
								},
							)),
						"Body": UASTType(uast.Block{}, Obj{
							"Statements": Var("body"),
						}),
					}),
				}),
			),
		},
	)),
}

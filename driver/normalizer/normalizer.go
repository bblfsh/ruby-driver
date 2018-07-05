package normalizer

import (
	"gopkg.in/bblfsh/sdk.v2/uast"
	"gopkg.in/bblfsh/sdk.v2/uast/role"
	. "gopkg.in/bblfsh/sdk.v2/uast/transformer"
)

var Preprocess = Transformers([][]Transformer{
	{Mappings(Preprocessors...)},
}...)

var Normalize = Transformers([][]Transformer{
	{Mappings(Normalizers...)},
}...)

var Preprocessors = []Mapping{}

func mapIdentifier(key string) Mapping {
	return MapSemantic(key, uast.Identifier{}, MapObj(
		Obj{uast.KeyToken: Var("val")},
		Obj{"Name": Var("val")},
	))
}

var Normalizers = []Mapping{
	MapSemantic("str", uast.String{}, MapObj(
		Obj{uast.KeyToken: Var("val")},
		Obj{
			"Value":  Var("val"),
			"Format": String(""),
		},
	)),

	mapIdentifier("splay"),
	mapIdentifier("lvar"),
	mapIdentifier("ivar"),
	mapIdentifier("gvar"),
	mapIdentifier("cvar"),
	mapIdentifier("Symbol"),
	mapIdentifier("Sym"),
	mapIdentifier("Const"),

	MapSemantic("send_qualified", uast.QualifiedIdentifier{}, MapObj(
		Obj{"qnames": Var("names")},
		Obj{"Names": Var("names")},
	)),

	// iflipflop / eflipflop, have selector but not names
	MapSemantic("flip_1", uast.Identifier{}, MapObj(
		Obj{uast.KeyToken: Var("ident")},
		Obj{"Name": Var("ident")},
	)),
	MapSemantic("flip_2", uast.Identifier{}, MapObj(
		Obj{uast.KeyToken: Var("ident")},
		Obj{"Name": Var("ident")},
	)),

	MapSemantic("comment", uast.Comment{}, MapObj(
		Obj{
			uast.KeyToken: CommentText([2]string{}, "comm"),
		},
		CommentNode(false, "comm", nil),
	)),


	AnnotateType("send_require", MapObj(
		Obj{
			"base": Var("path"),
			"values": Each("vals", Var("name")),
		},
		Obj{
			"Path": Var("path"),
			"Names": Each("vals", UASTType(uast.RuntimeImport{}, Obj{
				"Path": Var("name"),
			})),
		},
	), role.Expression, role.Import),

	//MapSemantic("send_require", uast.RuntimeImport{}, MapObj(
	//	Obj{
	//		"base": Var("path"),
	//		"values": Var("names"),
	//	},
	//	Obj{
	//		"Path": Var("path"),
	//		"Names": Var("names"),
	//	},
	//)),

	MapSemantic("arg", uast.Argument{}, MapObj(
		Obj{uast.KeyToken: Var("name")},
		Obj{"Name": UASTType(uast.Identifier{}, Obj{
			"Name": Var("name"),
		})},
	)),

	MapSemantic("optarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
			"default":     Var("init"),
		},
		Obj{
			"Name": Var("name"),
			"Init": Var("init"),
		},
	)),

	MapSemantic("kwoptarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
			"default":     Var("init"),
		},
		Obj{
			"Name": Var("name"),
			"Init": Var("init"),
		},
	)),

	MapSemantic("restarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
		},
		Obj{
			"Name":     Var("name"),
			"Variadic": Bool(true),
		},
	)),

	MapSemantic("restarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
			"default":     Var("init"),
		},
		Obj{
			"Name":     Var("name"),
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
			"Name":        Var("name"),
			"Init":        Var("init"),
			"MapVariadic": Bool(true),
		},
	)),

	MapSemantic("kwrestarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
		},
		Obj{
			"Name":        Var("name"),
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

	MapSemantic("def", uast.FunctionGroup{}, MapObj(
		Obj{
			"body":        Var("body"),
			uast.KeyToken: Var("name"),
			"args": Cases("case_args",
				Obj{
					uast.KeyType: String("args"),
					uast.KeyPos:  Var("_pos"),
					"children":   Var("args"),
				},
				Check(
					Not(Has{"children": Var("args")}),
					Var("_nochildren"),
				),
			),
		},
		Obj{
			"Nodes": Arr(
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

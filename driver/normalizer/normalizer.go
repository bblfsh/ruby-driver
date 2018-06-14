package normalizer

import (
	. "gopkg.in/bblfsh/sdk.v2/uast/transformer"
	"gopkg.in/bblfsh/sdk.v2/uast"
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
	mapIdentifier("class"),

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

	MapSemantic("send_require", uast.RuntimeImport{}, MapObj(
		Obj{uast.KeyToken: Var("path")},
		Obj{"Path": Var("path")},
	)),

	MapSemantic("arg", uast.Argument{}, MapObj(
		Obj{uast.KeyToken: Var("name")},
		Obj{"Name": Var("name")},
	)),

	MapSemantic("optarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
			"default": Var("init"),
		},
		Obj{
			"Name": Var("name"),
			"Init": Var("init"),
		},
	)),

	MapSemantic("kwoptarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
			"default": Var("init"),
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
			"Name": Var("name"),
			"Variadic": Bool(true),
		},
	)),

	MapSemantic("restarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
			"default": Var("init"),
		},
		Obj{
			"Name": Var("name"),
			"Init": Var("init"),
			"Variadic": Bool(true),
		},
	)),

	MapSemantic("kwrestarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
			"default": Var("init"),
		},
		Obj{
			"Name": Var("name"),
			"Init": Var("init"),
			"MapVariadic": Bool(true),
		},
	)),

	MapSemantic("kwrestarg", uast.Argument{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
		},
		Obj{
			"Name": Var("name"),
			"MapVariadic": Bool(true),
		},
	)),

	MapSemantic("def", uast.FunctionGroup{}, MapObj(
		Obj{
			"body": Var("body"),
			uast.KeyToken: Var("name"),
			"args": Var("args"),
		},
		Obj{
			"Nodes": Arr(
				UASTType(uast.Alias{}, Obj{
					"Name": Var("name"),
					"Node": UASTType(uast.Function{}, Obj{
						"Type": UASTType(uast.FunctionType{}, Obj{
							"Arguments": Var("args"),
						}),
						"Body": Var("body"),
					}),
				}),
			),
		},
	)),
}

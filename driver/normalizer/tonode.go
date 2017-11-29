package normalizer

import "gopkg.in/bblfsh/sdk.v1/uast"

// ToNode is an instance of `uast.ObjectToNode`, defining how to transform an
// into a UAST (`uast.Node`).
//
// https://godoc.org/gopkg.in/bblfsh/sdk.v1/uast#ObjectToNode
var ToNode = &uast.ObjectToNode{
	InternalTypeKey: "type",
	LineKey: "start_line",
	EndLineKey: "end_line",
	ColumnKey: "start_col",
	EndColumnKey: "end_col",

	TokenKeys: map[string]bool {
		"name": true,
		"token": true,
		"selector": true,
		"target": true,
	},
	Modifier: func(n map[string] interface{}) error {
		// Native parser uses a [) interval for columns, so add 1 to start_col
		if col, ok := n["start_col"].(float64); ok {
			n["start_col"] = col + 1
		}
		return nil
	},
}

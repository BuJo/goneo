// Package gcy contains a lexer and parser for a graph language.
/*
Language EBNF:

	Query := SearchQuery | DeleteQuery | CreateQuery
	SearchQuery := ( Roots [ Match ] | Match ) Returns
	Roots := "start" Root
	Root := name "=" NodeOrRel "(" id ")" ["," Root]
	NodeOrRel := "node" | "relation"
	Returns := "return" ReturnVal ["," Return]
	Return := name ["," Return]
	Match := "match" PathPart
	PathPart := PathAssignment | Path
	PathAssignment := name "=" Path
	Path := NodeRel
	NodeRel := "(" name ")" DirectionalRel "(" name ")"
	DirectionalRel := "<-" "[" name [":" name ] [ RelCount ] "]"
	RelCount := "*" | \d+ [".." * | \d+]

*/
package gcy

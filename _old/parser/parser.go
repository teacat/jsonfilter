package parser

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var pathLexer = lexer.MustSimple([]lexer.Rule{
	{`Ident`, `\*|[a-zA-Z][a-zA-Z_\d]*`, nil},
	{`Punct`, `[][=/,()]`, nil},
	{"whitespace", `\s+`, nil},
})

type Result struct {
	Props []*Prop `@@ ( "," @@ )*`
}

// Object | Group
type Prop struct {
	Group  *Group  `@@`
	Object *Object `| @@`
}

// @Ident | @Ident "/" Prop
type Object struct {
	Name string `@Ident`
	Prop *Prop  `( "/" @@ )*`
}

// @Ident "(" Prop ( "," Prop )* ")"
type Group struct {
	Name  string  `@Ident`
	Props []*Prop `"(" @@ ( "," @@ )* ")"`
}

func Parse(v string) (*Result, error) {
	res := &Result{}
	err := participle.MustBuild(&Result{}, participle.Lexer(pathLexer)).Parse("", strings.NewReader(v), res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

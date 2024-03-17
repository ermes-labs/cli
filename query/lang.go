package query

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Define the lexer for the query language.
var queryLexer = participle.Lexer(lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Whitespace", Pattern: `\s+`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Number", Pattern: `[-+]?[0-9]+`},
	{Name: "Punct", Pattern: `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
}))

var parser = participle.MustBuild[Query](
	(queryLexer),
	participle.Elide("Whitespace"),
	participle.UseLookahead(2),
)

// Query is the top-level query structure.
type Query struct {
	Set *Set `parser:"@@"`
}

// Set is a set of areas.
type Set struct {
	AreaId      *string      `parser:"  '#'@Ident"`
	NodesInArea *NodesInArea `parser:"| @@"`
	Sets        *Sets        `parser:"| '{' @@ '}'"`
}

// NodesInArea is a single area.
type NodesInArea struct {
	AreaId  string    `parser:"(@Ident | @('*'))"`
	Filters []*Filter `parser:"('(' @@ (',' @@)* ')')?"`
}

type Sets struct {
	InitialSet *Set           `parser:"@@"`
	NextSets   []*SetWithSign `parser:"@@*"`
}

type SetWithSign struct {
	Op  string `parser:"@('+' | '-')"`
	Set *Set   `parser:"@@"`
}

type Filter struct {
	Tag   *TagFilter   `parser:"  @@"`
	Level *LevelFilter `parser:"| @@"`
}

type TagFilter struct {
	Key   string `parser:"'tag' ':' @Ident '='"`
	Value string `parser:"@Ident"`
}

type LevelFilter struct {
	Comp string `parser:"'level' @( ('<' '>') | ('<' '=') | ('>' '=') | '=' | '<' | '>' | ('!' '=') )"`
	Numb string `parser:"@Number"`
}

func (q *Query) String() string {
	return q.Set.String()
}

func (s *Set) String() string {
	if s.AreaId != nil {
		return "#" + *s.AreaId
	}

	if s.NodesInArea != nil {
		return s.NodesInArea.String()
	}

	if s.Sets != nil {
		return s.Sets.String()
	}

	return ""
}

func (n *NodesInArea) String() string {
	return fmt.Sprintf("%s%v", n.AreaId, n.Filters)
}

func (g *Sets) String() string {
	result := g.InitialSet.String()

	for _, set := range g.NextSets {
		result += " " + set.String()
	}

	return "{ " + result + " }"
}

func (s *SetWithSign) String() string {
	return fmt.Sprintf("%s %s", s.Op, s.Set.String())
}

func (f *Filter) String() string {
	if f.Tag != nil {
		return f.Tag.String()
	}

	if f.Level != nil {
		return f.Level.String()
	}

	return ""
}

func (t *TagFilter) String() string {
	return fmt.Sprintf("tag:%s=%s", t.Key, t.Value)
}

func (l *LevelFilter) String() string {
	return fmt.Sprintf("level=%s", l.Numb)
}

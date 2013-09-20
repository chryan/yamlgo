package yaml

import "fmt"

type TokenType int

const (
	TOKEN_DIRECTIVE TokenType = iota
	TOKEN_DOC_START
	TOKEN_DOC_END
	TOKEN_BLOCK_SEQ_START
	TOKEN_BLOCK_MAP_START
	TOKEN_BLOCK_SEQ_END
	TOKEN_BLOCK_MAP_END
	TOKEN_BLOCK_ENTRY
	TOKEN_FLOW_SEQ_START
	TOKEN_FLOW_MAP_START
	TOKEN_FLOW_SEQ_END
	TOKEN_FLOW_MAP_END
	TOKEN_FLOW_MAP_COMPACT
	TOKEN_FLOW_ENTRY
	TOKEN_KEY
	TOKEN_VALUE
	TOKEN_ANCHOR
	TOKEN_ALIAS
	TOKEN_TAG
	TOKEN_PLAIN_SCALAR
	TOKEN_NON_PLAIN_SCALAR
)

type TokenStatus int

var tokenNames []string = []string{
	"DIRECTIVE",
	"DOC_START",
	"DOC_END",
	"BLOCK_SEQ_START",
	"BLOCK_MAP_START",
	"BLOCK_SEQ_END",
	"BLOCK_MAP_END",
	"BLOCK_ENTRY",
	"FLOW_SEQ_START",
	"FLOW_MAP_START",
	"FLOW_SEQ_END",
	"FLOW_MAP_END",
	"FLOW_MAP_COMPACT",
	"FLOW_ENTRY",
	"KEY",
	"VALUE",
	"ANCHOR",
	"ALIAS",
	"TAG",
	"PLAIN_SCALAR",
	"NON_PLAIN_SCALAR",
}

const (
	VALID TokenStatus = iota
	INVALID
	UNVERIFIED
)

type Token struct {
	Status TokenStatus
	Type   TokenType
	Mark   Mark
	Value  string
	Params []string
	Data   int
}

func NewToken(ttype TokenType, mark Mark) *Token {
	return &Token{Status: VALID, Type: ttype, Mark: mark, Params: make([]string, 0, 8)}
}

func (t Token) String() (out string) {
	out = fmt.Sprintf("%v: %v", tokenNames[t.Type], t.Value)
	//for i, val := range t.Params {
	//	out << std::string(" ") << token.params[i];
	//}

	return
}

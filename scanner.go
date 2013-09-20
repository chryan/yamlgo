package yaml

import (
	"io"
)

type indentType int
type indentStatus int
type flowMarker int

const (
	it_MAP indentType = iota
	it_SEQ
	it_NONE
)

const (
	is_VALID indentStatus = iota
	is_INVALID
	is_UNKNOWN
)

const (
	fm_FLOW_MAP flowMarker = iota
	fm_FLOW_SEQ
)

type Scanner struct {
	reader io.Reader
}

type indentMarker struct {
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{reader: reader}
}

func (s *Scanner) Empty() bool {
	return false
}

func (s *Scanner) Peek() *Token {
	return NewToken(TOKEN_VALUE, NullMark)
}

func (s *Scanner) Pop() {

}

func (s *Scanner) Mark() Mark {
	return NullMark
}

/********************/
/***** Scanning *****/
/********************/

func (s *Scanner) ensureTokensInQueue() {
}

func (s *Scanner) scanNextToken() {
}

func (s *Scanner) scanToNextToken() {
}

func (s *Scanner) startStream() {
}

func (s *Scanner) endStream() {
}

func (s *Scanner) pushToken(ttype TokenType) *Token {
	return nil
}

func (s *Scanner) inFlowContext() bool {
	//return !s.flows.empty()
	return true
}

func (s *Scanner) inBlockContext() bool {
	//return s.flows.empty()
	return true
}

func (s *Scanner) getFlowLevel() int {
	//return s.flows.size()
	return 0
}

func (s *Scanner) getStartTokenFor(itype indentType) TokenType {
	return TOKEN_DIRECTIVE
}

func (s *Scanner) pushIndentTo(column int, itype indentType) *indentMarker {
	return nil
}

func (s *Scanner) popIndentToHere() {
}

func (s *Scanner) popAllIndents() {
}

func (s *Scanner) popIndent() {
}

func (s *Scanner) getTopIndent() int {
	return 0
}

/**************************/
/***** Checking Input *****/
/**************************/

func (s *Scanner) canInsertPotentialSimpleKey() bool {
	return false
}

func (s *Scanner) existsActiveSimpleKey() bool {
	return false
}

func (s *Scanner) insertPotentialSimpleKey() {
}

func (s *Scanner) invalidateSimpleKey() {
}

func (s *Scanner) verifySimpleKey() bool {
	return true
}

func (s *Scanner) popAllSimpleKeys() {
}


func (s *Scanner) panicParserException(msg string) {
	/*
	mark := NullMark
	if(!m_tokens.empty()) {
		const Token& token = m_tokens.front()
		mark = token.mark
	}
	throw ParserException(mark, msg)
	*/
}

func (s *Scanner) isWhitespaceToBeEaten(ch byte) bool {
	return true
}

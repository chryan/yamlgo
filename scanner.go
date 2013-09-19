package yaml

import (
	"io"
)

type Scanner struct {
	reader io.Reader
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{reader: reader}
}

func (s *Scanner) Empty() bool {
	return false
}

func (s *Scanner) Peek() *Token {
	return NewToken(VALUE, NullMark)
}

func (s *Scanner) Pop() {

}

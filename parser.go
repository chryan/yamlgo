package yaml

import (
	"fmt"
	"io"
)

type Parser struct {
	scanner    *Scanner
	directives *Directives
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{scanner: NewScanner(reader), directives: NewDirectives()}
}

func (p *Parser) IsValid() bool {
	return !p.scanner.Empty()
}

func (p *Parser) Load(reader io.Reader) {
	p.scanner = NewScanner(reader)
	p.directives = NewDirectives()
}

func (p *Parser) HandleNextDocument(evtHandler EventHandler) bool {
	if p.scanner == nil {
		return false
	}

	p.parseDirectives()
	if p.scanner.Empty() {
		return false
	}

	sdp := &singleDocParser{p.scanner, p.directives}
	sdp.handleDocument(evtHandler)
	return true
}

func (p *Parser) PrintTokens() (output string) {
	if p.scanner == nil {
		return
	}

	for !p.scanner.Empty() {
		output += p.scanner.Peek().String() + "\n"
		p.scanner.Pop()
	}
	return
}

func (p *Parser) parseDirectives() {
	readDirective := false

	for !p.scanner.Empty() {
		token := p.scanner.Peek()
		if token.Type != DIRECTIVE {
			break
		} else if !readDirective {
			// we keep the directives from the last document if none are specified
			// but if any directives are specific, then we reset them
			p.directives = NewDirectives()
		}

		readDirective = true
		p.handleDirective(token)
		p.scanner.Pop()
	}
}

func (p *Parser) handleDirective(token *Token) {
	switch token.Value {
	case "YAML":
		p.handleYamlDirective(token)
	case "TAG":
		p.handleTagDirective(token)
	}
}

func (p *Parser) handleYamlDirective(token *Token) {
	if len(token.Params) != 1 {
		panic(fmt.Sprintf("%v - %v", token.Mark, ERR_YAML_DIRECTIVE_ARGS))
	} else if !p.directives.Version.IsDefault {
		panic(fmt.Sprintf("%v - %v", token.Mark, ERR_REPEATED_YAML_DIRECTIVE))
	} else if c, err := fmt.Sscanf(token.Params[0], "%d.%d", &p.directives.Version.Major, &p.directives.Version.Minor); c != 2 || err != nil {
		panic(fmt.Sprintf("%v - %v %v", token.Mark, ERR_YAML_VERSION, token.Params[0]))
	} else if p.directives.Version.Major > 1 {
		panic(fmt.Sprintf("%v - %v", token.Mark, ERR_YAML_MAJOR_VERSION))
	}

	p.directives.Version.IsDefault = false
}

func (p *Parser) handleTagDirective(token *Token) {
	if len(token.Params) != 2 {
		panic(fmt.Sprintf("%v - %v", token.Mark, ERR_TAG_DIRECTIVE_ARGS))
	}

	handle := token.Params[0]
	if _, ok := p.directives.Tags[handle]; !ok {
		panic(fmt.Sprintf("%v - %v", token.Mark, ERR_REPEATED_TAG_DIRECTIVE))
	} else {
		p.directives.Tags[handle] = token.Params[1] // token.Params[1] == prefix
	}
}

/********************************/
/* Single document parsing code */
/********************************/

type singleDocParser struct {
	scanner    *Scanner
	directives *Directives
}

func (s *singleDocParser) handleDocument(evtHandler EventHandler) {

}

package yaml

import (
	"fmt"
)

/*****************************************/
/* Single document collection stack code */
/*****************************************/

type collectionType int

const (
	ct_None collectionType = iota
	ct_BlockMap
	ct_BlockSeq
	ct_FlowMap
	ct_FlowSeq
	ct_CompactMap
)

type collectionstack struct {
	stack []collectionType
}

func newCollectionStack() *collectionstack {
	return &collectionstack{
		stack: make([]collectionType, 0, 8),
	}
}

func (c *collectionstack) push(ctype collectionType) {
	c.stack = append(c.stack, ctype)
}

func (c *collectionstack) pop(ctype collectionType) {
	if l := len(c.stack) - 1; l >= 0 {
		if c.stack[l] != ctype {
			panic(&ParseError{NullMark, fmt.Sprintf("Collection type mismatch: %v != %v", c.stack[l], ctype)})
		}
		c.stack = c.stack[:l]
	}
}

func (c *collectionstack) top() collectionType {
	l := len(c.stack)
	if l == 0 {
		return ct_None
	}
	return c.stack[l-1]
}

/********************************/
/* Single document parsing code */
/********************************/

type singleDocParser struct {
	scanner    *Scanner
	directives *Directives
	cstack     *collectionstack
	anchors    map[string]Anchor
	curranchor Anchor
}

func newSingleDocParser(scanner *Scanner, directives *Directives) *singleDocParser {
	return &singleDocParser{
		scanner: scanner,
		directives: directives,
		cstack: newCollectionStack(),
		anchors: make(map[string]Anchor),
		curranchor: NullAnchor,
	}
}

func (s *singleDocParser) handleDocument(evtHandler EventHandler) {
	if s.scanner.Empty() {
		panic(&ParseError{NullMark, "No tokens in scanner."})
	} else if s.curranchor != NullAnchor {
		panic(&ParseError{NullMark, "Anchor is not reset to 0!"})
	}

	evtHandler.DocumentStart(s.scanner.Peek().Mark)
	
	// eat doc start
	if s.scanner.Peek().Type == TOKEN_DOC_START {
		s.scanner.Pop()
	}
	
	// recurse!
	s.handleNode(evtHandler)
	evtHandler.DocumentEnd()
	
	// and finally eat any doc ends we see
	for !s.scanner.Empty() && s.scanner.Peek().Type == TOKEN_DOC_END {
		s.scanner.Pop()
	}
}

func (s *singleDocParser) handleNode(evtHandler EventHandler) {
	// an empty node *is* a possibility
	if s.scanner.Empty() {
		evtHandler.Null(s.scanner.Mark(), NullAnchor)
		return
	}
	
	// save location
	mark := s.scanner.Peek().Mark
	
	switch s.scanner.Peek().Type {
	// special case: a value node by itself must be a map, with no header
	case TOKEN_VALUE:
		evtHandler.MapStart(mark, "?", NullAnchor)
		s.handleMap(evtHandler)
		evtHandler.MapEnd()
		return
	// special case: an alias node
	case TOKEN_ALIAS:
		evtHandler.Alias(mark, s.lookupAnchor(mark, s.scanner.Peek().Value))
		s.scanner.Pop()
		return
	}
	
	tag, anchor := s.parseProperties()
	token := s.scanner.Peek()

	if token.Type == TOKEN_PLAIN_SCALAR && token.Value == "null" {
		evtHandler.Null(mark, anchor)
		s.scanner.Pop()
		return
	}
	
	// add non-specific tags
	if len(tag) == 0 {
		if token.Type == TOKEN_NON_PLAIN_SCALAR {
			tag = "!"
		} else {
			tag = "?"
		}
	}
	
	// now split based on what kind of node we should be
	switch token.Type {
		case TOKEN_PLAIN_SCALAR, TOKEN_NON_PLAIN_SCALAR:
			evtHandler.Scalar(mark, tag, anchor, token.Value)
			s.scanner.Pop()
			return
		case TOKEN_FLOW_SEQ_START, TOKEN_BLOCK_SEQ_START:
			evtHandler.SequenceStart(mark, tag, anchor)
			s.handleSequence(evtHandler)
			evtHandler.SequenceEnd()
			return
		case TOKEN_FLOW_MAP_START, TOKEN_BLOCK_MAP_START:
			evtHandler.MapStart(mark, tag, anchor)
			s.handleMap(evtHandler)
			evtHandler.MapEnd()
			return
		case TOKEN_KEY:
			// compact maps can only go in a flow sequence
			if s.cstack.top() == ct_FlowSeq {
				evtHandler.MapStart(mark, tag, anchor)
				s.handleMap(evtHandler)
				evtHandler.MapEnd()
				return
			}
			break
	}
	
	if tag == "?" {
		evtHandler.Null(mark, anchor)
	} else {
		evtHandler.Scalar(mark, tag, anchor, "")
	}
}

func (s *singleDocParser) handleSequence(evtHandler EventHandler) {
	// split based on start token
	switch s.scanner.Peek().Type {
	case TOKEN_BLOCK_SEQ_START:
		s.handleBlockSequence(evtHandler)
	case TOKEN_FLOW_SEQ_START:
		s.handleFlowSequence(evtHandler)
	}
}

func (s *singleDocParser) handleBlockSequence(evtHandler EventHandler) {
	// eat start token
	s.scanner.Pop()
	s.cstack.push(ct_BlockSeq)
	
	for {
		if s.scanner.Empty() {
			panic(&ParseError{s.scanner.Mark(), ERR_END_OF_SEQ})
		}
		
		// Make copy.
		token := *s.scanner.Peek()
		if token.Type != TOKEN_BLOCK_ENTRY && token.Type != TOKEN_BLOCK_SEQ_END {
			panic(&ParseError{token.Mark, ERR_END_OF_SEQ})
		}
		
		s.scanner.Pop()
		if token.Type == TOKEN_BLOCK_SEQ_END {
			break
		}
		
		// check for null
		if !s.scanner.Empty() {
			if token := s.scanner.Peek(); token.Type == TOKEN_BLOCK_ENTRY || token.Type == TOKEN_BLOCK_SEQ_END {
				evtHandler.Null(token.Mark, NullAnchor)
				continue
			}
		}
		
		s.handleNode(evtHandler)
	}
	
	s.cstack.pop(ct_BlockSeq)
}

func (s *singleDocParser) handleFlowSequence(evtHandler EventHandler) {
		// eat start token
	s.scanner.Pop()
	s.cstack.push(ct_FlowSeq)
	
	for {
		if s.scanner.Empty() {
			panic(&ParseError{s.scanner.Mark(), ERR_END_OF_SEQ_FLOW})
		}
		
		// first check for end
		if s.scanner.Peek().Type == TOKEN_FLOW_SEQ_END {
			s.scanner.Pop()
			break
		}
		
		// then read the node
		s.handleNode(evtHandler)
		
		if s.scanner.Empty() {
			panic(&ParseError{s.scanner.Mark(), ERR_END_OF_SEQ_FLOW})
		}
		
		// now eat the separator (or could be a sequence end, which we ignore - but if it's neither, then it's a bad node)
		if token := s.scanner.Peek(); token.Type == TOKEN_FLOW_ENTRY {
			s.scanner.Pop()
		} else if token.Type != TOKEN_FLOW_SEQ_END {
			panic(&ParseError{token.Mark, ERR_END_OF_SEQ_FLOW})
		}
	}
	
	s.cstack.pop(ct_FlowSeq)
}

func (s *singleDocParser) handleMap(evtHandler EventHandler) {
	// split based on start token
	switch s.scanner.Peek().Type {
		case TOKEN_BLOCK_MAP_START:
			s.handleBlockMap(evtHandler)
		case TOKEN_FLOW_MAP_START:
			s.handleFlowMap(evtHandler)
		case TOKEN_KEY:
			s.handleCompactMap(evtHandler)
		case TOKEN_VALUE:
			s.handleCompactMapWithNoKey(evtHandler)
	}
}

func (s *singleDocParser) handleBlockMap(evtHandler EventHandler) {
	// eat start token
	s.scanner.Pop()
	s.cstack.push(ct_BlockMap)
	
	for {
		if s.scanner.Empty() {
			panic(&ParseError{s.scanner.Mark(), ERR_END_OF_MAP})
		}
		
		token := s.scanner.Peek()
		if token.Type != TOKEN_KEY && token.Type != TOKEN_VALUE && token.Type != TOKEN_BLOCK_MAP_END {
			panic(&ParseError{token.Mark, ERR_END_OF_MAP})
		}
		
		if token.Type == TOKEN_BLOCK_MAP_END {
			s.scanner.Pop()
			break
		}
		
		// grab key (if non-null)
		if token.Type == TOKEN_KEY {
			s.scanner.Pop()
			s.handleNode(evtHandler)
		} else {
			evtHandler.Null(token.Mark, NullAnchor)
		}
		
		// now grab value (optional)
		if !s.scanner.Empty() && s.scanner.Peek().Type == TOKEN_VALUE {
			s.scanner.Pop()
			s.handleNode(evtHandler)
		} else {
			evtHandler.Null(token.Mark, NullAnchor)
		}
	}
	
	s.cstack.pop(ct_BlockMap)
}

func (s *singleDocParser) handleFlowMap(evtHandler EventHandler) {
	// eat start token
	s.scanner.Pop()
	s.cstack.push(ct_FlowMap)
	
	for {
		if s.scanner.Empty() {
			panic(&ParseError{s.scanner.Mark(), ERR_END_OF_MAP_FLOW})
		}
		
		token := s.scanner.Peek()
        mark := token.Mark
		// first check for end
		if token.Type == TOKEN_FLOW_MAP_END {
			s.scanner.Pop()
			break
		}
		
		// grab key (if non-null)
		if token.Type == TOKEN_KEY {
			s.scanner.Pop()
			s.handleNode(evtHandler)
		} else {
			evtHandler.Null(mark, NullAnchor)
		}
		
		// now grab value (optional)
		if !s.scanner.Empty() && s.scanner.Peek().Type == TOKEN_VALUE {
			s.scanner.Pop()
			s.handleNode(evtHandler)
		} else {
			evtHandler.Null(mark, NullAnchor)
		}
		
		if s.scanner.Empty() {
			panic(&ParseError{s.scanner.Mark(), ERR_END_OF_MAP_FLOW})
		}
		
		// now eat the separator (or could be a map end, which we ignore - but if it's neither, then it's a bad node)
		token = s.scanner.Peek()
		if token.Type == TOKEN_FLOW_ENTRY {
			s.scanner.Pop()
		} else if token.Type != TOKEN_FLOW_MAP_END {
			panic(&ParseError{token.Mark, ERR_END_OF_MAP_FLOW})
		}
	}
	
	s.cstack.pop(ct_FlowMap)
}

func (s *singleDocParser) handleCompactMap(evtHandler EventHandler) {
	s.cstack.push(ct_CompactMap)
	
	// grab key
	mark := s.scanner.Peek().Mark
	s.scanner.Pop()
	s.handleNode(evtHandler)
	
	// now grab value (optional)
	if !s.scanner.Empty() && s.scanner.Peek().Type == TOKEN_VALUE {
		s.scanner.Pop()
		s.handleNode(evtHandler)
	} else {
		evtHandler.Null(mark, NullAnchor)
	}
	
	s.cstack.pop(ct_CompactMap)
}

func (s *singleDocParser) handleCompactMapWithNoKey(evtHandler EventHandler) {
	s.cstack.push(ct_CompactMap)
	
	// null key
	evtHandler.Null(s.scanner.Peek().Mark, NullAnchor)
	
	// grab value
	s.scanner.Pop()
	s.handleNode(evtHandler)
	
	s.cstack.pop(ct_CompactMap)
}

func (s *singleDocParser) parseProperties() (tag string, anchor Anchor) {
	for !s.scanner.Empty() {
		switch s.scanner.Peek().Type {
		case TOKEN_TAG:
			s.parseTag(&tag)
		case TOKEN_ANCHOR:
			s.parseAnchor(&anchor)
		default:
			return
		}
	}
	return
}

func (s *singleDocParser) parseTag(tag *string) {
	token := s.scanner.Peek();
	if len(*tag) > 0 {
		panic(&ParseError{token.Mark, ERR_MULTIPLE_TAGS})
	}
	
	tagInfo := tagFromToken(token)
	*tag = tagInfo.Translate(s.directives)
	s.scanner.Pop();
}

func (s *singleDocParser) parseAnchor(anchor *Anchor) {
	token := s.scanner.Peek();
	if *anchor != NullAnchor {
		panic(&ParseError{token.Mark, ERR_MULTIPLE_ANCHORS})
	}
	
	*anchor = s.registerAnchor(token.Value)
	s.scanner.Pop();
}

func (s *singleDocParser) registerAnchor(name string) (ret Anchor) {
	if len(name) == 0 {
		return
	}
	
	s.curranchor++
	ret = s.curranchor
	s.anchors[name] = ret
	return
}

func (s *singleDocParser) lookupAnchor(mark Mark, name string) (ret Anchor) {
	if val, ok := s.anchors[name]; !ok {
		panic(&ParseError{mark, ERR_UNKNOWN_ANCHOR})
	} else {
		ret = val
	}
	return
}

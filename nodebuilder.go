package yaml

type nodeBuilder struct {
	root interface{}

	stack []interface{}
	anchors []interface{}

	// Pushed keys
	keys []struct {
		node interface{}
		flag bool
	}
	mapDepth uint
}

/*
func newNodeBuilder() *nodeBuilder {
	return &nodeBuilder{
		stack: make([]interface{}, 0),
		anchors: make([]interface{}, 1),
	}
}

func (n *nodeBuilder) OnDocumentStart(mark Mark) {

}

func (n *nodeBuilder) OnDocumentEnd() {

}

func (n *nodeBuilder) OnNull(mark Mark, anchor Anchor) {
	n.pushAnchor(anchor)
}

func (n *nodeBuilder) OnAlias(mark Mark, anchor Anchor) {

}

func (n *nodeBuilder) OnScalar(mark Mark, tag string, anchor Anchor, value string) {

}

func (n *nodeBuilder) OnSequenceStart(mark Mark, tag string, anchor Anchor) {

}

func (n *nodeBuilder) OnSequenceEnd() {

}

func (n *nodeBuilder) OnMapStart(mark Mark, tag string, anchor Anchor) {

}

func (n *nodeBuilder) OnMapEnd() {

}

func (n *nodeBuilder) pushAnchor(anchor Anchor) interface{} {
	return nil
}
/*
func (n *nodeBuilder) push(detail::node& node) {
}

func (n *nodeBuilder) pop() {

}

func (n *nodeBuilder) RegisterAnchor(anchor_t anchor, detail::node& node) {

}
*/
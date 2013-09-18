package yaml

type EventHandler interface {
	OnDocumentStart(mark Mark)
	OnDocumentEnd()
	
	OnNull(mark Mark, anchor Anchor)
	OnAlias(mark Mark, anchor Anchor)
	OnScalar(mark Mark, tag string, anchor Anchor, value string)

	OnSequenceStart(mark Mark, tag string, anchor Anchor)
	OnSequenceEnd()

	OnMapStart(mark Mark, tag string, anchor Anchor)
	OnMapEnd()
}
package yaml

type EventHandler interface {
	DocumentStart(mark Mark)
	DocumentEnd()

	Null(mark Mark, anchor Anchor)
	Alias(mark Mark, anchor Anchor)
	Scalar(mark Mark, tag string, anchor Anchor, value string)

	SequenceStart(mark Mark, tag string, anchor Anchor)
	SequenceEnd()

	MapStart(mark Mark, tag string, anchor Anchor)
	MapEnd()
}

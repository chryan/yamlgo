package yaml

type tagType int

const (
	tag_VERBATIM tagType = iota
	tag_PRIMARY_HANDLE
	tag_SECONDARY_HANDLE
	tag_NAMED_HANDLE
	tag_NON_SPECIFIC
)

type tag struct {
	tagtype tagType
	handle string
	value string
}

func tagFromToken(token *Token) tag {
	newTag := tag{
		tagtype: tagType(token.Data),
	}

	switch newTag.tagtype {
		case tag_VERBATIM:
			newTag.value = token.Value
		case tag_PRIMARY_HANDLE:
			newTag.value = token.Value
		case tag_SECONDARY_HANDLE:
			newTag.value = token.Value
		case tag_NAMED_HANDLE:
			newTag.handle = token.Value
			newTag.value = token.Params[0]
		case tag_NON_SPECIFIC:
		default:
			panic("Bad tag from token.")
	}

	return newTag
}

func (t tag) Translate(directives *Directives) string {
	switch t.tagtype {
		case tag_VERBATIM:
			return t.value
		case tag_PRIMARY_HANDLE:
			return directives.TranslateTagHandle("!") + t.value
		case tag_SECONDARY_HANDLE:
			return directives.TranslateTagHandle("!!") + t.value
		case tag_NAMED_HANDLE:
			return directives.TranslateTagHandle("!" + t.handle + "!") + t.value
		case tag_NON_SPECIFIC:
			return "!"
	}
	panic("internal error, bad tag type")
}
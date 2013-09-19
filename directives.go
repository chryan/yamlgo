package yaml

type Version struct {
	IsDefault bool
	Major     int
	Minor     int
}

type Directives struct {
	Version Version
	Tags    map[string]string
}

func NewDirectives() *Directives {
	return &Directives{Version: Version{true, 1, 2}}
}

func (d *Directives) TranslateTagHandle(handle string) string {
	if val, ok := d.Tags[handle]; ok {
		return val
	}

	if handle == "!!" {
		return "tag:yaml.org,2002:"
	}
	return handle
}

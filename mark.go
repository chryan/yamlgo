package yaml

type Mark struct {
	Pos    int
	Line   int
	Column int
}

var NullMark Mark = Mark{Pos: -1, Line: -1, Column: -1}

func (m Mark) String() string {
	return "Pos:%v Line:%v Col:%v"
}

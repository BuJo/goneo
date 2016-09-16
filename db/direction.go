package db

type Direction int

const (
	Both Direction = iota
	Incoming
	Outgoing
)

// Convert a string to a direction
func DirectionFromString(str string) Direction {
	if str == "out" || str == "outgoing" {
		return Outgoing
	} else if str == "in" || str == "incoming" {
		return Incoming
	}
	return Both
}

func (d Direction) String() string {
	switch d {
	case Both:
		return "-"
	case Incoming:
		return "<-"
	case Outgoing:
		return "->"
	default:
		return ""
	}
}

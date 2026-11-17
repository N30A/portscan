package main

type Status uint

const (
	StatusError Status = iota
	StatusClosed
	StatusFiltered
	StatusOpen
)

func (s Status) String() string {
	switch s {
	case StatusError:
		return "error"
	case StatusClosed:
		return "closed"
	case StatusFiltered:
		return "filtered"
	case StatusOpen:
		return "open"
	default:
		panic("unknown status")
	}
}

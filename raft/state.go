package raft

type State uint64

const (
	Leader State = iota
	Candidate
	Follower
	Unkown
)

func (r *Raft) State() State {
	switch r.state {
	case Leader:
		return Leader
	case Candidate:
		return Candidate
	case Follower:
		return Follower
	default:
		return Unkown
	}
}

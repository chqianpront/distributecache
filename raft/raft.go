package raft

type ServerId string
type CommitId int
type TermId int
type Commit struct {
	TermId   TermId
	CommitId CommitId
	Command  Command
}
type Raft struct {
	state          State
	stopCh         chan bool
	stoped         bool
	term           TermId
	replications   []Server
	peers          []Server
	leaderId       ServerId
	lastCommitId   CommitId
	commits        []Commit
	lastCommitedId CommitId
	commtied       []Commit
}

func StartServer() {
	raft := &Raft{
		stopCh: make(chan bool),
	}
	raft.restore()
	raft.Start()
}
func (r *Raft) Start() {
	switch r.State() {
	case Leader:
		r.LeaderLoop()
	case Candidate:
		r.CandidateLoop()
	case Follower:
		r.FollowerLoop()
	}
}

func (r *Raft) restore() {}

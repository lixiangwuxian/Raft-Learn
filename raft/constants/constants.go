package constants

type State int

const (
	Follower State = iota
	Leader
	Candidate
)

type PackegeType int

const (
	AppendEntries PackegeType = iota
	AppendEntriesReply
	RequestVote
	RequestVoteReply
)

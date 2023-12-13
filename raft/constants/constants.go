package constants

type State string

const (
	Follower  State = "follower"
	Leader    State = "leader"
	Candidate State = "candidate"
)

type PackegeType int

const (
	AppendEntries PackegeType = iota
	AppendEntriesReply
	RequestVote
	RequestVoteReply
	UserRequest
)

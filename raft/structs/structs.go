package structs

type Inform struct {
	KnownLeader      string
	KnownNodes       []string
	CurrentTerm      int
	FollowerTimeout  int //多长时间从follower转为leader，使用时需要加一个随机数
	CandidateTimeout int //多长时间选举超时
	CommitIndex      int
}

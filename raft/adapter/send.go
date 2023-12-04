package adapter

type Sender interface {
	AppendEntries(peerAddr string, data AppendEntries)
	AppendEntriesReply(peerAddr string, data AppendEntriesReply)
	RequestVote(peerAddr string, data RequestVote, currentTerm int)
	RequestVoteReply(peerAddr string, data RequestVoteReply, currentTerm int)
}

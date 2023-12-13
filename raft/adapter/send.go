package adapter

import "lxtend.com/m/packages"

type Sender interface {
	AppendEntries(peerAddr string, data packages.AppendEntries)
	AppendEntriesReply(peerAddr string, data packages.AppendEntriesReply)
	RequestVote(peerAddr string, data packages.RequestVote, currentTerm int)
	RequestVoteReply(peerAddr string, data packages.RequestVoteReply, currentTerm int)
}

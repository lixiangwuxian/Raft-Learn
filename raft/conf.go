package main

type Conf struct {
	Peers            map[int]string
	TotalNodes       int
	LeaderTimeout    int
	CandidateTimeout int
}

var conf *Conf

func initConf() *Conf {
	conf = new(Conf)
	conf.Peers = make(map[int]string)
	conf.Peers[0] = ""
	conf.Peers[1] = ""
	conf.Peers[2] = ""
	conf.TotalNodes = 3
	conf.LeaderTimeout = 1000
	conf.CandidateTimeout = 1000
	return conf
}

func getPeers() map[int]string {
	return conf.Peers
}

// func getLeaderTimeout() int {
// 	return conf.LeaderTimeout
// }

// func getCandidateTimeout() int {
// 	return conf.CandidateTimeout
// }

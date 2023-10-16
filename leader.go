package raft

func leader_loop() {
	for {
		var data = get_data_from_client()
		var log = create_log(data)
		var follower = get_follower()
		var success = send_log_to_followers(follower, log) //return true if log has been sent to over half of followers,others block
		if success {
			var commit = send_commit_to_followers(follower) //block until all followers has been committed
			if commit {
				add_log_to_stack(log)
			}
		}
	}
}

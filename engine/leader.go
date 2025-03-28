package engine

import "time"

//communicate with the leader

func (engine *UFollower) detectLeader() {
	for {
		if engine.Config.LeaderUrl != "" {

		} else {
			err := engine.broadCastToFindLeader()
			if err != nil {
				engine.Log.Error(err)
				time.Sleep(5 * time.Second)
				continue //retry
			}
		}

		break
	}
}

func (engine *UFollower) broadCastToFindLeader() error {
	return nil
}

package raft

import (
	"log"
	"time"
)

func (r *Raft) LeaderLoop() {
	timer := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-r.stopCh:
			return
		case <-timer.C:
			log.Printf("start looping\n")

		}
	}
}

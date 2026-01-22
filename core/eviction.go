package core

import "github.com/dicedb/dice/config"

func evictFirst() {
	for k := range store {
		delete(store, k)
		return
	}
}

func evict() {
	switch config.EvictionStrategy {
	case "simple-first":
		evictFirst()
	}
}

package main

import (
	"fmt"

	"github.com/xiaoweil/simple-ratelimiter/ratelimiter"
)

func main() {
	// Create a new user rate limiter with a rate of 1000 request/s,
	// and an initial bucket size of 3
	url := ratelimiter.NewUserRateLimiter(1000, 3)
	defer url.StopAll()

	users := []string{"user1", "user2", "user3"}
	for i := 0; i < 10; i++ {
		for _, u := range users {
			if url.Allow(u) {
				fmt.Printf("Request from %v is allowed.\n", u)
			} else {
				fmt.Printf("Request from %v is rejected.\n", u)
			}
		}
	}
}

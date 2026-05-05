// Package ratelimit implements a token-bucket rate limiter used to throttle
// outbound SSH connections and check executions in driftwatch.
//
// This prevents overwhelming remote hosts when scanning large inventories.
// The limiter is safe for concurrent use across goroutines.
//
// Basic usage:
//
//	l, err := ratelimit.New(ratelimit.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, host := range hosts {
//		if err := l.Wait(ctx); err != nil {
//			return err
//		}
//		go runChecks(host)
//	}
package ratelimit

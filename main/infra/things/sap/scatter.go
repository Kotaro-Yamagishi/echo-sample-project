package sap

import (
	"golang.org/x/sync/errgroup"
)

func scatter(n int, fn func(i int) error) error {
	var eg errgroup.Group

	for i := 0; i < n; i++ {
		i := i // Create a new variable scoped to the loop to avoid data race
		eg.Go(func() error {
			return fn(i)
		})
	}

	return eg.Wait()
}

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/mikloslorinczi/infra-exec/seeder"
	"github.com/mikloslorinczi/infra-exec/syncmap"
)

var myMap = syncmap.NewSafeMap()

func main() {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wg.Add(3)
	go populateMap(ctx, &wg)
	go populateMap(ctx, &wg)
	go printMap(ctx, &wg)
	wg.Wait()
}

func populateMap(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			myMap.SetKey(seeder.RandomHash(3), seeder.RandomHash(3))
		}
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}

func printMap(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("\nHash Map:\n%v\n", myMap.GetMap())
		}
		time.Sleep(time.Second)
	}
}

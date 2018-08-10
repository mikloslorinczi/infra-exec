package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	wg      sync.WaitGroup
	rwMutex sync.RWMutex
	hashMap = make(map[string]string)
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wg.Add(3)
	go populateMap(ctx, &wg)
	go populateMap(ctx, &wg)
	go printMap(ctx, &wg)
	wg.Wait()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getMap() map[string]string {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	return hashMap
}

func addRandomHash() {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	hashMap[randomHash(rand.Intn(8))] = randomHash(rand.Intn(8))
}

func randomHash(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func populateMap(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Populate Map canceled")
			return
		default:
			addRandomHash()
		}
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}

func printMap(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Print Map canceled")
			return
		default:
			fmt.Printf("\nHash Map:\n%v\n", getMap())
		}
		time.Sleep(time.Second)
	}
}

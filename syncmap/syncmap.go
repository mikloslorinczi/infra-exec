package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	wg      sync.WaitGroup
	mutex   sync.Mutex
	hashMap = make(map[string]string)
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func main() {
	wg.Add(1)
	go populateMap()
	go populateMap()
	go printMap()
	wg.Wait()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getMap() map[string]string {
	mutex.Lock()
	defer mutex.Unlock()
	return hashMap
}

func addRandomHash() {
	mutex.Lock()
	defer mutex.Unlock()
	hashMap[randomHash(rand.Intn(8))] = randomHash(rand.Intn(8))
}

func randomHash(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func populateMap() {
	for {
		addRandomHash()
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}

func printMap() {
	for {
		fmt.Printf("\nHash Map:\n%v\n", getMap())
		time.Sleep(time.Second)
	}
}

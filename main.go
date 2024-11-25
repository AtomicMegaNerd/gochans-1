package main

import (
	"log"
	"strings"
	"sync"
	"time"
)

type UserDB struct {
	users []string
	lock  sync.RWMutex
}

func (u *UserDB) GetUserChats() []string {
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.users
}

func (u *UserDB) AddUser(user string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.users = append(u.users, user)
}

func (u *UserDB) String() string {
	// Print all chats and friends
	sb := strings.Builder{}
	sb.WriteString("User Database: ")
	for _, user := range u.users {
		sb.WriteString(user)
		sb.WriteString(" ")
	}
	return sb.String()
}

type Message struct {
	users []string
}

func appendUsers(users []string, ch chan<- *Message, wg *sync.WaitGroup) {
	time.Sleep(2 * time.Second)
	defer wg.Done()
	ch <- &Message{users: users}
}

func processMessages(db *UserDB, ch <-chan *Message, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range ch {
		for _, user := range msg.users {
			db.AddUser(user)
		}
	}
}

func main() {

	timeStart := time.Now()

	// Create a new UserDB
	db := &UserDB{}

	ch := make(chan *Message)

	// Start the processMessages goroutine
	processWg := &sync.WaitGroup{}
	processWg.Add(1)
	go processMessages(db, ch, processWg)

	// Send messages to the channel
	messageWg := &sync.WaitGroup{}
	messageWg.Add(2)
	go appendUsers([]string{"Ted", "Alice"}, ch, messageWg)
	go appendUsers([]string{"Bob", "Fiona", "Wilma"}, ch, messageWg)

	// Wait for all messages to be sent before closing the channel
	messageWg.Wait()
	close(ch)

	// Wait for all messages to be processed
	processWg.Wait()

	// Print the final state of the database
	log.Println(db)
	log.Printf("Done in %.2f seconds\n", time.Since(timeStart).Seconds())
}

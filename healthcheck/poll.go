package healthcheck

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
)

var (
	pollsMutex sync.Mutex
	polls      map[string]*Poll
)

func init() {
	polls = make(map[string]*Poll)
}

type Poll struct {
}

func generateId() string {
	data := make([]byte, 12)
	rand.Read(data)
	return base64.RawURLEncoding.EncodeToString(data)
}

func CreatePoll() (key string, poll *Poll) {
	pollsMutex.Lock()
	defer pollsMutex.Unlock()

	key = generateId()

	poll = &Poll{}
	polls[key] = poll
	return
}

func GetPoll(key string) *Poll {
	pollsMutex.Lock()
	defer pollsMutex.Unlock()

	return polls[key]
}

func DeletePoll(key string) {
	pollsMutex.Lock()
	defer pollsMutex.Unlock()

	delete(polls, key)
}

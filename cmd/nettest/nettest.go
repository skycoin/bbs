package main

import (
	"github.com/skycoin/net/skycoin-messenger/factory"
	"sync"
	"errors"
)

const (
	MessengerServerAddress = ":8080"
)

func main() {
	var (
		quit = make(chan struct{})
		wg sync.WaitGroup
	)

	if e := runMessengerServer(quit, &wg); e != nil {
		panic(e)
	}

	if e := runTest(); e != nil {
		panic(e)
	}

	quit <- struct{}{}
	wg.Wait()
}

func runMessengerServer(quit chan struct{}, wg *sync.WaitGroup) error {
	f := factory.NewMessengerFactory()
	if e := f.Listen(MessengerServerAddress); e != nil {
		return e
	}
	go func(f *factory.MessengerFactory, quit chan struct{}, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()

		for {
			select {
			case <-quit:
				f.Close()
				return
			}
		}
	}(f, quit, wg)
	return nil
}

func runTest() error {
	f := factory.NewMessengerFactory()

	conn, e := f.Connect(MessengerServerAddress)
	if e != nil {
		return e
	}
	fromKey := conn.GetKey()

	if _, ok := f.GetConnection(fromKey); !ok {
		return errors.New("connection not found")
	}

	return nil
}
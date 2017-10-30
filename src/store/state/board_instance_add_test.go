package state

import (
	"github.com/skycoin/skycoin/src/cipher"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestBoardInstance_NewThread(t *testing.T) {
	const (
		boardSeed = "a" // Seed for board creation.
		userSeed  = "b" // Seed for user generation.
	)

	bi, quit := initInstance(t, boardSeed)
	defer quit()

	t.Run("add threads", func(t *testing.T) {
		const (
			threadCount = 30 // Number of threads to create.
		)
		var (
			threadsWG   sync.WaitGroup
			threadsChan = make(chan []uint64)
			threadsDone = make(chan struct{})
		)
		threadsWG.Add(threadCount)

		go func() {
			for {
				select {
				case signal := <-threadsChan:
					t.Logf("[%d] thread added. Expecting seq(%d).",
						signal[0], signal[1])
					threadsWG.Done()
				case <-threadsDone:
					return
				}
			}
		}()

		for i := 0; i < threadCount; i++ {
			actionAddThread := func(i int) {
				goal := addThread(t, bi, i, []byte(userSeed))
				threadsChan <- []uint64{uint64(i), goal}
			}
			if i%2 == 0 {
				go actionAddThread(i)
			} else {
				actionAddThread(i)
			}
		}

		threadsWG.Wait()
		threadsDone <- struct{}{}

		if e := bi.PublishChanges(); e != nil {
			t.Fatal("failed to publish changes:", e)
		}

		threadList := obtainThreadList(t, bi)
		if len(threadList) != threadCount {
			t.Fatalf("len(threadList) != addThreads_threadCount : got %d threads, expected %d.",
				len(threadList), threadCount)
		}

		for i, tHash := range threadList {
			t.Logf("[%d] thread '%s'.", i, tHash.Hex())
		}
	})

	t.Run("add threads with posts", func(t *testing.T) {
		const (
			threadCount = 30
			postCount   = 20
		)

		var wg sync.WaitGroup

		for i := 0; i < threadCount; i++ {
			addThread(t, bi, i, []byte(userSeed))

			threadList := obtainThreadList(t, bi)
			tHash := threadList[len(threadList)-1]

			for j := 0; j < postCount; j++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					addPost(t, bi, tHash, j, []byte(userSeed))
				}()
			}

			if e := bi.PublishChanges(); e != nil {
				t.Fatal("failed to publish changes:", e)
			}

			if i%2 == 0 {
				for j := postCount; j < postCount+postCount; j++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						addPost(t, bi, tHash, j, []byte(userSeed))
					}()
				}
			}
		}

		wg.Wait()
	})

	t.Run("add content while modifying board", func(t *testing.T) {
		const (
			threadCount   = 30
			postCount     = 100
			maxSubPKCount = 5
		)

		var wg sync.WaitGroup
		publishLoopBegin := make(chan struct{})
		quitChan := make(chan struct{})

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-publishLoopBegin
			ticker := time.NewTicker(time.Millisecond * 100)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if e := bi.PublishChanges(); e != nil {
						t.Fatal("failed to publish changes:", e)
					}
				case <-quitChan:
					return
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < threadCount; i++ {
				addThread(t, bi, i, []byte{byte(i)})
				time.Sleep(time.Millisecond * 50)
			}
			if e := bi.PublishChanges(); e != nil {
				t.Fatal("failed to publish changes:", e)
			}
			publishLoopBegin <- struct{}{}
			for i := 0; i < postCount; i++ {
				threads := obtainThreadList(t, bi)
				threadIndex := rand.Intn(int(len(threads)))
				addPost(t, bi, threads[threadIndex], i, []byte{byte(i)})
			}
			quitChan <- struct{}{}
			quitChan <- struct{}{}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			ticker := time.NewTicker(time.Millisecond * 100)
			defer ticker.Stop()

			for {
				select {
				case <-quitChan:
					return
				case <-ticker.C:
					pkCount := rand.Intn(maxSubPKCount)
					pks := make([]cipher.PubKey, pkCount)
					for i := 0; i < pkCount; i++ {
						pks[i], _ = cipher.GenerateKeyPair()
					}
					if _, e := bi.EnsureSubmissionKeys(pks); e != nil {
						t.Fatal("failed to change board submission keys:", e)
					}
				}
			}
		}()

		wg.Wait()
	})
}

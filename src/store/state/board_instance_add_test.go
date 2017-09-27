package state

import (
	"testing"
	"sync"
)

func TestBoardInstance_NewThread(t *testing.T) {
	const (
		boardSeed = "a" // Seed for board creation.
		userSeed = "b" // Seed for user generation.
	)

	bi, close := initInstance(t, boardSeed)
	defer close()

	t.Run("ADD_THREADS", func(t *testing.T) {
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

	t.Run("ADD_THREAD_AND_POST", func(t *testing.T) {
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
}
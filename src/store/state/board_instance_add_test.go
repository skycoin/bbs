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

		for i := 0; i < threadCount; i++ {
			for j := 0; j < postCount; j++ {
				//goal := addThread(t, bi, i, []byte(userSeed))

			}
		}
	})
}
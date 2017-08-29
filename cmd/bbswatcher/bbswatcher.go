package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

const (
	cmdStr        = "bbsnode"
	pauseDuration = 1
)

func main() {
	// Awaits Ctrl+C behaviour.
	keySig := CatchInterrupt()

	for {
		// Execute 'bbsnode', connecting outputs.
		cmd := exec.Command(cmdStr, os.Args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if e := cmd.Run(); e != nil {
			log.Println("FAILED HERE")
			return
		}

		// Await process exit.
		cmdSig := waiter(cmd.Process)

		select {
		case v := <-keySig:
			log.Println("Watcher received signal", v)
			log.Println("Process released:", cmd.Process.Release())
			return

		case v := <-cmdSig:
			if v != nil {
				log.Printf("Process (%d) exitied with state '%s'.",
					v.Pid(), v.String())
			}

			log.Printf("Restarting in %ds...", pauseDuration)
			time.Sleep(time.Second * time.Duration(pauseDuration))
		}
	}
	return
}

func waiter(process *os.Process) chan *os.ProcessState {
	out := make(chan *os.ProcessState)
	go func() {
		state, e := process.Wait()
		if e != nil {
			log.Println("Process exited with error:", e)
		}
		out <- state
	}()
	return out
}

// CatchInterrupt catches Ctrl+C behaviour.
func CatchInterrupt() chan int {
	quit := make(chan int)
	go func(q chan<- int) {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		signal.Stop(sigchan)
		q <- 1
	}(quit)
	return quit
}

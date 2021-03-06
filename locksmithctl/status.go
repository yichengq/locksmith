package main

import (
	"fmt"
	"os"

	"github.com/coreos/locksmith/lock"
)

var (
	cmdStatus = &Command{
		Name:        "status",
		Summary:     "Get the status of the cluster wide reboot lock.",
		Description: `Status will return the number of locks that are held and available and a list of the holders.`,
		Run:         runStatus,
	}
)

func printHolders(sem *lock.Semaphore) {
	fmt.Fprintln(out, "MACHINE ID")
	for _, h := range sem.Holders {
		fmt.Fprintln(out, h)
	}
}

func runStatus(args []string) (exit int) {
	elc, err := lock.NewEtcdLockClient(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing etcd client:", err)
		return 1
	}
	l := lock.New("", elc)

	sem, err := l.Get()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error setting value:", err)
	}

	fmt.Println("Available:", sem.Semaphore)
	fmt.Println("Max:", sem.Max)

	if len(sem.Holders) > 0 {
		fmt.Fprintln(out, "")
		printHolders(sem)
	}

	return
}

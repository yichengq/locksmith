package main

import (
	"fmt"
	"os"

	"github.com/coreos/locksmith/lock"
	"github.com/coreos/locksmith/pkg/machineid"
	"github.com/coreos/locksmith/third_party/github.com/coreos/go-systemd/login1"
)

var (
	cmdReboot = &Command{
		Name:        "reboot",
		Summary:     "Reboot honoring reboot locks.",
		Description: `Reboot will attempt to reboot immediatly after taking a reboot lock. The user is responsible for unlocking after a successful reboot.`,
		Run:         runReboot,
	}
)

func runReboot(args []string) int {
	if os.Geteuid() != 0 {
		fmt.Fprintln(os.Stderr, "Must be root to initiate reboot.")
		return 1
	}

	elc, _ := lock.NewEtcdLockClient(nil)
	lgn, err := login1.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot read machine-id")
		return 1
	}

	mID := machineid.MachineID("/")
	if mID == "" {
		fmt.Fprintln(os.Stderr, "Cannot read machine-id")
		return 1
	}

	l := lock.New(mID, elc)

	err = l.Lock()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error locking:", err)
		return 1
	}

	lgn.Reboot(false)

	// TODO(philips): Unlock if the reboot fails.

	return 0
}

package pcstat

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

const CLONE_NEWNS = 0x00020000

func SwitchMountNs(pid int) {
	myns := getMountNs(os.Getpid())
	pidns := getMountNs(pid)

	if myns != pidns {
		setns(pidns)
	}
}

func getMountNs(pid int) int {
	fname := fmt.Sprintf("/proc/%d/ns/mnt", pid)
	nss, err := os.Readlink(fname)
	if err != nil || len(nss) == 0 {
		return 0
	}
	nss = strings.TrimPrefix(nss, "mnt:[")
	nss = strings.TrimSuffix(nss, "]")
	ns, err := strconv.Atoi(nss)
	if err != nil {
		log.Fatalf("this namespace is not a number")
		return 0
	}

	return ns
}

// TODO: switch ns?
func setns(fd int) error {
	ret, _, err := unix.Syscall(unix.SYS_SETNS, uintptr(uint(fd)), uintptr(CLONE_NEWNS), 0)
	if ret != 0 {
		return fmt.Errorf("syscall SYS_SETNS failed: %v", err)
	}
	return nil
}

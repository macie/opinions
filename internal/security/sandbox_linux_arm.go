//go:build linux && arm && !openbsd && !unsafe

package security

import (
	"fmt"

	"github.com/elastic/go-seccomp-bpf"
)

// IsHardened reports whether security sandbox is enabled.
const IsHardened = true

// Sandbox restrict application access to necessary system calls needed by
// network connections and standard i/o.
func Sandbox() error {
	// How to create minimal whitelist:
	// 1. Create empty list of allowed syscalls
	// 2. Set `seccomp.ActionLog` as default filter action
	// 3. Compile and run program
	// 4. Use `dmesg` to find logged syscalls (started with _audit_)
	// 5. Translate syscalls numbers to names and add them to allowed list
	// 6. Go to point 3 and repeat until no new audit logs
	// 7. Reset default filter action to `seccomp.ActionKillProcess`
	whitelist := seccomp.SyscallGroup{
		Action: seccomp.ActionAllow,
		Names:  []string{
			// similar to stdio pledge
			// "clone3", "close", "epoll_create1", "epoll_ctl", "epoll_pwait",
			// "exit_group", "fcntl", "fstat", "futex", "getpid", "getrandom",
			// "getsockopt", "gettid", "madvise", "mmap", "mprotect", "munmap", "nanosleep",
			// "pipe2", "read", "rseq", "rt_sigprocmask", "rt_sigreturn",
			// "sched_getaffinity", "sched_yield", "set_robust_list", "setsockopt",
			// "sigaltstack", "tgkill", "uname", "write",

			// similar to inet pledge
			// "connect", "getpeername", "getsockname", "socket",

			// similar to rpath pledge
			// "getdents64", "newfstatat", "openat", "readlinkat",
		},
	}

	spec := seccomp.Filter{
		NoNewPrivs: true,
		Flag:       seccomp.FilterFlagTSync,
		Policy: seccomp.Policy{
			// By default goroutines don't play well with seccomp. Program will
			// hang when underlying thread is terminated silently.
			// We need to kill process - see:
			// <https://github.com/golang/go/issues/3405#issuecomment-750816828>
			DefaultAction: seccomp.ActionKillProcess,
			Syscalls:      []seccomp.SyscallGroup{whitelist},
		},
	}

	if err := seccomp.LoadFilter(spec); err != nil {
		return fmt.Errorf("%w: %w", ErrNoSandbox, err)
	}

	return nil
}

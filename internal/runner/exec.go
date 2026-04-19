package runner

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Options configures how the subprocess is run.
type Options struct {
	Args   []string
	Env    []string
	Stdin  bool
}

// Run executes the given command with the provided environment, forwarding
// signals and returning the exit code.
func Run(opts Options) (int, error) {
	if len(opts.Args) == 0 {
		return 1, errors.New("runner: no command specified")
	}

	path, err := exec.LookPath(opts.Args[0])
	if err != nil {
		return 1, err
	}

	cmd := exec.Command(path, opts.Args[1:]...)
	cmd.Env = opts.Env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if opts.Stdin {
		cmd.Stdin = os.Stdin
	}

	if err := cmd.Start(); err != nil {
		return 1, err
	}

	// Forward OS signals to the child process.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sigCh {
			_ = cmd.Process.Signal(sig)
		}
	}()

	err = cmd.Wait()
	signal.Stop(sigCh)
	close(sigCh)

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), nil
	}
	if err != nil {
		return 1, err
	}
	return 0, nil
}

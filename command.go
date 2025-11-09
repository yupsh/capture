package command

import (
	"context"
	"io"

	gloo "github.com/gloo-foo/framework"
)

type command struct {
	stdout io.Writer
	stderr io.Writer
}

// Capture creates a command that captures stdin to the provided writers.
// This is useful as a pipeline sink when you want to capture output instead of
// writing to os.Stdout/os.Stderr.
//
// Example:
//
//	var stdout, stderr bytes.Buffer
//	pipeline := gloo.Pipe(
//	    grep.Grep("ERROR"),
//	    sort.Sort(),
//	    capture.Capture(&stdout, &stderr),
//	)
//	gloo.MustRun(pipeline)
//	// Now stdout and stderr contain the captured output
func Capture(stdout, stderr io.Writer) gloo.Command {
	return command{
		stdout: stdout,
		stderr: stderr,
	}
}

func (c command) Executor() gloo.CommandExecutor {
	return func(ctx context.Context, stdin io.Reader, _, _ io.Writer) error {
		// Copy stdin to the provided stdout writer
		_, err := io.Copy(c.stdout, stdin)
		if err != nil {
			// If there's an error copying, write it to our stderr
			if c.stderr != nil {
				c.stderr.Write([]byte("capture: " + err.Error() + "\n"))
			}
			return err
		}
		return nil
	}
}

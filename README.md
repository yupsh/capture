# yup.capture

A command that captures pipeline output to custom `io.Writer` destinations instead of sending to stdout/stderr.

## Purpose

The `capture` command serves as a pipeline sink, allowing you to redirect command output to buffers, files, or any `io.Writer` implementation. This is particularly useful for:

- Testing pipelines and capturing their output
- Processing command output programmatically
- Storing results in memory for further manipulation
- Custom logging and output handling

## Usage

```go
import (
    "bytes"

    yup "github.com/yupsh/framework"
    "github.com/yupsh/capture"
    "github.com/yupsh/grep"
    "github.com/yupsh/pipe"
    "github.com/yupsh/sort"
)

func main() {
    // Create buffers to capture output
    var stdout, stderr bytes.Buffer

    // Build a pipeline that captures output
    pipeline := pipe.Pipeline(
        grep.Grep("ERROR"),
        sort.Sort(),
        capture.Capture(&stdout, &stderr),
    )

    // Run the pipeline
    err := yup.Run(pipeline)
    if err != nil {
        log.Fatal(err)
    }

    // Process the captured output
    fmt.Printf("Captured %d bytes: %s\n", stdout.Len(), stdout.String())
}
```

## API

### `Capture(stdout, stderr io.Writer) yup.Command`

Creates a capture command that routes:
- stdin → `stdout` writer
- error messages → `stderr` writer

**Parameters:**
- `stdout`: `io.Writer` - destination for captured standard output
- `stderr`: `io.Writer` - destination for captured error output

**Returns:**
- `yup.Command` - a command that can be used in pipelines

## Examples

### Capture to bytes.Buffer

```go
var out bytes.Buffer
pipeline := pipe.Pipeline(
    cat.Cat("file.txt"),
    capture.Capture(&out, io.Discard),
)
yup.MustRun(pipeline)
fmt.Println(out.String())
```

### Merge stdout and stderr to same writer

```go
// Equivalent to shell's 2>&1 - merge both streams to one writer
var combined bytes.Buffer
pipeline := pipe.Pipeline(
    grep.Grep("pattern"),
    capture.Capture(&combined, &combined),
)
yup.MustRun(pipeline)
// combined contains both stdout and stderr
```

### Capture to multiple destinations

```go
var buf bytes.Buffer
file, _ := os.Create("output.txt")
defer file.Close()

// Use io.MultiWriter to write to both
multi := io.MultiWriter(&buf, file)

pipeline := pipe.Pipeline(
    grep.Grep("pattern"),
    capture.Capture(multi, os.Stderr),
)
yup.MustRun(pipeline)
```

### Testing with capture

```go
func TestPipeline(t *testing.T) {
    var out, err bytes.Buffer

    pipeline := pipe.Pipeline(
        grep.Grep("test"),
        capture.Capture(&out, &err),
    )

    if err := yup.Run(pipeline); err != nil {
        t.Fatal(err)
    }

    if out.String() != "expected output\n" {
        t.Errorf("unexpected output: %s", out.String())
    }
}
```

## Behavior

- **Input**: Reads from stdin (the output of the previous command in the pipeline)
- **Output**: Writes to the provided `stdout` writer
- **Errors**: Writes to the provided `stderr` writer
- **Return**: Returns any error that occurs during copying

## Notes

- The capture command is always the last command in a pipeline (it's a sink)
- Both stdout and stderr writers must be non-nil
- Uses `io.Copy` for efficient streaming
- Does not buffer the entire output in memory unless you use a buffer like `bytes.Buffer`
- Compatible with any `io.Writer` implementation

## See Also

- [tee](../tee) - Write output to files while also passing it through
- [pipe](../pipe) - Pipeline composition
- [framework](../framework) - Core command framework

## License

GNU Affero General Public License v3.0 - see [LICENSE](LICENSE) file for details.


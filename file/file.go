package file

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/brentp/xopen"
)

// ReadFileWithBuffer reads file with buffer.
//
// the default value of line parser function fn is:
//
// 	fn = func(line string) (string, bool) {
// 		return line, true
// 	}
//
// A common one
//
// 	fn := func(line string) (string, bool) {
// 		line = strings.TrimSpace(line)
// 		if line == "" {
// 			return "", false
// 		}
// 		return line, ture
// 	}
func ReadFileWithBuffer(file string, batchSize int, bufferSize int, fn func(string) (string, bool)) (<-chan []string, error) {
	reader, err := xopen.Ropen(file)
	if err != nil {
		return nil, err
	}

	if batchSize <= 0 {
		batchSize = 1000000
	}
	if bufferSize < 0 {
		bufferSize = 0
	}

	if fn == nil {
		fn = func(line string) (string, bool) {
			return line, true
		}
	}

	ch := make(chan []string, bufferSize)
	batch := make([]string, batchSize)

	go func() {
		var (
			i    int
			line string
			err  error
		)
		for {
			line, err = reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					result, ok := fn(strings.TrimRight(line, "\n"))
					if ok {
						batch[i] = result
						i++
					}

					ch <- batch[0:i]
					// fmt.Println("sent", len(batch[0:i]))
					close(ch)
					break
				} else {
					fmt.Fprintln(os.Stderr, err)
					close(ch)
					os.Exit(-1)
					break
				}
			}

			result, ok := fn(strings.TrimRight(line, "\n"))
			if !ok {
				continue
			}
			batch[i] = result
			i++
			if i == batchSize {
				ch <- batch
				// fmt.Println("sent", len(batch))
				batch = make([]string, batchSize)
				i = 0
			}
		}
	}()

	return ch, nil
}

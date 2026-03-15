package cut

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

func Cut(reader io.Reader, writer io.Writer, opts Options) error {
	newReader := bufio.NewReader(reader)
	newWriter := bufio.NewWriter(writer)

	for {
		line, err := newReader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				if line != "" {
					if pErr := processLine(newWriter, strings.TrimRight(line, "\r\n"), opts); pErr != nil {
						return pErr
					}
				}
				break
			}
			return err
		}
		if pErr := processLine(newWriter, strings.TrimRight(line, "\r\n"), opts); pErr != nil {
			return pErr
		}
	}
	return newWriter.Flush()
}

func processLine(w *bufio.Writer, line string, opts Options) error {
	if !strings.Contains(line, opts.Delimiter) {
		if !opts.Separated {
			if _, err := w.WriteString(line); err != nil {
				return err
			}
			return w.WriteByte('\n')
		}
		return nil
	}

	cols := strings.Split(line, opts.Delimiter)
	var selected []string
	for _, f := range opts.Fields {
		if f >= 1 && f <= len(cols) {
			selected = append(selected, cols[f-1])
		}
	}
	if _, err := w.WriteString(strings.Join(selected, opts.Delimiter)); err != nil {
		return err
	}
	return w.WriteByte('\n')
}

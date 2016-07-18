// Package gotail provides high performing tail-like behavior for tailing files.
package gotail

import (
	"bufio"
	"strings"

	"io"
	"log"
	"os"
	"time"

	fsnotify "gopkg.in/fsnotify.v1"
)

type Tail struct {
	Lines chan string

	reader  *bufio.Reader
	watcher *fsnotify.Watcher
	fname   string
	file    *os.File
	config  Config
}

type Config struct {
	Timeout int
}

// NewTail creates a new Tail Object.  During initialization, it checks to see
// If the file exists.  If it doesn't, NewTail sleeps up to Config.Timeout seconds
// before returning an error.  If the file exists, then NewTail attaches an open file handle
// and a watcher to the file for new notifications.
func NewTail(fname string, config Config) (*Tail, error) {
	tail := &Tail{
		Lines:  make(chan string),
		fname:  fname,
		config: config,
	}

	err := tail.openAndWatch()

	if err != nil {
		return nil, err
	}

	tail.listenAndReadLines()

	return tail, nil
}

// Close closes the tail object when finished, closing the file handle and watcher
func (t *Tail) Close() {
	t.file.Close()
	t.watcher.Close()
}

// openAndWatch continually polls the target file to try to set an open file handler and watcher.
// If the timeout is reached, it sends the error back to the timeout signal
// and the function returns an error.  If no error is detected, it returns immediately.
func (t *Tail) openAndWatch() error {
	var err error
	var newFile bool

	timeout := make(chan error, 1)

	go func() {
		for {
			err = t.openFile(newFile)
			if err != nil {
				if os.IsNotExist(err) && newFile == false {
					newFile = true
				}

				if t.config.Timeout == 0 {
					timeout <- err
					return
				} else {
					continue
				}

			}
			newFile = false

			err = t.watchFile()

			if err == nil {
				timeout <- nil
				return
			}
		}
	}()

	if t.config.Timeout != 0 {
		go func() {
			time.Sleep(time.Duration(t.config.Timeout) * time.Second)

			timeout <- err
		}()
	}

	select {
	case err := <-timeout:
		if err != nil {
			return err
		}
	}

	return nil
}

// openFile opens a file and finds the offset byte to start reading.
// If it's a new file that has been created after Tail is following,
// then it processes the entire file first.
// This is because sometimes, a new file is considered "MODIFY" and
// file.Seek will automatically point to the last byte of the file,
// causing it to skip the first line.
func (t *Tail) openFile(newFile bool) (err error) {
	if t.file != nil {
		t.file.Close()
	}

	t.file, err = os.Open(t.fname)
	if err != nil {
		return err
	}

	if !newFile {
		_, err = t.file.Seek(0, 2)
	}

	if err != nil {
		return err
	}

	t.reader = bufio.NewReader(t.file)

	return nil
}

// watchFile assigns a new fsnotify watcher to the file if possible.
// It it watches for any signals that lead to file change, and responds accordingly.
func (t *Tail) watchFile() error {
	if t.watcher != nil {
		t.watcher.Close()
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	t.watcher = watcher

	err = t.watcher.Add(t.fname)

	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case evt := <-t.watcher.Events:
				if evt.Op == fsnotify.Create || evt.Op == fsnotify.Rename || evt.Op == fsnotify.Remove {
					if err = t.openAndWatch(); err != nil {
						log.Fatalln("open and watch failed:", err)
					}
				}
			case err := <-t.watcher.Errors:
				if err != nil {
					log.Fatalln("Watcher err:", err)
				}
			}
		}
	}()

	return nil
}

// listenAndReadLines continually polls the file in question,
// reading any new lines that gets added to the file.
func (t *Tail) listenAndReadLines() {
	go func() {
		for {
			if t.reader == nil {
				continue
			}

			line, err := t.reader.ReadString('\n')

			if err == io.EOF {
				continue
			}

			t.Lines <- strings.TrimRight(line, "\n")
		}
	}()
}

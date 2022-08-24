package util

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

// Exist check file or dir exist
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// ReadFileLines read by line
func ReadFileLines(filename string) (lines []string, err error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	bio := bufio.NewReader(file)
	for {
		var line []byte

		line, _, err = bio.ReadLine()
		if err != nil {
			if err == io.EOF {
				file.Close()
				return lines, nil
			}
			return
		}

		lines = append(lines, string(line))
	}

}

// CurrentDir for working directory
func CurrentDir(path ...string) string {
	_, currentFilePath, _, _ := runtime.Caller(1)
	if len(path) == 0 {
		return filepath.Dir(currentFilePath)
	}
	return filepath.Join(filepath.Dir(currentFilePath), filepath.Join(path...))
}

func WaitTerm(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sig)
	<-sig
	cancel()
}

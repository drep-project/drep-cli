package common

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

	panic("not reached")
}

func IsFileExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return !fi.IsDir()
	}

	panic("not reached")
}

// EachChildFile get child fi  and process ,if get error after processing stop, if get a stop flag , stop
func EachChildFile(directory string, process func(path string) (bool, error)) error {
	fds, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, fi := range fds {
		if !fi.IsDir() {
			isContinue, err := process(path.Join(directory, fi.Name()))
			if !isContinue {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func EachDirectory(directory string, process func(path string) (bool, error)) error {
	fds, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, fi := range fds {
		if fi.IsDir() {
			isContinue, err := process(path.Join(directory, fi.Name()))
			if !isContinue {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func AbsolutePath(datadir string, filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(datadir, filename)
}

package util

import "os"

func FileOrPathIsExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil || os.IsExist(err)
}

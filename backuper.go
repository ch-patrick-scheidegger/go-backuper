package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	cp "github.com/otiai10/copy"
	//"os/dir.DirEntry"
	//"path/filepath"
	//"io/ioutil"
)

// https://pkg.go.dev/github.com/spf13/cobra?GOOS=windows
// https://github.com/etcd-io/etcd/blob/main/etcdctl/go.mod
// https://golang.google.cn/pkg/os/#DirEntry
func main() {
	fmt.Print("--- Hello World ---\n")
	directoryPath := "C:/Users/Pat/Desktop/testFolder/"
	entries, _ := os.ReadDir(directoryPath)
	for _, entry := range entries {
		fmt.Printf("Is dir: {%v}\n", entry.IsDir())
		fmt.Printf("Name: {%v}\n", entry.Name())
		if !entry.IsDir() {
			signature, _ := hashFile(directoryPath + entry.Name())
			fmt.Printf("SHA256 signature: {%v}\n", signature)
		}
	}
	// items, _ := ioutil.ReadDir(directoryPath)
	// os.Open(directoryPath)
}

// Cases to backup files
// 1. fromFile exists in toFile but fromHash != toHash
// 2. fromFile does not exists in toFile
// Cases to delete files
// 1. toFile does not exists in fromFile -> check if some files have same hash
// 2.
func backupDirectory(srcPath string, dstPath string) {
	fmt.Printf("--- %v ---", srcPath)
	srcFilesMap := createFileNameHashValueMap(srcPath)
	dstFilesMap := createFileNameHashValueMap(dstPath)
	// Copy non existing files and overwritte existing files where hash is not the same
	for srcFileName, srcHash := range srcFilesMap {
		if dstHash, exists := dstFilesMap[srcFileName]; !exists {
			//TODO copy srcFile dstFile
		} else if dstHash != srcHash {
			//TODO overwrite dstFile
		}
	}
	// Delete files that only exist in dst
	for dstFileName := range dstFilesMap {
		if _, exists := srcFilesMap[dstFileName]; !exists {
			toRemove := dstPath + dstFileName
			err := os.Remove(toRemove)
			if err != nil {
				fmt.Printf("Could not delete file %v because of %v", toRemove, err)
			}
		}
	}

	srcSubDirs, _ := getSubDirs(srcPath)
	dstSubDirs, _ := getSubDirs(dstPath)
	// Delete sub dirs that exist in dst but not in src
	for dstSubDir := range dstSubDirs {
		if !srcSubDirs.Contains(dstSubDir) {
			// TODO remove dstSubDir
			toRemove := dstPath + dstSubDir + "/"
			err := os.RemoveAll(toRemove)
			if err != nil {
				fmt.Printf("Could not delete dst dir %v because of %v", toRemove, err)
			} else {
				fmt.Print("Deleted ", toRemove)
			}
		}
	}
	// Copy sub dirs that do not exist in dst dir
	// Recursively call this function to backup sub dirs
	for srcSubDir := range srcSubDirs {
		if !dstSubDirs.Contains(srcSubDir) {
			err := cp.Copy(srcSubDir, dstPath)
			if err != nil {
				fmt.Printf("Could not copy dir %v because of %v", srcSubDir, err)
			} else {
				fmt.Print("Copied ", srcSubDir)
			}
		} else {
			subDirToCheck := srcSubDir + "/"
			backupDirectory(srcPath+subDirToCheck, dstPath+subDirToCheck)
		}
	}
}

func getSubDirs(dirPath string) (Set, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var set Set = NewSet()
	for _, entry := range entries {
		if entry.IsDir() {
			set.Add(entry.Name())
		}
	}
	return set, nil
}

func createFileNameHashValueMap(dirPath string) map[string][32]byte {
	dict := make(map[string][32]byte)
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Print()
	}
	for _, entry := range files {
		if !entry.IsDir() {
			hash, err := hashFile(dirPath + entry.Name())
			if err != nil {
				fmt.Print("Failed to hash file " + dirPath + entry.Name() + " error: " + err.Error())
			} else {
				dict[entry.Name()] = hash
			}
		}
	}
	return dict
}

// Returned hash is of size 32 Byte
func hashFile(filePath string) ([32]byte, error) {
	var result *[32]byte
	cryptoHash := sha256.New()
	file, err := os.Open(filePath)
	if err != nil {
		return *result, err
	}

	buffer := make([]byte, 32*1024)
	for {
		if n, err := file.Read(buffer); err == nil {
			cryptoHash.Write(buffer[:n])
		} else if err == io.EOF {
			break
		} else {
			return *result, err
		}
	}
	// Cast the slice to an array according to https://stackoverflow.com/a/30285971/12828371
	result = (*[32]byte)(cryptoHash.Sum(nil))
	//var result [32]byte
	//copy(result[:], cryptoHash.Sum(nil))
	return *result, nil
}

package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// TODO remove me
func toBeRemoved() {
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
}

// https://pkg.go.dev/github.com/spf13/cobra?GOOS=windows
// https://github.com/etcd-io/etcd/blob/main/etcdctl/go.mod
// https://golang.google.cn/pkg/os/#DirEntry
func main() {
	fmt.Println("\n--- Backuper ---")
	srcPath := "C:/Users/Pat/Desktop/src/"
	dstPath := "C:/Users/Pat/Desktop/dst/"
	backupDirectory(srcPath, dstPath)
}

func backupDirectory(srcPath string, dstPath string) {
	fmt.Printf("Dir %v\n", srcPath)
	srcFilesMap := createFileNameHashValueMap(srcPath)
	dstFilesMap := createFileNameHashValueMap(dstPath)
	// Copy non existing files and overwritte existing files where hash is not the same
	for srcFileName, srcHash := range srcFilesMap {
		src := srcPath + srcFileName
		dst := dstPath + srcFileName
		if dstHash, exists := dstFilesMap[srcFileName]; !exists {
			err := copyFile(src, dst)
			if err != nil {
				fmt.Printf("Could not copy %v because of %v\n", src, err)
			} else {
				fmt.Print("Copied ", src, "\n")
			}
		} else if dstHash != srcHash {
			err := copyFile(src, dst)
			if err != nil {
				fmt.Printf("Could not update %v because of %v\n", src, err)
			} else {
				fmt.Print("Updated ", src, "\n")
			}
		}
	}
	// Delete files that only exist in dst
	for dstFileName := range dstFilesMap {
		if _, exists := srcFilesMap[dstFileName]; !exists {
			toRemove := dstPath + dstFileName
			err := os.Remove(toRemove)
			if err != nil {
				fmt.Printf("Could not delete file %v because of %v\n", toRemove, err)
			} else {
				fmt.Println("Deleted ", toRemove)
			}
		}
	}

	srcSubDirs, err := getSubDirs(srcPath)
	if err != nil {
		fmt.Printf("Could not read src subdirs of %v because of %v", srcPath, err)
	}
	dstSubDirs, err := getSubDirs(dstPath)
	if err != nil {
		fmt.Printf("Could not read dst subdirs of %v because of %v", dstPath, err)
	}
	// Delete sub dirs that exist in dst but not in src
	for dstSubDir := range dstSubDirs {
		if !srcSubDirs.Contains(dstSubDir) {
			toRemove := dstPath + dstSubDir + "/"
			err := os.RemoveAll(toRemove)
			if err != nil {
				fmt.Printf("Could not delete dst dir %v because of %v\n", toRemove, err)
			} else {
				fmt.Print("Deleted dir ", toRemove, "\n")
			}
		}
	}
	// Copy sub dirs that do not exist in dst dir
	// Recursively call this function to backup sub dirs
	for srcSubDir := range srcSubDirs {
		if !dstSubDirs.Contains(srcSubDir) {
			err := copyDir(srcPath+srcSubDir+"/", dstPath+srcSubDir+"/")
			if err != nil {
				fmt.Printf("Could not copy dir %v because of %v\n", srcSubDir, err)
			} else {
				fmt.Print("Copied dir ", srcPath+srcSubDir, "\n")
			}
		} else {
			subDirToCheck := srcSubDir + "/"
			backupDirectory(srcPath+subDirToCheck, dstPath+subDirToCheck)
		}
	}
}

func copyDir(srcDirPath, dstDirPath string) error {
	err := os.Mkdir(dstDirPath, 0777)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(srcDirPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			if err := copyFile(srcDirPath+entry.Name(), dstDirPath+entry.Name()); err != nil {
				return err
			}
		} else {
			if err := copyDir(srcDirPath+entry.Name()+"/", dstDirPath+entry.Name()+"/"); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(srcFileName, dstFileName string) error {
	srcFile, err := os.Open(srcFileName)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstFileName)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
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
			fileToHash := dirPath + entry.Name()
			hash, err := hashFile(fileToHash)
			if err != nil {
				fmt.Printf("Failed to hash file %v because of %v\n", fileToHash, err)
			} else {
				dict[entry.Name()] = hash
			}
		}
	}
	return dict
}

// Returned hash is of size 32 Byte
func hashFile(filePath string) ([32]byte, error) {
	var resultHash *[32]byte
	cryptoHash := sha256.New()
	file, err := os.Open(filePath)
	if err != nil {
		return *resultHash, err
	}
	defer file.Close()
	buffer := make([]byte, 32*1024)
	for {
		if n, err := file.Read(buffer); err == nil {
			cryptoHash.Write(buffer[:n])
		} else if err == io.EOF {
			break
		} else {
			return *resultHash, err
		}
	}
	// Cast the slice to an array according to https://stackoverflow.com/a/30285971/12828371
	resultHash = (*[32]byte)(cryptoHash.Sum(nil))
	return *resultHash, nil
}

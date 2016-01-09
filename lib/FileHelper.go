package lib

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"strings"
)

type FileHelper struct {
}

//读文件
func (helper *FileHelper) CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//写文件
func (helper *FileHelper) WriteFile(filename string, content string) error {
	lock := new(sync.RWMutex)
	lock.Lock()

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0x644)

	if err != nil {
		fmt.Println("file open error:", err)
		return err
	}

	_, err = file.WriteString(content + "\n")
	file.Close()

	lock.Unlock()
	return err
}

//读文件
func (helper *FileHelper) ReadFile(filename string) (string, error) {
	lock := new(sync.RWMutex)
	lock.RLock()

	var content string
	file, err := os.Open(filename)

	if err != nil {
		fmt.Println("file open error:", err)
		return content, err
	}

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
        
        content+=string(line)+","
	}
    
    content=strings.TrimRight(content,",")
	file.Close()
	lock.RUnlock()
    
    fmt.Println(content)
	return content, err
}

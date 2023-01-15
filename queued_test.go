package queued

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var expectedDownloadFileCount = 1000
var downloadFactor = 11

var downloadedLocker sync.Mutex
var filesDownloaded = map[string]bool{}
var downloadedCount = 0

func init() {
	rand.Seed(time.Now().UnixNano())
}

func download(index int, url string) (string, error) {
	log.Printf("%03d: start downloading %s\n", index, url)

	if _, ok := filesDownloaded[url]; ok {
		return url, nil
	}

	downloadedLocker.Lock()
	downloadedCount += 1
	filesDownloaded[url] = true
	downloadedLocker.Unlock()

	log.Printf("%03d: downloaded %s\n", index, url)

	time.Sleep(3 * time.Second)

	return url, nil
}

func TestLock(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < expectedDownloadFileCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			fileName := fmt.Sprintf("test.%03d.txt", i%downloadFactor)
			_, err := Queued(fileName, func() (string, error) {
				return download(i, fileName)
			})
			if err != nil {
				panic(err)
			}
		}(i)
	}

	wg.Wait()

	log.Printf("actually downloaded %d, expected %d", downloadedCount, downloadFactor)

	if downloadedCount != downloadFactor {
		t.Fatalf("download multiple times\n")
	}
}

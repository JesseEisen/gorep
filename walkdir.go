package main  

import (
    "fmt"
    "os"
    "sync"
    "runtime"
)

var wg sync.WaitGroup
var semaphoreChan = make(chan struct{}, runtime.GOMAXPROCS(runtime.NumCPU()))

func main(){
    wg.Add(1)
    lsFiles(".")
    wg.Wait()
}

func lsFiles(dir string) {
    semaphoreChan <- struct{} {}
    file, err := os.Open(dir)
    if err != nil {
        fmt.Println("error opening directory")
    }
    defer file.Close()

    files, err := file.Readdir(-1)
    if err != nil {
        fmt.Println("error readding directory")
    }

    for _, f := range files {
        if f.IsDir() {
            wg.Add(1)
            go lsFiles(dir + "/" + f.Name())
        }
        fmt.Println(dir + "/" + f.Name())
    }

    defer func() {
        <-semaphoreChan
        wg.Done()
    }()
}

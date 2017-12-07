package main  

import (
    "fmt"
    "os"
    "sync"
    "runtime"
    "bufio"
    "strings"
)

var wg sync.WaitGroup
var semaphoreChan = make(chan struct{}, runtime.GOMAXPROCS(runtime.NumCPU()))
var magicNumber = []byte{0x7F,0x45,0x4C,0x46}

func main(){
    fmt.Println("\n")
    wg.Add(1)
    lsFiles(os.Args[1], os.Args[2])
    wg.Wait()
    fmt.Println("\n")
}

func lsFiles(dir string, pattern string) {
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
            go lsFiles(dir + "/" + f.Name(), pattern)
        } else if f.Mode().IsRegular() {
            SearchFile(dir + "/" + f.Name(), pattern)
        }
    }

    defer func() {
        <-semaphoreChan
        wg.Done()
    }()
}

func SearchFile(filename, pattern string) {
    handler, err:= os.Open(filename)
    if err != nil {
        fmt.Println("error in open file", filename, err)
    }

    reader := bufio.NewReader(handler)
    var lineno  uint32 = 1

    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        first4 := []rune(line)
        if strings.Compare(string(first4[:4]), string(magicNumber)) == 0 {
            break;
        }

        lineno += 1
    
        if strings.Contains(line, pattern) {
            fmt.Printf("%s - %d: %s",filename, lineno, line)
        }
    }

    handler.Close()
    
}

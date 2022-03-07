package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/briandowns/spinner"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type ResultStruct struct {
	File string
	Line int
}

var noExec bool = false

func searchTextInFile(path, text string) ([]int, error) {
	var result []int
	f, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	line := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), text) {
			result = append(result, line)
		}
		line++
	}
	if err := scanner.Err(); err != nil {
		return result, err
	}
	return result, nil
}

func isExec(info fs.DirEntry) bool {
	nfo, err := info.Info()
	if err != nil {
		return false
	}
	return nfo.Mode().Perm()&0111 != 0
}

func readPath(path string, filesChan chan string, except *[]string) {
	err := filepath.WalkDir(path, func(f string, info fs.DirEntry, err error) error {
		if noExec && isExec(info) {
			return nil
		}
		if !info.IsDir() && !stringInSlice(f, except) {
			// fmt.Printf("%s\n", f)
			filesChan <- f
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func searchInFiles(out chan ResultStruct, c <-chan string, str string) {
	opened := true
	file := ""
	for opened {
		file, opened = <-c
		lines, err := searchTextInFile(file, str)
		if err == nil {
			for _, n := range lines {
				if n > 0 {
					out <- ResultStruct{file, n}
				}
			}
		}
	}
	close(out)
}

func main() {
	fmt.Print(
		"                                                       _\n" +
			" ___ _   _ _ __   ___ _ __      ___  ___  __ _ _ __ ___| |__\n" +
			"/ __| | | | '_ \\ / _ \\ '__|____/ __|/ _ \\/ _` | '__/ __| '_ \\\n" +
			"\\__ \\ |_| | |_) |  __/ | |_____\\__ \\  __/ (_| | | | (__| | | |\n" +
			"|___/\\__,_| .__/ \\___|_|       |___/\\___|\\__,_|_|  \\___|_| |_|\n" +
			"          |_|\n")
	var wgSearchers, wgPrinter sync.WaitGroup
	except := ""
	//var exceptList *[]string

	exceptList := readExcludeFile()

	flag.BoolVar(&noExec, "no-exec", false, "exclude executables")
	flag.StringVar(&except, "except", "", "list of files/folders pattern to avoid")
	flag.Parse()
	flag.Usage = func() {
		fmt.Printf("usage : %s [--no-exec] [--except=e1,e2..] searchText path1 path2 ... pathN\n", os.Args[0])
		os.Exit(0)
	}

	if except != "" {
		tmp := cleanDuplicatesAndEmpties(append(*exceptList, strings.Split(except, ",")...))
		exceptList = &tmp
	}
	fmt.Printf("%s\n%s\n%d\n", except, *exceptList, len(*exceptList))

	if flag.NArg() < 2 {
		flag.Usage()
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	time.Sleep(4 * time.Second)
	defer s.Stop()

	search := ""
	filesChan := make(chan string)
	for i, arg := range flag.Args() {
		if i == 0 {
			search = arg
		} else {
			wgSearchers.Add(1)
			go func(path string) {
				defer wgSearchers.Done()
				colorize(ColorRed, fmt.Sprintf("Collecting files in path %s ...\n", path))
				readPath(path, filesChan, exceptList)
			}(arg)
		}
	}
	resultChan := make(chan ResultStruct)
	go searchInFiles(resultChan, filesChan, search)

	wgPrinter.Add(1)
	go func() {
		defer wgPrinter.Done()
		ok := true
		var i ResultStruct
		for ok {
			i, ok = <-resultChan
			if !ok {
				break
			}
			colorize(ColorGreen, fmt.Sprintf("\t* Found %s in %s at line %d\n", search, i.File, i.Line))
		}
	}()

	wgSearchers.Wait()
	close(filesChan)
	wgPrinter.Wait()
}

package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var problems string
	var timeLimit time.Duration
	flag.StringVar(&problems, "problems", "problems.csv", "CSV formatted file with quiz questions")
	flag.DurationVar(&timeLimit, "limit", 30*time.Second, "Time limit for the quiz")
	flag.Parse()

	quiz, _ := ioutil.ReadFile(problems)
	reader := csv.NewReader(strings.NewReader(string(quiz)))
	inputReader := bufio.NewReader(os.Stdin)

	var nCorrect, nIncorrect int = 0, 0

	fmt.Print("PRESS ENTER TO BEGIN")
	inputReader.ReadString('\n')

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(timeLimit)
		timeout <- true
	}()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var question, answer string = record[0], record[1]
		fmt.Print(question + "? ")

		inCh := make(chan string, 1)
		go func() {
			in, _ := inputReader.ReadString('\n')
			in = strings.TrimSpace(in)
			inCh <- in
		}()

		var quit bool = false
		select {
		case offer := <-inCh:
			if offer == answer {
				nCorrect++
			} else {
				nIncorrect++
			}
		case <-timeout:
			quit = true
		}

		if quit {
			fmt.Println()
			break
		}
	}

	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		}
		nIncorrect++
	}

	fmt.Printf("Answered %d of %d questions correctly!\n", nCorrect, nCorrect+nIncorrect)
}

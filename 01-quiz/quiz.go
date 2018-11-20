package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type Question struct {
	Text   string
	Answer string
}

func main() {
	csvPath := flag.String("csv", "problems.csv", `a csv file in the format of 'question,answer' (default "problems.csv")`)
	shuffle := flag.Bool("shuffle", true, "a boolean value to shuffle the questions from the csv (default false).")
	timeLimit := flag.Int("limit", 30, "the time limit for te quiz in seconds (default 30)")
	flag.Parse()

	quiz, err := getQuestions(*csvPath)
	if err != nil {
		log.Fatal(err)
	}

	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(quiz), func(i, j int) { quiz[i], quiz[j] = quiz[j], quiz[i] })
	}

	var score int
	var wg sync.WaitGroup
	wg.Add(1)

	fmt.Println("PRESS ENTER TO START")
	fmt.Scanln()

	go func() {
		timeOver := time.After(time.Duration(*timeLimit) * time.Second)
		select {
		case <-timeOver:
			fmt.Printf("\nTime is OVER!!\n")
			wg.Done()
			break
		}
	}()

	go func() {
		for i, q := range quiz {
			fmt.Printf("Problem #%d: %s = ", i+1, q.Text)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(strings.TrimSpace(response)) == strings.ToLower(q.Answer) {
				score++
			}
		}
		wg.Done()
	}()

	wg.Wait()

	fmt.Printf("You scored %d out of %d.\n", score, len(quiz))
}

func getQuestions(path string) ([]Question, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	quiz := []Question{}

	for _, row := range records {
		q := Question{
			Text:   row[0],
			Answer: row[1],
		}

		quiz = append(quiz, q)
	}

	return quiz, nil
}

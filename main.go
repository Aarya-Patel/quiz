package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

/* ----- Structs ----- */
type Quiz struct {
	questions   []QnA
	submissions []QuestionSubmission
}
type QnA struct {
	prompt, answer string
}

type QuestionSubmission struct {
	question   *QnA
	userAnswer string
	isCorrect  bool
}

/* ----- Globals ----- */
var (
	filename  = flag.String("filename", "problems.csv", "Filename of the problem set.")
	timeLimit = flag.Int("time", 30, "Time limit of the quiz.")
	quiz      Quiz
	SCANNER   *bufio.Scanner = bufio.NewScanner(os.Stdin)
)

func main() {
	flag.Parse()

	csvFile, err := openCSV(*filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	questions, err := readCSVAndFormatQuestions(csvFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	quiz.questions = questions

	fmt.Println("--------------------------------")
	fmt.Println("|------------ Quiz ------------|")
	fmt.Println("--------------------------------")
	promptQuizStart()

	done := make(chan interface{})
	time.AfterFunc(time.Duration(*timeLimit)*time.Second, func() {
		close(done)
	})

	go func() {
		for idx := 0; idx < len(quiz.questions); idx++ {
			// User finished the quiz
			if idx >= len(quiz.questions) {
				close(done)
				return
			}

			userSubmission, err := askUserQuestion(quiz.questions[idx])
			if err != nil {
				fmt.Println(err)
				return
			}
			quiz.submissions = append(quiz.submissions, userSubmission)

		}
	}()

awaitingFinish:
	for {
		select {
		case <-done:
			break awaitingFinish
		default:
		}
	}

	var totalScore uint8 = 0
	for _, submission := range quiz.submissions {
		if submission.isCorrect {
			totalScore += 1
		}
	}
	fmt.Printf("\nYou scored %d/%d on this quiz.\n", totalScore, len(quiz.questions))
}

func promptQuizStart() (bool, error) {
	fmt.Print("Please press enter to begin the quiz.")
	ok := SCANNER.Scan()
	if !ok {
		return false, SCANNER.Err()
	}
	return true, nil

}

func askUserQuestion(qna QnA) (QuestionSubmission, error) {
	fmt.Print(qna.prompt + ": ")
	ok := SCANNER.Scan()
	if !ok {
		return QuestionSubmission{}, SCANNER.Err()
	}

	userAns := strings.Trim(SCANNER.Text(), " ")
	isCorrect := userAns == qna.answer

	return QuestionSubmission{&qna, userAns, isCorrect}, nil
}

func readCSVAndFormatQuestions(csvFile io.Reader) ([]QnA, error) {
	questions, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, err
	}

	if len(questions) == 0 {
		return nil, fmt.Errorf("There are no questions in the quiz!")
	}

	var quizQuestions []QnA
	for _, questionAndAnswer := range questions {
		prompt := questionAndAnswer[0]
		answer := questionAndAnswer[1]

		quizQ := QnA{prompt, answer}

		quizQuestions = append(quizQuestions, quizQ)
	}
	return quizQuestions, nil
}

func openCSV(filename string) (io.Reader, error) {
	return os.Open(filename)
}

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
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

/* ----- IO Scanner ----- */
var SCANNER *bufio.Scanner = bufio.NewScanner(os.Stdin)

func main() {
	fmt.Println("Golang is awesome!")

	csvFile, err := openCSV("problems.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	questions, err := readCSVAndFormatQuestions(csvFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	var quiz Quiz
	quiz.questions = questions
	var totalScore uint8 = 0

	for _, question := range quiz.questions {
		userSubmission, err := askUserQuestion(question)
		if err != nil {
			fmt.Println(err)
			return
		}

		quiz.submissions = append(quiz.submissions, userSubmission)
		if userSubmission.isCorrect {
			totalScore += 1
		}
	}

	fmt.Printf("You scored %d/%d on this quiz.\n", totalScore, len(quiz.questions))

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

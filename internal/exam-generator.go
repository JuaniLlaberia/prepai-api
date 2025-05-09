package internal

import (
	"encoding/json"
	"fmt"

	"google.golang.org/genai"
	"prepai.app/configs"
)

type ExamQuestion struct {
	Question    string   `json:"question"`
	Options     []string `json:"options"`
	Correct     int64    `json:"correct"`
	Explanation string   `json:"explanation"`
}

type ExamResponse struct {
	Title     string         `json:"title"`
	Questions []ExamQuestion `json:"questions"`
}

func GenerateExam(subject string, difficulty string, examType string) (ExamResponse, error) {
	prompt := fmt.Sprintf(`
		Generate a %v exam on the topic %v, with %v difficulty.
		- If the exam type is multiple choice, generate 4 options per question.
		- If the exam type is true/false, generate only 2 options: "True" and "False".

		Based on the difficulty level:
		- "easy": generate 10 questions
		- "medium": generate 15 questions
		- "hard": generate 20 questions

		For each question:
		- Randomly shuffle the answer options so the correct one is not always in the same index.
		- Provide the correct answer's index (0-based).
		- Make sure the correct answer value matches the position of the correct option after shuffling.
		- Provide an explanation (Explain in 3-4 lines why the correct answer is correct)
		- Format the output in the following JSON schema:
		{
			"title": string,
			"questions": [
				{
				"question": string,
				"options": [string],
				"correct": int64
				"explanation": string
				}
			]
		}
`, examType, subject, difficulty)

	result, err := configs.Gemini(genai.Text(prompt))
	if err != nil {
		return ExamResponse{}, err
	}

	var questions ExamResponse

	err = json.Unmarshal([]byte(result), &questions)
	if err != nil {
		return ExamResponse{}, err
	}

	return questions, nil
}

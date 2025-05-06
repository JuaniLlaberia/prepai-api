package internal

import (
	"encoding/json"
	"fmt"

	"prepai.app/configs"
)

type InterviewQuestion struct {
	Question string `json:"question"`
	Hint     string `json:"hint"`
	Type     string `json:"type"`
}

type InterviewResponse struct {
	Title     string              `json:"title"`
	Questions []InterviewQuestion `json:"questions"`
}

func GenerateInterview(jobRole string, jobLevel string, topics []string) (InterviewResponse, error) {
	prompt := fmt.Sprintf(`
		Generate 5 job interview questions for a role of %v with a %v. And a title for the interview.
		The interview topics are: %v.
		For each question provide:
		- The question.
		- A hint (Short text to help the interviewee).
		- Question type ("Behavioral", "Technical", "HR", etc)

		Follow this JSON schema:
		{
			"title": string,
			"questions": [
				{
					"question": string,
					"hint": string,
					"type": string
				}
			]
		}
	`, jobRole, jobLevel, topics)

	result, err := configs.Gemini(prompt)
	if err != nil {
		return InterviewResponse{}, err
	}

	var questions InterviewResponse

	err = json.Unmarshal([]byte(result), &questions)
	if err != nil {
		return InterviewResponse{}, err
	}

	return questions, nil
}

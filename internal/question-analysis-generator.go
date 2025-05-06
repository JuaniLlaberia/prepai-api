package internal

import (
	"encoding/json"
	"fmt"

	"prepai.app/configs"
)

type QuestionAnswer struct {
	Structure string   `json:"structure"`
	KeyPoints []string `json:"key_points"`
	Example   string   `json:"example"`
}

type QuestionAnalysisResponse struct {
	Type           string         `json:"type"`
	Difficulty     string         `json:"difficulty"`
	Explanation    string         `json:"explanation"`
	ExpectedLength string         `json:"expected_length"`
	IdealAnswer    QuestionAnswer `json:"ideal_answer"`
}

func GenerateQuestionAnalysis(question string) (QuestionAnalysisResponse, error) {
	prompt := fmt.Sprintf(`
		Analyze the following interview question: %v.

		Return a JSON object with the following fields:
		- "type": The type of the question. Choose one of: "Behavioral", "Technical", "HR", "Situational", or "Other".
		- "difficulty": One of: "Easy", "Medium", or "Hard".
		- "explanation": A 2-3 sentence explanation of what the question evaluates and why interviewers ask it.
		- "expected_length": How long in minutes should the interviewee take to answer.
		- "ideal_answer": An object that contains:
			- "structure": Describe the best format or method to answer the question (e.g., STAR, technical breakdown, etc.).
			- "key_points": A list of the most important concepts, points, or themes the answer should include.
			- "example": A sample ideal answer (5-8 lines) that would score highly in a real interview.

		Respond only in the following JSON format:
		{
			"type": string,
			"difficulty": string,
			"explanation": string,
			"expected_length": string,
			"ideal_answer": {
				"structure": string,
				"key_points": [string],
				"example": string
			}
		}
	`, question)

	result, err := configs.Gemini(prompt)
	if err != nil {
		return QuestionAnalysisResponse{}, err
	}

	var questionAnalysis QuestionAnalysisResponse

	err = json.Unmarshal([]byte(result), &questionAnalysis)
	if err != nil {
		return QuestionAnalysisResponse{}, err
	}

	return questionAnalysis, nil
}

package internal

import (
	"encoding/json"
	"fmt"

	"google.golang.org/genai"
	"prepai.app/configs"
)

type UserInterviewResponse struct {
	Question string
	Answer   string
}

type InterviewFeedback struct {
	Feedback   string  `json:"feedback"`
	Score      float64 `json:"score"`
	Suggestion string  `json:"suggestion"`
}

type InterviewFeedbackResponse struct {
	Feedbacks      []InterviewFeedback `json:"feedbacks"`
	Analysis       string              `json:"analysis"`
	Strengths      []string            `json:"strengths"`
	AreasToImprove []string            `json:"areas_to_improve"`
}

func GenerateInterviewFeedback(responses []UserInterviewResponse) (InterviewFeedbackResponse, error) {
	prompt := fmt.Sprintf(`
		Generate feedback on how the interviewee answered the following questions.

		This is the JSON containing the questions and answers: %v

		For each question, provide:
		- Feedback on how well the interviewee answered the question, considering vocabulary, technical terminology, structure, depth of knowledge, and relevance to the question.
		- The feedback must be between 3 to 5 sentences.
		- If the response is empty or missing, state clearly: "This question was not answered."
		- Provide a score from 1 to 10 (1 = very poor, 10 = excellent) based on the quality of the response.
		- Suggestion must give a direct and practical advice for how to improve the answer.

		Then, generate an overall interview analysis, taking into account:
		- Use of vocabulary and domain-specific terminology.
		- Clarity and confidence in communication.
		- Word repetition or redundancy.
		- Excessive use of filler words (e.g., "um", "like", "you know").
		- Overall ability to communicate thoughts effectively and professionally.
		- The overall analysis must be between 5 to 8 sentences.

		In addition provide:
		- Strengths and areas to improve.

		Format the output in the following JSON schema:
		{
		"feedbacks": [
			{
			"feedback": string,
			"score": int,
			"suggestion": string
			}
		],
		"analysis": string,
		"strengths": [string],
  		"areas_to_improve": [string]
		}
	`, responses)

	result, err := configs.Gemini(genai.Text(prompt))
	if err != nil {
		return InterviewFeedbackResponse{}, err
	}

	var feedback InterviewFeedbackResponse

	err = json.Unmarshal([]byte(result), &feedback)
	if err != nil {
		return InterviewFeedbackResponse{}, err
	}

	return feedback, nil
}

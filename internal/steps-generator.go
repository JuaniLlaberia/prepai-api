package internal

import (
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"
	"prepai.app/configs"
)

type Step struct {
	Title      string `json:"title"`
	Type       string `json:"type"`
	Order      int64  `json:"order"`
	Difficulty string `json:"difficulty"`
}

type StepsResponse struct {
	Steps []Step `json:"steps"`
}

func GenerateSteps(moduleTitle string, moduleDescription string, topics []string) (StepsResponse, error) {
	topicsStr := strings.Join(topics, ", ")

	prompt := fmt.Sprintf(`
		Create structured, gamified steps for the following module: %v. The module description is: "%v", with this topics: %v.
		The steps must help the candidate understand and practice for this module.

			- Include between 8 and 10 steps in total.
			- Each module should have:
				- Title (Descriptive of the step).
				- Type (What type of activity is: mock exam, open question, mock interview, lesson, etc).
				- Order (To sort the steps)
			- The last module needs to be an exam or interview to sum up this module.
			- Sort them from easier to harder in terms of difficulty (1: easiest and X: hardesr). This number is the "order" field and add the appropiate difficulty.

		Follow this JSON schema:
		{
			"steps": [
				{
					"title": string,
					"type": string,
					"order": int64,
					"difficulty": string
				}
			]
		}
	`, moduleTitle, moduleDescription, topicsStr)

	result, err := configs.Gemini(genai.Text(prompt))
	if err != nil {
		return StepsResponse{}, err
	}

	var steps StepsResponse

	err = json.Unmarshal([]byte(result), &steps)
	if err != nil {
		return StepsResponse{}, err
	}

	return steps, nil
}

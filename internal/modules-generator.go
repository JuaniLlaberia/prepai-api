package internal

import (
	"encoding/json"
	"fmt"
	"strings"

	"prepai.app/configs"
)

type Module struct {
	Title     string `json:"title"`
	Objective string `json:"objective"`
	Order     int64  `json:"order"`
	Topic     string `json:"topic"`
}

type ModuleResponse struct {
	Modules []Module `json:"modules"`
}

func GenerateModules(jobRole string, jobLevel string, jobDescription string, topics []string) (ModuleResponse, error) {
	topicsStr := strings.Join(topics, ", ")
	if len(jobDescription) > 500 {
		jobDescription = jobDescription[:497] + "..."
	}

	prompt := fmt.Sprintf(`
		Create structured, gamified modules for the following role: %v at a %v-level. The job description is: "%v".
		The modules must help the candidate prepare for a job interview for this role.

			- Include between 10 and 12 modules in total.
			- Each module should have:
				- Title (Descriptive of the module).
				- Objective (What is the aim of that module).
				- Topic (Topics included in the module separated by ",").
				- Order (To sort the modules)
			- Some modules should focus on technical skills, others on soft skills.
			- The last module needs to be: "Final challenge".
			- Sort them from easier to harder in terms of difficulty (1: easiest and X: hardesr). This number is the "order" field.
			- In addition to the role, level and description, this are topics that the interviewee needs to know: %v

		Follow this JSON schema:
		{
			"modules": [
				{
					"title": string,
					"objective": string,
					"topic": string,
					"order": int64
				}
			]
		}
	`, jobRole, jobLevel, jobDescription, topicsStr)

	result, err := configs.Gemini(prompt)
	if err != nil {
		return ModuleResponse{}, err
	}

	var modules ModuleResponse

	err = json.Unmarshal([]byte(result), &modules)
	if err != nil {
		return ModuleResponse{}, err
	}

	return modules, nil
}

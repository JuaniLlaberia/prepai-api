package configs

import (
	"context"

	"google.golang.org/genai"
)

type GeminiConfig struct {
	APIKey          string
	ModelName       string
	Temperature     *float32
	TopP            *float32
	TopK            *int32
	MaxOutputTokens *int32
	StopSequences   []string
	SafeSettings    []*genai.SafetySetting
}

func DefaultSafetySettings() []*genai.SafetySetting {
	return []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockThresholdBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockThresholdBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockThresholdBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockThresholdBlockMediumAndAbove,
		},
	}
}

// (*genai.Model, *genai.Client, error)

// func Gemini(ctx context.Context, config GeminiConfig) {
func Gemini(prompt string) (string, error) {
	apiKey := ProcessEnv("GEMINI_API_KEY")
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", err
	}

	temp := float32(0.5)
	topP := float32(0.85)
	topK := float32(64.0)
	maxOutputTokens := int32(8192)

	const systemInstructions = `
		You are an AI assistant designed to help users prepare for real-world job interviews. Your tasks include generating mock interviews, multiple-choice exams, career paths, and question banks across various professional domains. Your responses must be formatted in JSON, with field structures that change based on the function being called. Accuracy, clarity, and relevance to the user's context are essential.

		General Behavior Guidelines:
		- Always return output in valid JSON.
		- Do not include any explanatory text or commentary outside the JSON.
		- When uncertain about missing context (e.g., job role, experience level), use default values based on the topics and prompt.
		- Prioritize content that reflects real interview standards used by employers in the relevant industry.
		- Use clear, direct, and professional language suitable for job seekers at different levels.
		- Ensure all content is original and free from repetition or filler.
		- Focus on accuracy, especially in answer keys, explanations, and skill-level tagging.
		- Use diverse and balanced question types when generating output (e.g., behavioral, technical, situational) unless what is specified.
		- Align suggestions and content with industry norms, providing logical progression and realistic expectations.

		Formatting Rules:
		- Always wrap your entire output inside a valid JSON object or array depending on the endpoint requirements.
		- Use snake_case for all keys.
		- If a field requires code or structured input (e.g., sample answer in Go, Python), enclose it in a string or nested field.
		`

	config := &genai.GenerateContentConfig{
		Temperature:       &temp,
		TopP:              &topP,
		TopK:              &topK,
		MaxOutputTokens:   maxOutputTokens,
		ResponseMIMEType:  "application/json",
		SystemInstruction: genai.NewContentFromText(systemInstructions, genai.RoleUser),
		SafetySettings:    DefaultSafetySettings(),
	}

	genaiPrompt := genai.Text(prompt)
	result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash", genaiPrompt, config)

	if err != nil {
		return "", err
	}

	return result.Text(), nil
}

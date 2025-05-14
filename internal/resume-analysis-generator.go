package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"

	"google.golang.org/genai"
	"prepai.app/configs"
)

type Metrics struct {
	AtsMatchScore          int64  `json:"ats_match_score"`
	ClarityScore           int64  `json:"clarity_score"`
	GrammarIssues          int64  `json:"grammar_issues"`
	SoftVsHardSkillBalance string `json:"soft_vs_hard_skill_balance"`
	ResumeLengthFeedback   string `json:"resume_length_feedback"`
	FillerWordUsage        string `json:"filler_word_usage"`
}

type ResumeAnalyzerResponse struct {
	Title                  string  `json:"title"`
	OverallScore           int64   `json:"overall_score"`
	AnalysisSummary        string  `json:"analysis_summary"`
	ImprovementSuggestions string  `json:"improvement_suggestions"`
	Metrics                Metrics `json:"metrics"`
}

func ResumeAnalyzer(file multipart.File, header *multipart.FileHeader, jobDescription string) (ResumeAnalyzerResponse, error) {
	bs := make([]byte, header.Size)

	_, err := bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		return ResumeAnalyzerResponse{}, err
	}

	prompt := fmt.Sprintf(`
		You are an expert technical recruiter and resume reviewer. Analyze the following resume in relation to the provided job description.
		Evaluate and return your analysis using the JSON format described below. Be objective, precise, and explain each metric when necessary.
		You need to talk/address as if you were talking to the candidate.

		Job description:
		"%v"

		Return the results using this JSON schema:
		{
			"title": string, // A title for this analysis (no more than one line)
			"overall_score": int, // From 1 to 100: how well the resume fits the job,
			"analysis_summary": string, // A brief 5-8 line summary of the resume quality and fit,
			"metrics": {
				"ats_match_score": int, // (1-100) How well the resume matches keywords/structure from the job,
				"clarity_score": int, // (1-10) Based on grammar, conciseness, and readability,
				"grammar_issues": int, // Total number of grammar or spelling problems,
				"soft_vs_hard_skill_balance": string, // e.g., "Balanced", "Too much soft", "Too technical",
				"resume_length_feedback": string, // e.g., "Appropriate", "Too long for a junior", "Too short",
				"filler_word_usage": string, // e.g., "Minimal", "Moderate", "Heavy use of vague language",
			},
			"improvement_suggestions":
				string // Improvements to enhance the resume
		}

	`, jobDescription)
	parts := []*genai.Part{
		{
			InlineData: &genai.Blob{
				MIMEType: "application/pdf",
				Data:     bs,
			},
		},
		genai.NewPartFromText(prompt),
	}
	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := configs.Gemini(contents)
	if err != nil {
		return ResumeAnalyzerResponse{}, err
	}

	var analysis ResumeAnalyzerResponse

	err = json.Unmarshal([]byte(result), &analysis)
	if err != nil {
		fmt.Print(err)
		return ResumeAnalyzerResponse{}, err
	}

	return analysis, nil
}

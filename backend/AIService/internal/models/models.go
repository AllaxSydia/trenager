package models

type HintRequest struct {
	TaskID    string
	UserCode  string
	Language  string
	HintLevel int
}

type HintResponse struct {
	Hint        string
	HintLevel   int
	ExampleCode string
}

type CodeReviewRequest struct {
	Code            string
	Language        string
	TaskDescription string
}

type CodeReviewResponse struct {
	QualityScore    int
	Issues          []string
	Suggestions     []string
	OverallFeedback string
}

type Recommendation struct {
	TaskID          string
	Title           string
	Reason          string
	DifficultyScore float64
}

type Question struct {
	Text        string
	CodeContext string
	Language    string
}

type Answer struct {
	Text         string
	CodeExamples []string
	HasCode      bool
}

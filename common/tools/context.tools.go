package tools

import "context"

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	// ContextKeySubjectID var
	ContextKeySubjectID = contextKey("subjectID")
	// ContextKeyJobID var
	ContextKeyJobID contextKey
)

// GetSubjectFromContext gets the caller value from the context.
func GetSubjectFromContext(ctx context.Context) (string, bool) {
	caller, ok := ctx.Value(ContextKeySubjectID).(string)
	return caller, ok
}

// GetJobIDFromContext gets the jobID value from the context.
func GetJobIDFromContext(ctx context.Context) (string, bool) {
	jobID, ok := ctx.Value(ContextKeyJobID).(string)
	return jobID, ok
}

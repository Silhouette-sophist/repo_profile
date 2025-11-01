package ast_type

import (
	"context"
	"testing"
)

func TestAstTypeAnalyzer_AnalyzeRepo(t *testing.T) {
	ctx := context.Background()
	r := &AstTypeAnalyzer{
		RepoPath: "/Users/silhouette/work-practice/repo_profile",
		RootPkg:  "github.com/Silhouette-sophist/repo_profile",
	}
	r.AnalyzeRepo(ctx)
}

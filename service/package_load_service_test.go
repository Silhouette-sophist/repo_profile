package service

import (
	"context"
	"testing"
)

func TestLoadPackages(t *testing.T) {
	ctx := context.Background()
	repoPath := "/Users/silhouette/work-practice/repo_profile"
	pkgPath := "github.com/Silhouette-sophist/repo_profile"
	// LoadPackages(ctx, &LoadConfig{
	// 	RepoPath: repoPath,
	// 	PkgPath:  pkgPath,
	// 	LoadEnum: LoadAllPkg,
	// })
	packages, err := LoadPackages(ctx, &LoadConfig{
		RepoPath: repoPath,
		PkgPath:  pkgPath,
		LoadEnum: LoadSpecificPkgWithChild,
	})
	if err != nil {
		t.Fatalf("LoadPackages err: %v", err)
	}
	t.Logf("packages: %v", packages)
}

func TestLoadAllPackages(t *testing.T) {
	ctx := context.Background()
	repoPath := "/Users/silhouette/work-practice/repo_profile"
	pkgPath := "github.com/Silhouette-sophist/repo_profile"
	packages, err := LoadPackages(ctx, &LoadConfig{
		RepoPath: repoPath,
		PkgPath:  pkgPath,
		LoadEnum: LoadCurrentRepo,
	})
	if err != nil {
		t.Fatalf("LoadAllPackages err: %v", err)
	}
	t.Logf("packages: %v", packages)
}

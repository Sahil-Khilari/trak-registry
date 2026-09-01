package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Allowed Categories for Tracks
var validCategories = map[string]bool{
	"lang":  true,
	"os":    true,
	"cloud": true,
	"db":    true,
	"tool":  true,
}

// AST Node
type Node struct {
	Name     string `json:"name"`
	Type     string `json:"type"` // "file" or "directory"
	Content  string `json:"content,omitempty"`
	Children []Node `json:"children,omitempty"`
}

// Blueprint Structure
type Blueprint struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Root        Node   `json:"root"`
}

func main() {
	fmt.Println("=======================================================")
	fmt.Println("  🔍 Trak Registry - Blueprint AST Schema Validator")
	fmt.Println("=======================================================")
	fmt.Println()

	changedFiles := getChangedFiles()
	eventName := strings.TrimSpace(os.Getenv("GITHUB_EVENT_NAME"))
	actor := strings.TrimSpace(os.Getenv("GITHUB_ACTOR"))
	repoOwner := strings.TrimSpace(os.Getenv("REPO_OWNER"))

	if eventName == "pull_request" {
		fmt.Printf("📋 PR Event Detected | Author: @%s | Target Repo Owner: @%s\n", actor, repoOwner)
		if len(changedFiles) > 0 {
			fmt.Printf("   Detected %d changed file(s) in this PR.\n", len(changedFiles))
		}
		fmt.Println()

		// If PR author is not repo owner, prevent tampering with scripts/ or .github/
		if actor != "" && repoOwner != "" && !strings.EqualFold(actor, repoOwner) {
			for file := range changedFiles {
				if strings.HasPrefix(file, "scripts/") || strings.HasPrefix(file, ".github/") {
					fmt.Printf("❌ Security Violation: PR author '@%s' cannot modify CI workflows or validation scripts ('%s')\n", actor, file)
					os.Exit(1)
				}
			}
		}
	}

	errorsCount := 0
	filesCount := 0

	// Walk through templates/ and users/ directories
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Normalize path to forward slashes and trim leading ./
		normalizedPath := filepath.ToSlash(path)
		normalizedPath = strings.TrimPrefix(normalizedPath, "./")

		// Only check .json files inside templates/ or users/
		if !strings.HasSuffix(normalizedPath, ".json") {
			return nil
		}

		// Skip top-level config files
		if normalizedPath == "registry.json" || normalizedPath == "package.json" {
			return nil
		}

		isOfficial := strings.HasPrefix(normalizedPath, "templates/")
		isUser := strings.HasPrefix(normalizedPath, "users/")

		if !isOfficial && !isUser {
			return nil
		}

		filesCount++
		fmt.Printf("🔍 [%s] Checking: %-45s ", formatType(isOfficial), normalizedPath)

		if err := validateFile(normalizedPath, isOfficial, changedFiles); err != nil {
			fmt.Printf("❌ FAILED\n   Error: %v\n", err)
			errorsCount++
		} else {
			fmt.Println("✔ PASSED")
		}

		return nil
	})

	if err != nil {
		fmt.Printf("\n❌ Error walking registry directories: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("=======================================================")
	if errorsCount > 0 {
		fmt.Printf("  ❌ Validation FAILED with %d error(s) across %d files.\n", errorsCount, filesCount)
		fmt.Println("=======================================================")
		os.Exit(1)
	}

	fmt.Printf("  ✨ All %d Blueprints PASSED Validation Successfully! 🎉\n", filesCount)
	fmt.Println("=======================================================")
}

func getChangedFiles() map[string]bool {
	changed := make(map[string]bool)

	// 1. Check if CHANGED_FILES env var is populated by CI
	if envList := os.Getenv("CHANGED_FILES"); envList != "" {
		for _, f := range strings.Split(envList, "\n") {
			f = strings.TrimSpace(filepath.ToSlash(f))
			f = strings.TrimPrefix(f, "./")
			if f != "" {
				changed[f] = true
			}
		}
		return changed
	}

	// 2. If running under pull_request event, attempt to use git diff
	if os.Getenv("GITHUB_EVENT_NAME") == "pull_request" {
		baseRef := os.Getenv("GITHUB_BASE_REF")
		if baseRef == "" {
			baseRef = "main"
		}
		// Try git diff --name-only origin/<baseRef>...HEAD
		out, err := exec.Command("git", "diff", "--name-only", "origin/"+baseRef+"...HEAD").Output()
		if err == nil && len(out) > 0 {
			for _, line := range strings.Split(string(out), "\n") {
				f := strings.TrimSpace(filepath.ToSlash(line))
				f = strings.TrimPrefix(f, "./")
				if f != "" {
					changed[f] = true
				}
			}
			return changed
		}
		// Fallback: git diff --name-only HEAD~1
		out, err = exec.Command("git", "diff", "--name-only", "HEAD~1").Output()
		if err == nil && len(out) > 0 {
			for _, line := range strings.Split(string(out), "\n") {
				f := strings.TrimSpace(filepath.ToSlash(line))
				f = strings.TrimPrefix(f, "./")
				if f != "" {
					changed[f] = true
				}
			}
			return changed
		}
	}

	return changed
}

func formatType(isOfficial bool) string {
	if isOfficial {
		return "OFFICIAL"
	}
	return "COMMUNITY"
}

func validateFile(filePath string, isOfficial bool, changedFiles map[string]bool) error {
	eventName := strings.TrimSpace(os.Getenv("GITHUB_EVENT_NAME"))
	actor := strings.TrimSpace(os.Getenv("GITHUB_ACTOR"))
	repoOwner := strings.TrimSpace(os.Getenv("REPO_OWNER"))

	// In a pull request, only files that were ACTUALLY modified/added in this PR
	// are subject to author ownership / security verification.
	isChangedInPR := true
	if eventName == "pull_request" && len(changedFiles) > 0 {
		isChangedInPR = changedFiles[filePath]
	}

	// 1. Path structure verification
	parts := strings.Split(filePath, "/")
	if isOfficial {
		// Expect: templates/<category>/<slug>.json (3 parts)
		if len(parts) != 3 {
			return fmt.Errorf("official template path must be 'templates/<category>/<slug>.json', got '%s'", filePath)
		}
		category := parts[1]
		if !validCategories[category] {
			return fmt.Errorf("invalid category '%s'. Must be one of: lang, os, cloud, db, tool", category)
		}

		// Security Check for Official templates in PRs (only if PR changed/added this file)
		if eventName == "pull_request" && isChangedInPR && actor != "" && repoOwner != "" {
			if !strings.EqualFold(actor, repoOwner) {
				return fmt.Errorf("security violation: PR author '@%s' cannot modify official templates/. Please submit your curriculum under 'users/%s/'", actor, actor)
			}
		}
	} else {
		// Expect: users/<username>/<category>/<slug>.json (4 parts)
		if len(parts) != 4 {
			return fmt.Errorf("community template path must be 'users/<username>/<category>/<slug>.json', got '%s'", filePath)
		}
		folderUser := parts[1]
		category := parts[2]
		if !validCategories[category] {
			return fmt.Errorf("invalid category '%s'. Must be one of: lang, os, cloud, db, tool", category)
		}

		// Security Check: PR Author authorization (only if PR changed/added this file)
		if eventName == "pull_request" && isChangedInPR && actor != "" {
			if !strings.EqualFold(folderUser, actor) && (repoOwner == "" || !strings.EqualFold(actor, repoOwner)) {
				return fmt.Errorf("security violation: PR author '@%s' cannot modify namespace 'users/%s/'. You can only add/edit templates under 'users/%s/'", actor, folderUser, actor)
			}
		}
	}

	// 2. Read and parse JSON
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	// Check max size (5MB limit per blueprint)
	if len(data) > 5*1024*1024 {
		return fmt.Errorf("file size (%d bytes) exceeds maximum 5MB limit", len(data))
	}

	var bp Blueprint
	if err := json.Unmarshal(data, &bp); err != nil {
		return fmt.Errorf("invalid JSON syntax: %w", err)
	}

	// 3. Metadata Validation
	if strings.TrimSpace(bp.ID) == "" {
		return fmt.Errorf("missing mandatory field 'id'")
	}
	if strings.TrimSpace(bp.Name) == "" {
		return fmt.Errorf("missing mandatory field 'name'")
	}
	if strings.TrimSpace(bp.Version) == "" {
		return fmt.Errorf("missing mandatory field 'version'")
	}

	// 4. Root AST validation
	if bp.Root.Type != "directory" {
		return fmt.Errorf("root node type must be 'directory', got '%s'", bp.Root.Type)
	}
	if strings.TrimSpace(bp.Root.Name) == "" {
		return fmt.Errorf("root node must have a non-empty name")
	}

	// 5. Recursive Node Validation
	return validateNode(&bp.Root, "/")
}

func validateNode(n *Node, parentPath string) error {
	if strings.TrimSpace(n.Name) == "" {
		return fmt.Errorf("node at '%s' has an empty name", parentPath)
	}

	// Prevent path traversal and illegal characters
	if strings.ContainsAny(n.Name, "/\\:\x00") || n.Name == ".." || n.Name == "." {
		return fmt.Errorf("illegal node name '%s' at '%s'", n.Name, parentPath)
	}

	currentPath := parentPath + n.Name

	switch n.Type {
	case "file":
		if len(n.Children) > 0 {
			return fmt.Errorf("file node '%s' cannot have children", currentPath)
		}
		// Disallow dangerous executable binaries inside templates
		lowerName := strings.ToLower(n.Name)
		if strings.HasSuffix(lowerName, ".exe") || strings.HasSuffix(lowerName, ".dll") || strings.HasSuffix(lowerName, ".so") || strings.HasSuffix(lowerName, ".dylib") {
			return fmt.Errorf("forbidden binary file '%s' detected inside template", currentPath)
		}
	case "directory":
		for i := range n.Children {
			if err := validateNode(&n.Children[i], currentPath+"/"); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown node type '%s' at '%s' (must be 'file' or 'directory')", n.Type, currentPath)
	}

	return nil
}

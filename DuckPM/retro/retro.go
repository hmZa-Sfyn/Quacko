package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

const (
	supabaseURL    = "https://ugykvbblboigrjdzlfte.supabase.co"
	supabaseAPIKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InVneWt2YmJsYm9pZ3JqZHpsZnRlIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDA5MDk3MjksImV4cCI6MjA1NjQ4NTcyOX0.0mtbFGKycclZr1TjG6NYlOTuniiTWyHoW0DG2IAIOXw"
)

var (
	globalInstallPath string
)

func init() {
	// Set global install path based on OS
	if runtime.GOOS == "windows" {
		globalInstallPath = `C:\Program Files\Quacko`
	} else {
		globalInstallPath = "/usr/local/quacko"
	}
}

type Library struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	GithubRepoLink string   `json:"github_repo_link"`
	AuthorID       string   `json:"author_id"`
	License        string   `json:"license"`
	Tags           []string `json:"tags"`
	Version        string   `json:"version"`
	NumOfLikes     int      `json:"num_of_likes"`
	NumOfDownloads int      `json:"num_of_downloads"`
}

func main() {
	var rootCmd = &cobra.Command{Use: "retro"}
	var globalFlag bool

	// Color setup
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	var installCmd = &cobra.Command{
		Use:   "install [library_name]",
		Short: "Install a QuackoLang library",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			libName := args[0]
			if !isValidLibraryName(libName) {
				fmt.Printf("%s: Invalid library name '%s'\n", red("Error"), libName)
				return
			}
			installLibrary(libName, globalFlag, green, red)
		},
	}

	var uninstallCmd = &cobra.Command{
		Use:   "uninstall [library_name]",
		Short: "Uninstall a QuackoLang library",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			libName := args[0]
			if !isValidLibraryName(libName) {
				fmt.Printf("%s: Invalid library name '%s'\n", red("Error"), libName)
				return
			}
			uninstallLibrary(libName, globalFlag, green, red)
		},
	}

	var infoCmd = &cobra.Command{
		Use:   "info [library_name]",
		Short: "Display information about a QuackoLang library",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			libName := args[0]
			if !isValidLibraryName(libName) {
				fmt.Printf("%s: Invalid library name '%s'\n", red("Error"), libName)
				return
			}
			getLibraryInfo(libName, blue, red)
		},
	}

	var repoCmd = &cobra.Command{
		Use:   "repo [library_name]",
		Short: "Print the GitHub repository URL for a QuackoLang library",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			libName := args[0]
			if !isValidLibraryName(libName) {
				fmt.Printf("%s: Invalid library name '%s'\n", red("Error"), libName)
				return
			}
			getLibraryRepoURL(libName, blue, red)
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List installed QuackoLang libraries",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			listInstalledLibraries(globalFlag, blue, red, green)
		},
	}

	var pushCmd = &cobra.Command{
		Use:   "push",
		Short: "Push changes to the git repository and optionally update version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			version, _ := cmd.Flags().GetString("version")
			if version != "" && !isValidVersion(version) {
				fmt.Printf("%s: Invalid version format '%s'\n", red("Error"), version)
				return
			}
			pushChanges(version, green, red, yellow)
		},
	}

	installCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, fmt.Sprintf("Install the library globally in %s", globalInstallPath))
	uninstallCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "Uninstall the library from the global directory")
	listCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "List libraries installed globally")
	pushCmd.Flags().StringP("version", "v", "", "Update the library version in Supabase")

	rootCmd.AddCommand(installCmd, uninstallCmd, infoCmd, repoCmd, listCmd, pushCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%s: %v\n", red("Error"), err)
		os.Exit(1)
	}
}

func isValidLibraryName(name string) bool {
	// Basic validation: non-empty, no slashes, no special characters
	return len(name) > 0 && !strings.ContainsAny(name, "/\\:") && !strings.HasPrefix(name, ".")
}

func isValidVersion(version string) bool {
	// Basic semantic version check (e.g., 1.0.0)
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return false
	}
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return false
		}
	}
	return true
}

func fetchLibrary(libName string) (*Library, error) {
	url := fmt.Sprintf("%s/rest/v1/retro_libs?select=*&name=ilike.%s", supabaseURL, strings.ReplaceAll(libName, "%", ""))
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+supabaseAPIKey)
	req.Header.Set("apikey", supabaseAPIKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch library: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var libraries []Library
	if err := json.Unmarshal(body, &libraries); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(libraries) == 0 {
		return nil, fmt.Errorf("library '%s' not found", libName)
	}

	return &libraries[0], nil
}

func listInstalledLibraries(global bool, blue, red, green func(a ...interface{}) string) {
	installPath := "."
	if global {
		installPath = globalInstallPath
	}

	dir, err := os.Open(installPath)
	if err != nil {
		fmt.Printf("%s: Failed to open directory %s: %v\n", red("Error"), installPath, err)
		return
	}
	defer dir.Close()

	libs, err := dir.Readdir(-1)
	if err != nil {
		fmt.Printf("%s: Failed to read directory %s: %v\n", red("Error"), installPath, err)
		return
	}

	if len(libs) == 0 {
		fmt.Printf("%s: No libraries installed in %s\n", red("Warning"), installPath)
		return
	}

	fmt.Printf("%s in %s:\n", blue("Installed libraries"), installPath)
	for _, lib := range libs {
		if lib.IsDir() {
			// Try to fetch version from Supabase
			libName := lib.Name()
			library, err := fetchLibrary(libName)
			version := "unknown"
			if err == nil {
				version = library.Version
			}
			fmt.Printf(" - %s %s\n", libName, green("(v"+version+")"))
		}
	}
}

func installLibrary(libName string, global bool, green, red func(a ...interface{}) string) {
	library, err := fetchLibrary(libName)
	if err != nil {
		fmt.Printf("%s: %v\n", red("Error"), err)
		return
	}

	installPath := "."
	if global {
		installPath = globalInstallPath
	}

	// Check if library already exists
	repoPath := filepath.Join(installPath, library.Name)
	if _, err := os.Stat(repoPath); !os.IsNotExist(err) {
		fmt.Printf("%s: Library %s already installed in %s\n", red("Error"), library.Name, installPath)
		return
	}

	// Ensure the install directory exists
	if err := os.MkdirAll(installPath, 0755); err != nil {
		fmt.Printf("%s creating directory %s: %v\n", red("Error"), installPath, err)
		return
	}

	// Validate GitHub URL
	if !strings.HasPrefix(library.GithubRepoLink, "https://github.com/") {
		fmt.Printf("%s: Invalid GitHub repository URL for %s\n", red("Error"), library.Name)
		return
	}

	// Clone the repository
	_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:      library.GithubRepoLink,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Printf("%s cloning repository: %v\n", red("Error"), err)
		return
	}

	fmt.Printf("%s %s (v%s) to %s\n", green("Successfully installed"), library.Name, library.Version, repoPath)

	// Update download count
	updateDownloadCount(library.ID, red)
}

func uninstallLibrary(libName string, global bool, green, red func(a ...interface{}) string) {
	installPath := "."
	if global {
		installPath = globalInstallPath
	}

	repoPath := filepath.Join(installPath, libName)
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fmt.Printf("%s: Library %s not found in %s\n", red("Error"), libName, installPath)
		return
	}

	if err := os.RemoveAll(repoPath); err != nil {
		fmt.Printf("%s uninstalling library %s: %v\n", red("Error"), libName, err)
		return
	}

	fmt.Printf("%s %s from %s\n", green("Successfully uninstalled"), libName, installPath)
}

func getLibraryInfo(libName string, blue, red func(a ...interface{}) string) {
	library, err := fetchLibrary(libName)
	if err != nil {
		fmt.Printf("%s: %v\n", red("Error"), err)
		return
	}

	fmt.Printf("%s: %s\n", blue("Library"), library.Name)
	fmt.Printf("%s: %s\n", blue("Description"), library.Description)
	fmt.Printf("%s: %s\n", blue("Version"), library.Version)
	fmt.Printf("%s: %s\n", blue("License"), library.License)
	fmt.Printf("%s: %s\n", blue("Tags"), strings.Join(library.Tags, ", "))
	fmt.Printf("%s: %d\n", blue("Downloads"), library.NumOfDownloads)
	fmt.Printf("%s: %d\n", blue("Likes"), library.NumOfLikes)
	fmt.Printf("%s: %s\n", blue("GitHub Repository"), library.GithubRepoLink)
}

func getLibraryRepoURL(libName string, blue, red func(a ...interface{}) string) {
	library, err := fetchLibrary(libName)
	if err != nil {
		fmt.Printf("%s: %v\n", red("Error"), err)
		return
	}
	fmt.Println(blue(library.GithubRepoLink))
}

func pushChanges(version string, green, red, yellow func(a ...interface{}) string) {
	// Open the current directory as a git repository
	repo, err := git.PlainOpen(".")
	if err != nil {
		fmt.Printf("%s: Failed to open git repository: %v\n", red("Error"), err)
		return
	}

	// Check repository status
	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Printf("%s: Failed to get worktree: %v\n", red("Error"), err)
		return
	}

	status, err := worktree.Status()
	if err != nil {
		fmt.Printf("%s: Failed to check repository status: %v\n", red("Error"), err)
		return
	}
	if status.IsClean() {
		fmt.Printf("%s: No changes to commit\n", yellow("Warning"))
		return
	}

	// Add all changes
	if _, err := worktree.Add("."); err != nil {
		fmt.Printf("%s: Failed to add files: %v\n", red("Error"), err)
		return
	}

	// Commit changes
	commitMsg := "retro push!!!"
	commit, err := worktree.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Retro User",
			Email: "retro@quacko.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		fmt.Printf("%s: Failed to commit changes: %v\n", red("Error"), err)
		return
	}

	// Verify main branch exists
	_, err = repo.Branch("main")
	if err != nil {
		fmt.Printf("%s: Main branch does not exist: %v\n", red("Error"), err)
		return
	}

	// Push to origin main
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/heads/main:refs/heads/main")},
	})
	if err != nil {
		fmt.Printf("%s: Failed to push to origin main: %v\n", red("Error"), err)
		return
	}

	fmt.Printf("%s: Changes pushed to origin main (commit %s)\n", green("Success"), commit.String()[:7])

	if version != "" {
		// Update version in Supabase
		dir, err := os.Getwd()
		if err != nil {
			fmt.Printf("%s: Failed to get current directory: %v\n", red("Error"), err)
			return
		}
		libName := filepath.Base(dir)

		library, err := fetchLibrary(libName)
		if err != nil {
			fmt.Printf("%s: Failed to fetch library %s: %v\n", red("Error"), libName, err)
			return
		}

		oldVersion := library.Version
		updateLibraryVersion(library.ID, version, red)
		fmt.Printf("%s: %s ==> %s\n", yellow("Version updated"), red(oldVersion), green(version))
	}
}

func updateDownloadCount(libraryID string, red func(a ...interface{}) string) {
	library, err := fetchLibraryByID(libraryID)
	if err != nil {
		fmt.Printf("%s: Failed to fetch library for download count update: %v\n", red("Error"), err)
		return
	}

	url := fmt.Sprintf("%s/rest/v1/retro_libs?id=eq.%s", supabaseURL, libraryID)
	newDownloads := library.NumOfDownloads + 1
	data := fmt.Sprintf(`{"num_of_downloads": %d}`, newDownloads)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("PATCH", url, strings.NewReader(data))
	if err != nil {
		fmt.Printf("%s: Error creating update request: %v\n", red("Error"), err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+supabaseAPIKey)
	req.Header.Set("apikey", supabaseAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s: Error updating download count: %v\n", red("Error"), err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("%s: Failed to update download count: status code %d, response: %s\n", red("Error"), resp.StatusCode, string(body))
	}
}

func updateLibraryVersion(libraryID, version string, red func(a ...interface{}) string) {
	url := fmt.Sprintf("%s/rest/v1/retro_libs?id=eq.%s", supabaseURL, libraryID)
	data := fmt.Sprintf(`{"version": "%s"}`, version)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("PATCH", url, strings.NewReader(data))
	if err != nil {
		fmt.Printf("%s: Error creating version update request: %v\n", red("Error"), err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+supabaseAPIKey)
	req.Header.Set("apikey", supabaseAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s: Error updating library version: %v\n", red("Error"), err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("%s: Failed to update library version: status code %d, response: %s\n", red("Error"), resp.StatusCode, string(body))
	}
}

func fetchLibraryByID(libraryID string) (*Library, error) {
	url := fmt.Sprintf("%s/rest/v1/retro_libs?select=*&id=eq.%s", supabaseURL, libraryID)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+supabaseAPIKey)
	req.Header.Set("apikey", supabaseAPIKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch library: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var libraries []Library
	if err := json.Unmarshal(body, &libraries); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(libraries) == 0 {
		return nil, fmt.Errorf("library with ID '%s' not found", libraryID)
	}

	return &libraries[0], nil
}

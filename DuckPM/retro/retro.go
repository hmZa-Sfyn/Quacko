package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

const (
	supabaseURL       = "https://ugykvbblboigrjdzlfte.supabase.co"
	supabaseAPIKey    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InVneWt2YmJsYm9pZ3JqZHpsZnRlIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDA5MDk3MjksImV4cCI6MjA1NjQ4NTcyOX0.0mtbFGKycclZr1TjG6NYlOTuniiTWyHoW0DG2IAIOXw"
	globalInstallPath = `C:\Program Files\Quacko`
)

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

type APIResponse struct {
	Data  []Library `json:"data"`
	Error string    `json:"message"`
}

func main() {
	var rootCmd = &cobra.Command{Use: "retro"}
	var globalFlag bool

	var installCmd = &cobra.Command{
		Use:   "install [library_name]",
		Short: "Install a QuackoLang library",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			libName := args[0]
			installLibrary(libName, globalFlag)
		},
	}

	var uninstallCmd = &cobra.Command{
		Use:   "uninstall [library_name]",
		Short: "Uninstall a QuackoLang library",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			libName := args[0]
			uninstallLibrary(libName, globalFlag)
		},
	}

	var infoCmd = &cobra.Command{
		Use:   "info [library_name]",
		Short: "Display information about a QuackoLang library",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			libName := args[0]
			getLibraryInfo(libName)
		},
	}

	var repoCmd = &cobra.Command{
		Use:   "repo [library_name]",
		Short: "Print the GitHub repository URL for a QuackoLang library",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			libName := args[0]
			getLibraryRepoURL(libName)
		},
	}

	installCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "Install the library globally in C:\\Program Files\\Quacko")
	uninstallCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "Uninstall the library from the global directory")

	rootCmd.AddCommand(installCmd, uninstallCmd, infoCmd, repoCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func fetchLibrary(libName string) (*Library, error) {
	url := fmt.Sprintf("%s/rest/v1/retro_libs?select=*&name=ilike.%s", supabaseURL, libName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+supabaseAPIKey)
	req.Header.Set("apikey", supabaseAPIKey)

	client := &http.Client{}
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

func installLibrary(libName string, global bool) {
	library, err := fetchLibrary(libName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	installPath := "."
	if global {
		installPath = globalInstallPath
	}

	// Ensure the install directory exists
	if err := os.MkdirAll(installPath, 0755); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", installPath, err)
		return
	}

	// Clone the repository
	repoPath := filepath.Join(installPath, library.Name)
	_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:      library.GithubRepoLink,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Printf("Error cloning repository: %v\n", err)
		return
	}

	fmt.Printf("Successfully installed %s to %s\n", library.Name, repoPath)

	// Update download count
	updateDownloadCount(library.ID)
}

func uninstallLibrary(libName string, global bool) {
	installPath := "."
	if global {
		installPath = globalInstallPath
	}

	repoPath := filepath.Join(installPath, libName)
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fmt.Printf("Library %s not found in %s\n", libName, installPath)
		return
	}

	if err := os.RemoveAll(repoPath); err != nil {
		fmt.Printf("Error uninstalling library %s: %v\n", libName, err)
		return
	}

	fmt.Printf("Successfully uninstalled %s from %s\n", libName, installPath)
}

func getLibraryInfo(libName string) {
	library, err := fetchLibrary(libName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Library: %s\n", library.Name)
	fmt.Printf("Description: %s\n", library.Description)
	fmt.Printf("Version: %s\n", library.Version)
	fmt.Printf("License: %s\n", library.License)
	fmt.Printf("Tags: %s\n", strings.Join(library.Tags, ", "))
	fmt.Printf("Downloads: %d\n", library.NumOfDownloads)
	fmt.Printf("Likes: %d\n", library.NumOfLikes)
	fmt.Printf("GitHub Repository: %s\n", library.GithubRepoLink)
}

func getLibraryRepoURL(libName string) {
	library, err := fetchLibrary(libName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(library.GithubRepoLink)
}

func updateDownloadCount(libraryID string) {
	url := fmt.Sprintf("%s/rest/v1/retro_libs?id=eq.%s", supabaseURL, libraryID)
	data := []byte(`{"num_of_downloads": "increment"}`)

	req, err := http.NewRequest("PATCH", url, strings.NewReader(string(data)))
	if err != nil {
		fmt.Printf("Error creating update request: %v\n", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+supabaseAPIKey)
	req.Header.Set("apikey", supabaseAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error updating download count: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		fmt.Printf("Failed to update download count: status code %d\n", resp.StatusCode)
	}
}

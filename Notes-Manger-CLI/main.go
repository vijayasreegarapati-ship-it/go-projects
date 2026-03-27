package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"
)

const timeFormat = "2006-01-02-15-04-05"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	initStorage()
	command := strings.ToLower(os.Args[1])

	switch command {
	case "add":
		handleAdd()
	case "new":
		handleNew()
	case "list":
		handleList()
	case "view":
		handleView()
	case "search":
		handleSearch()
	case "edit":
		handleEdit()
	case "delete":
		handleDelete()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

// initStorage ensures the local ~/.gonotes directory exists.
func initStorage() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
		os.Exit(1)
	}

	notesDir := filepath.Join(homeDir, ".gonotes")
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating notes directory: %v\n", err)
		os.Exit(1)
	}

	return notesDir
}

// handleAdd captures a quick one-line note from the CLI args.
func handleAdd() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: Please provide some text for your note.\nExample: gonotes add \"This is my note\"")
		os.Exit(1)
	}

	noteText := strings.Join(os.Args[2:], " ")
	filename := fmt.Sprintf("%s.md", time.Now().Format(timeFormat))
	fullPath := filepath.Join(initStorage(), filename)

	if err := os.WriteFile(fullPath, []byte(noteText), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving note: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Note saved: %s\n", filename)
}

// handleNew drops the user into their default $EDITOR to write a note.
func handleNew() {
	filename := fmt.Sprintf("%s.md", time.Now().Format(timeFormat))
	fullPath := filepath.Join(initStorage(), filename)

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	cmd := exec.Command(editor, fullPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error opening editor: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf(" Note saved: %s\n", filename)
}

// handleList outputs all notes in a formatted terminal table.
func handleList() {
	notesDir := initStorage()
	files, err := os.ReadDir(notesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading notes directory: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No notes found. Create one with 'gonotes add'!")
		return
	}

	fmt.Println("\nYour Notes:")
	fmt.Println(strings.Repeat("-", 40))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tFilename")

	for i, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
			fmt.Fprintf(w, "%d\t%s\n", i+1, f.Name())
		}
	}
	w.Flush()

	fmt.Println(strings.Repeat("-", 40))
}

// handleView prints the contents of a specific file to stdout.
func handleView() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: Please provide a filename.\nExample: gonotes view 2026-03-27-15-04-05.md")
		os.Exit(1)
	}

	filename := os.Args[2]
	fullPath := filepath.Join(initStorage(), filename)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading note '%s': %v\n", filename, err)
		os.Exit(1)
	}

	fmt.Printf("\n--- %s ---\n\n", filename)
	fmt.Println(string(content))
}

// handleSearch scans all markdown notes for a case-insensitive keyword match.
func handleSearch() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: Please provide a keyword to search for.")
		os.Exit(1)
	}

	query := strings.Join(os.Args[2:], " ")
	queryLower := strings.ToLower(query)
	notesDir := initStorage()

	files, err := os.ReadDir(notesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nSearching for '%s'...\n\n", query)
	foundMatch := false

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}

		fullPath := filepath.Join(notesDir, f.Name())
		file, err := os.Open(fullPath)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(file)
		lineNumber := 1
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(strings.ToLower(line), queryLower) {
				fmt.Printf("[%s: Line %d] %s\n", f.Name(), lineNumber, strings.TrimSpace(line))
				foundMatch = true
			}
			lineNumber++
		}
		file.Close()
	}

	if !foundMatch {
		fmt.Println("No matches found.")
	}
	fmt.Println()
}

// handleDelete removes a note file from the local filesystem.
func handleDelete() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: Please provide a filename to delete.\nExample: gonotes delete 2026-03-27-15-04-05.md")
		os.Exit(1)
	}

	filename := os.Args[2]
	fullPath := filepath.Join(initStorage(), filename)

	if err := os.Remove(fullPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting note '%s': %v\n", filename, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully deleted note: %s\n", filename)
}

// handleEdit opens an existing note in the default text editor.
func handleEdit() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: Please provide a filename to edit.\nExample: gonotes edit 2026-03-27-15-04-05.md")
		os.Exit(1)
	}

	filename := os.Args[2]
	fullPath := filepath.Join(initStorage(), filename)

	// Safety check: Ensure the file actually exists before we try to edit it
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Note '%s' does not exist. Run 'list' to see your notes.\n", filename)
		os.Exit(1)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	cmd := exec.Command(editor, fullPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error opening editor: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Note updated successfully: %s\n", filename)
}

func printUsage() {
	usage := `
gonotes - A simple CLI notes manager

Usage:
  gonotes <command> [arguments]

Commands:
  add <text>      Add a quick one-line note
  new             Open default editor to write a note
  edit <File>     Edit an existing note
  list            List all saved notes
  view <File>     View the full text of a note
  search <query>  Search all notes for a keyword
  delete <File>   Permanently delete a note
`
	fmt.Println(usage)
}

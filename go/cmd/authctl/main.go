package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/tehnerd/vatran/go/server/auth"
	"golang.org/x/term"
)

const (
	defaultDBPath    = "/var/lib/katran/auth.db"
	defaultBcryptCost = 12
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "add":
		cmdAdd(args)
	case "list":
		cmdList(args)
	case "delete":
		cmdDelete(args)
	case "passwd":
		cmdPasswd(args)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`authctl - Katran authentication user management tool

Usage:
  authctl <command> [options]

Commands:
  add      Add a new user
  list     List all users
  delete   Delete a user
  passwd   Change user password
  help     Show this help message

Examples:
  authctl add -username admin -db /var/lib/katran/auth.db
  authctl list -db /var/lib/katran/auth.db
  authctl delete -username olduser -db /var/lib/katran/auth.db
  authctl passwd -username admin -db /var/lib/katran/auth.db`)
}

func cmdAdd(args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	username := fs.String("username", "", "Username to create (required)")
	password := fs.String("password", "", "Password (will prompt if not provided)")
	dbPath := fs.String("db", defaultDBPath, "Path to SQLite database")
	bcryptCost := fs.Int("bcrypt-cost", defaultBcryptCost, "bcrypt cost factor")
	fs.Parse(args)

	if *username == "" {
		fmt.Fprintln(os.Stderr, "Error: -username is required")
		fs.Usage()
		os.Exit(1)
	}

	// Prompt for password if not provided
	pwd := *password
	if pwd == "" {
		var err error
		pwd, err = promptPassword("Enter password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			os.Exit(1)
		}
		confirm, err := promptPassword("Confirm password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			os.Exit(1)
		}
		if pwd != confirm {
			fmt.Fprintln(os.Stderr, "Error: passwords do not match")
			os.Exit(1)
		}
	}

	if len(pwd) < 8 {
		fmt.Fprintln(os.Stderr, "Error: password must be at least 8 characters")
		os.Exit(1)
	}

	store, err := auth.NewStore(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	userRepo := auth.NewUserRepository(store, *bcryptCost)
	user, err := userRepo.Create(*username, pwd)
	if err != nil {
		if err == auth.ErrUserExists {
			fmt.Fprintf(os.Stderr, "Error: user '%s' already exists\n", *username)
		} else {
			fmt.Fprintf(os.Stderr, "Error creating user: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("Created user '%s' (id: %d)\n", user.Username, user.ID)
}

func cmdList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	dbPath := fs.String("db", defaultDBPath, "Path to SQLite database")
	fs.Parse(args)

	store, err := auth.NewStore(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	rows, err := store.DB().Query("SELECT id, username, created_at, updated_at FROM users ORDER BY id")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error querying users: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tUSERNAME\tCREATED\tUPDATED")
	fmt.Fprintln(w, "--\t--------\t-------\t-------")

	count := 0
	for rows.Next() {
		var id int64
		var username, createdAt, updatedAt string
		if err := rows.Scan(&id, &username, &createdAt, &updatedAt); err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning row: %v\n", err)
			continue
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", id, username, createdAt, updatedAt)
		count++
	}
	w.Flush()

	if count == 0 {
		fmt.Println("No users found.")
	} else {
		fmt.Printf("\nTotal: %d user(s)\n", count)
	}
}

func cmdDelete(args []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	username := fs.String("username", "", "Username to delete (required)")
	dbPath := fs.String("db", defaultDBPath, "Path to SQLite database")
	force := fs.Bool("force", false, "Skip confirmation prompt")
	fs.Parse(args)

	if *username == "" {
		fmt.Fprintln(os.Stderr, "Error: -username is required")
		fs.Usage()
		os.Exit(1)
	}

	if !*force {
		fmt.Printf("Are you sure you want to delete user '%s'? [y/N]: ", *username)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Aborted.")
			return
		}
	}

	store, err := auth.NewStore(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	result, err := store.DB().Exec("DELETE FROM users WHERE username = ?", *username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting user: %v\n", err)
		os.Exit(1)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		fmt.Fprintf(os.Stderr, "Error: user '%s' not found\n", *username)
		os.Exit(1)
	}

	fmt.Printf("Deleted user '%s'\n", *username)
}

func cmdPasswd(args []string) {
	fs := flag.NewFlagSet("passwd", flag.ExitOnError)
	username := fs.String("username", "", "Username to update (required)")
	password := fs.String("password", "", "New password (will prompt if not provided)")
	dbPath := fs.String("db", defaultDBPath, "Path to SQLite database")
	bcryptCost := fs.Int("bcrypt-cost", defaultBcryptCost, "bcrypt cost factor")
	fs.Parse(args)

	if *username == "" {
		fmt.Fprintln(os.Stderr, "Error: -username is required")
		fs.Usage()
		os.Exit(1)
	}

	// Prompt for password if not provided
	pwd := *password
	if pwd == "" {
		var err error
		pwd, err = promptPassword("Enter new password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			os.Exit(1)
		}
		confirm, err := promptPassword("Confirm new password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			os.Exit(1)
		}
		if pwd != confirm {
			fmt.Fprintln(os.Stderr, "Error: passwords do not match")
			os.Exit(1)
		}
	}

	if len(pwd) < 8 {
		fmt.Fprintln(os.Stderr, "Error: password must be at least 8 characters")
		os.Exit(1)
	}

	store, err := auth.NewStore(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	userRepo := auth.NewUserRepository(store, *bcryptCost)

	// Find user first
	user, err := userRepo.FindByUsername(*username)
	if err != nil {
		if err == auth.ErrUserNotFound {
			fmt.Fprintf(os.Stderr, "Error: user '%s' not found\n", *username)
		} else {
			fmt.Fprintf(os.Stderr, "Error finding user: %v\n", err)
		}
		os.Exit(1)
	}

	if err := userRepo.UpdatePassword(user.ID, pwd); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating password: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Password updated for user '%s'\n", *username)
}

// promptPassword prompts for a password without echoing input.
func promptPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // newline after password input
	if err != nil {
		return "", err
	}
	return string(password), nil
}

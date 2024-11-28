package utils

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// StringPrompt prompts the user for a string input.
func StringPrompt(label string, allowEmpty bool, hasOptions bool, options []string) (string, error) {
	var input string
	for {
		fmt.Fprint(os.Stderr, label)
		if _, err := fmt.Fscanln(os.Stdin, &input); err != nil {
			return "", err
		}

		input = strings.TrimSpace(input)

		if !allowEmpty && len(input) == 0 {
			fmt.Fprintln(os.Stderr, "Input cannot be empty")
			continue
		}

		if hasOptions {
			validOption := false
			for _, opt := range options {
				if strings.EqualFold(input, opt) {
					validOption = true
					break
				}
			}
			if !validOption {
				fmt.Fprintf(os.Stderr, "Invalid input. Must be one of: %s\n", strings.Join(options, ", "))
				continue
			}
		}

		break
	}
	return input, nil
}

// IntegerPrompt prompts the user for an integer input.
func IntegerPrompt(label string) (int, error) {
	var input int
	for {
		fmt.Fprint(os.Stderr, label)
		if _, err := fmt.Fscanln(os.Stdin, &input); err != nil {
			return 0, err
		}
		if input > 0 {
			break
		}
		fmt.Fprintln(os.Stderr, "Input must be a positive integer")
	}
	return input, nil
}

func PasswordPrompt(label string, toConfirm bool) (string, error) {
	var password string
	for {
		fmt.Fprint(os.Stderr, label)
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}
		password = string(bytePassword)
		fmt.Fprintln(os.Stderr)

		if !isValidPassword(password) {
			fmt.Fprintln(os.Stderr, "Password does not meet the requirements. Please try again.")
			continue
		}

		if toConfirm {
			fmt.Fprint(os.Stderr, "Confirm password: ")
			byteConfirmPassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return "", err
			}
			confirmPassword := string(byteConfirmPassword)
			fmt.Fprintln(os.Stderr)

			if password != confirmPassword {
				fmt.Fprintln(os.Stderr, "Passwords do not match. Please try again.")
				continue
			}
		}

		break
	}
	return password, nil
}

// TODO: Add password validation logic
func isValidPassword(password string) bool {
	// if len(password) < 8 {
	// 	return false
	// }

	// var (
	// 	hasUpper   bool
	// 	hasLower   bool
	// 	hasNumber  bool
	// 	hasSpecial bool
	// )

	// for _, char := range password {
	// 	switch {
	// 	case 'A' <= char && char <= 'Z':
	// 		hasUpper = true
	// 	case 'a' <= char && char <= 'z':
	// 		hasLower = true
	// 	case '0' <= char && char <= '9':
	// 		hasNumber = true
	// 	case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
	// 		hasSpecial = true
	// 	}
	// }

	// return hasUpper && hasLower && hasNumber && hasSpecial
	return len(password) > 0
}

// ConfirmPrompt prompts the user for a confirmation input.
func ConfirmPrompt(label string) (bool, error) {
	var input string
	for {
		fmt.Fprint(os.Stderr, label)
		if _, err := fmt.Fscanln(os.Stdin, &input); err != nil {
			return false, err
		}
		switch strings.ToLower(strings.Trim(input, " ")) {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			fmt.Fprintln(os.Stderr, "Invalid input. Please enter 'y' or 'n'")
		}
	}
}

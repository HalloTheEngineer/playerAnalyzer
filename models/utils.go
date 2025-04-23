package models

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var numbers = regexp.MustCompile(`^[0-9]+$`)

func RequireNumbers(str string) bool {
	return numbers.MatchString(str)
}

func GetInput(q string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(q)

	key, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(key), nil
}

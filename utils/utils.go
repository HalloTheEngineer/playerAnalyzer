package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
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

func RandomStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func OpenFile(dir string) {
	if runtime.GOOS == "darwin" {
		cmd := `open "` + dir + `"`
		_ = exec.Command("/bin/bash", "-c", cmd).Start()
	} else {
		dir = strings.ReplaceAll(dir, "/", "\\")
		_ = exec.Command("explorer", dir).Start()
	}
}

func FetchToStruct[T any](url string) (*T, error) {
	replayResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(replayResp.Body)
	if err != nil {
		return nil, err
	}
	if replayResp.StatusCode != 200 {
		return nil, errors.New(strconv.Itoa(replayResp.StatusCode))
	}

	var str T

	err = json.Unmarshal(bytes, &str)
	if err != nil {
		return nil, err
	}

	return &str, nil
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

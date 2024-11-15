package logger

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/MagicalCrawler/RealEstateApp/utils"
)

func TestLog6(t *testing.T) {
	logsDir := t.TempDir()
	os.Setenv("LOG_PATH", logsDir)

	l1 := utils.NewLogger("c1")
	l1.Info("info1 in c1")
	l1.Info("info2 in c1")

	l2 := utils.NewLogger("c2")
	l2.Info("info1 in c2")
	l2.Info("info2 in c2")

	file1Path := filepath.Join(logsDir, "c1.log")
	output1Content, err := os.ReadFile(file1Path)
	if err != nil {
		t.Errorf("Could not open %v", file1Path)
	}
	file1Lines := strings.Split(strings.TrimSpace(string(output1Content)), "\n")

	file2Path := filepath.Join(logsDir, "c2.log")
	output2Content, err := os.ReadFile(file2Path)
	if err != nil {
		t.Errorf("Could not open %v", file2Path)
	}
	file2Lines := strings.Split(strings.TrimSpace(string(output2Content)), "\n")

	yearRegex := `[0-9][0-9][0-9][0-9]`
	monthRegex := `([0-9]|0[0-9]|1[0-2])`
	dayRegex := `([0-9]|0[0-9]|1[0-9]|2[0-9]|3[0-1])`
	dateRegex := yearRegex + `-` + monthRegex + `-` + dayRegex

	hourRegex := `([0-9]|0[0-9]|1[0-9]|2[0-3])`
	minuteRegex := `[0-5][0-9]`
	secondRegex := `[0-5][0-9]`
	millisecondRegex := `([0-9]|[0-9][0-9]|[0-9][0-9][0-9])`
	timezoneRegex := hourRegex + `:` + minuteRegex
	timeRegex := hourRegex + `:` + minuteRegex + `:` + secondRegex + `\.` + millisecondRegex + `\+` + timezoneRegex

	expectedFile1Lines := [2]*regexp.Regexp{}
	for index := range 2 {
		expectedFile1Lines[index] = regexp.MustCompile(`time=` + dateRegex + `T` + timeRegex + ` level=INFO msg="info` + strconv.Itoa(index+1) + ` in c1"`)
	}
	expectedFile2Lines := [2]*regexp.Regexp{}
	for index := range 2 {
		expectedFile2Lines[index] = regexp.MustCompile(`time=` + dateRegex + `T` + timeRegex + ` level=INFO msg="info` + strconv.Itoa(index+1) + ` in c2"`)
	}

	for index, expected := range expectedFile1Lines {
		if !expected.MatchString(file1Lines[index]) {
			t.Errorf("Missing line %v", expected)
		}
	}
	for index, expected := range expectedFile2Lines {
		if !expected.MatchString(file2Lines[index]) {
			t.Errorf("Could not open %v", file2Path)
		}
	}
}

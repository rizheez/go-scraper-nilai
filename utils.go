package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func log(level, message string) {
	fmt.Printf("%s %s\n", level, message)
}

func logf(level, format string, args ...interface{}) {
	fmt.Printf("%s %s\n", level, fmt.Sprintf(format, args...))
}

func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func printHeader(title string, lines []string) {
	width := 60
	if len(lines) > 0 {
		for _, l := range lines {
			if n := len(l); n > width {
				width = n
			}
		}
		width += 4
	}
	if width < len(title)+4 {
		width = len(title) + 4
	}
	pad := (width - len(title) - 2) / 2
	left := strings.Repeat("=", pad)
	right := strings.Repeat("=", width-len(left)-len(title)-2)
	fmt.Printf("%s %s %s\n", left, title, right)
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "-", "\\", "-", ":", " ", "*", "",
		"?", "", "\"", "", "<", "", ">", "", "|", "-",
	)
	return strings.TrimSpace(replacer.Replace(name))
}

func updateProgress(label string, done, total int, lastPercent *int) {
	if total == 0 {
		return
	}
	percent := done * 100 / total
	if percent-*lastPercent >= 5 || percent == 100 {
		barLen := ProgressBarLength
		pos := percent * barLen / 100

		// pakai blok penuh █ dan filler spasi
		bar := strings.Repeat("█", pos) + strings.Repeat(" ", barLen-pos)

		fmt.Printf("\r[PROGRESS] %s: [%s] %d%% (%d/%d)", label, bar, percent, done, total)
		*lastPercent = percent
	}
}

func formatDuration(d time.Duration) string {
	seconds := int(d.Seconds())
	if seconds < 60 {
		return fmt.Sprintf("%d detik", seconds)
	}
	minutes := seconds / 60
	sec := seconds % 60
	if minutes < 60 {
		return fmt.Sprintf("%d menit %d detik", minutes, sec)
	}
	hours := minutes / 60
	min := minutes % 60
	return fmt.Sprintf("%d jam %d menit %d detik", hours, min, sec)
}

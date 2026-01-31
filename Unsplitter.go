package main

import (
	"bufio"
	"io"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// --- helpers ---

var ignorePrefixes []string

func splitLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.Split(s, "\n")
}

func joinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

// get substring using delimiter (Pascal-style Explode)
func explode(str, delim string, n int) string {
	if !strings.Contains(str, delim) && (n == 0 || n == -1) {
		return str
	}

	str = str + delim
	for n > 0 {
		i := strings.Index(str, delim)
		if i == -1 {
			return ""
		}
		str = str[i+len(delim):]
		n--
	}

	i := strings.Index(str, delim)
	if i == -1 {
		return str
	}
	return str[:i]
}

// clean a line like Pascal CleanLine
func cleanLine(s string) string {
	s = explode(s, ";", 0)              // strip comments
	s = strings.ReplaceAll(s, ":", " ") // colons -> spaces
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "'", "\"")
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return s
}

func loadIgnoreList(rootFile string) {
    data, err := os.ReadFile("Unsplitter_ignore.txt")
    if err != nil {
        return
    }

    lines := splitLines(string(data))
    for _, l := range lines {
        l = strings.TrimSpace(l)
        if l != "" {
            ignorePrefixes = append(ignorePrefixes, l)
        }
    }
}

// --- core processing ---
func processFile(path string, log *bufio.Writer) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(log, path+" failed to load")
		return nil, err
	}

	fmt.Fprint(log, path+" opened")

	origLines := splitLines(string(data))
	fmt.Fprintln(log, " -->", len(origLines), "lines in included")

	out := make([]string, 0, len(origLines))

	for _, line := range origLines {
		if !strings.Contains(line, "include") {
			out = append(out, line)
			continue
		}

		cur := cleanLine(line)
		pos := strings.Index(cur, " include \"")
		if pos == -1 {
			out = append(out, line)
			continue
		}

		// filename between quotes
		incName := explode(cur, "\"", 1)
		if incName == "" {
			out = append(out, line)
			continue
		}

		// only .asm
		if !strings.HasSuffix(strings.ToLower(incName), ".asm") {
			out = append(out, line)
			continue
		}

		// ignore based on prefix match
		ignore := false
		for _, p := range ignorePrefixes {
		    if strings.HasPrefix(incName, p) {
			fmt.Fprintln(log, incName+" ignored")
			out = append(out, line) // keep original include line
			ignore = true
			break
		    }
		}
		if ignore {
		    continue
		}


		// keep the original include line, commented out
		copyLabel := strings.TrimSpace(cur[:pos])
		out = append(out, "; "+line)

		if copyLabel != "" {
		    out = append(out, copyLabel+":")
		}

		incLines, err := processFile(incName, log)
		if err != nil {
		    out = append(out, line) // fallback
		    continue
		}

		out = append(out, incLines...)

	}

	return out, nil
}


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: Drop root .asm file into this tool (example: sonic.asm)")
		fmt.Println("Ignore: You can create a file called Unsplitter_ignore.txt with a list of prefixes to ignore.")
		fmt.Println("\nPress Enter to exit...")
		os.Stdin.Read(make([]byte, 1))
		return
	}

	root := os.Args[1]

	loadIgnoreList(root)

	logFile, err := os.Create("Unsplitter.log")
	if err != nil {
		fmt.Println("Error creating Unsplit.log:", err)
		return
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log := bufio.NewWriter(mw)
	defer log.Flush()

	data, err := os.ReadFile(root)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}
	rootLines := splitLines(string(data))
	fmt.Fprintln(log, len(rootLines), "lines in original file")

	// process root file (same rules as includes)
	finalLines, err := processFile(root, log)
	if err != nil {
		fmt.Println("Error processing root file:", err)
		return
	}

	outPath := root + ".unsplit.asm"
	_ = os.Remove(outPath)

	err = os.WriteFile(outPath, []byte(joinLines(finalLines)), 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}

	fmt.Fprintln(log, len(finalLines), "lines written")
	fmt.Println("\nUnsplitting done! Output:", filepath.Base(outPath))

	fmt.Println("\nPress Enter to exit...")
	os.Stdin.Read(make([]byte, 1))
}

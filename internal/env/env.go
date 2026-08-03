package env

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Report contains redacted comparison details between environment files.
// Only variable names are stored in slices; values are never included.
type Report struct {
	Missing     []string
	Extra       []string
	Duplicates  []string
	Empty       []string
	MissingFile bool
}

type parsed struct {
	values map[string]string
	dups   []string
}

func parseFile(path string) (parsed, error) {
	f, err := os.Open(path)
	if err != nil {
		return parsed{values: map[string]string{}}, err
	}
	defer f.Close()

	p := parsed{values: map[string]string{}}
	seen := map[string]bool{}
	dupSeen := map[string]bool{}

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(kv[0])
		if key == "" {
			continue
		}
		val := ""
		if len(kv) == 2 {
			val = strings.TrimSpace(kv[1])
		}
		if seen[key] {
			if !dupSeen[key] {
				p.dups = append(p.dups, key)
				dupSeen[key] = true
			}
		}
		seen[key] = true
		p.values[key] = val
	}
	if err := sc.Err(); err != nil {
		return p, err
	}
	return p, nil
}

// Compare checks envPath against examplePath.
// If envPath is missing, Report.MissingFile is set to true and nil error is returned.
// If examplePath is invalid or unreadable, an error is returned.
func Compare(envPath, examplePath string) (Report, error) {
	example, err := parseFile(examplePath)
	if err != nil {
		return Report{}, fmt.Errorf("failed to parse example env file: %w", err)
	}

	env, err := parseFile(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Report{MissingFile: true}, nil
		}
		return Report{}, fmt.Errorf("failed to parse env file: %w", err)
	}

	r := Report{MissingFile: false}
	for k := range example.values {
		if _, present := env.values[k]; !present {
			r.Missing = append(r.Missing, k)
		}
	}
	for k, v := range env.values {
		if _, present := example.values[k]; !present {
			r.Extra = append(r.Extra, k)
		}
		if v == "" {
			r.Empty = append(r.Empty, k)
		}
	}
	r.Duplicates = env.dups

	sort.Strings(r.Missing)
	sort.Strings(r.Extra)
	sort.Strings(r.Empty)
	sort.Strings(r.Duplicates)

	return r, nil
}

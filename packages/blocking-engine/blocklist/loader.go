package blocklist

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/domain"
)

// LoadFromReader reads domain entries from an io.Reader and populates the given Blocklist.
// Blank lines and lines starting with '#' are ignored. Inline comments are stripped.
// Invalid entries are safely skipped. Returns the count of newly added domains.
func LoadFromReader(r io.Reader, b *Blocklist) (int, error) {
	scanner := bufio.NewScanner(r)
	initialSize := b.Size()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if commentIdx := strings.Index(line, "#"); commentIdx != -1 {
			line = strings.TrimSpace(line[:commentIdx])
			if line == "" {
				continue
			}
		}

		norm, err := domain.Normalize(line)
		if err != nil {
			// Skip malformed domains safely
			continue
		}

		_ = b.Add(norm)
	}

	if err := scanner.Err(); err != nil {
		return b.Size() - initialSize, err
	}

	return b.Size() - initialSize, nil
}

// LoadFromFile opens the specified local file and populates the given Blocklist.
func LoadFromFile(filePath string, b *Blocklist) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return LoadFromReader(file, b)
}

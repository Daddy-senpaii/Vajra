package blocklist

import (
	_ "embed"
	"strings"
)

//go:embed data/adult-domains.txt
var defaultDomains string

// NewDefault creates a blocklist from the adult-domain snapshot shipped with
// the engine. The returned count is the number of valid, unique domains.
func NewDefault() (*Blocklist, int, error) {
	list := New()
	count, err := LoadFromReader(strings.NewReader(defaultDomains), list)
	if err != nil {
		return nil, count, err
	}
	return list, count, nil
}

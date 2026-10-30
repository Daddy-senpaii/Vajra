package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/allowlist"
	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/blocklist"
	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/engine"
)

func main() {
	blocklistPath := flag.String("blocklist", "", "Path to blocklist file")
	allowlistPath := flag.String("allowlist", "", "Path to allowlist file")
	domainArg := flag.String("domain", "", "Single domain to test")
	flag.Parse()

	al := allowlist.New()
	bl := blocklist.New()

	if *blocklistPath != "" {
		count, err := blocklist.LoadFromFile(*blocklistPath, bl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading blocklist: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Loaded %d domains into blocklist from %s\n", count, *blocklistPath)
	}

	if *allowlistPath != "" {
		file, err := os.Open(*allowlistPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading allowlist: %v\n", err)
			os.Exit(1)
		}
		count, err := al.Load(file)
		_ = file.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading allowlist: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Loaded %d domains into allowlist from %s\n", count, *allowlistPath)
	}

	eng := engine.New(al, bl)

	if *domainArg != "" {
		res := eng.Evaluate(*domainArg)
		printResult(*domainArg, res)
		return
	}

	// Interactive / Stdin streaming mode
	fmt.Println("Enter domain names to evaluate (Ctrl+D to exit):")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		res := eng.Evaluate(line)
		printResult(line, res)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}

func printResult(input string, res engine.Result) {
	fmt.Printf("[%s] %s (Normalized: %q) - %s\n",
		res.Decision, input, res.NormalizedDomain, res.Reason)
}

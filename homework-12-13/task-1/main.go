package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func countURL(client *http.Client, address, query string) (int, error) {
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, address, nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}
	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return 0, fmt.Errorf("HTTP %d", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("read response: %w", err)
	}
	return strings.Count(string(body), query), nil
}

func processURLs(input io.Reader, stdout, stderr io.Writer, client *http.Client, limit int, query string) error {
	slots := make(chan struct{}, limit)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var total int64
	var outputErr error

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		address := strings.TrimSpace(scanner.Text())
		if address == "" {
			continue
		}
		slots <- struct{}{}
		wg.Go(func() {
			defer func() { <-slots }()
			count, err := countURL(client, address, query)

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				_, err = fmt.Fprintf(stderr, "Error for %s: %v\n", address, err)
			} else {
				total += int64(count)
				_, err = fmt.Fprintf(stdout, "Count for %s: %d\n", address, count)
			}
			if outputErr == nil {
				outputErr = err
			}
		})
	}
	wg.Wait()
	_, err := fmt.Fprintf(stdout, "Total: %d\n", total)
	if scanErr := scanner.Err(); scanErr != nil {
		return fmt.Errorf("read stdin: %w", scanErr)
	}
	if outputErr != nil {
		return fmt.Errorf("write output: %w", outputErr)
	}
	if err != nil {
		return fmt.Errorf("write total: %w", err)
	}
	return nil
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("url-count", flag.ContinueOnError)
	flags.SetOutput(stderr)
	limit := flags.Int("k", 5, "maximum number of concurrent requests")
	query := flags.String("q", "Go", "case-sensitive search string")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "Read URLs from stdin, one per line.")
		return 2
	}
	if *limit < 1 {
		fmt.Fprintln(stderr, "-k must be positive")
		return 2
	}
	if *query == "" {
		fmt.Fprintln(stderr, "-q must not be empty")
		return 2
	}
	client := &http.Client{Timeout: 10 * time.Second}
	if err := processURLs(stdin, stdout, stderr, client, *limit, *query); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

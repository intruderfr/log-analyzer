package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Raw       string `json:"raw"`
}

type Pattern struct {
	Name        string `json:"name"`
	Pattern     string `json:"pattern"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type AnalysisResult struct {
	TotalLines     int                    `json:"total_lines"`
	LogsByLevel    map[string]int         `json:"logs_by_level"`
	PatternMatches map[string]int         `json:"pattern_matches"`
	TopErrors      []string               `json:"top_errors"`
	TimeRange      map[string]string      `json:"time_range"`
	Summary        string                 `json:"summary"`
}

var (
	logFile     = flag.String("file", "", "Path to log file")
	outputFile  = flag.String("output", "", "Output JSON report file")
	configFile  = flag.String("config", "", "Custom patterns config file")
	verbose     = flag.Bool("verbose", false, "Verbose output")
	watch       = flag.Bool("watch", false, "Watch file for changes")
	threshold   = flag.Int("threshold", 10, "Alert threshold for pattern matches")
)

// Default patterns for common log analysis
var defaultPatterns = []Pattern{
	{
		Name:        "HTTP_ERROR",
		Pattern:     `HTTP/1\.\d+"\s+[45]\d{2}`,
		Description: "HTTP 4xx/5xx error responses",
		Severity:    "high",
	},
	{
		Name:        "MEMORY_ERROR",
		Pattern:     `(?i)(out of memory|memory leak|heap|oom)`,
		Description: "Memory-related errors",
		Severity:    "critical",
	},
	{
		Name:        "CONNECTION_ERROR",
		Pattern:     `(?i)(connection (refused|reset|timeout)|network unreachable)`,
		Description: "Network connection issues",
		Severity:    "medium",
	},
	{
		Name:        "AUTH_FAILURE",
		Pattern:     `(?i)(authentication (failed|denied)|unauthorized|invalid (token|credentials))`,
		Description: "Authentication failures",
		Severity:    "high",
	},
	{
		Name:        "DATABASE_ERROR",
		Pattern:     `(?i)(database (connection|error)|sql (error|exception)|deadlock)`,
		Description: "Database-related errors",
		Severity:    "high",
	},
}

func main() {
	flag.Parse()

	if *logFile == "" {
		fmt.Println("Usage: log-analyzer -file=<path> [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	patterns := defaultPatterns
	if *configFile != "" {
		customPatterns, err := loadCustomPatterns(*configFile)
		if err != nil {
			fmt.Printf("Error loading custom patterns: %v\n", err)
			os.Exit(1)
		}
		patterns = append(patterns, customPatterns...)
	}

	if *watch {
		watchFile(*logFile, patterns)
	} else {
		result, err := analyzeLogFile(*logFile, patterns)
		if err != nil {
			fmt.Printf("Error analyzing log file: %v\n", err)
			os.Exit(1)
		}
		
		displayResults(result)
		
		if *outputFile != "" {
			saveResults(result, *outputFile)
		}
	}
}

func analyzeLogFile(filename string, patterns []Pattern) (*AnalysisResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	result := &AnalysisResult{
		LogsByLevel:    make(map[string]int),
		PatternMatches: make(map[string]int),
		TimeRange:      make(map[string]string),
		TopErrors:      []string{},
	}

	scanner := bufio.NewScanner(file)
	errorMessages := make(map[string]int)
	
	var firstTime, lastTime time.Time
	
	// Compile regex patterns once
	compiledPatterns := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		compiledPatterns[i] = regexp.MustCompile(pattern.Pattern)
	}
	
	// Common log level patterns
	logLevelRegex := regexp.MustCompile(`(?i)\b(DEBUG|INFO|WARN|ERROR|FATAL|TRACE)\b`)
	timestampRegex := regexp.MustCompile(`\d{4}-\d{2}-\d{2}[\sT]\d{2}:\d{2}:\d{2}`)

	for scanner.Scan() {
		line := scanner.Text()
		result.TotalLines++
		
		if *verbose && result.TotalLines%1000 == 0 {
			fmt.Printf("Processed %d lines...\n", result.TotalLines)
		}

		// Extract timestamp
		if timeMatch := timestampRegex.FindString(line); timeMatch != "" {
			if parsedTime, err := time.Parse("2006-01-02 15:04:05", timeMatch); err == nil {
				if firstTime.IsZero() {
					firstTime = parsedTime
				}
				lastTime = parsedTime
			} else if parsedTime, err := time.Parse("2006-01-02T15:04:05", timeMatch); err == nil {
				if firstTime.IsZero() {
					firstTime = parsedTime
				}
				lastTime = parsedTime
			}
		}

		// Extract log level
		if levelMatch := logLevelRegex.FindString(line); levelMatch != "" {
			level := strings.ToUpper(levelMatch)
			result.LogsByLevel[level]++
			
			// Track error messages for top errors
			if level == "ERROR" || level == "FATAL" {
				// Extract meaningful part of error message
				errorPart := extractErrorMessage(line)
				if errorPart != "" {
					errorMessages[errorPart]++
				}
			}
		}

		// Check against patterns
		for i, pattern := range patterns {
			if compiledPatterns[i].MatchString(line) {
				result.PatternMatches[pattern.Name]++
				
				if result.PatternMatches[pattern.Name] >= *threshold && pattern.Severity == "critical" {
					fmt.Printf("ALERT: Pattern '%s' exceeded threshold (%d matches)\n", pattern.Name, result.PatternMatches[pattern.Name])
				}
			}
		}
	}

	// Set time range
	if !firstTime.IsZero() {
		result.TimeRange["start"] = firstTime.Format("2006-01-02 15:04:05")
		result.TimeRange["end"] = lastTime.Format("2006-01-02 15:04:05")
		result.TimeRange["duration"] = lastTime.Sub(firstTime).String()
	}

	// Get top errors
	result.TopErrors = getTopErrors(errorMessages, 10)
	
	// Generate summary
	result.Summary = generateSummary(result)

	return result, scanner.Err()
}

func extractErrorMessage(line string) string {
	// Remove timestamp and log level prefix
	cleaned := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}[\sT]\d{2}:\d{2}:\d{2}[^\]]*\]\s*`).ReplaceAllString(line, "")
	cleaned = regexp.MustCompile(`(?i)^(ERROR|FATAL)[\s:]*`).ReplaceAllString(cleaned, "")
	
	// Take first 100 chars as error summary
	if len(cleaned) > 100 {
		cleaned = cleaned[:100] + "..."
	}
	
	return strings.TrimSpace(cleaned)
}

func getTopErrors(errorMessages map[string]int, limit int) []string {
	type errorCount struct {
		message string
		count   int
	}
	
	var errors []errorCount
	for msg, count := range errorMessages {
		errors = append(errors, errorCount{msg, count})
	}
	
	sort.Slice(errors, func(i, j int) bool {
		return errors[i].count > errors[j].count
	})
	
	var topErrors []string
	for i, err := range errors {
		if i >= limit {
			break
		}
		topErrors = append(topErrors, fmt.Sprintf("%s (count: %d)", err.message, err.count))
	}
	
	return topErrors
}

func generateSummary(result *AnalysisResult) string {
	summary := fmt.Sprintf("Analyzed %d log lines", result.TotalLines)
	
	if len(result.TimeRange) > 0 {
		summary += fmt.Sprintf(" spanning %s", result.TimeRange["duration"])
	}
	
	totalErrors := result.LogsByLevel["ERROR"] + result.LogsByLevel["FATAL"]
	if totalErrors > 0 {
		errorRate := float64(totalErrors) / float64(result.TotalLines) * 100
		summary += fmt.Sprintf(". Error rate: %.2f%% (%d errors)", errorRate, totalErrors)
	}
	
	criticalPatterns := 0
	for _, count := range result.PatternMatches {
		if count > 0 {
			criticalPatterns++
		}
	}
	
	if criticalPatterns > 0 {
		summary += fmt.Sprintf(". %d critical patterns detected", criticalPatterns)
	}
	
	return summary
}

func displayResults(result *AnalysisResult) {
	fmt.Printf("\n📊 Log Analysis Results\n")
	fmt.Printf("═══════════════════════\n\n")
	
	fmt.Printf("📈 Summary: %s\n\n", result.Summary)
	
	if len(result.TimeRange) > 0 {
		fmt.Printf("⏰ Time Range:\n")
		fmt.Printf("   Start: %s\n", result.TimeRange["start"])
		fmt.Printf("   End:   %s\n", result.TimeRange["end"])
		fmt.Printf("   Duration: %s\n\n", result.TimeRange["duration"])
	}
	
	if len(result.LogsByLevel) > 0 {
		fmt.Printf("📋 Logs by Level:\n")
		for level, count := range result.LogsByLevel {
			percentage := float64(count) / float64(result.TotalLines) * 100
			fmt.Printf("   %s: %d (%.1f%%)\n", level, count, percentage)
		}
		fmt.Printf("\n")
	}
	
	if len(result.PatternMatches) > 0 {
		fmt.Printf("🔍 Pattern Matches:\n")
		for pattern, count := range result.PatternMatches {
			if count > 0 {
				status := "✅"
				if count >= *threshold {
					status = "⚠️"
				}
				fmt.Printf("   %s %s: %d matches\n", status, pattern, count)
			}
		}
		fmt.Printf("\n")
	}
	
	if len(result.TopErrors) > 0 {
		fmt.Printf("🚨 Top Errors:\n")
		for i, error := range result.TopErrors {
			fmt.Printf("   %d. %s\n", i+1, error)
		}
		fmt.Printf("\n")
	}
}

func saveResults(result *AnalysisResult, filename string) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling results: %v\n", err)
		return
	}
	
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("Error writing results to file: %v\n", err)
		return
	}
	
	fmt.Printf("💾 Results saved to %s\n", filename)
}

func loadCustomPatterns(filename string) ([]Pattern, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	var patterns []Pattern
	err = json.Unmarshal(data, &patterns)
	return patterns, err
}

func watchFile(filename string, patterns []Pattern) {
	fmt.Printf("👀 Watching %s for changes...\n", filename)
	
	lastSize := int64(0)
	
	for {
		stat, err := os.Stat(filename)
		if err != nil {
			fmt.Printf("Error watching file: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
		
		if stat.Size() > lastSize {
			fmt.Printf("📄 File changed, analyzing new content...\n")
			
			result, err := analyzeLogFile(filename, patterns)
			if err != nil {
				fmt.Printf("Error analyzing file: %v\n", err)
			} else {
				displayResults(result)
			}
			
			lastSize = stat.Size()
		}
		
		time.Sleep(2 * time.Second)
	}
}
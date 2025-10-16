package main

import (
	"fmt"
	"log"
	"path/filepath"

	"smart-image-consolidators/analyzer"
	"smart-image-consolidators/ci"
	"smart-image-consolidators/metrics"
	"smart-image-consolidators/scanner"
	"smart-image-consolidators/security"
)

func main() {
	// Path to Dockerfiles in PR
	dockerfileDir := "./test_dockerfiles"

	// 1. Scan Dockerfiles
	dockerfiles, err := scanner.FindDockerfiles(dockerfileDir)
	if err != nil {
		log.Fatalf("Failed to scan Dockerfiles: %v", err)
	}
	if len(dockerfiles) == 0 {
		fmt.Println("No Dockerfiles found in PR")
		return
	}

	results := ""

	for _, df := range dockerfiles {
		fmt.Printf("Processing Dockerfile: %s\n", df)

		// 2. Analyze image for canonical base suggestion
		image, suggestion := analyzer.AnalyzeDockerfile(df)

		results += fmt.Sprintf("Dockerfile: %s\n", filepath.Base(df))
		results += fmt.Sprintf("Current Image: %s\n", image)
		results += fmt.Sprintf("Suggested Canonical Base: %s\n", suggestion)

		// 3. Measure performance metrics
		perf, err := metrics.MeasurePerformanceAndCost(image)
		if err != nil {
			results += fmt.Sprintf("Performance measurement failed: %v\n", err)
		} else {
			results += fmt.Sprintf("CPU: %s, Memory: %s, Storage: %s, Estimated Cost: $%.2f\n",
				perf.CPUUsage, perf.MemoryUsage, perf.Storage, perf.EstimatedCost)
		}

		// 4. Run AnchoreCTL security scan
		vulns, err := security.RunAnchoreScan(image)
		if err != nil {
			results += fmt.Sprintf("Security scan failed: %v\n", err)
		} else {
			results += fmt.Sprintf("Vulnerabilities found: %d\n", len(vulns))
		}

		results += "\n-----------------------------------\n"
	}

	// 5. Post results as PR comment
	if err := ci.CommentOnPR(results); err != nil {
		log.Printf("Failed to post PR comment: %v\n", err)
	}

	fmt.Println("✅ Smart Image Consolidator completed for all Dockerfiles")
}
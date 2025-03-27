package repository

import (
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

func ExtractChart(repo, chart string) (chartPath string, err error) {
	// Create a temporary directory for the chart
	tempDir, err := os.MkdirTemp("", "helm-chart-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Set up the Helm settings
	settings := cli.New()

	// Create a chart downloader
	chartDownloader := downloader.ChartDownloader{
		Out:              os.Stdout,
		Getters:          getter.All(settings),
		RepositoryConfig: settings.RepositoryConfig,
		RepositoryCache:  settings.RepositoryCache,
	}

	// Download the chart
	chartRef := fmt.Sprintf("%s/%s", repo, chart)
	savedChart, _, err := chartDownloader.DownloadTo(chartRef, "", tempDir)
	if err != nil {
		// Clean up temp directory on error
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to download chart %s: %w", chartRef, err)
	}

	// Return the path to the downloaded chart
	return savedChart, nil
}

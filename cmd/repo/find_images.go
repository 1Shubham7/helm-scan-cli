package repository

import (
	"context"
	"fmt"
	"io"

	"io/ioutil"
	"os"

	"path/filepath"
	"regexp"

	"strings"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type ImageInfo struct {
	Name   string `json:"name"`
	Size   string `json:"size"`
	Layers int    `json:"layers"`
}

var failed = "Failed to fetch"

func FindImages(extractedChartPath string) ([]ImageInfo, error) {
	patterns := []*regexp.Regexp{
		// Matches: image: "nginx:1.21"  OR  image: postgres:13
		// Explanation:
		// - `image:\s*` → Matches "image:" followed by any spaces.
		// - `["']?` → Optionally matches a single or double quote around the image name.
		// - `([^"'\s]+)` → Captures the actual image name (anything except spaces or quotes).
		// - `["']?` → Optionally matches the closing quote.
		// regexp.MustCompile(`image:\s*["']?([^"'\s]+)["']?`),

		// Matches: repository: "myrepo/custom-app:v2"  OR  repository: myrepo/custom-app:v2
		// Explanation:
		// - `repository:\s*` → Matches "repository:" followed by any spaces.
		// - `["']?` → Optionally matches a single or double quote around the repository name.
		// - `([^"'\s]+)` → Captures the repository name.
		// - `["']?` → Optionally matches the closing quote.
		regexp.MustCompile(`(?:image|repository):\s*["']?([a-zA-Z0-9.-]+(?:/[a-zA-Z0-9.-]+)*(?::[a-zA-Z0-9.-]+)?)["']?`),
	}

	images := make(map[string]ImageInfo)

	err := filepath.Walk(extractedChartPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip dirs and only processs yaml files
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatch(string(content), -1)
			for _, match := range matches {
				if len(match) > 1 {
					imageName := match[1]

					if isValidImageReference(imageName) {

						// Normalize image reference
						normalizedImage := normalizeImage(imageName)

						// Store unique image
						images[normalizedImage.Name] = normalizedImage
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk chart directory: %w", err)
	}

	// converting map to slice
	respImages := make([]ImageInfo, 0, len(images))
	for _, i := range images {
		respImages = append(respImages, i)
	}

	return respImages, nil
}


// normalizeImage standardizes image references
func normalizeImage(rawImage string) ImageInfo {
	// Split image into parts
	parts := strings.Split(rawImage, ":")
	name := parts[0]

	// // rawImage is image + tag
	size, layers, _ := getSizeAndLayers(rawImage)

	return ImageInfo{
		Name:       name,
		Size:       size,
		Layers:     layers,
	}
}

func getSizeAndLayers(imageWithTag string) (size string, layers int, err error) {
	ctx := context.Background()
	// my Docker daemon supports 1.47 at most, and client was of latest 1.48 version.
	// cli, err := client.NewClientWithOpts(client.FromEnv)

	// this give client highest API version that both the client and daemon support.
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return failed, -1, fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	imageCreateOptions := image.CreateOptions{}
	createReader, err := cli.ImageCreate(ctx, imageWithTag, imageCreateOptions)
	if err != nil {
		fmt.Println("one", err)
		return failed, -1, fmt.Errorf("failed to create image %s: %w", imageWithTag, err)
	}
	// Ensure we read and close the reader
	if createReader != nil {
		io.Copy(io.Discard, createReader)
		createReader.Close()
	}

	// Now inspect the image to get size and layers
	inspect, err := cli.ImageInspect(ctx, imageWithTag)
	if err != nil {
		return failed, -1, fmt.Errorf("failed to inspect image %s: %w", imageWithTag, err)
	}

	sizeInBytes := inspect.Size

	// Format size based on the condition
	if sizeInBytes > 1024*1024 { // More than 1 MB
		size = fmt.Sprintf("%.2f MB", float64(sizeInBytes)/(1024*1024))
	} else {
		size = fmt.Sprintf("%d bytes", sizeInBytes)
	}

	layers = len(inspect.RootFS.Layers)

	fmt.Printf("Image %s - Layers: %d, Size: %s\n", imageWithTag, layers, size)

	// Delete the image after usage
	_, err = cli.ImageRemove(ctx, imageWithTag, image.RemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
	if err != nil {
		fmt.Printf("Warning: failed to remove image %s: %v\n", imageWithTag, err)
	}
	fmt.Printf("image removed %s\n", imageWithTag)

	return size, layers, nil
}

// isValidImageReference checks if the image reference is potentially valid
func isValidImageReference(image string) bool {
	// Checks to filter out invalid references
	invalidPatterns := []string{
		`^\{\{.*\}\}$`,                       // Helm template variables
		`^["']+$`,                            // Just quotes
		`^[{}]+$`,                            // Just braces
		`^\s*$`,                              // Empty or whitespace
		`^[/:.]+$`,                           // Just delimiters
		`^(your-image|placeholder|example)$`, // Common placeholder names
		`^(oci|registry)$`,                   // Generic registry names
	}

	for _, pattern := range invalidPatterns {
		match, _ := regexp.MatchString(pattern, image)
		if match {
			return false
		}
	}

	// Additional check: must contain alphanumeric characters
	return regexp.MustCompile(`[a-zA-Z0-9]`).MatchString(image)
}
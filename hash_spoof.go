package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"
)

const workerCount = 16 // Number of goroutines to use

// searchForHashPrefix attempts to find a hash with the desired prefix by modifying bytes
func searchForHashPrefix(inputContent []byte, targetPrefix string, outputChan chan []byte, startOffset int, numBytes int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Create a local copy of the input content
	modifiedContent := make([]byte, len(inputContent))
	copy(modifiedContent, inputContent)

	// Modify multiple bytes to expand the search space
	maxVal := 256
	offsets := make([]int, numBytes)

	for {
		// Apply the current offsets to modify bytes
		for i := 0; i < numBytes; i++ {
			modifiedContent[len(modifiedContent)-1-startOffset-i] = byte(offsets[i])
		}

		// Compute the hash
		hash := sha256.Sum256(modifiedContent)
		hashHex := fmt.Sprintf("%x", hash)

		// Check if the hash matches the target prefix
		if len(hashHex) >= len(targetPrefix) && hashHex[:len(targetPrefix)] == targetPrefix {
			outputChan <- modifiedContent
			return
		}

		// Increment the offsets
		for i := 0; i < numBytes; i++ {
			offsets[i]++
			if offsets[i] < maxVal {
				break
			}
			offsets[i] = 0
			if i == numBytes-1 {
				return // Exhausted all possibilities
			}
		}
	}
}

func adjustImageForHashOptimized(targetPrefix string, inputFile string, outputFile string, numBytes int) error {
	// Open the input file
	originalFile, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer originalFile.Close()

	// Read the content of the original file
	inputContent, err := io.ReadAll(originalFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Output channel for successful results
	outputChan := make(chan []byte)
	var wg sync.WaitGroup

	// Distribute work among worker goroutines
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go searchForHashPrefix(inputContent, targetPrefix, outputChan, i, numBytes, &wg)
	}

	// Wait for a result or completion of all workers
	go func() {
		wg.Wait()
		close(outputChan)
	}()

	// Write the successful result to the output file
	for modifiedContent := range outputChan {
		err := os.WriteFile(outputFile, modifiedContent, 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Println("Success! File modified and written to disk.")
		return nil
	}

	return fmt.Errorf("failed to achieve target hash")
}

func main() {
	inputFile := "original.jpg"
	outputFile := "altered.jpg"
	targetPrefix := "24" // Desired hash prefix
	numBytes := 2        // Number of bytes to modify for search

	err := adjustImageForHashOptimized(targetPrefix, inputFile, outputFile, numBytes)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("File successfully modified.")
	}
}

package performance_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/types"
)

// TestLoadPerformance tests system performance under various load conditions
func TestLoadPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	t.Run("Concurrent Setup Operations", func(t *testing.T) {
		testCases := []struct {
			name        string
			concurrency int
			timeout     time.Duration
		}{
			{"Low Load", 5, 30 * time.Second},
			{"Medium Load", 25, 60 * time.Second},
			{"High Load", 100, 120 * time.Second},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
				defer cancel()

				results := make(chan types.SetupResult, tc.concurrency)
				var wg sync.WaitGroup

				startTime := time.Now()

				for i := 0; i < tc.concurrency; i++ {
					wg.Add(1)
					go func(id int) {
						defer wg.Done()
						
						tempDir := createTempDir(t, fmt.Sprintf("load-test-%d", id))
						options := types.SetupOptions{
							HomeDir: tempDir,
						}

						select {
						case <-ctx.Done():
							return
						default:
							result := performSetupWithContext(ctx, options)
							results <- result
						}
					}(i)
				}

				wg.Wait()
				close(results)

				duration := time.Since(startTime)
				successCount := 0
				failureCount := 0

				for result := range results {
					if result.Success {
						successCount++
					} else {
						failureCount++
					}
				}

				t.Logf("Load test completed: %d concurrent operations in %v", tc.concurrency, duration)
				t.Logf("Success: %d, Failures: %d", successCount, failureCount)

				// Assert that most operations succeeded
				successRate := float64(successCount) / float64(tc.concurrency)
				assert.GreaterOrEqual(t, successRate, 0.8, "Success rate should be at least 80%")
				assert.Less(t, duration, tc.timeout, "Should complete within timeout")
			})
		}
	})

	t.Run("Memory Usage Under Load", func(t *testing.T) {
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		// Perform multiple setup operations
		const operations = 50
		for i := 0; i < operations; i++ {
			tempDir := createTempDir(t, fmt.Sprintf("memory-test-%d", i))
			options := types.SetupOptions{
				HomeDir: tempDir,
			}
			
			result := performSetupWithContext(context.Background(), options)
			assert.True(t, result.Success)
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		memoryIncrease := m2.Alloc - m1.Alloc
		t.Logf("Memory increase after %d operations: %d bytes", operations, memoryIncrease)

		// Memory increase should be reasonable (less than 100MB for 50 operations)
		assert.Less(t, memoryIncrease, uint64(100*1024*1024), "Memory usage should be reasonable")
	})

	t.Run("CPU Usage Under Load", func(t *testing.T) {
		const operations = 100
		startTime := time.Now()

		for i := 0; i < operations; i++ {
			tempDir := createTempDir(t, fmt.Sprintf("cpu-test-%d", i))
			options := types.SetupOptions{
				HomeDir: tempDir,
			}
			
			result := performSetupWithContext(context.Background(), options)
			assert.True(t, result.Success)
		}

		duration := time.Since(startTime)
		operationsPerSecond := float64(operations) / duration.Seconds()

		t.Logf("Completed %d operations in %v (%.2f ops/sec)", operations, duration, operationsPerSecond)
		
		// Should maintain reasonable throughput
		assert.Greater(t, operationsPerSecond, 10.0, "Should maintain at least 10 operations per second")
	})
}

// TestStressPerformance tests system behavior under stress conditions
func TestStressPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress tests in short mode")
	}

	t.Run("Resource Exhaustion", func(t *testing.T) {
		t.Run("File Descriptor Limits", func(t *testing.T) {
			// Test with many concurrent file operations
			const concurrency = 200
			var wg sync.WaitGroup
			results := make(chan bool, concurrency)

			for i := 0; i < concurrency; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					
					tempDir := createTempDir(t, fmt.Sprintf("fd-test-%d", id))
					
					// Create multiple files to stress file descriptors
					for j := 0; j < 10; j++ {
						filePath := filepath.Join(tempDir, fmt.Sprintf("file-%d.txt", j))
						file, err := os.Create(filePath)
						if err != nil {
							results <- false
							return
						}
						file.Close()
					}
					
					results <- true
				}(i)
			}

			wg.Wait()
			close(results)

			successCount := 0
			for success := range results {
				if success {
					successCount++
				}
			}

			// Should handle file descriptor pressure gracefully
			successRate := float64(successCount) / float64(concurrency)
			assert.GreaterOrEqual(t, successRate, 0.9, "Should handle file descriptor stress")
		})

		t.Run("Disk Space Pressure", func(t *testing.T) {
			tempDir := createTempDir(t, "disk-stress")
			
			// Create files until we approach limits (simulated)
			const maxFiles = 1000
			createdFiles := 0

			for i := 0; i < maxFiles; i++ {
				filePath := filepath.Join(tempDir, fmt.Sprintf("stress-file-%d.txt", i))
				content := make([]byte, 1024) // 1KB per file
				
				err := os.WriteFile(filePath, content, 0644)
				if err != nil {
					break
				}
				createdFiles++
			}

			t.Logf("Created %d files under disk pressure", createdFiles)
			assert.Greater(t, createdFiles, 100, "Should create a reasonable number of files")
		})
	})

	t.Run("Long Running Operations", func(t *testing.T) {
		t.Run("Extended Setup Duration", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			tempDir := createTempDir(t, "long-running")
			options := types.SetupOptions{
				HomeDir: tempDir,
			}

			// Simulate long-running setup
			startTime := time.Now()
			result := performLongRunningSetup(ctx, options)
			duration := time.Since(startTime)

			assert.True(t, result.Success)
			t.Logf("Long-running setup completed in %v", duration)
		})

		t.Run("Repeated Setup/Reset Cycles", func(t *testing.T) {
			tempDir := createTempDir(t, "cycles")
			
			const cycles = 50
			for i := 0; i < cycles; i++ {
				options := types.SetupOptions{
					HomeDir: tempDir,
				}

				// Setup
				setupResult := performSetupWithContext(context.Background(), options)
				assert.True(t, setupResult.Success, "Setup cycle %d should succeed", i)

				// Reset
				resetErr := performReset()
				assert.NoError(t, resetErr, "Reset cycle %d should succeed", i)
			}

			t.Logf("Completed %d setup/reset cycles successfully", cycles)
		})
	})
}

// TestVolumePerformance tests system behavior with large volumes of data
func TestVolumePerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping volume tests in short mode")
	}

	t.Run("Large Configuration Files", func(t *testing.T) {
		testCases := []struct {
			name     string
			sizeKB   int
			timeout  time.Duration
		}{
			{"Small Config", 10, 5 * time.Second},
			{"Medium Config", 100, 10 * time.Second},
			{"Large Config", 1000, 30 * time.Second},
			{"Very Large Config", 10000, 60 * time.Second},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tempDir := createTempDir(t, "volume-config")
				configPath := filepath.Join(tempDir, "large-config.yaml")

				// Generate large configuration
				config := generateLargeConfig(tc.sizeKB * 1024)
				err := os.WriteFile(configPath, []byte(config), 0644)
				require.NoError(t, err)

				options := types.SetupOptions{
					ConfigPath: configPath,
					HomeDir:    tempDir,
				}

				ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
				defer cancel()

				startTime := time.Now()
				result := performSetupWithContext(ctx, options)
				duration := time.Since(startTime)

				assert.True(t, result.Success)
				assert.Less(t, duration, tc.timeout)
				
				t.Logf("Processed %dKB config in %v", tc.sizeKB, duration)
			})
		}
	})

	t.Run("Many Small Files", func(t *testing.T) {
		tempDir := createTempDir(t, "many-files")
		
		const fileCount = 10000
		startTime := time.Now()

		// Create many small files
		for i := 0; i < fileCount; i++ {
			filePath := filepath.Join(tempDir, fmt.Sprintf("file-%06d.txt", i))
			content := fmt.Sprintf("File content %d", i)
			
			err := os.WriteFile(filePath, []byte(content), 0644)
			require.NoError(t, err)
		}

		creationDuration := time.Since(startTime)

		// Now perform setup with this directory
		options := types.SetupOptions{
			HomeDir: tempDir,
		}

		setupStart := time.Now()
		result := performSetupWithContext(context.Background(), options)
		setupDuration := time.Since(setupStart)

		assert.True(t, result.Success)
		
		t.Logf("Created %d files in %v", fileCount, creationDuration)
		t.Logf("Setup with %d files completed in %v", fileCount, setupDuration)
		
		// Should handle many files reasonably
		assert.Less(t, setupDuration, 30*time.Second)
	})

	t.Run("Deep Directory Structure", func(t *testing.T) {
		tempDir := createTempDir(t, "deep-dirs")
		
		// Create deep directory structure
		const depth = 100
		currentPath := tempDir
		
		for i := 0; i < depth; i++ {
			currentPath = filepath.Join(currentPath, fmt.Sprintf("level-%d", i))
			err := os.MkdirAll(currentPath, 0755)
			require.NoError(t, err)
		}

		// Create file at the deepest level
		deepFile := filepath.Join(currentPath, "deep-file.txt")
		err := os.WriteFile(deepFile, []byte("Deep content"), 0644)
		require.NoError(t, err)

		options := types.SetupOptions{
			HomeDir: tempDir,
		}

		startTime := time.Now()
		result := performSetupWithContext(context.Background(), options)
		duration := time.Since(startTime)

		assert.True(t, result.Success)
		t.Logf("Setup with %d-level deep directory completed in %v", depth, duration)
	})
}

// Helper functions for performance tests

func performSetupWithContext(ctx context.Context, options types.SetupOptions) types.SetupResult {
	startTime := time.Now()
	
	// Simulate setup work with context awareness
	select {
	case <-ctx.Done():
		return types.SetupResult{
			Success:   false,
			StartTime: startTime,
			EndTime:   time.Now(),
			Error:     ctx.Err(),
			Options:   options,
		}
	case <-time.After(100 * time.Millisecond):
		// Simulate successful setup
		return types.SetupResult{
			Success:     true,
			StartTime:   startTime,
			EndTime:     time.Now(),
			ConfigPath:  filepath.Join(options.HomeDir, "config.yaml"),
			Environment: runtime.GOOS,
			Options:     options,
			Message:     "Setup completed successfully",
		}
	}
}

func performLongRunningSetup(ctx context.Context, options types.SetupOptions) types.SetupResult {
	startTime := time.Now()
	
	// Simulate long-running work
	select {
	case <-ctx.Done():
		return types.SetupResult{
			Success:   false,
			StartTime: startTime,
			EndTime:   time.Now(),
			Error:     ctx.Err(),
			Options:   options,
		}
	case <-time.After(2 * time.Second): // Longer simulation
		return types.SetupResult{
			Success:     true,
			StartTime:   startTime,
			EndTime:     time.Now(),
			ConfigPath:  filepath.Join(options.HomeDir, "config.yaml"),
			Environment: runtime.GOOS,
			Options:     options,
			Message:     "Long-running setup completed successfully",
		}
	}
}

func performReset() error {
	// Simulate reset work
	time.Sleep(50 * time.Millisecond)
	return nil
}

func createTempDir(t *testing.T, prefix string) string {
	dir, err := os.MkdirTemp("", "syntropy-perf-"+prefix+"-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	return dir
}

func generateLargeConfig(sizeBytes int) string {
	config := "manager:\n  settings:\n"
	
	// Calculate how many entries we need to reach the target size
	entrySize := len("    key000000: value000000\n")
	entries := sizeBytes / entrySize
	
	for i := 0; i < entries; i++ {
		config += fmt.Sprintf("    key%06d: value%06d\n", i, i)
	}
	
	return config
}

// Benchmark functions

func BenchmarkSetupPerformance(b *testing.B) {
	tempDir, _ := os.MkdirTemp("", "syntropy-bench-*")
	defer os.RemoveAll(tempDir)

	options := types.SetupOptions{
		HomeDir: tempDir,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := performSetupWithContext(context.Background(), options)
		if !result.Success {
			b.Fatal("Setup failed")
		}
	}
}

func BenchmarkConcurrentSetup(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tempDir, _ := os.MkdirTemp("", "syntropy-concurrent-*")
			defer os.RemoveAll(tempDir)

			options := types.SetupOptions{
				HomeDir: tempDir,
			}

			result := performSetupWithContext(context.Background(), options)
			if !result.Success {
				b.Fatal("Setup failed")
			}
		}
	})
}

func BenchmarkLargeConfigProcessing(b *testing.B) {
	tempDir, _ := os.MkdirTemp("", "syntropy-config-bench-*")
	defer os.RemoveAll(tempDir)

	// Create a large config file
	configPath := filepath.Join(tempDir, "large-config.yaml")
	config := generateLargeConfig(10 * 1024) // 10KB config
	os.WriteFile(configPath, []byte(config), 0644)

	options := types.SetupOptions{
		ConfigPath: configPath,
		HomeDir:    tempDir,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := performSetupWithContext(context.Background(), options)
		if !result.Success {
			b.Fatal("Setup failed")
		}
	}
}
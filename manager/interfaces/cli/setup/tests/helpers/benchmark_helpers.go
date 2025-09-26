// Package helpers provides benchmark utilities for performance testing
package helpers

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/types"
)

// BenchmarkHelper provides utilities for benchmark tests
type BenchmarkHelper struct {
	StartTime    time.Time
	EndTime      time.Time
	MemoryBefore runtime.MemStats
	MemoryAfter  runtime.MemStats
}

// NewBenchmarkHelper creates a new benchmark helper
func NewBenchmarkHelper() *BenchmarkHelper {
	return &BenchmarkHelper{}
}

// StartBenchmark starts benchmark measurement
func (bh *BenchmarkHelper) StartBenchmark() {
	runtime.GC() // Force garbage collection
	runtime.ReadMemStats(&bh.MemoryBefore)
	bh.StartTime = time.Now()
}

// EndBenchmark ends benchmark measurement
func (bh *BenchmarkHelper) EndBenchmark() {
	bh.EndTime = time.Now()
	runtime.ReadMemStats(&bh.MemoryAfter)
}

// GetDuration returns the benchmark duration
func (bh *BenchmarkHelper) GetDuration() time.Duration {
	return bh.EndTime.Sub(bh.StartTime)
}

// GetMemoryUsage returns memory usage in bytes
func (bh *BenchmarkHelper) GetMemoryUsage() uint64 {
	return bh.MemoryAfter.TotalAlloc - bh.MemoryBefore.TotalAlloc
}

// GetMemoryUsageMB returns memory usage in megabytes
func (bh *BenchmarkHelper) GetMemoryUsageMB() float64 {
	return float64(bh.GetMemoryUsage()) / (1024 * 1024)
}

// BenchmarkSetupOperation benchmarks a setup operation
func BenchmarkSetupOperation(b *testing.B, setupFunc func() (types.SetupResult, error)) {
	b.Helper()
	
	for i := 0; i < b.N; i++ {
		helper := NewBenchmarkHelper()
		helper.StartBenchmark()
		
		result, err := setupFunc()
		
		helper.EndBenchmark()
		
		if err != nil {
			b.Fatalf("Setup operation failed: %v", err)
		}
		
		if !result.Success {
			b.Fatalf("Setup operation was not successful")
		}
		
		// Report custom metrics
		b.ReportMetric(float64(helper.GetDuration().Nanoseconds()), "ns/op")
		b.ReportMetric(helper.GetMemoryUsageMB(), "MB/op")
	}
}

// BenchmarkConcurrentSetup benchmarks concurrent setup operations
func BenchmarkConcurrentSetup(b *testing.B, concurrency int, setupFunc func() (types.SetupResult, error)) {
	b.Helper()
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			helper := NewBenchmarkHelper()
			helper.StartBenchmark()
			
			result, err := setupFunc()
			
			helper.EndBenchmark()
			
			if err != nil {
				b.Errorf("Concurrent setup operation failed: %v", err)
				continue
			}
			
			if !result.Success {
				b.Errorf("Concurrent setup operation was not successful")
				continue
			}
		}
	})
}

// BenchmarkWithTimeout benchmarks an operation with timeout
func BenchmarkWithTimeout(b *testing.B, timeout time.Duration, operation func(ctx context.Context) error) {
	b.Helper()
	
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		
		helper := NewBenchmarkHelper()
		helper.StartBenchmark()
		
		err := operation(ctx)
		
		helper.EndBenchmark()
		cancel()
		
		if err != nil {
			if err == context.DeadlineExceeded {
				b.Fatalf("Operation timed out after %v", timeout)
			} else {
				b.Fatalf("Operation failed: %v", err)
			}
		}
		
		// Report timing
		b.ReportMetric(float64(helper.GetDuration().Nanoseconds()), "ns/op")
	}
}

// MeasureMemoryUsage measures memory usage of a function
func MeasureMemoryUsage(b *testing.B, operation func()) {
	b.Helper()
	
	var memBefore, memAfter runtime.MemStats
	
	runtime.GC()
	runtime.ReadMemStats(&memBefore)
	
	operation()
	
	runtime.GC()
	runtime.ReadMemStats(&memAfter)
	
	memUsed := memAfter.TotalAlloc - memBefore.TotalAlloc
	b.ReportMetric(float64(memUsed)/(1024*1024), "MB")
}

// BenchmarkFileOperations benchmarks file I/O operations
func BenchmarkFileOperations(b *testing.B, fileSize int, operation func([]byte) error) {
	b.Helper()
	
	data := GenerateTestData(fileSize)
	
	b.ResetTimer()
	b.SetBytes(int64(fileSize))
	
	for i := 0; i < b.N; i++ {
		if err := operation(data); err != nil {
			b.Fatalf("File operation failed: %v", err)
		}
	}
}

// BenchmarkNetworkOperations benchmarks network operations
func BenchmarkNetworkOperations(b *testing.B, operation func() error) {
	b.Helper()
	
	for i := 0; i < b.N; i++ {
		helper := NewBenchmarkHelper()
		helper.StartBenchmark()
		
		err := operation()
		
		helper.EndBenchmark()
		
		if err != nil {
			b.Fatalf("Network operation failed: %v", err)
		}
		
		// Report latency
		b.ReportMetric(float64(helper.GetDuration().Milliseconds()), "ms/op")
	}
}

// BenchmarkConfigurationParsing benchmarks configuration parsing
func BenchmarkConfigurationParsing(b *testing.B, configData []byte, parseFunc func([]byte) (types.SetupConfig, error)) {
	b.Helper()
	
	b.SetBytes(int64(len(configData)))
	
	for i := 0; i < b.N; i++ {
		helper := NewBenchmarkHelper()
		helper.StartBenchmark()
		
		_, err := parseFunc(configData)
		
		helper.EndBenchmark()
		
		if err != nil {
			b.Fatalf("Configuration parsing failed: %v", err)
		}
		
		b.ReportMetric(helper.GetMemoryUsageMB(), "MB/op")
	}
}

// BenchmarkValidation benchmarks validation operations
func BenchmarkValidation(b *testing.B, validateFunc func() (types.ValidationResult, error)) {
	b.Helper()
	
	for i := 0; i < b.N; i++ {
		helper := NewBenchmarkHelper()
		helper.StartBenchmark()
		
		result, err := validateFunc()
		
		helper.EndBenchmark()
		
		if err != nil {
			b.Fatalf("Validation failed: %v", err)
		}
		
		if !result.Valid {
			b.Fatalf("Validation result was invalid")
		}
		
		b.ReportMetric(float64(helper.GetDuration().Microseconds()), "μs/op")
	}
}

// BenchmarkResourceUsage benchmarks resource usage patterns
func BenchmarkResourceUsage(b *testing.B, operation func() error) {
	b.Helper()
	
	var totalCPU time.Duration
	var maxMemory uint64
	
	for i := 0; i < b.N; i++ {
		helper := NewBenchmarkHelper()
		helper.StartBenchmark()
		
		err := operation()
		
		helper.EndBenchmark()
		
		if err != nil {
			b.Fatalf("Operation failed: %v", err)
		}
		
		totalCPU += helper.GetDuration()
		if helper.GetMemoryUsage() > maxMemory {
			maxMemory = helper.GetMemoryUsage()
		}
	}
	
	avgCPU := totalCPU / time.Duration(b.N)
	b.ReportMetric(float64(avgCPU.Nanoseconds()), "ns/op")
	b.ReportMetric(float64(maxMemory)/(1024*1024), "peak-MB")
}

// BenchmarkScalability benchmarks scalability with different loads
func BenchmarkScalability(b *testing.B, loads []int, operation func(load int) error) {
	b.Helper()
	
	for _, load := range loads {
		b.Run(fmt.Sprintf("load-%d", load), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				helper := NewBenchmarkHelper()
				helper.StartBenchmark()
				
				err := operation(load)
				
				helper.EndBenchmark()
				
				if err != nil {
					b.Fatalf("Operation failed with load %d: %v", load, err)
				}
				
				b.ReportMetric(float64(helper.GetDuration().Nanoseconds()), "ns/op")
				b.ReportMetric(helper.GetMemoryUsageMB(), "MB/op")
			}
		})
	}
}

// SetBenchmarkDefaults sets default benchmark parameters
func SetBenchmarkDefaults(b *testing.B) {
	b.Helper()
	
	// Set reasonable defaults for benchmarks
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}
	
	// Ensure consistent benchmark environment
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GC()
}
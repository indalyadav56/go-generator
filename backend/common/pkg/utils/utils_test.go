package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetServiceAddress(t *testing.T) {

	address := GetServiceAddress()

	fmt.Println("address", address)
}

func TestFormatMemorySize(t *testing.T) {
	// Test table structure
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		// Test cases for bytes (B)
		{
			name:     "Zero bytes",
			input:    0,
			expected: "0 B",
		},
		{
			name:     "Small bytes",
			input:    512,
			expected: "512 B",
		},
		{
			name:     "Bytes upper limit",
			input:    1023,
			expected: "1023 B",
		},

		// Test cases for kilobytes (KB)
		{
			name:     "One KB exact",
			input:    1024,
			expected: "1.00 KB",
		},
		{
			name:     "Middle range KB",
			input:    1500,
			expected: "1.46 KB",
		},
		{
			name:     "KB upper limit",
			input:    1024*1024 - 1,
			expected: "1024.00 KB",
		},

		// Test cases for megabytes (MB)
		{
			name:     "One MB exact",
			input:    1024 * 1024,
			expected: "1.00 MB",
		},
		{
			name:     "Middle range MB",
			input:    1024 * 1024 * 50,
			expected: "50.00 MB",
		},
		{
			name:     "MB upper limit",
			input:    1024*1024*1024 - 1,
			expected: "1024.00 MB",
		},

		// Test cases for gigabytes (GB)
		{
			name:     "One GB exact",
			input:    1024 * 1024 * 1024,
			expected: "1.00 GB",
		},
		{
			name:     "Multiple GB",
			input:    1024 * 1024 * 1024 * 2,
			expected: "2.00 GB",
		},

		// Test cases for []byte input
		{
			name:     "Empty byte slice",
			input:    []byte{},
			expected: "0 B",
		},
		{
			name:     "Small byte slice",
			input:    make([]byte, 100),
			expected: "100 B",
		},
		{
			name:     "KB range byte slice",
			input:    make([]byte, 1500),
			expected: "1.46 KB",
		},
		{
			name:     "MB range byte slice",
			input:    make([]byte, 1024*1024*2),
			expected: "2.00 MB",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			switch v := tt.input.(type) {
			case int:
				result = FormatMemorySize(v)
			case []byte:
				result = FormatMemorySize(v)
			}

			if result != tt.expected {
				t.Errorf("FormatMemorySize(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkFormatMemorySize(b *testing.B) {
	benchmarks := []struct {
		name  string
		input int
	}{
		{"Bytes", 500},
		{"Kilobytes", 1500},
		{"Megabytes", 1024 * 1024 * 5},
		{"Gigabytes", 1024 * 1024 * 1024 * 2},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				FormatMemorySize(bm.input)
			}
		})
	}
}

// Example test for documentation
func ExampleFormatMemorySize() {
	fmt.Println(FormatMemorySize(1500))
	fmt.Println(FormatMemorySize(2 * 1024 * 1024))
	fmt.Println(FormatMemorySize([]byte{1, 2, 3, 4, 5}))
	// Output:
	// 1.46 KB
	// 2.00 MB
	// 5 B
}

// Test edge cases and potential panics
func TestFormatMemorySizeEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{
			name:     "Max int",
			input:    int(^uint(0) >> 1),
			expected: "8589934592.00 GB", // This may vary based on 32/64 bit systems
		},
		{
			name:     "Negative number",
			input:    -1024,
			expected: "-1.00 KB",
		},
		{
			name:     "Very small negative",
			input:    -1,
			expected: "-1 B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatMemorySize(tt.input)
			if result != tt.expected {
				t.Errorf("FormatMemorySize(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Test helper function to create large byte slices
func createByteSlice(size int) []byte {
	return make([]byte, size)
}

// Test memory handling with large inputs
func TestFormatMemorySizeLargeInputs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping large input test in short mode")
	}

	largeSizes := []int{
		1024 * 1024 * 100,      // 100 MB
		1024 * 1024 * 1024,     // 1 GB
		1024 * 1024 * 1024 * 2, // 2 GB
	}

	for _, size := range largeSizes {
		t.Run(fmt.Sprintf("Size-%d", size), func(t *testing.T) {
			data := createByteSlice(size)
			result := FormatMemorySize(data)

			// Verify the result contains the expected unit
			expected := "GB"
			if size < 1024*1024*1024 {
				expected = "MB"
			}

			if !strings.Contains(result, expected) {
				t.Errorf("Expected result to contain %s, got %s", expected, result)
			}
		})
	}
}

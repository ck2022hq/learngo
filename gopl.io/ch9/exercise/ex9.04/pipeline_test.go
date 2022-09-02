package pipeline

import (
	"testing"
)

func benchmarkPipeline(b *testing.B, stages int) {
	// b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline(stages)
	}
}

func BenchmarkPipeline10(b *testing.B)     { benchmarkPipeline(b, 10) }
func BenchmarkPipeline100(b *testing.B)    { benchmarkPipeline(b, 100) }
func BenchmarkPipeline1000(b *testing.B)   { benchmarkPipeline(b, 1000) }
func BenchmarkPipeline10000(b *testing.B)  { benchmarkPipeline(b, 10000) }
func BenchmarkPipeline100000(b *testing.B) { benchmarkPipeline(b, 100000) }

func BenchmarkPipeline1000000(b *testing.B) { benchmarkPipeline(b, 1000000) }

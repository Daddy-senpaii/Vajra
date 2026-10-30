package benchmarks_test

import (
	"fmt"
	"testing"

	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/allowlist"
	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/blocklist"
	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/domain"
	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/engine"
)

func setupBenchmarkEngine(numBlocked int) *engine.Engine {
	al := allowlist.New()
	bl := blocklist.New()

	_ = al.Add("allowed.com")
	_ = al.Add("safe.blocked-domain.com")

	for i := 0; i < numBlocked; i++ {
		_ = bl.Add(fmt.Sprintf("adult-site-%d.com", i))
	}
	_ = bl.Add("blocked-domain.com")

	return engine.New(al, bl)
}

func BenchmarkDomain_Normalize(b *testing.B) {
	input := "  HTTPS://WWW.Example-Site.COM:8080/path?query=1  "
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = domain.Normalize(input)
	}
}

func BenchmarkEngine_AllowlistHit(b *testing.B) {
	eng := setupBenchmarkEngine(10000)
	input := "allowed.com"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = eng.Evaluate(input)
	}
}

func BenchmarkEngine_BlocklistExactHit(b *testing.B) {
	eng := setupBenchmarkEngine(10000)
	input := "adult-site-500.com"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = eng.Evaluate(input)
	}
}

func BenchmarkEngine_BlocklistSubdomainHit(b *testing.B) {
	eng := setupBenchmarkEngine(10000)
	input := "deep.sub.adult-site-500.com"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = eng.Evaluate(input)
	}
}

func BenchmarkEngine_UnblockedMiss(b *testing.B) {
	eng := setupBenchmarkEngine(10000)
	input := "innocent-news-portal.org"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = eng.Evaluate(input)
	}
}

func BenchmarkEngine_ParallelLookups(b *testing.B) {
	eng := setupBenchmarkEngine(10000)
	domains := []string{
		"allowed.com",
		"adult-site-500.com",
		"deep.sub.adult-site-500.com",
		"safe.blocked-domain.com",
		"innocent-news-portal.org",
	}

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		idx := 0
		for pb.Next() {
			_ = eng.Evaluate(domains[idx%len(domains)])
			idx++
		}
	})
}

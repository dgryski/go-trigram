package trigram

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func mktri(s string) T { return T(uint32(s[0])<<16 | uint32(s[1])<<8 | uint32(s[2])) }

func mktris(ss ...string) []T {
	var ts []T
	for _, s := range ss {
		ts = append(ts, mktri(s))
	}
	return ts
}

func TestExtract(t *testing.T) {

	tests := []struct {
		s    string
		want []T
	}{
		{"", nil},
		{"a", nil},
		{"ab", nil},
		{"abc", mktris("abc")},
		{"abcabc", mktris("abc", "bca", "cab")},
		{"abcd", mktris("abc", "bcd")},
	}

	for _, tt := range tests {
		if got := Extract(tt.s, nil); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Extract(%q)=%+v, want %+v", tt.s, got, tt.want)
		}
	}
}

func TestQuery(t *testing.T) {

	s := []string{
		"foo",
		"foobar",
		"foobfoo",
		"quxzoot",
		"zotzot",
		"azotfoba",
	}

	idx := NewIndex(s)

	tests := []struct {
		q   string
		ids []DocID
	}{
		{"", []DocID{0, 1, 2, 3, 4, 5}},
		{"foo", []DocID{0, 1, 2}},
		{"foob", []DocID{1, 2}},
		{"zot", []DocID{4, 5}},
		{"oba", []DocID{1, 5}},
	}

	for _, tt := range tests {
		if got := idx.Query(tt.q); !reflect.DeepEqual(got, tt.ids) {
			t.Errorf("Query(%q)=%+v, want %+v", tt.q, got, tt.ids)
		}
	}

	idx.Add("zlot")
	docs := idx.Query("lot")
	if len(docs) != 1 || docs[0] != 6 {
		t.Errorf("Query(`lot`)=%+v, want []DocID{6}", docs)
	}

	idx.Delete("foobar", 1)
	docs = idx.Query("fooba")
	if len(docs) != 0 {
		t.Errorf("Query(`fooba`)=%+v, want []DocID{}", docs)
	}
}

func TestFullPrune(t *testing.T) {

	s := []string{
		"foo",
		"foobar",
		"foobfoo",
		"quxzoot",
		"zotzot",
		"azotfoba",
	}

	idx := NewIndex(s)
	idx.Prune(0)

	tests := []struct {
		q   string
		ids []DocID
	}{
		{"", []DocID{0, 1, 2, 3, 4, 5}},
		{"foo", []DocID{0, 1, 2, 3, 4, 5}},
		{"foob", []DocID{0, 1, 2, 3, 4, 5}},
		{"zot", []DocID{0, 1, 2, 3, 4, 5}},
		{"oba", []DocID{0, 1, 2, 3, 4, 5}},
	}

	for _, tt := range tests {
		if got := idx.Query(tt.q); !reflect.DeepEqual(got, tt.ids) {
			t.Errorf("Query(%q)=%+v, want %+v", tt.q, got, tt.ids)
		}
	}

	idx.Add("ahafoo")
	tests = []struct {
		q   string
		ids []DocID
	}{
		{"", []DocID{0, 1, 2, 3, 4, 5, 6}},
		{"foo", []DocID{0, 1, 2, 3, 4, 5, 6}},
		{"foob", []DocID{0, 1, 2, 3, 4, 5, 6}},
		{"zot", []DocID{0, 1, 2, 3, 4, 5, 6}},
		{"oba", []DocID{0, 1, 2, 3, 4, 5, 6}},
	}

	for _, tt := range tests {
		if got := idx.Query(tt.q); !reflect.DeepEqual(got, tt.ids) {
			t.Errorf("Query(%q)=%+v, want %+v", tt.q, got, tt.ids)
		}
	}
}

var result int
var podNames = getPodNames()
var f1 = getFileNames(100, 5000)
var f2 = getFileNames(100, 100000)

var idx1 = NewIndex(f1)
var idx2 = NewIndex(f2)


func BenchmarkQuery1_1(b *testing.B) {
	var r []DocID
	format := "general.tuning.%s.metric-*"
	q := fmt.Sprintf(format, podNames[rand.Intn(100)])

	for n := 0; n < b.N; n++ {
		ts := extractTrigrams(q)
		r = idx1.QueryTrigrams(ts)
	}
	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = len(r)
}

func BenchmarkQuery1_2(b *testing.B) {
	var r []DocID
	format := "general.tuning.%s.metric-1"
	q := fmt.Sprintf(format, "*")

	for n := 0; n < b.N; n++ {
		ts := extractTrigrams(q)
		r = idx1.QueryTrigrams(ts)
	}
	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = len(r)
}

func BenchmarkQuery1_3(b *testing.B) {
	var r []DocID
	format := "general.tuning.%s.metric-1*"
	q := fmt.Sprintf(format, "*")

	for n := 0; n < b.N; n++ {
		ts := extractTrigrams(q)
		r = idx1.QueryTrigrams(ts)
	}
	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = len(r)
}

func BenchmarkQuery2_1(b *testing.B) {
	var r []DocID
	format := "general.tuning.%s.metric-*"
	q := fmt.Sprintf(format, podNames[rand.Intn(100)])

	for n := 0; n < b.N; n++ {
		ts := extractTrigrams(q)
		r = idx2.QueryTrigrams(ts)
	}
	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = len(r)
}

func BenchmarkQuery2_2(b *testing.B) {
	var r []DocID
	format := "general.tuning.%s.metric-1"
	q := fmt.Sprintf(format, "*")

	for n := 0; n < b.N; n++ {
		ts := extractTrigrams(q)
		r = idx2.QueryTrigrams(ts)
	}
	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = len(r)
}

func BenchmarkQuery2_3(b *testing.B) {
	var r []DocID
	format := "general.tuning.%s.metric-1*"
	q := fmt.Sprintf(format, "*")

	for n := 0; n < b.N; n++ {
		ts := extractTrigrams(q)
		r = idx2.QueryTrigrams(ts)
	}
	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = len(r)
}


func getFileNames(pods int, perPodMetrics int) []string {
	var fileNames []string
	format := "general.tuning.%s.metric-%d"


	for i := 0; i < pods; i++ {
		podName := podNames[i]
		for j := 0; j < perPodMetrics; j++{
			x := fmt.Sprintf(format, podName, i)
			fileNames = append(fileNames, x)
		}

	}
	return fileNames
}

func getPodNames() [100]string {
	var pNames [100]string
	for i := 0; i < 100; i++ {
		pNames[i] = RandStringRunes(8)
	}
	return pNames
}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func extractTrigrams(query string) []T {

	if len(query) < 3 {
		return nil
	}

	var start int
	var i int

	var trigrams []T

	for i < len(query) {
		if query[i] == '[' || query[i] == '*' || query[i] == '?' {
			trigrams = Extract(query[start:i], trigrams)

			if query[i] == '[' {
				for i < len(query) && query[i] != ']' {
					i++
				}
			}

			start = i + 1
		}
		i++
	}

	if start < i {
		trigrams = Extract(query[start:i], trigrams)
	}

	return trigrams
}

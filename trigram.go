// Package trigram is a dumb trigram index
package trigram

// T is a trigram
type T uint32

func (t T) String() string {
	b := [3]byte{byte(t >> 16), byte(t >> 8), byte(t)}
	return string(b[:])
}

// DocID is a document ID
type DocID uint32

// Index is a trigram index
type Index map[T][]DocID

// a special (and invalid) trigram that holds all the document IDs
const tAllDocIDs T = 0xFFFFFFFF

// Extract returns a list of trigrams in s
func Extract(s string, trigrams []T) []T {

	for i := 0; i <= len(s)-3; i++ {
		t := T(uint32(s[i])<<16 | uint32(s[i+1])<<8 | uint32(s[i+2]))
		trigrams = appendIfUnique(trigrams, t)
	}

	return trigrams
}

func appendIfUnique(t []T, n T) []T {
	for _, v := range t {
		if v == n {
			return t
		}
	}

	return append(t, n)
}

// NewIndex returns an index for the strings in docs
func NewIndex(docs []string) Index {

	idx := make(Index)

	var allDocIDs []DocID

	var trigrams []T

	for id, d := range docs {
		ts := Extract(d, trigrams)
		docid := DocID(id)
		allDocIDs = append(allDocIDs, docid)
		for _, t := range ts {
			idx[t] = append(idx[t], docid)
		}
		trigrams = trigrams[:0]
	}

	idx[tAllDocIDs] = allDocIDs

	return idx
}

// Add adds a new string to the search index
func (idx Index) Add(s string) {

	id := DocID(len(idx))

	ts := Extract(s, nil)
	for _, t := range ts {
		idx[t] = append(idx[t], id)
	}

	idx[tAllDocIDs] = append(idx[tAllDocIDs], id)
}

// Query returns a list of document IDs that match the trigrams in the query s
func (idx Index) Query(s string) []DocID {
	ts := Extract(s, nil)
	return idx.QueryTrigrams(ts)
}

// QueryTrigrams returns a list of document IDs that match the trigram set ts
func (idx Index) QueryTrigrams(ts []T) []DocID {

	if len(ts) == 0 {
		return idx[tAllDocIDs]
	}

	midx := 0
	mtri := ts[midx]

	for i, t := range ts {
		if len(idx[t]) < len(idx[mtri]) {
			midx = i
			mtri = t
		}
	}

	ts[0], ts[midx] = ts[midx], ts[0]

	return idx.Filter(idx[mtri], ts[1:]...)
}

// Filter removes documents that don't contain the specified trigrams
func (idx Index) Filter(docs []DocID, ts ...T) []DocID {
	for _, t := range ts {
		docs = intersect(docs, idx[t])
	}

	return docs
}

func intersect(a, b []DocID) []DocID {

	// TODO(dgryski): reduce allocations by reusing A

	var aidx, bidx int

	var result []DocID

scan:
	for aidx < len(a) && bidx < len(b) {
		if a[aidx] == b[bidx] {
			result = append(result, a[aidx])
			aidx++
			bidx++
			if aidx >= len(a) || bidx >= len(b) {
				break scan
			}
		}

		for a[aidx] < b[bidx] {
			aidx++
			if aidx >= len(a) {
				break scan
			}
		}

		for bidx < len(b) && a[aidx] > b[bidx] {
			bidx++
			if bidx >= len(b) {
				break scan
			}
		}
	}

	return result
}

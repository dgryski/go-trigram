/* Package trigram is a dumb trigram index */
package trigram

type tri uint32

type Index map[tri][]int

func MakeTrigrams(s string) []tri {

	var trigrams []tri

	for i := 0; i <= len(s)-3; i++ {
		t := tri(uint32(s[i])<<16 | uint32(s[i+1])<<8 | uint32(s[i+2]))
		trigrams = appendIfUnique(trigrams, t)
	}

	return trigrams
}

func appendIfUnique(t []tri, n tri) []tri {
	for _, v := range t {
		if v == n {
			return t
		}
	}

	return append(t, n)
}

func NewIndex(docs []string) Index {

	idx := make(Index)

	for id, d := range docs {
		ts := MakeTrigrams(d)
		for _, t := range ts {
			idx[t] = append(idx[t], id)
		}
	}

	return idx
}

func (idx Index) Add(s string) {

	id := len(idx)

	ts := MakeTrigrams(s)
	for _, t := range ts {
		idx[t] = append(idx[t], id)
	}
}

func (idx Index) Query(s string) []int {
	ts := MakeTrigrams(s)
	return idx.QueryTrigrams(ts)
}

func (idx Index) QueryTrigrams(ts []tri) []int {

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

func (idx Index) Filter(docs []int, ts ...tri) []int {
	for _, t := range ts {
		docs = intersect(docs, idx[t])
	}

	return docs
}

func intersect(a, b []int) []int {

	// TODO(dgryski): reduce allocations by reusing A

	var aidx, bidx int

	var result []int

	for aidx < len(a) && bidx < len(b) {
		switch {
		case a[aidx] == b[bidx]:
			result = append(result, a[aidx])
			aidx++
			bidx++
		case a[aidx] < b[bidx]:
			aidx++
		case a[aidx] > b[bidx]:
			bidx++
		}
	}

	return result
}

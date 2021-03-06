// Copyright (c) 2018-2019 Burak Sezer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sorted

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkSortedMapSet(b *testing.B) {
	m := NewSortedMap(0)
	defer m.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err := m.Set(bkey(n), bval(n))
		if err != nil {
			b.Fatalf("Expected nil. Got %v", err)
		}
	}
}

func BenchmarkSortedMapGet(b *testing.B) {
	m := NewSortedMap(0)
	defer m.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err := m.Set(bkey(n), bval(n))
		if err != nil {
			b.Fatalf("Expected nil. Got %v", err)
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := m.Get(bkey(n))
		if err != nil {
			b.Fatalf("Expected nil. Got %v", err)
		}
	}
}

func parallelKey(threadID int, counter int) string {
	return fmt.Sprintf("%04d-%09d", threadID, counter)
}

func BenchmarkSortedMapSetParallel(b *testing.B) {
	m := NewSortedMap(0)
	defer m.Close()

	rand.Seed(time.Now().Unix())
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000000)
		counter := 0
		for pb.Next() {
			key := []byte(parallelKey(id, counter))
			err := m.Set(key, key)
			if err != nil {
				b.Errorf("Expected nil. Got %s", err)
			}
			counter = counter + 1
		}
	})
}

func BenchmarkSortedMapGetParallel(b *testing.B) {
	m := NewSortedMap(0)
	defer m.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err := m.Set(bkey(n), bval(n))
		if err != nil {
			b.Fatalf("Expected nil. Got %v", err)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		n := 0
		for pb.Next() {
			_, err := m.Get(bkey(n))
			if err != nil {
				b.Errorf("Expected nil. Got %s", err)
			}
			n++
		}
	})
}

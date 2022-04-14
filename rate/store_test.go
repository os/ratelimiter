package rate

import (
	"reflect"
	"testing"
)

func TestMemoryStore_Increment(t *testing.T) {
	for _, tc := range []struct {
		Description    string
		IncrementKeys  []string
		ExpectedCounts map[string]int
	}{
		{
			Description: "Counters match increment calls",
			IncrementKeys: []string{
				"user1",
				"user2",
				"user3",
				"user4",
				"user5",
				"user1",
				"user2",
				"user3",
				"user1",
			},
			ExpectedCounts: map[string]int{
				"user1": 3,
				"user2": 2,
				"user3": 2,
				"user4": 1,
				"user5": 1,
			},
		},
	} {
		t.Run(tc.Description, func(t *testing.T) {
			store := MemoryStore{
				store: map[string]*memoryRecord{},
			}

			counts := map[string]int{}
			for _, key := range tc.IncrementKeys {
				value, err := store.Increment(key, 10)
				if err != nil {
					t.Errorf("Expected no errors, got: %s", err)
				}
				counts[key] = value
			}

			if !reflect.DeepEqual(counts, tc.ExpectedCounts) {
				t.Errorf("Counts don't match")
			}
		})
	}
}

package pokecache

import
(	
	"testing"
	"time"
	"fmt"
)

func TestAddGet(t *testing.T) {	
	const interval = 5 * time.Second
	type test struct {
		key string
		val []byte
	}

	tests := []test{
		{key: "Toontown", val: []byte("toontown")},
		{key: "Route66", val: []byte("cerulean")},
	}
	fmt.Println("Running Add/Get.")
	for i, testCase := range tests {
		t.Run( fmt.Sprintf("Test case %v", i), func (t *testing.T) {
			cache := NewCache(interval)
			cache.Add(testCase.key, testCase.val)
			val, ok := cache.Get(testCase.key)

			if !ok {
				t.Errorf("expected to find key")
				return
			}

			if string(val) != string(testCase.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
	fmt.Println("Add/Get finished.")
	return
}

func TestReapLoop(t *testing.T) {
	fmt.Println("Running Reaploop")
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + time.Second
	cache := NewCache(baseTime)
	cache.Add( "Toontown", []byte("toontown") )

	_, ok := cache.Get("Toontown")
	if !ok {
		t.Errorf("Expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("Toontown")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
	fmt.Println("Reaploop finished")
	return 
}
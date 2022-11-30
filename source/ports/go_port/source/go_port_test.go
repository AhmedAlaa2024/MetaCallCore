package metacall

import (
	"log"
	"os"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	if err := Initialize(); err != nil {
		log.Fatal(err)
	}

	// if benchmark {
		buffer := "module.exports = { benchmark: async x => x }"

		if err := LoadFromMemory("node", buffer); err != nil {
			log.Fatal(err)
			return
		}
	// }

	code := m.Run()
	Destroy()
	os.Exit(code)
}

func TestMock(t *testing.T) {
	scripts := []string{"test.mock"}

	if err := LoadFromFile("mock", scripts); err != nil {
		t.Fatal(err)
		return
	}

	ret, err := Call("three_str", "e", "f", "g")

	if err != nil {
		t.Fatal(err)
		return
	}

	str, ok := ret.(string)

	if !ok || str != "Hello World" {
		t.Fatalf("expected 'Hello World', received %s", str)
	}
}

func TestNodeJSArray(t *testing.T) {
	buffer := "module.exports = { g: () => [0, 1, 2] }"

	if err := LoadFromMemory("node", buffer); err != nil {
		t.Fatal(err)
		return
	}

	array, err := Call("g")

	if err != nil {
		t.Fatal(err)
		return
	}

	iArray := array.([]interface{})

	if len(iArray) != 3 {
		t.Fatal("Invalid size of array")
		return
	}

	for i := 0; i < len(iArray); i++ {
		switch iArray[i].(type) {
		case float64:
			{
				if iArray[i] != float64(i) {
					t.Fatalf("Invalid value of array at position %d: %f (current) != %f (expected)", i, iArray[i], float64(i))
					return
				}
			}
		default:
			{
				t.Fatalf("Invalid type of array at position %d", i)
				return
			}
		}
	}
}

func TestNodeJSAwait(t *testing.T) {
	buffer := "module.exports = { f: async (ms) => await new Promise(resolve => setTimeout(resolve, ms)) }"

	if err := LoadFromMemory("node", buffer); err != nil {
		t.Fatal(err)
		return
	}

	var wg sync.WaitGroup

	wg.Add(1)

	_, err := Await("f",
		func(interface{}, interface{}) interface{} {
			log.Println("from go resolve")
			wg.Done()
			return nil
		},
		func(interface{}, interface{}) interface{} {
			log.Println("from go reject")
			wg.Done()
			return nil
		},
		nil,
		300, // 0.3 seconds
	)

	if err != nil {
		t.Fatal(err)
		return
	}

	wg.Wait()
}

func benchmarkNodeJS(b *testing.B, n int) {
	var wg sync.WaitGroup

	wg.Add(n)

	for i := 0; i < n; i++ {
		_, err := Await("benchmark",
			func(interface{}, interface{}) interface{} {
				wg.Done()
				return nil
			},
			func(interface{}, interface{}) interface{} {
				wg.Done()
				return nil
			},
			nil,
		)

		if err != nil {
			b.Fatal(err)
			return
		}
	}

	wg.Wait()
}

func BenchmarkNodeJSSequential(b *testing.B) {
	benchmarkNodeJS(b, 5)
}

func BenchmarkNodeJSParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchmarkNodeJS(b, 5)
		}
	})
}

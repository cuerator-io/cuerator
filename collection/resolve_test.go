package collection_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cuerator-io/cuerator/collection"
	"github.com/google/go-cmp/cmp"
)

func TestResolve(t *testing.T) {
	t.Parallel()

	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range entries {
		if entry.Name()[0] == '_' {
			continue
		}

		dir := filepath.Join("testdata", entry.Name())

		t.Run(dir, func(t *testing.T) {
			t.Parallel()

			instances, err := collection.Load(dir)
			if err != nil {
				t.Fatal(err)
			}

			inputs := loadInputs(dir)
			outputs, err := collection.Resolve(inputs, instances)
			if err != nil {
				t.Fatal(err)
			}

			got := actualOutputsAsJSON(outputs)
			want := expectedOutputsAsJSON(dir)

			if diff := cmp.Diff(want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func loadInputs(dir string) map[string]any {
	filename := filepath.Join(dir, "inputs.json")

	data, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		panic(err)
	}

	var inputs map[string]any
	if err := json.Unmarshal(data, &inputs); err != nil {
		panic(err)
	}

	return inputs
}

func expectedOutputsAsJSON(dir string) map[string]string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	want := map[string]string{}

	for _, entry := range entries {
		var name string
		if entry.Name() == "outputs.json" {
			name = "_"
		} else if n, ok := strings.CutSuffix(entry.Name(), ".outputs.json"); ok {
			name = n
		} else {
			continue
		}

		filename := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		want[name] = normalizeJSON(data)
	}

	return want
}

func actualOutputsAsJSON(outputs map[string]map[string]any) map[string]string {
	got := map[string]string{}

	for name, output := range outputs {
		data, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			panic(err)
		}

		got[name] = string(data)
	}

	return got
}

func normalizeJSON(data []byte) string {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		panic(err)
	}

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(data)
}

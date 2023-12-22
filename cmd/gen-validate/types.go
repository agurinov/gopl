package main

type (
	TargetField struct {
		tags map[string]string
		name string
		typ  string
	}
	Target struct {
		imports  map[string]string
		name     string
		generics []string
		fields   []TargetField
	}
)

func (t Target) InvolvedInGenerate() bool {
	involved := false

	for _, f := range t.fields {
		if f.tags == nil {
			continue
		}

		if f.tags["validate"] != "" {
			involved = true
		}
	}

	return involved
}

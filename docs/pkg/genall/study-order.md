## Study Order For `pkg/genall` And Related Packages

This guide explains a recommended study order for the `*_study_test.go` files added around `pkg/genall` and its supporting packages.

### Goal

Follow the code in the same order that controller-tools uses it:

`version -> yaml helpers -> markers -> loader -> genall -> help -> pretty`

This keeps the learning path moving from small pure functions to the higher-level generation and documentation pipeline.

### Recommended Order

1. `pkg/version/version_study_test.go`
2. `pkg/schemapatcher/internal/yaml/yaml_study_test.go`
3. `pkg/markers/markers_study_test.go`
4. `pkg/loader/loader_study_test.go`
5. `pkg/genall/genall_study_test.go`
6. `pkg/genall/help/help_study_test.go`
7. `pkg/genall/help/help_flow_study_test.go`
8. `pkg/genall/help/pretty/pretty_study_test.go`

### Why This Order Works

1. `version`
   Start with a tiny package and learn the study-test style.
2. `yaml`
   Practice reading and mutating tree-shaped data.
3. `markers`
   Learn the marker definition and parsing model.
4. `loader`
   Learn the AST, struct tag, and error helper layer that supports marker collection.
5. `genall`
   See how generation options and output rules are orchestrated.
6. `help`
   See how marker definitions become documentation data.
7. `help_flow`
   Follow the end-to-end path from definition input to grouped help output.
8. `pretty`
   Finish with the terminal rendering rules and exact output strings.

### Study Loop

Use the same loop for each test file.

1. Read only the test names.
2. Predict the expected behavior.
3. Read the test body and confirm the input and expected output.
4. Open the implementation file and find the code path.
5. Explain in your own words why the assertion is correct.
6. Run only that study test.
7. Change one small value, rerun, and observe what breaks.

### Useful Commands

Run all study tests:

```bash
go test ./pkg/... -run Study -v
```

Run package-by-package:

```bash
go test ./pkg/version -run Study -v
go test ./pkg/schemapatcher/internal/yaml -run Study -v
go test ./pkg/markers -run Study -v
go test ./pkg/loader -run Study -v
go test ./pkg/genall -run Study -v
go test ./pkg/genall/help -run Study -v
go test ./pkg/genall/help/pretty -run Study -v
```

### What To Focus On

For each package, answer these questions:

1. What contract is the test protecting?
2. Is this package doing data transformation or orchestration?
3. What data shape enters this package?
4. What data shape leaves this package?
5. Which lower-level package does it rely on?

### Big Picture

The most useful mental model is this pipeline:

1. `loader` reads packages and AST information.
2. `markers` defines and parses marker comments.
3. `genall` coordinates generators and output behavior.
4. `help` converts marker definitions into documentation models.
5. `pretty` renders those models into terminal-friendly text.

The `help_flow` and `pretty` study tests are the best place to confirm that this pipeline makes sense end-to-end.

### Suggested Passes

1. First pass
   Read tests first, then implementations.
2. Second pass
   Read implementations first, then explain the tests.
3. Third pass
   Make a tiny change, predict the failure, run the test, and revert.

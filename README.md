# Trak Registry 📦

Official template registry and catalog for **Trak**.

## Structure

```text
trak-registry/
├── registry.json             # Small index catalog with categories and metadata
└── templates/
    ├── lang/
    │   ├── go.json          # Go workspace blueprint
    │   └── python.json
    └── tool/
        ├── jenkins.json     # Jenkins workspace blueprint
        └── docker.json
```

## How Trak CLI Resolves Templates

When a user runs:
```bash
trak init lang/go
```

1. Trak CLI fetches `registry.json`.
2. Locates category `lang` and template `go`.
3. Reads the `source` field (`templates/lang/go.json`).
4. Fetches the template blueprint and recursively builds the workspace locally.

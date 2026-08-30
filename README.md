# Trak Registry 📦

> Central template catalog and workspace blueprints for the **[Trak CLI](https://github.com/ndk123-web/trak)** learning tool.

[![Schema Version](https://img.shields.io/badge/Schema-v1.0.0-blue.svg)](registry.json)
[![Templates](https://img.shields.io/badge/Templates-Go%20%7C%20Jenkins-orange.svg)](templates/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 🗂️ Registry Architecture

The Trak registry acts as a lightweight, distributed index for learning curricula. It separates **catalog discovery** from **workspace blueprints**:

```text
trak-registry/
├── registry.json                 # Fast index catalog with categories and metadata
└── templates/
    ├── lang/
    │   ├── go.json              # Go workspace blueprint (tree structure)
    │   └── python.json
    └── tool/
        ├── jenkins.json         # Jenkins workspace blueprint
        └── docker.json
```

---

## 📄 Schema Contracts

### 1. `registry.json` (Catalog Index)
Provides categorized listings and points to the raw source of each template.

```json
{
  "schema_version": "1.0.0",
  "updated_at": "2026-08-30T12:00:00Z",
  "categories": {
    "lang": {
      "title": "Programming Languages",
      "description": "Interactive learning workspaces for programming languages",
      "templates": {
        "go": {
          "name": "Go",
          "description": "Comprehensive Go fundamentals, concurrency, and hands-on exercises",
          "version": "1.0.0",
          "source": "templates/lang/go.json",
          "tags": ["golang", "backend", "concurrency"]
        }
      }
    },
    "tool": {
      "title": "DevOps & Developer Tools",
      "description": "Hands-on environments and labs for DevOps & engineering tools",
      "templates": {
        "jenkins": {
          "name": "Jenkins",
          "description": "Jenkins declarative pipelines, automation jobs, and CI/CD labs",
          "version": "1.0.0",
          "source": "templates/tool/jenkins.json",
          "tags": ["cicd", "devops", "automation"]
        }
      }
    }
  }
}
```

---

### 2. `templates/<category>/<name>.json` (Workspace Blueprint)
Defines the hierarchical directory and file tree that Trak CLI materializes locally.

```json
{
  "id": "lang/go",
  "name": "Go",
  "version": "1.0.0",
  "description": "Go fundamentals and exercises",
  "root": {
    "name": "go-workspace",
    "type": "directory",
    "children": [
      {
        "name": "01-hello-world",
        "type": "directory",
        "children": [
          {
            "name": "main.go",
            "type": "file",
            "content": "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, Trak!\")\n}\n"
          },
          {
            "name": "README.md",
            "type": "file",
            "content": "# 01 - Hello World\n"
          }
        ]
      }
    ]
  }
}
```

---

## 🛠️ How to Contribute a Template

We welcome new languages, frameworks, and tools!

1. **Fork** this repository.
2. **Create a template JSON** under `templates/<category>/<template_name>.json`.
3. **Register your template** in `registry.json` under the appropriate category.
4. **Submit a Pull Request**.

---

## 🔗 Live Raw URLs

- **Registry Index:** `https://raw.githubusercontent.com/ndk123-web/trak-registry/main/registry.json`
- **Go Blueprint:** `https://raw.githubusercontent.com/ndk123-web/trak-registry/main/templates/lang/go.json`
- **Jenkins Blueprint:** `https://raw.githubusercontent.com/ndk123-web/trak-registry/main/templates/tool/jenkins.json`

---

## 📜 License

MIT License. See [LICENSE](LICENSE) for details.

<p align="center">
  <img src="public/trak.png" width="120" alt="Trak Registry Logo" style="border-radius: 20px;" />
</p>

<h1 align="center">Trak Registry</h1>

<p align="center">
  <strong>Decoupled GitOps Curriculum Catalog & Blueprint Engine for Trak CLI</strong>
</p>

<p align="center">
  <a href="registry.json"><img src="https://img.shields.io/badge/Schema-v1.1.0-blue?style=flat-square" alt="Schema Version" /></a>
  <a href="templates/"><img src="https://img.shields.io/badge/Official-19%20Tracks%20(350%2B%20Modules)-emerald?style=flat-square" alt="Official Tracks" /></a>
  <a href="users/"><img src="https://img.shields.io/badge/Community-GitOps%20Namespaces-cyan?style=flat-square" alt="Community Namespaces" /></a>
  <a href="https://github.com/ndk123-web/trak"><img src="https://img.shields.io/badge/CLI-v1.1.0-00ADD8?style=flat-square&logo=go" alt="CLI" /></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-green.svg?style=flat-square" alt="License" /></a>
</p>

---

## ⚡ Overview

**Trak Registry** is the central, open-source curriculum repository for the [Trak CLI](https://github.com/ndk123-web/trak) ecosystem. It hosts declarative Abstract Syntax Tree (AST) blueprints that define comprehensive, multi-module learning environments with real source code, build configs, and hands-on exercises.

By decoupling curriculum blueprints from the CLI binary, new learning tracks and community updates are **streamed instantly** to users globally without requiring CLI recompilation, database migrations, or server maintenance.

---

## 🎬 Demo

<p align="center">
  <video src="https://github.com/user-attachments/assets/4210baaf-ef0d-469b-9a8a-f0e244d9b9a3" controls="controls" width="100%" style="max-width: 900px; border-radius: 12px;"></video>
</p>

---

## 🗂️ 2-Tier Registry Hierarchy

Trak Registry separates verified core curricula from decentralized community tracks using two clear filesystem tiers:

```text
trak-registry/
├── registry.json                 # Master catalog index for official tracks
├── templates/                    # 📦 Tier 1: Official Verified Curriculums (Maintained by Core Team)
│   ├── lang/                     # Programming Languages (go, rust, python, typescript...)
│   ├── os/                       # Operating Systems (linux, macos, windows)
│   ├── cloud/                    # Cloud Platforms (aws)
│   ├── db/                       # Databases & Storage (postgres, redis, sql)
│   └── tool/                     # DevOps Tools (docker, k8s, terraform, git, jenkins, ansible)
│
└── users/                        # 🌐 Tier 2: Community GitOps Blueprints (Open to All Creators)
    ├── [github-username-1]/      # Isolated Creator Namespace (e.g. ndk123-web/)
    │   └── lang/
    │       └── go.json
    └── [github-username-2]/
        └── db/
            └── postgres.json
```

---

## 🔄 End-to-End GitOps Contribution & Verification Flow

Publishing custom learning tracks to Trak Registry requires **zero custom accounts, zero passwords, and zero API keys**. The entire lifecycle is governed by automated GitOps and GitHub Actions CI:

```mermaid
flowchart TD
    A["1. Creator: Design Blueprint in Studio"] -->|"Export AST JSON"| B["2. Fork trak-registry on GitHub"]
    
    B -->|"Save file to users/:username/:category/:tool.json"| C["3. Submit Pull Request (PR)"]
    
    C --> D["4. GitHub Actions CI Trigger: validate.yml"]
    
    subgraph CI_Checks ["Automated CI AST & Security Validator (scripts/validate.go)"]
        E1["1. Actor Isolation: GITHUB_ACTOR == folderUsername"]
        E2["2. File Size Limit: Max 5MB per JSON blueprint"]
        E3["3. Category Taxonomy: lang, os, cloud, db, tool"]
        E4["4. AST Structure: root.type == directory & non-empty children"]
        E5["5. Path Traversal Guard: Reject .., /, :, null bytes in names"]
        E6["6. Security Check: Reject .exe, .dll, .so, .dylib binaries"]
    end
    
    D --> CI_Checks
    
    CI_Checks -->|"Any Check Fails"| F["❌ PR Rejected with Detailed Security/Schema Error"]
    CI_Checks -->|"All Checks Pass"| G["✔ PR Approved & Merged to main"]
    
    G --> H["5. Instantly Available Worldwide!"]
    H --> I["Developer Runs: trak init :username/:category/:tool"]
```

---

## 🛡️ Comprehensive Validator Rules (`scripts/validate.go`)

Every blueprint submitted via Pull Request is strictly validated across 6 automated checks before merging:

| Rule | Enforcement | Details |
| :--- | :--- | :--- |
| **1. Actor Isolation** | Mandatory | In PRs, `GITHUB_ACTOR` must match the parent folder in `users/<username>/`. No contributor can touch another creator's namespace or `templates/`. |
| **2. Max File Size (5MB)** | `< 5,242,880 bytes` | Blueprint JSON files cannot exceed 5MB to ensure fast HTTPS downloads over CLI. |
| **3. Category Taxonomy** | Strict Enum | Must reside under one of the 5 valid categories: `lang`, `os`, `cloud`, `db`, `tool`. |
| **4. Metadata Schema** | Mandatory Fields | `id`, `name`, `version`, and `root` directory object must be present and non-empty. |
| **5. Path Traversal Safety** | Security Sanitization | Node names cannot contain `/`, `\`, `:`, null bytes, `.`, or `..` to prevent filesystem escapes. |
| **6. Binary Prohibition** | Security Shield | Rejects compiled executable extensions (`.exe`, `.dll`, `.so`, `.dylib`, `.bin`) inside templates. |

---

## 🛠️ Step-by-Step Community Publishing Guide

### Step 1: Design in Blueprint Studio
Instead of writing hundreds of lines of nested JSON manually with escaped newlines and quotes, use **Blueprint Studio** directly in your browser:
1. Open [Blueprint Studio](https://trak-web.vercel.app/studio) (or run `trak-web` locally).
2. Scaffold your directories, add starter files, and write code in Monaco Editor.
3. Click **"Download AST JSON"** to export your valid template.

### Step 2: Fork the Repository
Fork [`github.com/ndk123-web/trak-registry`](https://github.com/ndk123-web/trak-registry) to your own GitHub account.

### Step 3: Place in Your Isolated User Namespace
Create your blueprint file using the exact directory structure:
```text
users/<your-github-username>/<category>/<tool>.json
```

> **Allowed Categories:** `lang` (languages), `os` (operating systems), `cloud` (cloud infra), `db` (databases), `tool` (DevOps tools).

### Step 4: Submit Pull Request (PR)
1. Push your branch and open a PR against `main`.
2. GitHub Actions will automatically run `go run scripts/validate.go`.
3. Once merged, anyone across the globe can immediately run:
   ```bash
   trak init <your-github-username>/<category>/<tool>
   ```

---

## 📄 Blueprint AST Schema Specification

Every blueprint JSON file represents a recursive Abstract Syntax Tree (AST) defining the materialized workspace:

```json
{
  "id": "lang/go",
  "name": "Go (Golang) Comprehensive Mastery Track",
  "version": "1.2.0",
  "description": "Complete Go curriculum from fundamentals to production concurrency",
  "root": {
    "name": "go-workspace",
    "type": "directory",
    "children": [
      {
        "name": "go.mod",
        "type": "file",
        "content": "module learn-go\n\ngo 1.22\n"
      },
      {
        "name": "00-setup-and-prerequisites",
        "type": "directory",
        "children": [
          {
            "name": "README.md",
            "type": "file",
            "content": "# 00 - Setup & Toolchain\n\n## 🎯 Learning Objectives\n..."
          },
          {
            "name": "main.go",
            "type": "file",
            "content": "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, Trak!\")\n}\n"
          }
        ]
      }
    ]
  }
}
```

### Node Schema Fields:
| Field | Type | Description |
| :--- | :--- | :--- |
| `name` | `string` | File or directory name (e.g. `main.go`, `01-basics`). Cannot contain path traversal characters. |
| `type` | `string` | Node type: `"directory"` or `"file"`. |
| `content` | `string` | *(Files only)* Raw UTF-8 string content with escaped newlines. |
| `children` | `array[Node]` | *(Directories only)* Nested list of child file and directory nodes. |

---

## 📚 Master Blueprint Matrix (19 Official Tracks)

| Category | Identifier | Name | Modules | Version | Source File |
| :--- | :--- | :--- | :---: | :---: | :--- |
| **`lang/`** | `lang/go` | Go (Golang) | 20 | `1.2.0` | [`templates/lang/go.json`](templates/lang/go.json) |
| **`lang/`** | `lang/rust` | Rust Systems | 22 | `1.0.0` | [`templates/lang/rust.json`](templates/lang/rust.json) |
| **`lang/`** | `lang/typescript` | TypeScript Fullstack | 18 | `1.0.0` | [`templates/lang/typescript.json`](templates/lang/typescript.json) |
| **`lang/`** | `lang/python` | Python Architecture | 20 | `1.0.0` | [`templates/lang/python.json`](templates/lang/python.json) |
| **`lang/`** | `lang/javascript` | JavaScript & Node.js | 18 | `1.0.0` | [`templates/lang/javascript.json`](templates/lang/javascript.json) |
| **`lang/`** | `lang/java` | Java & JVM Systems | 19 | `1.0.0` | [`templates/lang/java.json`](templates/lang/java.json) |
| **`lang/`** | `lang/cpp` | Modern C++ (C++23) | 19 | `1.0.0` | [`templates/lang/cpp.json`](templates/lang/cpp.json) |
| **`lang/`** | `lang/c` | C Low-Level Systems | 18 | `1.0.0` | [`templates/lang/c.json`](templates/lang/c.json) |
| **`os/`** | `os/linux` | Linux Systems & Bash | 21 | `1.0.0` | [`templates/os/linux.json`](templates/os/linux.json) |
| **`os/`** | `os/macos` | macOS & Darwin XNU | 16 | `1.0.0` | [`templates/os/macos.json`](templates/os/macos.json) |
| **`os/`** | `os/windows` | Windows & PowerShell | 18 | `1.0.0` | [`templates/os/windows.json`](templates/os/windows.json) |
| **`cloud/`** | `cloud/aws` | AWS Cloud Architecture | 22 | `1.0.0` | [`templates/cloud/aws.json`](templates/cloud/aws.json) |
| **`db/`** | `db/postgres` | PostgreSQL & DBA | 20 | `1.0.0` | [`templates/db/postgres.json`](templates/db/postgres.json) |
| **`db/`** | `db/redis` | Redis In-Memory Engine | 18 | `1.0.0` | [`templates/db/redis.json`](templates/db/redis.json) |
| **`db/`** | `db/sql` | Comprehensive SQL | 18 | `1.0.0` | [`templates/db/sql.json`](templates/db/sql.json) |
| **`tool/`** | `tool/docker` | Docker Containers | 20 | `1.0.0` | [`templates/tool/docker.json`](templates/tool/docker.json) |
| **`tool/`** | `tool/k8s` | Kubernetes (CKA/CKAD) | 22 | `1.0.0` | [`templates/tool/k8s.json`](templates/tool/k8s.json) |
| **`tool/`** | `tool/terraform` | Terraform Infrastructure | 18 | `1.0.0` | [`templates/tool/terraform.json`](templates/tool/terraform.json) |
| **`tool/`** | `tool/ansible` | Ansible Automation | 18 | `1.0.0` | [`templates/tool/ansible.json`](templates/tool/ansible.json) |
| **`tool/`** | `tool/git` | Git Internals & Workflows | 16 | `1.0.0` | [`templates/tool/git.json`](templates/tool/git.json) |
| **`tool/`** | `tool/jenkins` | Jenkins CI/CD Pipelines | 18 | `1.0.0` | [`templates/tool/jenkins.json`](templates/tool/jenkins.json) |

---

## 🧪 Local Schema Validation

You can run the official registry validator locally before opening a PR:

```bash
# From trak-registry repository root
go run scripts/validate.go
```

---

## 📜 Ecosystem & License

- ⚡ **[Trak CLI](https://github.com/ndk123-web/trak)** — The official Go binary workspace generator.
- 🌐 **[Trak Web](https://github.com/ndk123-web/trak-web)** — The interactive web portal and Blueprint Studio.

This project is licensed under the **MIT License**.

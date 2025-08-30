# compoundassign

A Go static analysis tool that suggests using compound assignment operators (`+=`, `-=`, etc.) instead of verbose assignments like `x = x + y`.

---

## 📌 Overview

Writing `x = x + y` works, but it’s more idiomatic and concise in Go to use `x += y`.  
This analyzer scans your Go code and reports such patterns, suggesting the replacement with a compound assignment.  

---

## ✨ Features

- Detects verbose self-assignments like `x = x + y`
- Suggests and auto-fixes them to compound form (`x += y`)
- Supports identifiers and selectors (e.g., `obj.field = obj.field + y` → `obj.field += y`)
- Provides suggested fixes that can be applied automatically by supporting tools

---

## 🔧 Supported Operators

The analyzer converts the following:

| Operator | Replacement |
|----------|-------------|
| `+`      | `+=`        |
| `-`      | `-=`        |
| `*`      | `*=`        |
| `/`      | `/=`        |
| `%`      | `%=`        |
| `&`      | `&=`        |
| `|`      | `|=`        |
| `^`      | `^=`        |
| `<<`     | `<<=`       |
| `>>`     | `>>=`       |
| `&^`     | `&^=`       |

---

## 🚀 Installation

You can install the package with:

```bash
go install github.com/ajinbaby01/compoundassign@latest

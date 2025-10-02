# README.md Generation Rules for Directory Navigation

## Version 1.0.0

## Purpose Statement

This ruleset defines standardized rules for creating README.md files throughout a project structure. Each README.md serves as a directory guide, providing developers with immediate context about the contents and purpose of files and subdirectories within the current location.

## Core Principles

1. **Navigation First**: README files are wayfinding tools, not documentation
2. **Brevity is Essential**: Maximum 1 page when rendered
3. **Descriptive, Not Technical**: Explain "what is here" not "how it works"
4. **Consistency Across Project**: Every directory follows identical format
5. **Hierarchical Awareness**: Each README knows its place in the structure
6. **LLM Optimized**: Clear patterns for automated generation and parsing

## README.md Scope and Purpose

### What README.md IS:
- A directory map for developers
- A quick orientation tool
- A file/folder purpose descriptor
- A navigation aid

### What README.md IS NOT:
- Technical documentation
- API reference
- Implementation details
- Tutorial or guide
- Changelog
- Contributing guidelines

## Required Structure Template

```markdown
# [Directory Name]

## Overview
[One sentence describing this directory's purpose within the parent context]

## Contents

### Files
| File | Purpose |
|------|---------|
| [filename.ext] | [One-line description of what this file contains/does] |
| [filename.ext] | [One-line description of what this file contains/does] |

### Directories
| Directory | Purpose |
|-----------|---------|
| [dirname/] | [One-line general description of directory contents] |
| [dirname/] | [One-line general description of directory contents] |

## Navigation
- **Parent**: [../](../) - [Parent directory purpose]
- **Related**: [Optional - Link to related directories at same level]
```

## Section Rules

### Title Section
- Must match the directory name exactly
- Use proper case formatting
- No additional decorations or emojis

### Overview Section
- EXACTLY one sentence
- Maximum 100 characters
- Describes purpose within parent context
- No technical details
- No implementation specifics

### Contents Section

#### Files Subsection
- List ONLY files in current directory
- Do NOT list files in subdirectories
- One-line descriptions only
- Maximum 80 characters per description
- Sort alphabetically by filename
- Include file extensions
- Skip if no files present

#### Directories Subsection
- List ONLY immediate child directories
- Add trailing slash (/) to directory names
- General description only - no file specifics
- Maximum 80 characters per description
- Sort alphabetically
- Skip if no subdirectories present

### Navigation Section
- ALWAYS include parent directory link
- Parent link uses relative path (../)
- Include one-line parent purpose description
- Related directories are optional
- Use relative paths for all links

## File Description Guidelines

### Description Language Patterns

#### For Source Code Files
```markdown
| main.py | Application entry point |
| utils.py | Utility functions collection |
| config.py | Configuration settings |
```

#### For Configuration Files
```markdown
| .env | Environment variables |
| settings.json | Application settings |
| docker-compose.yml | Container orchestration setup |
```

#### For Documentation Files
```markdown
| README.md | Directory navigation guide |
| LICENSE | Project license terms |
| CHANGELOG.md | Version history record |
```

#### For Asset Files
```markdown
| logo.png | Project logo image |
| styles.css | Stylesheet definitions |
| data.csv | Sample data set |
```

### Directory Description Guidelines

#### DO Use Generic Descriptions
```markdown
| src/ | Source code files |
| tests/ | Test suites and fixtures |
| docs/ | Documentation files |
| assets/ | Static resources and media |
```

#### DON'T Include Specifics
❌ WRONG:
```markdown
| src/ | Contains main.py, utils.py, and 5 model files |
```

✅ CORRECT:
```markdown
| src/ | Application source code |
```

## Formatting Standards

### Table Formatting
- Use pipe tables with headers
- Include separator row (|------|---------|)
- Align pipes for readability
- No empty cells - use "—" if needed

### Path Formatting
- Files: include extension (file.ext)
- Directories: include trailing slash (dir/)
- Links: use relative paths
- Parent: always ../

### Text Formatting
- No bold/italic in descriptions
- No code blocks
- No inline code markers
- No links in descriptions

## Special Cases

### Empty Directories
```markdown
## Contents

*This directory is currently empty or contains only subdirectories.*

### Directories
| Directory | Purpose |
|-----------|---------|
| [subdirs if any] | [descriptions] |
```

### Single Purpose Directories
```markdown
## Contents

*This directory contains [type] files exclusively.*

### Files
[Table of files]
```

### Root Directory README
```markdown
# [Project Name]

## Overview
[One sentence project description]

## Contents

### Files
| File | Purpose |
|------|---------|
| [files in root] | [descriptions] |

### Directories
| Directory | Purpose |
|-----------|---------|
| [root directories] | [descriptions] |

## Quick Start
See [docs/](docs/) for detailed documentation.
```

## Prohibited Content

### Never Include:
1. Code examples
2. Installation instructions  
3. API documentation
4. Implementation details
5. TODO items
6. Version numbers
7. Dates or timestamps
8. Author information
9. Links to external resources
10. Markdown beyond tables and headers

### Never Describe:
1. How something works
2. Why something exists
3. When to use something
4. Who should use something
5. Technical specifications

## Directory Hierarchy Awareness

### Depth Level Indicators

#### Level 0 (Root)
```markdown
# Project Name
## Overview
Main project repository
```

#### Level 1
```markdown
# src
## Overview
Source code for the application
```

#### Level 2+
```markdown
# components
## Overview
Reusable UI components
```

### Path Context
Each README should be aware of its location:
- Root: Project level description
- Level 1: Major section description
- Level 2+: Specific functionality description

## File Filtering Rules

### Always List:
- Source code files
- Configuration files
- Documentation files
- Asset files
- Data files
- Script files

### Never List:
- Hidden files starting with . (except .env, .gitignore)
- Temporary files
- Cache files
- Build artifacts
- Auto-generated files
- IDE configuration files

### Optional Listings:
- Lock files (package-lock.json, etc.)
- Compiled files (only if checked in)
- Vendor directories (summarize as single entry)

## Examples

### Example 1: Source Code Directory

```markdown
# src

## Overview
Application source code and core business logic

## Contents

### Files
| File | Purpose |
|------|---------|
| app.js | Application initialization |
| server.js | Server configuration |
| routes.js | Request routing definitions |

### Directories
| Directory | Purpose |
|-----------|---------|
| controllers/ | Request handlers |
| models/ | Data models |
| utils/ | Helper functions |
| middleware/ | Request processing pipeline |

## Navigation
- **Parent**: [../](../) - Project root
```

### Example 2: Test Directory

```markdown
# tests

## Overview
Test suites and testing utilities

## Contents

### Files
| File | Purpose |
|------|---------|
| setup.js | Test environment configuration |
| helpers.js | Testing utility functions |

### Directories
| Directory | Purpose |
|-----------|---------|
| unit/ | Unit test cases |
| integration/ | Integration test suites |
| fixtures/ | Test data and mocks |
| e2e/ | End-to-end test scenarios |

## Navigation
- **Parent**: [../](../) - Project root
- **Related**: [../src/](../src/) - Source code being tested
```

### Example 3: Configuration Directory

```markdown
# config

## Overview
Application and environment configuration files

## Contents

### Files
| File | Purpose |
|------|---------|
| default.json | Default configuration values |
| production.json | Production environment settings |
| development.json | Development environment settings |
| database.json | Database connection settings |

### Directories
| Directory | Purpose |
|-----------|---------|
| env/ | Environment-specific configurations |
| keys/ | Security keys and certificates |

## Navigation
- **Parent**: [../](../) - Project root
```

## LLM Generation Instructions

### Generation Process
```
FUNCTION generate_readme(directory_path):
  1. SCAN directory for files and subdirectories
  2. FILTER according to filtering rules
  3. SORT alphabetically
  4. GENERATE overview from directory name and parent context
  5. CREATE files table if files exist
  6. CREATE directories table if subdirectories exist
  7. ADD navigation section with parent link
  8. VALIDATE against ruleset
  9. OUTPUT formatted markdown
```

### Context Requirements
When generating README.md, LLM needs:
1. Current directory name
2. Parent directory name and purpose
3. List of files in current directory
4. List of immediate subdirectories
5. Project-level context (if available)

### Validation Checklist
- [ ] Title matches directory name
- [ ] Overview is one sentence
- [ ] Files listed are only from current directory
- [ ] Directories listed are only immediate children
- [ ] All descriptions are under 80 characters
- [ ] Tables are properly formatted
- [ ] Navigation section includes parent
- [ ] No technical details included
- [ ] No prohibited content present
- [ ] Total length under 1 page

## Common Patterns

### Backend Directories
```
| api/ | REST API endpoints |
| services/ | Business logic services |
| database/ | Database schemas and migrations |
| middleware/ | Request/response processors |
```

### Frontend Directories
```
| components/ | UI components |
| pages/ | Page components |
| styles/ | CSS and styling files |
| assets/ | Images and static files |
```

### DevOps Directories
```
| scripts/ | Automation scripts |
| docker/ | Container definitions |
| k8s/ | Kubernetes manifests |
| ci/ | CI/CD pipeline configurations |
```

### Documentation Directories
```
| docs/ | Project documentation |
| api-docs/ | API documentation |
| guides/ | User and developer guides |
| examples/ | Usage examples |
```

## Quality Metrics

### Good README Indicators
- Can understand directory purpose in 5 seconds
- Can locate any file/directory instantly
- No need to open files to understand structure
- Clear navigation path
- Consistent with sibling directories

### Poor README Indicators
- Contains implementation details
- Longer than one screen
- Includes code examples
- Duplicates documentation
- Missing or incorrect navigation

## Maintenance Rules

### When to Update README.md
- New file added to directory
- File removed from directory
- New subdirectory created
- Subdirectory removed
- Directory purpose changes
- File purpose changes significantly

### When NOT to Update README.md
- File content changes
- Implementation changes
- Bug fixes
- Refactoring
- Documentation updates
- Version updates

## Anti-Patterns to Avoid

### 1. Documentation Duplication
❌ Don't copy from actual documentation
✅ Just point to where things are

### 2. Technical Specifications
❌ Don't explain how code works
✅ Just say what the file contains

### 3. Nested Descriptions
❌ Don't describe subdirectory contents
✅ Only describe immediate children

### 4. Tutorial Content
❌ Don't include how-to instructions
✅ Just describe what is present

### 5. Historical Information
❌ Don't include changelog or history
✅ Only current state

## Template for Quick Use

```markdown
# [DIRECTORY_NAME]

## Overview
[ONE_SENTENCE_PURPOSE]

## Contents

### Files
| File | Purpose |
|------|---------|
| [FILE_NAME] | [ONE_LINE_DESCRIPTION] |

### Directories
| Directory | Purpose |
|-----------|---------|
| [DIR_NAME/] | [ONE_LINE_DESCRIPTION] |

## Navigation
- **Parent**: [../](../) - [PARENT_PURPOSE]
```

## Final Notes

### Remember:
- README.md is a map, not a manual
- Every directory deserves a README.md
- Consistency across all directories
- Brevity over completeness
- Navigation over documentation

### Success Criteria:
A developer should be able to:
1. Understand what's in a directory without opening files
2. Navigate to related directories easily
3. Locate specific functionality quickly
4. Orient themselves in the project structure
5. All in under 30 seconds

---

*END OF README.md GENERATION RULES v1.0.0*

*Use these rules to ensure consistent, helpful navigation throughout your project structure.*
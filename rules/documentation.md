# LLM Documentation Rules for Application Development - Language Agnostic

## Core Principles

1. **Consistency First**: All documentation must follow identical structure patterns across all components
2. **LLM Optimization**: Use clear hierarchies, explicit markers, and structured formatting for optimal LLM parsing
3. **Three-Layer Perspective**: Every documentation must include macro (project), meso (module), and micro (component) views
4. **Progressive Detail**: Information depth increases from README.md → DEV.md → API.md → TEST.md
5. **Language Agnostic**: Rules apply to any programming language or framework

## Documentation File Structure

Every component MUST have exactly four documentation files in `component/docs/`:
```
component/
├── docs/
│   ├── README.md  - Quick start and overview for users
│   ├── DEV.md     - Architecture and implementation for developers
│   ├── API.md     - Complete interface documentation
│   └── TEST.md    - Testing strategies and execution
└── src/
```

## File Purpose Matrix

| File | Primary Audience | Focus | Technical Depth | Use When |
|------|-----------------|-------|-----------------|----------|
| README.md | New Users | What & Why | Low | First contact with component |
| DEV.md | Developers | How it works | High | Modifying/extending component |
| API.md | Integrators | How to use | Medium | Integrating component |
| TEST.md | QA/DevOps | How to verify | Medium-High | Testing/debugging component |

---

## README.md Rules

### Purpose
**First-contact documentation providing immediate understanding and quick start capabilities.**

### Target Audience
- New users exploring the component
- Stakeholders evaluating adoption
- Developers browsing for solutions

### Required Sections Structure

```markdown
# [Component Name]

## What is This?
[MACRO_VIEW]
One-sentence description of component's role in the entire project.
[/MACRO_VIEW]

[MESO_VIEW]
How this component relates to its parent module and sibling components.
[/MESO_VIEW]

[MICRO_VIEW]
The specific problem this component solves.
[/MICRO_VIEW]

## Why Use This?
### Problems It Solves
- Problem 1: [Description]
- Problem 2: [Description]

### Key Benefits
- Benefit 1: [User-focused benefit]
- Benefit 2: [User-focused benefit]

## Quick Start
### Prerequisites
- Prerequisite 1 with minimum version
- Prerequisite 2 with minimum version

### Installation
```bash
[Platform-agnostic installation command]
```

### Basic Usage
```
[Simplest possible working example with comments]
Input: [Example input]
Output: [Example output]
```

## Features
| Feature | Description | Status |
|---------|-------------|--------|
| Feature 1 | What it does for the user | Stable |
| Feature 2 | What it does for the user | Beta |

## Component Structure
```
component/
├── docs/           # Documentation
├── src/            # Source code
├── tests/          # Test files
├── examples/       # Usage examples
└── config/         # Configuration files
```

## Next Steps
- [Explore the API](./API.md) - Detailed usage instructions
- [Developer Guide](./DEV.md) - Understanding the internals
- [Testing Guide](./TEST.md) - Running and writing tests
- [Examples](../examples/) - More usage scenarios

## Support
- Issue Tracker: [Link]
- Discussion Forum: [Link]
- Contact: [Method]

## License
[License type] - See LICENSE file for details
```

### README.md Unique Characteristics
1. **Language**: Simple, jargon-free, benefit-focused
2. **Examples**: Minimal, immediately runnable
3. **Depth**: Surface-level, linking to details
4. **Tone**: Welcoming and encouraging
5. **Length**: Maximum 2 pages when rendered

---

## DEV.md Rules

### Purpose
**Technical deep-dive into architecture, design decisions, and implementation details.**

### Target Audience
- Developers maintaining the component
- Contributors adding features
- Architects evaluating design

### Required Sections Structure

```markdown
# [Component Name] - Developer Documentation

## Architecture Overview
[MACRO_VIEW]
Component's architectural pattern and its role in system architecture.
[/MACRO_VIEW]

[MESO_VIEW]
Integration patterns with module architecture and data flow.
[/MESO_VIEW]

[MICRO_VIEW]
Internal component architecture and processing pipeline.
[/MICRO_VIEW]

## Design Decisions

### Architectural Pattern
**Pattern Used**: [Pattern name]
**Justification**: [Why this pattern was chosen]
**Trade-offs**: 
- Pros: [List]
- Cons: [List]

### Core Abstractions
| Abstraction | Purpose | Design Principle |
|-------------|---------|------------------|
| [Name] | [What it abstracts] | [SOLID principle, etc.] |

## Component Internals

### Directory Structure Deep Dive
```
src/
├── core/           # [Purpose and responsibility]
│   ├── [file1]     # [Specific responsibility]
│   └── [file2]     # [Specific responsibility]
├── interfaces/     # [Purpose and responsibility]
├── implementations/# [Purpose and responsibility]
└── utilities/      # [Purpose and responsibility]
```

### Core Components

#### Component: [Name]
##### Responsibility
[Single responsibility description]

##### Collaborators
- [Component A]: [How they interact]
- [Component B]: [How they interact]

##### Key Algorithms
| Algorithm | Complexity | Use Case |
|-----------|------------|----------|
| [Name] | Time: O(n), Space: O(1) | [When used] |

##### State Management
```
[State diagram or description]
Initial State → Processing → Final State
```

### Data Flow Architecture
```
[Input Source] 
    ↓ [Transformation A]
[Intermediate State]
    ↓ [Transformation B]
[Output Sink]
```

### Dependency Graph
```
[ASCII or text representation of dependencies]
Component A
├── depends on → Component B
└── depends on → Component C
    └── depends on → Component D
```

## Extension Points

### How to Add New Features
1. **Identify Extension Point**
   - [Where to look]
   - [What to consider]

2. **Implement Interface/Contract**
   - [Required methods/properties]
   - [Constraints to respect]

3. **Register Component**
   - [How to register]
   - [Configuration needed]

### Plugin Architecture
[If applicable, describe how plugins work]

## Performance Characteristics

### Resource Usage
| Resource | Typical Usage | Maximum Usage | Scaling Factor |
|----------|--------------|---------------|----------------|
| Memory | [Amount] | [Amount] | [O(n) notation] |
| CPU | [Percentage] | [Percentage] | [O(n) notation] |
| I/O | [Operations/sec] | [Operations/sec] | [Pattern] |

### Optimization Strategies
1. **Strategy Name**
   - Implementation: [How it's done]
   - Impact: [Measured improvement]
   - Trade-off: [What we sacrifice]

## Security Considerations

### Threat Model
| Threat | Mitigation | Residual Risk |
|--------|------------|---------------|
| [Threat type] | [How we handle it] | [What remains] |

### Security Boundaries
```
[Trusted Zone] | [Security Boundary] | [Untrusted Zone]
```

## Development Workflow

### Setting Up Development Environment
```bash
# Step 1: [Description]
[command]

# Step 2: [Description]
[command]

# Step 3: Verify setup
[verification command]
```

### Code Organization Principles
- **Separation of Concerns**: [How applied]
- **Dependency Injection**: [How implemented]
- **Error Handling**: [Strategy used]

### Debugging Techniques
| Scenario | Technique | Tools |
|----------|-----------|-------|
| [Problem type] | [Approach] | [Tool names] |

## Monitoring and Observability

### Key Metrics
| Metric | Purpose | Alert Threshold |
|--------|---------|-----------------|
| [Metric name] | [What it indicates] | [Value] |

### Logging Strategy
- **Debug Level**: [What's logged]
- **Info Level**: [What's logged]
- **Error Level**: [What's logged]

### Debugging Hooks
[How to enable verbose debugging]

## Maintenance Guidelines

### Code Health Metrics
- Cyclomatic Complexity: Maximum [number]
- Coupling: Maximum [metric]
- Cohesion: Minimum [metric]

### Refactoring Triggers
1. [Condition] → [Recommended action]
2. [Condition] → [Recommended action]

## Migration Guide

### Breaking Changes Policy
[How breaking changes are handled]

### Version Compatibility Matrix
| Component Version | Compatible With | Migration Required |
|-------------------|-----------------|-------------------|
| 2.x | Module 3.x | No |
| 1.x | Module 2.x | Yes - see guide |

## Troubleshooting Development Issues

### Common Problems
| Symptom | Likely Cause | Solution |
|---------|--------------|----------|
| [What goes wrong] | [Why it happens] | [How to fix] |

## Contributing

### Code Review Checklist
- [ ] Follows architectural patterns
- [ ] Maintains abstraction boundaries
- [ ] Includes performance impact analysis
- [ ] Updates relevant documentation
```

### DEV.md Unique Characteristics
1. **Language**: Technical, precise, detailed
2. **Focus**: How and why, not what
3. **Diagrams**: ASCII architecture diagrams
4. **Metrics**: Quantifiable measures
5. **Decisions**: Documented trade-offs

---

## API.md Rules

### Purpose
**Complete interface documentation for component integration and usage.**

### Target Audience
- Developers integrating the component
- API consumers
- Technical writers

### Required Sections Structure

```markdown
# [Component Name] - API Documentation

## API Overview
[MACRO_VIEW]
Component's API philosophy and its place in the project's API ecosystem.
[/MACRO_VIEW]

[MESO_VIEW]
How this API relates to module-level APIs and contracts.
[/MESO_VIEW]

[MICRO_VIEW]
The specific capabilities exposed through this API.
[/MICRO_VIEW]

## API Principles
- **Principle 1**: [e.g., Consistency - All methods follow verb-noun pattern]
- **Principle 2**: [e.g., Idempotency - Safe retry semantics]
- **Principle 3**: [e.g., Versioning - Backward compatibility guaranteed]

## Authentication & Authorization
### Authentication Methods
| Method | Use Case | Example |
|--------|----------|---------|
| [Method type] | [When to use] | [Code example] |

### Required Permissions
| Endpoint/Method | Permission Level | Scope |
|-----------------|------------------|--------|
| [Name] | [Level] | [Scope description] |

## Core API Reference

### [Category 1 - e.g., Data Operations]

#### Method/Endpoint: [Name]
**Purpose**: [One-line description]

**Signature**:
```
[Method signature with types]
```

**Parameters**:
| Parameter | Type | Required | Default | Description | Constraints |
|-----------|------|----------|---------|-------------|-------------|
| param1 | [type] | Yes/No | [value] | [description] | [min/max, regex, etc.] |

**Returns**:
```
[Return type structure]
{
  field1: [type], // [description]
  field2: [type]  // [description]
}
```

**Errors**:
| Error Code | Condition | Resolution |
|------------|-----------|------------|
| [Code] | [When thrown] | [How to fix] |

**Example**:
```
// Request
[Example request]

// Response - Success
[Example successful response]

// Response - Error
[Example error response]
```

**Notes**:
- [Important consideration]
- [Performance note]
- [Security note]

### [Category 2 - e.g., Configuration]
[Repeat structure for each category]

## Event API / Callbacks

### Event: [Event Name]
**Triggered When**: [Condition]

**Payload**:
```
{
  eventType: [string],
  timestamp: [format],
  data: {
    // Event-specific data
  }
}
```

**Subscribe Example**:
```
[How to subscribe to this event]
```

## Streaming API
[If applicable]

### Stream: [Stream Name]
**Data Format**: [Format description]
**Frequency**: [Update frequency]
**Buffer Size**: [Size limits]

## Batch Operations

### Batch Processing
**Maximum Batch Size**: [Number]
**Timeout**: [Duration]
**Retry Policy**: [Description]

**Example**:
```
[Batch operation example]
```

## Rate Limiting

| Endpoint Category | Requests/Minute | Burst Limit | Retry-After |
|-------------------|-----------------|-------------|-------------|
| [Category] | [Number] | [Number] | [Seconds] |

## Pagination

### Pagination Strategy
**Type**: [Cursor/Offset/Page-based]
**Default Page Size**: [Number]
**Maximum Page Size**: [Number]

**Example**:
```
[Pagination example with request/response]
```

## Filtering & Sorting

### Available Filters
| Field | Operators | Example |
|-------|-----------|---------|
| [Field name] | [=, >, <, contains, etc.] | [Example usage] |

### Sorting Options
| Field | Direction | Default |
|-------|-----------|---------|
| [Field name] | ASC/DESC | [Yes/No] |

## Webhooks
[If applicable]

### Webhook Configuration
**Endpoint Requirements**: [URL format, HTTPS, etc.]
**Retry Policy**: [Attempts and backoff]
**Security**: [Signature verification]

## API Versioning

### Version Strategy
**Current Version**: [Version]
**Supported Versions**: [List]
**Deprecation Policy**: [Timeline]

### Version Differences
| Feature | v1 | v2 | v3 |
|---------|----|----|-----|
| [Feature name] | [Support status] | [Support status] | [Support status] |

## SDK/Client Libraries

### Official Libraries
| Language | Package | Version | Documentation |
|----------|---------|---------|---------------|
| [Language] | [Package name] | [Version] | [Link] |

### Code Generation
```
[How to generate client code from API spec]
```

## Testing the API

### Test Environment
**Base URL**: [URL]
**Test Credentials**: [How to obtain]
**Limitations**: [Any differences from production]

### Postman/OpenAPI
**Collection**: [Link to collection]
**OpenAPI Spec**: [Link to spec]

### Example Test Flow
```
1. [Step 1 - Setup]
2. [Step 2 - Execute]
3. [Step 3 - Verify]
```

## Migration Guides

### Migrating from v1 to v2
#### Breaking Changes
- [Change 1]: [Migration path]
- [Change 2]: [Migration path]

#### Deprecated Features
| Feature | Deprecated In | Removed In | Alternative |
|---------|---------------|------------|-------------|
| [Feature] | [Version] | [Version] | [New approach] |

## API Best Practices

### DO
- [Best practice with example]
- [Best practice with example]

### DON'T
- [Anti-pattern with explanation]
- [Anti-pattern with explanation]

## Troubleshooting

### Common Integration Issues
| Issue | Symptoms | Solution |
|-------|----------|----------|
| [Issue name] | [What you'll see] | [How to fix] |

## API Metrics

### SLA
| Metric | Target | Measurement |
|--------|--------|-------------|
| Availability | [99.9%] | [How measured] |
| Response Time (p95) | [<200ms] | [How measured] |

## Changelog
### [Version] - [Date]
#### Added
- [New feature]

#### Changed
- [Modified behavior]

#### Deprecated
- [Feature to be removed]

#### Removed
- [Removed feature]
```

### API.md Unique Characteristics
1. **Language**: Technical but accessible
2. **Examples**: Every endpoint/method has examples
3. **Completeness**: All parameters, returns, errors documented
4. **Practical**: Includes testing and troubleshooting
5. **Versioned**: Clear version strategy and migration paths

---

## TEST.md Rules

### Purpose
**Comprehensive guide for testing strategy, execution, and quality assurance.**

### Target Audience
- QA Engineers
- Developers writing tests
- DevOps setting up CI/CD
- Anyone debugging issues

### Required Sections Structure

```markdown
# [Component Name] - Testing Documentation

## Testing Philosophy
[MACRO_VIEW]
How testing this component contributes to overall project quality goals.
[/MACRO_VIEW]

[MESO_VIEW]
Integration with module-level testing strategy and shared testing infrastructure.
[/MESO_VIEW]

[MICRO_VIEW]
Component-specific testing approach and quality metrics.
[/MICRO_VIEW]

## Testing Strategy

### Test Pyramid
```
        /\
       /E2E\      5%  - Critical user journeys
      /------\
     /Integration\ 25% - Component interactions  
    /------------\
   /     Unit     \ 70% - Individual functions
  /----------------\
```

### Testing Dimensions
| Dimension | Coverage Goal | Current Coverage | Priority |
|-----------|--------------|------------------|----------|
| Functional | 90% | [%] | High |
| Performance | Key paths | [Status] | Medium |
| Security | All inputs | [Status] | High |
| Accessibility | WCAG 2.1 AA | [Status] | Medium |

## Test Environment Setup

### Prerequisites
#### System Requirements
| Component | Minimum Version | Recommended | Notes |
|-----------|-----------------|-------------|-------|
| Runtime | [Version] | [Version] | [Special requirements] |
| Database | [Version] | [Version] | [Test instance needed] |
| Cache | [Version] | [Version] | [Optional for unit tests] |

#### Environment Variables
```bash
# Required
TEST_ENV=[value]         # Description
TEST_TIMEOUT=[value]     # Description

# Optional
TEST_VERBOSE=[value]     # Description
TEST_PARALLEL=[value]    # Description
```

### Installation Steps
```bash
# Step 1: Install test dependencies
[command]

# Step 2: Setup test database
[command]

# Step 3: Generate test fixtures
[command]

# Step 4: Verify installation
[command]
Expected output: [description]
```

## Test Organization

### Directory Structure
```
tests/
├── unit/              # Isolated component tests
│   ├── core/          # Core logic tests
│   └── utilities/     # Utility function tests
├── integration/       # Component interaction tests
│   ├── api/          # API integration tests
│   └── database/     # Database integration tests
├── e2e/              # End-to-end user journeys
│   └── scenarios/    # User scenario tests
├── performance/      # Performance benchmarks
├── security/         # Security test cases
├── fixtures/         # Test data
│   ├── valid/       # Valid test cases
│   └── invalid/     # Invalid test cases
├── mocks/           # Mock implementations
└── helpers/         # Test utilities
```

## Running Tests

### All Tests
```bash
# Standard run
[command to run all tests]

# With coverage report
[command with coverage]

# In watch mode
[command with watch]
```

### Test Suites by Type

#### Unit Tests
```bash
# Run all unit tests
[command]

# Run specific unit test file
[command with file path]

# Run with debugging
[command with debug flag]
```

#### Integration Tests
```bash
# Prerequisites
[setup command if needed]

# Run integration tests
[command]

# Cleanup after tests
[cleanup command]
```

#### End-to-End Tests
```bash
# Start test environment
[command to start services]

# Run E2E tests
[command]

# Run specific scenario
[command with scenario]
```

#### Performance Tests
```bash
# Run performance benchmarks
[command]

# Generate performance report
[command for report]
```

### Platform-Specific Instructions

#### Windows
```powershell
# PowerShell specific commands
[PowerShell commands]

# CMD specific adjustments
[CMD commands]

# Common issues and solutions
[Windows-specific troubleshooting]
```

#### macOS
```bash
# macOS specific setup
[macOS commands]

# Permission requirements
[Permission setup]
```

#### Linux
```bash
# Distribution-specific notes
[Linux commands]

# Container-based testing
[Docker commands if applicable]
```

#### CI/CD Environment
```yaml
# Example pipeline configuration
test:
  stage: test
  script:
    - [command 1]
    - [command 2]
```

## Test Data Management

### Fixtures
#### Structure
```
fixtures/
├── users/           # User test data
│   ├── valid.json   # Valid user data
│   └── invalid.json # Invalid user data
└── config/          # Configuration fixtures
```

#### Loading Fixtures
```
[Example of loading fixture data]
```

#### Generating Fixtures
```bash
# Generate new fixtures
[command to generate]

# Update existing fixtures
[command to update]
```

### Mocks and Stubs

#### Available Mocks
| Mock Name | Purpose | Configuration |
|-----------|---------|---------------|
| [Service Mock] | [What it mocks] | [How to configure] |

#### Creating Custom Mocks
```
[Example of creating a mock]
```

#### Mock Best Practices
- Always reset mocks between tests
- Use realistic mock data
- Document mock limitations

## Test Categories Explained

### Unit Tests
**Purpose**: Test individual components in isolation
**Characteristics**:
- No external dependencies
- Fast execution (<100ms per test)
- Deterministic results

**Example Structure**:
```
describe: [Component Name]
  context: [Specific scenario]
    it: [Expected behavior]
      // Arrange
      // Act  
      // Assert
```

### Integration Tests
**Purpose**: Test component interactions
**Characteristics**:
- May use real dependencies
- Slower than unit tests
- Test data flow between components

**Key Test Scenarios**:
| Scenario | Components Involved | Validation |
|----------|-------------------|------------|
| [Scenario name] | [Component A + B] | [What to verify] |

### E2E Tests
**Purpose**: Validate complete user workflows
**Characteristics**:
- Test from user perspective
- Slowest test type
- Catch integration issues

**Critical User Journeys**:
1. [Journey name]: [Start] → [Steps] → [Expected outcome]
2. [Journey name]: [Start] → [Steps] → [Expected outcome]

## Code Coverage

### Coverage Requirements
| Metric | Minimum | Target | Current |
|--------|---------|--------|---------|
| Line Coverage | 70% | 85% | [%] |
| Branch Coverage | 65% | 80% | [%] |
| Function Coverage | 75% | 90% | [%] |
| Statement Coverage | 70% | 85% | [%] |

### Generating Coverage Reports
```bash
# Generate HTML report
[command]

# Generate JSON report for CI
[command]

# View coverage report
[command to open report]
```

### Coverage Exclusions
```
// Patterns excluded from coverage
- Generated files: [pattern]
- Config files: [pattern]
- Test files: [pattern]
```

## Continuous Integration

### Test Pipeline
```
1. Lint/Format Check
   └── Unit Tests
       └── Integration Tests
           └── E2E Tests
               └── Performance Tests
                   └── Coverage Report
```

### CI Configuration
```yaml
[Example CI configuration]
```

### Test Reports
**Location**: [Where reports are stored]
**Format**: [HTML, JSON, XML]
**Retention**: [How long kept]

## Performance Testing

### Benchmarks
| Operation | Target | Current | Variance Allowed |
|-----------|--------|---------|------------------|
| [Operation name] | [Target metric] | [Current value] | [+/- %] |

### Load Testing
```bash
# Run load test
[command]

# Parameters
Users: [number]
Duration: [time]
Ramp-up: [time]
```

### Performance Regression Detection
**Threshold**: [% change that triggers alert]
**Baseline**: [How baseline is established]

## Security Testing

### Security Test Suite
| Test Type | Tool/Method | Frequency |
|-----------|-------------|-----------|
| SAST | [Tool name] | Every commit |
| DAST | [Tool name] | Weekly |
| Dependency Scan | [Tool name] | Daily |

### Security Test Scenarios
1. Input validation testing
2. Authentication bypass attempts
3. Authorization boundary testing
4. Injection attack testing

## Debugging Tests

### Debug Strategies
| Problem Type | Debug Approach | Tools |
|--------------|----------------|-------|
| Flaky tests | [Approach] | [Tools] |
| Slow tests | [Approach] | [Tools] |
| False failures | [Approach] | [Tools] |

### Interactive Debugging
```bash
# Run tests in debug mode
[command]

# Attach debugger
[instructions]
```

### Logging During Tests
```
# Enable verbose logging
[environment variable or flag]

# Log locations
Test logs: [path]
Application logs: [path]
```

## Test Maintenance

### Test Health Metrics
| Metric | Good | Warning | Critical |
|--------|------|---------|----------|
| Flaky Test Rate | <1% | 1-5% | >5% |
| Test Runtime | <5min | 5-15min | >15min |
| Test Failures | <2% | 2-5% | >5% |

### Test Review Checklist
- [ ] Test name clearly describes scenario
- [ ] Test is independent (no order dependencies)
- [ ] Test data is properly cleaned up
- [ ] Assertions are specific and meaningful
- [ ] Test covers edge cases
- [ ] Performance impact is acceptable

### Updating Tests
**When to update**:
- API changes
- Business logic changes
- Bug fixes (add regression test)
- Performance improvements

**Update process**:
1. Run existing tests to establish baseline
2. Make changes
3. Update affected tests
4. Verify all tests pass
5. Update documentation

## Troubleshooting

### Common Issues
| Issue | Platform | Symptoms | Solution |
|-------|----------|----------|----------|
| [Issue name] | All | [What happens] | [How to fix] |
| [Issue name] | Windows | [What happens] | [How to fix] |
| [Issue name] | macOS | [What happens] | [How to fix] |
| [Issue name] | Linux | [What happens] | [How to fix] |

### FAQ
**Q: Tests pass locally but fail in CI**
A: [Common causes and solutions]

**Q: How to run a single test?**
A: [Platform-specific commands]

**Q: Tests are running slowly**
A: [Performance optimization tips]

## Test Documentation

### Writing Test Descriptions
```
GOOD: "should return error when user ID is invalid"
BAD: "test error case"

GOOD: "should process 1000 items in under 2 seconds"  
BAD: "performance test"
```

### Test Comments
```
// GIVEN: [Initial state]
// WHEN: [Action taken]
// THEN: [Expected outcome]
```

## Resources
- [Testing Best Practices Guide]
- [Mock Data Generator]
- [Test Report Dashboard]
- [Performance Baseline History]
```

### TEST.md Unique Characteristics
1. **Language**: Procedural, command-focused
2. **Platform Coverage**: OS-specific instructions
3. **Metrics**: Quantifiable quality measures
4. **Troubleshooting**: Extensive problem-solving guides
5. **Maintenance**: Test health monitoring

---

## General Documentation Rules

### File Location
All documentation MUST be located in `component/docs/`:
```
component/
├── docs/
│   ├── README.md
│   ├── DEV.md
│   ├── API.md
│   └── TEST.md
```

### Perspective Markers
Every document MUST include three perspective levels:
- `[MACRO_VIEW]...[/MACRO_VIEW]` - Project-wide context
- `[MESO_VIEW]...[/MESO_VIEW]` - Module-level context  
- `[MICRO_VIEW]...[/MICRO_VIEW]` - Component-specific context

### Formatting Standards
1. **Headers**: Use ATX-style headers (#), never skip levels
2. **Code Blocks**: Always use triple backticks with language hint
3. **Tables**: Always include header row with separator
4. **Lists**: Use `-` for unordered, `1.` for ordered
5. **Emphasis**: Use **bold** for important terms, *italic* for emphasis
6. **Line Length**: Maximum 100 characters for readability

### Language-Agnostic Requirements
1. Never use language-specific syntax in rules
2. Use generic terms: "component", "module", "function", "method"
3. Provide structure examples without language syntax
4. Focus on concepts over implementations

### LLM Optimization Techniques
1. **Consistent Structure**: Same section order across all components
2. **Explicit Sections**: Clear, searchable section names
3. **Structured Data**: Tables over paragraphs for specifications
4. **Clear Hierarchy**: Proper nesting with consistent indentation
5. **Unique Identifiers**: Section names that are globally unique
6. **Context Markers**: [TAG] markers for context switching
7. **Cross-References**: Explicit links between related sections

### Documentation Distinctions

| Aspect | README | DEV | API | TEST |
|--------|---------|-----|-----|------|
| **Audience** | End users | Developers | Integrators | QA/DevOps |
| **Questions** | What? Why? | How does it work? | How to use? | How to verify? |
| **Depth** | Surface | Deep internals | Interface details | Execution details |
| **Code** | Minimal examples | Architecture code | Full API examples | Test examples |
| **Tone** | Friendly | Technical | Precise | Procedural |
| **Length** | 1-2 pages | 5-10 pages | 10+ pages | 5-8 pages |

### Version Control Rules
1. Documentation updates MUST accompany code changes
2. Documentation reviews are part of code review
3. Breaking changes require migration guides in API.md
4. Version documentation separately from code

### Quality Validation Checklist
- [ ] All four files present in `component/docs/`
- [ ] Perspective markers present and correctly tagged
- [ ] No language-specific examples in rules
- [ ] Tables properly formatted with headers
- [ ] Internal links are relative and working
- [ ] Code blocks have language hints
- [ ] Platform-specific sections in TEST.md
- [ ] Clear distinction between file purposes
- [ ] Examples are self-contained
- [ ] Sections follow prescribed order

### Documentation Maintenance Triggers
| Event | README | DEV | API | TEST |
|-------|---------|-----|-----|------|
| New feature | Update features | Update architecture | Add endpoints | Add test cases |
| Bug fix | If user-facing | If design change | If behavior change | Add regression test |
| Performance improvement | If user-visible | Update metrics | Update SLA | Update benchmarks |
| Security patch | Update if relevant | Update threat model | Update auth | Add security test |
| Dependency update | Update prerequisites | Update dependencies | If API affected | Update environment |
| Refactoring | No change | Update internals | No change | Update if needed |

## Documentation Anti-Patterns to Avoid

1. **Mixing Concerns**: Don't put API details in README
2. **Duplication**: Don't repeat content across files
3. **Language Bias**: Don't assume specific programming paradigms
4. **Missing Context**: Always include all three view levels
5. **Outdated Examples**: Examples must work with current version
6. **Unclear Audience**: Each file has ONE primary audience
7. **Wall of Text**: Use structure, tables, and formatting
8. **Hidden Knowledge**: Document all assumptions explicitly

## LLM Processing Instructions

When generating documentation following these rules:

1. **Start with template structure** for consistency
2. **Fill sections in order** to maintain flow
3. **Verify perspective tags** are present and properly closed
4. **Use language-agnostic terms** throughout
5. **Check table formatting** for consistency
6. **Validate cross-references** between documents
7. **Ensure examples are generic** and conceptual
8. **Apply the correct tone** for each document type

### Document Generation Order
1. Generate README.md first (defines what)
2. Generate API.md second (defines how to use)
3. Generate DEV.md third (explains how it works)
4. Generate TEST.md last (verifies everything)

### Validation Steps for LLMs
```
FOR each document:
  CHECK file location is component/docs/
  CHECK all required sections present
  CHECK [MACRO_VIEW], [MESO_VIEW], [MICRO_VIEW] tags exist
  CHECK no language-specific syntax in examples
  CHECK target audience alignment
  CHECK cross-references to other docs
  VALIDATE formatting standards
  VALIDATE table structures
END FOR
```

## Example Component Documentation Structure

### Complete File Tree
```
my-component/
├── docs/
│   ├── README.md          # 2 pages - User introduction
│   ├── DEV.md            # 8 pages - Technical deep-dive
│   ├── API.md            # 12 pages - Complete interface
│   └── TEST.md           # 6 pages - Testing guide
├── src/
│   ├── core/
│   ├── interfaces/
│   └── implementations/
├── tests/
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── examples/
│   ├── basic/
│   └── advanced/
└── config/
```

### Document Relationship Map
```
README.md (Entry Point)
    ↓ "Learn to use" → API.md
    ↓ "Learn internals" → DEV.md
    ↓ "Verify it works" → TEST.md
    
API.md ←→ DEV.md (Architecture informs API)
API.md ←→ TEST.md (Tests verify API contract)
DEV.md ←→ TEST.md (Architecture guides testing)
```

## Enforcement Rules for Consistency

### Mandatory Elements Per Document

#### README.md Must Have:
- [ ] "What is This?" section
- [ ] "Why Use This?" section with benefits
- [ ] Quick Start with immediate value
- [ ] Links to all other docs
- [ ] Simple directory structure
- [ ] License information

#### DEV.md Must Have:
- [ ] Architecture overview with diagrams
- [ ] Design decisions with trade-offs
- [ ] Directory deep-dive with purposes
- [ ] Performance characteristics
- [ ] Security considerations
- [ ] Development workflow

#### API.md Must Have:
- [ ] Complete parameter documentation
- [ ] Return value specifications
- [ ] Error codes and handling
- [ ] Examples for every endpoint/method
- [ ] Rate limiting information
- [ ] Version compatibility matrix

#### TEST.md Must Have:
- [ ] Test pyramid/strategy diagram
- [ ] Platform-specific instructions (Windows/macOS/Linux)
- [ ] Coverage requirements and current status
- [ ] Troubleshooting guide
- [ ] Performance benchmarks
- [ ] CI/CD configuration

### Prohibited Practices

1. **Never mix document purposes**
   - ❌ Don't put setup instructions in API.md
   - ❌ Don't put API details in README.md
   - ❌ Don't put architecture in TEST.md

2. **Never use language-specific examples in rules**
   - ❌ `public class MyClass`
   - ✅ `[Component initialization]`

3. **Never skip perspective levels**
   - ❌ Only micro-view
   - ✅ All three views with proper tags

4. **Never assume prior knowledge**
   - ❌ "As everyone knows..."
   - ✅ "This component handles..."

## Quick Reference Card for LLMs

### Document Purpose in One Line
- **README.md**: What it does and why you'd want it
- **DEV.md**: How it's built and why it's built that way
- **API.md**: How to integrate and use it completely
- **TEST.md**: How to verify it works correctly

### Tone Guide
- **README.md**: Friendly teacher explaining benefits
- **DEV.md**: Senior engineer explaining architecture
- **API.md**: Technical writer documenting contract
- **TEST.md**: QA lead explaining verification

### Audience Needs
- **README.md reader asks**: "Should I use this?"
- **DEV.md reader asks**: "How do I modify this?"
- **API.md reader asks**: "How do I call this?"
- **TEST.md reader asks**: "How do I verify this?"

## Final Validation Checklist for LLMs

Before completing documentation generation:

### Structure Validation
- [ ] All documents in `component/docs/` directory
- [ ] All four documents present (README, DEV, API, TEST)
- [ ] Each document has ALL required sections
- [ ] Sections appear in prescribed order

### Content Validation
- [ ] Three perspective views in each document
- [ ] Perspective tags properly opened and closed
- [ ] No language-specific syntax (Java, Python, etc.)
- [ ] Examples are conceptual, not implementation-specific
- [ ] Cross-references use relative paths

### Formatting Validation
- [ ] ATX headers (#) used consistently
- [ ] Tables have header separator rows
- [ ] Code blocks use triple backticks
- [ ] Lists use consistent markers
- [ ] Line length ≤ 100 characters

### Audience Alignment
- [ ] README speaks to new users
- [ ] DEV speaks to developers
- [ ] API speaks to integrators
- [ ] TEST speaks to QA/DevOps

### Completeness Check
- [ ] README has working quick start
- [ ] DEV has architecture diagrams
- [ ] API has examples for all endpoints
- [ ] TEST has OS-specific instructions

---

## Meta: About This Ruleset

**Version**: 1.0.0
**Purpose**: Ensure consistent, high-quality documentation across all components
**Optimization**: Designed for LLM processing and generation
**Agnostic**: Works with any programming language or framework
**Location**: This ruleset should be referenced by LLMs when generating component documentation

### When to Update This Ruleset
- New documentation needs discovered through usage
- New sections required by multiple components
- Clarity improvements based on confusion
- Additional LLM optimization techniques discovered

### Ruleset Principles
1. **Clarity** over brevity
2. **Consistency** over creativity
3. **Completeness** over conciseness
4. **Structure** over prose
5. **Examples** over explanations

---

*END OF RULESET - Use these rules for all component documentation generation to ensure consistency and quality across the application.*
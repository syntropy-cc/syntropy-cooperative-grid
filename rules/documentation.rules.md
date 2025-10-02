# LLM Documentation Rules for Application Development - Language Agnostic v2.0

## Core Principles

1. **Consistency First**: All documentation must follow identical structure patterns across all components
2. **LLM Optimization**: Use clear hierarchies, explicit markers, and structured formatting for optimal LLM parsing
3. **Three-Layer Perspective**: Every documentation must include macro (project), meso (module), and micro (component) views
4. **Progressive Detail**: Information depth increases from DOC.md → API.md → DEV.md → TEST.md → LEARN.md
5. **Language Agnostic**: Rules apply to any programming language or framework
6. **Knowledge Transfer**: Documentation serves both operational and educational purposes

## Documentation File Structure

Every component MUST have exactly five documentation files in `component/docs/`:
```
component/
├── docs/
│   ├── DOC.md     - Quick start and overview for users
│   ├── DEV.md     - Architecture and implementation for developers
│   ├── API.md     - Complete interface documentation
│   ├── TEST.md    - Testing strategies and execution
│   └── LEARN.md   - Educational journey and knowledge synthesis
└── src/
```

## File Purpose Matrix

| File | Primary Audience | Focus | Technical Depth | Use When |
|------|-----------------|-------|-----------------|----------|
| DOC.md | New Users | What & Why | Low | First contact with component |
| DEV.md | Developers | How it works | High | Modifying/extending component |
| API.md | Integrators | How to use | Medium | Integrating component |
| TEST.md | QA/DevOps | How to verify | Medium-High | Testing/debugging component |
| LEARN.md | Learners/Educators | Why decisions matter | Very High | Teaching/mastering concepts |

---

## DOC.md Rules

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
- [Learning Path](./LEARN.md) - Deep dive into concepts and theory
- [Examples](../examples/) - More usage scenarios

## Support
- Issue Tracker: [Link]
- Discussion Forum: [Link]
- Contact: [Method]

## License
[License type] - See LICENSE file for details
```

### DOC.md Unique Characteristics
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

## Further Learning
For theoretical foundations and pedagogical insights, see [LEARN.md](./LEARN.md).
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

## Educational Resources
For conceptual understanding and learning exercises, see [LEARN.md](./LEARN.md).
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

## Learning Resources
For test-driven development principles and testing philosophy, see [LEARN.md](./LEARN.md).
```

### TEST.md Unique Characteristics
1. **Language**: Procedural, command-focused
2. **Platform Coverage**: OS-specific instructions
3. **Metrics**: Quantifiable quality measures
4. **Troubleshooting**: Extensive problem-solving guides
5. **Maintenance**: Test health monitoring

---

## LEARN.md Rules

### Purpose
**Comprehensive educational documentation synthesizing theoretical foundations, practical learnings, and pedagogical strategies for deep component understanding.**

### Target Audience
- Developers seeking mastery
- Educators teaching the concepts
- Team members onboarding
- Researchers studying implementations
- Students learning software engineering

### Required Sections Structure

```markdown
# [Component Name] - Learning Documentation

## Learning Overview
[MACRO_VIEW]
The fundamental computer science concepts and software engineering principles this component embodies within the larger system.
[/MACRO_VIEW]

[MESO_VIEW]
How mastering this component contributes to understanding the module's architecture and design patterns.
[/MESO_VIEW]

[MICRO_VIEW]
Specific technical skills, patterns, and problem-solving techniques learned through this component.
[/MICRO_VIEW]

## Theoretical Foundations

### Computer Science Fundamentals
#### Core Concepts
| Concept | CS Domain | Theoretical Basis | Practical Application |
|---------|-----------|------------------|----------------------|
| [Concept] | [e.g., Algorithms] | [Academic foundation] | [How applied here] |

#### Mathematical Foundations
```
Relevant Mathematics:
- [Mathematical concept]: [Why needed]
- [Formula/Theorem]: [Application]
- [Proof technique]: [Where used]
```

#### Computational Complexity
| Operation | Time Complexity | Space Complexity | Trade-offs |
|-----------|-----------------|------------------|------------|
| [Operation] | O(n) | O(1) | [Discussion] |

### Software Engineering Principles

#### Design Patterns Applied
##### Pattern: [Name]
**Gang of Four Classification**: [Creational/Structural/Behavioral]
**Intent**: [Academic definition]
**Structure**:
```
[UML-style diagram in ASCII]
```
**Our Implementation**:
- Participants: [Components involved]
- Collaborations: [How they interact]
- Consequences: [Benefits and liabilities]
- Implementation Notes: [Specific considerations]

#### SOLID Principles Demonstration
| Principle | How Applied | Code Location | Learning Value |
|-----------|-------------|---------------|----------------|
| Single Responsibility | [Example] | [File/Line] | [What it teaches] |
| Open/Closed | [Example] | [File/Line] | [What it teaches] |
| Liskov Substitution | [Example] | [File/Line] | [What it teaches] |
| Interface Segregation | [Example] | [File/Line] | [What it teaches] |
| Dependency Inversion | [Example] | [File/Line] | [What it teaches] |

#### Architectural Patterns
**Pattern Used**: [e.g., MVC, Hexagonal, Event-Driven]
**Theoretical Background**: [Academic sources]
**Implementation Choices**: [Our specific decisions]
**Learning Objectives**:
1. Understand [concept]
2. Apply [technique]
3. Evaluate [trade-offs]

## Pedagogical Approach

### Learning Theory Applied
#### Bloom's Taxonomy Levels
| Level | Cognitive Process | Component Application | Assessment Method |
|-------|------------------|---------------------|------------------|
| Remember | Recall facts | [Component terms] | [Quiz questions] |
| Understand | Explain concepts | [Core concepts] | [Explanations] |
| Apply | Use in new situations | [Practical tasks] | [Exercises] |
| Analyze | Draw connections | [Pattern recognition] | [Comparisons] |
| Evaluate | Justify decisions | [Trade-off analysis] | [Design choices] |
| Create | Design new solutions | [Extensions] | [Projects] |

#### Constructivist Learning Path
```
1. Prior Knowledge Activation
   └── Connect to known concepts
2. Cognitive Conflict
   └── Challenge assumptions
3. Construction
   └── Build new understanding
4. Reflection
   └── Internalize learning
5. Application
   └── Transfer to new contexts
```

### Learning Objectives

#### Knowledge Objectives
After studying this component, learners will know:
- [ ] [Factual knowledge item]
- [ ] [Conceptual understanding]
- [ ] [Procedural knowledge]
- [ ] [Metacognitive awareness]

#### Skill Objectives
After practicing with this component, learners will be able to:
- [ ] [Technical skill]
- [ ] [Problem-solving ability]
- [ ] [Analytical capability]
- [ ] [Design skill]

#### Attitude Objectives
After mastering this component, learners will appreciate:
- [ ] [Engineering principle]
- [ ] [Best practice value]
- [ ] [Quality attribute importance]

## Development Journey

### Problem Evolution
#### Initial Problem Statement
**As Presented**: [Original problem description]
**Assumptions Made**: [What we initially believed]
**Constraints Identified**: [Known limitations]

#### Problem Reframing
**Discovered Complexity**: [What we learned]
**Hidden Requirements**: [What wasn't obvious]
**Real Problem**: [Actual challenge to solve]

### Solution Evolution

#### Iteration 1: Naive Approach
**Hypothesis**: [What we thought would work]
**Implementation**:
```
[Pseudocode/diagram of first attempt]
```
**Result**: [What happened]
**Learning**: [Key insight gained]
**Cognitive Bias Revealed**: [What assumption was wrong]

#### Iteration 2: Informed Attempt
**New Understanding**: [What we now knew]
**Revised Approach**:
```
[Pseudocode/diagram of second attempt]
```
**Improvement**: [What got better]
**Remaining Issues**: [What still didn't work]
**Design Pattern Recognized**: [Pattern that emerged]

#### Iteration 3: Mature Solution
**Synthesis**: [Combined learnings]
**Final Design**:
```
[Pseudocode/diagram of solution]
```
**Success Criteria Met**: [How we knew it worked]
**Trade-offs Accepted**: [What we consciously sacrificed]

### Decision Tree Documentation
```
Decision: [Technology/Pattern Choice]
├── Option A: [Choice]
│   ├── Pros: [Benefits]
│   ├── Cons: [Drawbacks]
│   └── Rejected because: [Reason]
├── Option B: [Choice]
│   ├── Pros: [Benefits]
│   ├── Cons: [Drawbacks]
│   └── Selected because: [Reason]
└── Option C: [Choice]
    ├── Pros: [Benefits]
    ├── Cons: [Drawbacks]
    └── Rejected because: [Reason]
```

## Learning Through Failure

### Failure Analysis
#### Failed Approach: [Name]
**What We Tried**: [Description]
**Why It Should Have Worked**: [Theory]
**How It Failed**: [Specific failure mode]
**Root Cause Analysis**:
```
Symptom
└── Why? → [Immediate cause]
    └── Why? → [Underlying cause]
        └── Why? → [Deeper cause]
            └── Why? → [Root cause]
                └── Why? → [System cause]
```
**Prevention Strategy**: [How to avoid]
**Generalizable Lesson**: [Broader application]

### Anti-Pattern Catalog
| Anti-Pattern | How We Used It | Why It Failed | Correct Pattern | Learning |
|--------------|----------------|---------------|-----------------|----------|
| [Name] | [Our implementation] | [Failure mode] | [Better approach] | [Insight] |

## Cognitive Models

### Mental Models
#### Model: [Component Abstraction]
**Simplified View**:
```
[ASCII diagram of mental model]
```
**Accurate View**:
```
[ASCII diagram of actual implementation]
```
**Common Misconceptions**:
- Misconception: [Wrong model]
  - Why it's appealing: [Reason]
  - Why it's wrong: [Explanation]
  - Correct understanding: [Right model]

### System Thinking
#### Component Interactions
```
System Boundary
├── Input Boundaries
│   ├── Valid ranges: [Specifications]
│   └── Edge cases: [Boundary conditions]
├── Processing Core
│   ├── Transformations: [What changes]
│   └── Invariants: [What stays constant]
└── Output Boundaries
    ├── Guarantees: [What we promise]
    └── Limitations: [What we don't promise]
```

#### Emergent Properties
| Property | Component Behavior | System Behavior | Emergence Explanation |
|----------|-------------------|-----------------|----------------------|
| [Property] | [Local behavior] | [Global behavior] | [How it emerges] |

## Knowledge Transfer Strategies

### For Self-Directed Learners
#### Week 1: Foundation Building
**Learning Goals**:
- Understand problem domain
- Identify key concepts
- Run basic examples

**Activities**:
1. Read [DOC.md](./DOC.md) - 30 min
2. Complete Exercise Set A - 2 hours
3. Implement toy version - 4 hours

**Success Criteria**:
- [ ] Can explain component purpose
- [ ] Can run all examples
- [ ] Can modify simple behaviors

#### Week 2: Deep Understanding
**Learning Goals**:
- Grasp architecture
- Understand design decisions
- Debug common issues

**Activities**:
1. Study [DEV.md](./DEV.md) - 1 hour
2. Trace execution flow - 3 hours
3. Complete Exercise Set B - 3 hours

**Success Criteria**:
- [ ] Can explain architecture choices
- [ ] Can debug typical problems
- [ ] Can predict behavior changes

#### Week 3: Mastery
**Learning Goals**:
- Internalize patterns
- Transfer knowledge
- Create extensions

**Activities**:
1. Complete all exercises - 4 hours
2. Build extension - 6 hours
3. Teach someone else - 2 hours

**Success Criteria**:
- [ ] Can design similar components
- [ ] Can optimize performance
- [ ] Can teach core concepts

### For Instructors

#### Course Integration
**Prerequisites Coverage**:
- Data Structures: [Which ones needed]
- Algorithms: [Which ones applied]
- Design Patterns: [Which ones used]

**Learning Outcomes Mapping**:
| Course Outcome | Component Contribution | Assessment |
|---------------|----------------------|------------|
| [Outcome] | [How this helps] | [How to test] |

#### Lesson Plans
##### Lesson 1: Problem Introduction
**Duration**: 50 minutes
**Objectives**:
- Motivate the problem
- Explore naive solutions
- Identify challenges

**Activities**:
1. (10 min) Problem presentation
2. (20 min) Group brainstorming
3. (15 min) Solution attempts
4. (5 min) Reflection

**Materials**:
- Slides: [Topics to cover]
- Handout: [Exercises]
- Code: [Starter template]

##### Lesson 2: Pattern Recognition
[Similar structure for additional lessons]

#### Common Student Difficulties
| Difficulty | Indicators | Intervention | Prevention |
|------------|-----------|--------------|------------|
| [Concept confusion] | [How to spot] | [How to help] | [How to avoid] |

## Exercises and Challenges

### Conceptual Exercises
#### Exercise 1: Predict Behavior
**Setup**: Given configuration X
**Task**: Predict output for input Y
**Concepts Tested**: [List concepts]
**Solution Approach**:
1. [Step 1]
2. [Step 2]
**Common Mistakes**: [What to watch for]

### Implementation Exercises
#### Exercise 2: Extend Functionality
**Current State**: [What exists]
**Goal**: Add feature X
**Constraints**: [Limitations]
**Hints**:
- Hint 1 (after 10 min): [Gentle nudge]
- Hint 2 (after 20 min): [More specific]
- Hint 3 (after 30 min): [Strong hint]
**Solution**: [Complete implementation]
**Learning Points**: [What this teaches]

### Design Exercises
#### Exercise 3: Redesign Component
**Challenge**: Optimize for [different requirement]
**Trade-offs to Consider**:
- Performance vs. Readability
- Flexibility vs. Simplicity
- Memory vs. Speed
**Evaluation Criteria**: [How to judge solutions]

### Research Exercises
#### Exercise 4: Literature Review
**Task**: Find three alternative approaches
**Deliverable**: Comparison table
**Questions to Answer**:
1. What are the theoretical foundations?
2. What are the practical trade-offs?
3. When would each be preferred?

## Assessment Rubrics

### Understanding Assessment
| Level | Indicator | Example Evidence |
|-------|-----------|-----------------|
| Novice | Can use component | Runs examples successfully |
| Advanced Beginner | Can modify component | Makes simple changes |
| Competent | Can debug component | Fixes common issues |
| Proficient | Can extend component | Adds new features |
| Expert | Can redesign component | Proposes improvements |

### Skill Assessment
| Skill | Basic | Intermediate | Advanced |
|-------|-------|--------------|----------|
| Implementation | Follows patterns | Adapts patterns | Creates patterns |
| Debugging | Uses tools | Reads state | Predicts behavior |
| Optimization | Measures performance | Identifies bottlenecks | Implements improvements |
| Documentation | Reads docs | Updates docs | Writes docs |

## Knowledge Synthesis

### Cross-Component Patterns
| Pattern | This Component | Related Component | Similarity | Difference |
|---------|---------------|-------------------|------------|------------|
| [Pattern] | [How implemented] | [Other example] | [Common aspects] | [Unique aspects] |

### Transferable Skills
#### Skill: [Problem Decomposition]
**Learned Here**: [How this component teaches it]
**Applied Elsewhere**: [Where else useful]
**Practice Progression**:
1. Simple: [Easy application]
2. Medium: [Moderate application]
3. Complex: [Advanced application]

### System-Level Understanding
```
Component Knowledge Graph:
[This Component]
├── Depends on understanding:
│   ├── [Prerequisite 1]
│   └── [Prerequisite 2]
├── Enables understanding:
│   ├── [Advanced Topic 1]
│   └── [Advanced Topic 2]
└── Relates to:
    ├── [Parallel Concept 1]
    └── [Parallel Concept 2]
```

## Reflection and Metacognition

### Learning Reflection Questions
#### After Implementation
1. What surprised you most?
2. What was harder than expected?
3. What was easier than expected?
4. What would you do differently?

#### After Debugging
1. What assumption was wrong?
2. How did you discover the issue?
3. What debugging technique helped most?
4. How will you prevent similar issues?

#### After Optimization
1. What was the bottleneck?
2. Why wasn't it obvious initially?
3. What measurement proved the improvement?
4. What did you sacrifice for performance?

### Metacognitive Strategies
#### Learning How to Learn
**Strategy**: [e.g., Rubber Duck Debugging]
**When to Use**: [Situation]
**How It Helps**: [Cognitive benefit]
**Practice Exercise**: [How to develop skill]

## Research Directions

### Open Problems
| Problem | Current Limitation | Research Needed | Potential Impact |
|---------|-------------------|-----------------|------------------|
| [Problem] | [What doesn't work] | [What to investigate] | [Why it matters] |

### Literature Connections
#### Foundational Papers
1. [Paper Title] - [Author, Year]
   - Key Contribution: [What it established]
   - Relevance: [How it relates]
   - Further Reading: [Related papers]

#### Recent Developments
1. [Paper/Article Title] - [Author, Year]
   - Innovation: [What's new]
   - Application: [How to use]
   - Open Questions: [What's unsolved]

### Future Learning Paths
```
After This Component:
├── Breadth Path:
│   ├── Study: [Related Component A]
│   ├── Study: [Related Component B]
│   └── Project: [Integration Project]
├── Depth Path:
│   ├── Study: [Advanced Theory]
│   ├── Research: [Paper Implementation]
│   └── Contribute: [Open Source Project]
└── Application Path:
    ├── Build: [Real-world Project]
    ├── Optimize: [Performance Tuning]
    └── Deploy: [Production System]
```

## Learning Resources

### Primary Resources
- **Essential Reading**: [Book/Paper with why it's essential]
- **Video Lectures**: [Course/Series with what it covers]
- **Interactive Tutorials**: [Resource with what you'll learn]

### Supplementary Resources
- **Alternative Explanations**: [Resource for different perspective]
- **Practice Problems**: [Source with difficulty level]
- **Community Forums**: [Where to discuss and ask questions]

### Advanced Resources
- **Research Papers**: [Cutting-edge developments]
- **Conference Talks**: [Industry insights]
- **Open Source Studies**: [Real-world implementations]

## Conclusion

### Mastery Checklist
- [ ] Can implement from scratch
- [ ] Can explain to beginners
- [ ] Can debug without documentation
- [ ] Can optimize for different constraints
- [ ] Can identify appropriate use cases
- [ ] Can propose meaningful improvements
- [ ] Can teach the patterns
- [ ] Can transfer knowledge to new domains

### Final Reflection
**The One Key Insight**: [Most important learning]
**The Unexpected Discovery**: [What surprised us]
**The Lasting Principle**: [What transcends this component]

### Next Steps
1. **Immediate**: [What to do now]
2. **Short-term**: [What to pursue next]
3. **Long-term**: [Where this leads]
```

### LEARN.md Unique Characteristics
1. **Language**: Educational, theoretical, reflective
2. **Structure**: Progressive from theory to practice
3. **Content**: Emphasizes learning science and pedagogy
4. **Depth**: Connects CS theory to practical implementation
5. **Exercises**: Scaffolded learning with rubrics
6. **Theory**: Grounded in educational psychology
7. **Assessment**: Clear learning objectives and outcomes
8. **Metacognition**: Promotes learning about learning

---

## General Documentation Rules

### File Location
All documentation MUST be located in `component/docs/`:
```
component/
├── docs/
│   ├── DOC.md
│   ├── DEV.md
│   ├── API.md
│   ├── TEST.md
│   └── LEARN.md
```

### Perspective Markers
Every document MUST include three perspective levels:
- `[MACRO_VIEW]...[/MACRO_VIEW]` - Project-wide context
- `[MESO_VIEW]...[/MESO_VIEW]` - Module-level context  
- `[MICRO_VIEW]...[/MICRO_VIEW]` - Component-specific context

### Document Navigation Flow
```
Entry Point (DOC.md)
    ├→ Usage Path (API.md)
    ├→ Development Path (DEV.md)
    ├→ Quality Path (TEST.md)
    └→ Mastery Path (LEARN.md)

Knowledge Progression:
• DOC.md: "I need to use this" (5 min read)
• API.md: "I need to integrate this" (20 min read)
• DEV.md: "I need to modify this" (30 min read)
• TEST.md: "I need to verify this" (25 min read)
• LEARN.md: "I need to master this" (2+ hours study)
```

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

| Aspect | DOC | DEV | API | TEST | LEARN |
|--------|-----|-----|-----|------|-------|
| **Audience** | End users | Developers | Integrators | QA/DevOps | Learners/Educators |
| **Questions** | What? Why? | How built? | How to use? | How to verify? | Why designed this way? |
| **Depth** | Surface | Deep internals | Interface details | Execution details | Complete understanding |
| **Code** | Minimal | Architecture | Full API | Test examples | Evolution examples |
| **Tone** | Friendly | Technical | Precise | Procedural | Educational |
| **Length** | 1-2 pages | 5-10 pages | 10+ pages | 5-8 pages | 15+ pages |
| **Outcome** | Can use it | Can modify it | Can integrate it | Can test it | Can teach it |

### Version Control Rules
1. Documentation updates MUST accompany code changes
2. Documentation reviews are part of code review
3. Breaking changes require migration guides in API.md
4. Version documentation separately from code
5. LEARN.md captures evolution history

### Quality Validation Checklist
- [ ] All five files present in `component/docs/`
- [ ] Perspective markers present and correctly tagged
- [ ] No language-specific examples in rules
- [ ] Tables properly formatted with headers
- [ ] Internal links are relative and working
- [ ] Code blocks have language hints
- [ ] Platform-specific sections in TEST.md
- [ ] Clear distinction between file purposes
- [ ] Examples are self-contained
- [ ] Sections follow prescribed order
- [ ] LEARN.md includes theoretical foundations

### Documentation Maintenance Triggers
| Event | DOC | DEV | API | TEST | LEARN |
|-------|-----|-----|-----|------|-------|
| New feature | Update features | Update architecture | Add endpoints | Add test cases | Document learning journey |
| Bug fix | If user-facing | If design change | If behavior change | Add regression test | Add to failure analysis |
| Performance improvement | If user-visible | Update metrics | Update SLA | Update benchmarks | Document optimization process |
| Security patch | Update if relevant | Update threat model | Update auth | Add security test | Add to security lessons |
| Dependency update | Update prerequisites | Update dependencies | If API affected | Update environment | Note in evolution |
| Refactoring | No change | Update internals | No change | Update if needed | Document refactoring decision |
| Failed experiment | No change | No change | No change | No change | Document in failure archive |
| New understanding | No change | Possibly update | No change | No change | Update theoretical foundations |

## Documentation Anti-Patterns to Avoid

1. **Mixing Concerns**: Don't put API details in DOC, learning theory in TEST
2. **Duplication**: Don't repeat content across files
3. **Language Bias**: Don't assume specific programming paradigms
4. **Missing Context**: Always include all three view levels
5. **Outdated Examples**: Examples must work with current version
6. **Unclear Audience**: Each file has ONE primary audience
7. **Wall of Text**: Use structure, tables, and formatting
8. **Hidden Knowledge**: Document all assumptions explicitly
9. **Missing Theory**: LEARN.md must connect to CS foundations
10. **No Progression**: Ensure clear learning path through documents

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
9. **Include educational theory** in LEARN.md
10. **Connect to CS fundamentals** where appropriate

### Document Generation Order
1. Generate DOC.md first (defines what)
2. Generate API.md second (defines how to use)
3. Generate DEV.md third (explains how it works)
4. Generate TEST.md fourth (verifies everything)
5. Generate LEARN.md last (synthesizes all learning)

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
  IF document is LEARN.md:
    CHECK theoretical foundations present
    CHECK pedagogical approach defined
    CHECK exercises with solutions included
    CHECK assessment rubrics provided
  END IF
END FOR

VALIDATE cross-document consistency
VALIDATE progressive complexity increase
VALIDATE no content duplication
```

## Example Component Documentation Structure

### Complete File Tree
```
my-component/
├── docs/
│   ├── DOC.md      # 2 pages - User introduction
│   ├── DEV.md      # 8 pages - Technical deep-dive
│   ├── API.md      # 12 pages - Complete interface
│   ├── TEST.md     # 6 pages - Testing guide
│   └── LEARN.md    # 20+ pages - Educational journey
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
│   ├── advanced/
│   └── educational/    # Examples for LEARN.md
└── config/
```

### Document Relationship Map
```
DOC.md (Entry Point)
    ↓ "Learn to use" → API.md
    ↓ "Learn internals" → DEV.md
    ↓ "Verify it works" → TEST.md
    ↓ "Master concepts" → LEARN.md
    
API.md ←→ DEV.md (Architecture informs API)
API.md ←→ TEST.md (Tests verify API contract)
DEV.md ←→ TEST.md (Architecture guides testing)
LEARN.md ← ALL (Synthesizes everything)
LEARN.md → ALL (Provides theoretical foundation)
```

## Enforcement Rules for Consistency

### Mandatory Elements Per Document

#### DOC.md Must Have:
- [ ] "What is This?" section with three views
- [ ] "Why Use This?" section with benefits
- [ ] Quick Start with immediate value
- [ ] Links to all other docs including LEARN.md
- [ ] Simple directory structure
- [ ] License information

#### DEV.md Must Have:
- [ ] Architecture overview with diagrams
- [ ] Design decisions with trade-offs
- [ ] Directory deep-dive with purposes
- [ ] Performance characteristics
- [ ] Security considerations
- [ ] Development workflow
- [ ] Link to LEARN.md for theory

#### API.md Must Have:
- [ ] Complete parameter documentation
- [ ] Return value specifications
- [ ] Error codes and handling
- [ ] Examples for every endpoint/method
- [ ] Rate limiting information
- [ ] Version compatibility matrix
- [ ] Link to LEARN.md for deeper understanding

#### TEST.md Must Have:
- [ ] Test pyramid/strategy diagram
- [ ] Platform-specific instructions (Windows/macOS/Linux)
- [ ] Coverage requirements and current status
- [ ] Troubleshooting guide
- [ ] Performance benchmarks
- [ ] CI/CD configuration
- [ ] Link to LEARN.md for testing philosophy

#### LEARN.md Must Have:
- [ ] Theoretical foundations section
- [ ] Pedagogical approach with Bloom's taxonomy
- [ ] Development journey with failures
- [ ] Exercises with solutions
- [ ] Assessment rubrics
- [ ] Mental models and misconceptions
- [ ] Literature connections
- [ ] Metacognitive reflection questions
- [ ] Progressive learning paths

### Prohibited Practices

1. **Never mix document purposes**
   - ❌ Don't put learning theory in API.md
   - ❌ Don't put API details in DOC.md
   - ❌ Don't put pedagogical content in TEST.md
   - ✅ Keep learning content in LEARN.md

2. **Never use language-specific examples in rules**
   - ❌ `public class MyClass`
   - ✅ `[Component initialization]`

3. **Never skip perspective levels**
   - ❌ Only micro-view
   - ✅ All three views with proper tags

4. **Never assume prior knowledge**
   - ❌ "As everyone knows..."
   - ✅ "This component handles..."

5. **Never ignore theoretical foundations in LEARN.md**
   - ❌ Only practical examples
   - ✅ Connect to CS theory and learning science

## Quick Reference Card for LLMs

### Document Purpose in One Line
- **DOC.md**: What it does and why you'd want it
- **DEV.md**: How it's built and why it's built that way
- **API.md**: How to integrate and use it completely
- **TEST.md**: How to verify it works correctly
- **LEARN.md**: How to master it and teach others

### Tone Guide
- **DOC.md**: Friendly teacher explaining benefits
- **DEV.md**: Senior engineer explaining architecture
- **API.md**: Technical writer documenting contract
- **TEST.md**: QA lead explaining verification
- **LEARN.md**: Professor combining theory and practice

### Audience Needs
- **DOC.md reader asks**: "Should I use this?"
- **DEV.md reader asks**: "How do I modify this?"
- **API.md reader asks**: "How do I call this?"
- **TEST.md reader asks**: "How do I verify this?"
- **LEARN.md reader asks**: "How do I truly understand and teach this?"

## Final Validation Checklist for LLMs

Before completing documentation generation:

### Structure Validation
- [ ] All documents in `component/docs/` directory
- [ ] All five documents present (DOC, DEV, API, TEST, LEARN)
- [ ] Each document has ALL required sections
- [ ] Sections appear in prescribed order

### Content Validation
- [ ] Three perspective views in each document
- [ ] Perspective tags properly opened and closed
- [ ] No language-specific syntax (Java, Python, etc.)
- [ ] Examples are conceptual, not implementation-specific
- [ ] Cross-references use relative paths

### LEARN.md Specific Validation
- [ ] Includes CS theoretical foundations
- [ ] Maps to Bloom's taxonomy
- [ ] Contains exercises with solutions
- [ ] Provides assessment rubrics
- [ ] Documents failure analysis
- [ ] Includes metacognitive elements
- [ ] Connects to research literature

### Formatting Validation
- [ ] ATX headers (#) used consistently
- [ ] Tables have header separator rows
- [ ] Code blocks use triple backticks
- [ ] Lists use consistent markers
- [ ] Line length ≤ 100 characters

### Audience Alignment
- [ ] DOC speaks to new users
- [ ] DEV speaks to developers
- [ ] API speaks to integrators
- [ ] TEST speaks to QA/DevOps
- [ ] LEARN speaks to learners and educators

### Completeness Check
- [ ] DOC has working quick start
- [ ] DEV has architecture diagrams
- [ ] API has examples for all endpoints
- [ ] TEST has OS-specific instructions
- [ ] LEARN has complete learning journey

---

## Meta: About This Ruleset v2.0

**Version**: 2.0.0
**Purpose**: Ensure consistent, high-quality, educational documentation across all components
**Optimization**: Designed for LLM processing and generation
**Agnostic**: Works with any programming language or framework
**Educational**: Emphasizes learning theory and knowledge transfer
**Location**: This ruleset should be referenced by LLMs when generating component documentation

### Version 2.0 Changes
1. Added comprehensive LEARN.md specification
2. Emphasized theoretical foundations and pedagogy
3. Included learning science principles
4. Added educational assessment frameworks
5. Integrated metacognitive elements
6. Connected to CS academic literature

### When to Update This Ruleset
- New documentation needs discovered through usage
- New pedagogical approaches identified
- LLM processing improvements discovered
- Feedback from learners and educators
- New theoretical frameworks emerge
- Better assessment methods developed

### Ruleset Principles
1. **Clarity** over brevity
2. **Consistency** over creativity
3. **Completeness** over conciseness
4. **Structure** over prose
5. **Examples** over explanations
6. **Theory** grounding practice
7. **Learning** as primary goal

---

*END OF RULESET v2.0 - Use these rules for all component documentation generation to ensure consistency, quality, and effective knowledge transfer across the application.*
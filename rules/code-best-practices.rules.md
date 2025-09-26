# Code Best Practices Rules for LLM Development

## Purpose
This document provides strict, language-agnostic guidelines for generating clean, maintainable, and production-ready code following industry standards and best practices.

## Pre-Development Analysis Requirements

### 0. MANDATORY: Project Context Analysis

#### 0.1 Existing Code Discovery
```
BEFORE WRITING ANY CODE, MUST:
1. Search for existing files in this order:
   a. Current directory (./)
   b. Parent directory (../)
   c. Sibling directories (../*)
   d. Project root directory
   e. Common directories (src/, lib/, utils/, common/, shared/)
2. Identify naming patterns and conventions used
3. Map existing functionality to avoid duplication
4. Document findings before proceeding
```

#### 0.2 Codebase Analysis Protocol
```
REQUIRED ANALYSIS STEPS:
1. FILE STRUCTURE ANALYSIS:
   - List all relevant existing files
   - Understand directory organization
   - Identify module boundaries
   - Note import/dependency patterns

2. CODE PATTERN RECOGNITION:
   - Coding style and formatting
   - Naming conventions in use
   - Error handling patterns
   - Testing approach
   - Documentation style

3. FUNCTIONALITY MAPPING:
   - List existing functions/classes
   - Understand data models
   - Map API endpoints or interfaces
   - Identify utility functions
   - Note design patterns used
```

#### 0.3 Context Coherence Rules
```
PRIORITY ORDER:
1. CORRECTNESS: Zero errors, no breaking changes
2. CONSISTENCY: Match existing patterns exactly
3. COMPATIBILITY: Ensure seamless integration
4. COMPLETENESS: Address all requirements
5. CLEANLINESS: Follow best practices

COHERENCE CHECKLIST:
- [ ] Analyzed all related existing files
- [ ] Understood project structure
- [ ] Matched existing naming conventions
- [ ] Compatible with current architecture
- [ ] No functionality duplication
- [ ] Consistent error handling
- [ ] Follows established patterns
```

## Core Principles

### 1. File Organization Rules

#### 1.1 File Size Limits
```
STRICT RULES:
- Maximum 300-500 lines of EXECUTABLE CODE per file
- Comments and documentation DO NOT count toward this limit
- Blank lines DO NOT count toward this limit
- Import/include statements DO NOT count toward this limit
- When approaching 300 lines, evaluate for splitting
- At 500 lines, MUST split into multiple files
```

#### 1.2 Horizontal Scaling Strategy
```
WHEN FILES EXCEED LIMITS:
1. Identify logical boundaries and responsibilities
2. Extract related functions into separate modules
3. Create specialized utility files
4. Implement service/repository pattern for separation
5. Use composition over monolithic files

EXAMPLE STRUCTURE:
Instead of: UserController (800 lines)
Split into:
- UserController (150 lines) - routing and coordination
- UserService (200 lines) - business logic
- UserValidator (100 lines) - validation rules
- UserTransformer (80 lines) - data transformation
- UserRepository (150 lines) - data access
```

#### 1.3 Documentation Requirements
```
MANDATORY DOCUMENTATION:
- File header: purpose, author, dependencies
- Function/method documentation for ALL public interfaces
- Complex algorithm explanations
- Business logic reasoning
- Integration points and external dependencies
- Usage examples for non-trivial functions

DOCUMENTATION DOES NOT COUNT TOWARD LINE LIMITS:
- Write comprehensive documentation
- Include examples liberally
- Document edge cases and assumptions
- Add inline comments for complex logic
- Maintain README files for modules
```

### 2. Code Quality Standards

#### 2.1 Readability First
- **ALWAYS** prioritize code readability over clever optimizations
- **USE** descriptive names that express intent clearly
- **AVOID** abbreviations except universally understood ones
- **MAINTAIN** consistent indentation throughout the codebase
- **LIMIT** line length to 80-120 characters for readability

#### 2.2 Universal Naming Conventions
```
LANGUAGE-AGNOSTIC PATTERNS:
- Actions/Functions: verb + noun (calculateTotal, getUserData)
- Variables: descriptive nouns (userEmail, totalAmount)
- Constants: clearly indicate immutability (MAX_RETRY, API_URL)
- Classes/Types: noun representing entity (UserService, Database)
- Booleans: state or capability (isLoading, hasPermission, canEdit)
- Collections: plural nouns (users, items, connections)
```

### 3. SOLID Principles (Language-Agnostic)

#### 3.1 Single Responsibility Principle
```
UNIVERSAL APPLICATION:
- Each module/class/function has ONE clear purpose
- Split code when it handles multiple concerns
- Ideal function length: 5-20 lines
- Maximum function length: 50 lines
- If describing the function requires "AND", split it
```

#### 3.2 Open/Closed Principle
```
IMPLEMENTATION STRATEGY:
- Use abstraction layers for extensibility
- Implement plugin architectures where appropriate
- Prefer configuration over modification
- Use dependency injection patterns
```

#### 3.3 Liskov Substitution Principle
```
CONSISTENCY RULES:
- Subtypes must be substitutable for base types
- Maintain behavioral consistency in hierarchies
- Avoid breaking contracts in derived implementations
```

#### 3.4 Interface Segregation Principle
```
INTERFACE DESIGN:
- Create focused, specific interfaces
- Avoid "fat" interfaces with unnecessary methods
- Split large interfaces into role-specific ones
```

#### 3.5 Dependency Inversion Principle
```
DEPENDENCY MANAGEMENT:
- Depend on abstractions, not implementations
- Inject dependencies rather than creating them
- Use factory patterns for object creation
```

### 4. Universal Clean Code Rules

#### 4.1 Function Design
```
UNIVERSAL FUNCTION RULES:
- Maximum 20 lines of executable code
- Single purpose and responsibility
- Maximum 3-4 parameters
- Return early to reduce nesting
- Avoid side effects when possible
- Name should fully describe behavior
```

#### 4.2 Error Handling
```
LANGUAGE-AGNOSTIC ERROR STRATEGIES:
- Fail fast and explicitly
- Handle errors at appropriate abstraction levels
- Never ignore or suppress errors silently
- Provide context in error messages
- Use error codes or types for different failures
- Implement retry logic for transient failures
- Log errors with sufficient detail for debugging
```

#### 4.3 Code Duplication
```
DRY PRINCIPLE ENFORCEMENT:
- Extract common code after 2nd duplication
- Create utility modules for shared functionality
- Use configuration for varying parameters
- Implement template/strategy patterns for similar flows
- Maximum 3-5 lines of acceptable duplication
```

### 5. Module Architecture

#### 5.1 Module Cohesion
```
MODULE ORGANIZATION RULES:
- Group related functionality together
- Separate by feature/domain, not technical layer
- Maintain clear module boundaries
- Minimize inter-module dependencies
- Export only necessary interfaces
```

#### 5.2 Dependency Management
```
DEPENDENCY RULES:
- No circular dependencies
- Minimize coupling between modules
- Use dependency injection
- Prefer composition over inheritance
- Implement facade pattern for complex subsystems
```

### 6. Code Complexity Management

#### 6.1 Cyclomatic Complexity
```
COMPLEXITY LIMITS:
- Maximum cyclomatic complexity: 10 per function
- Ideal complexity: 1-5
- Split complex conditionals into separate functions
- Use early returns to reduce nesting
- Extract complex boolean logic into named functions
```

#### 6.2 Nesting Depth
```
NESTING RULES:
- Maximum nesting depth: 3 levels
- Prefer guard clauses and early returns
- Extract nested loops into separate functions
- Use functional approaches when available
```

### 7. Data Handling

#### 7.1 Input Validation
```
UNIVERSAL VALIDATION RULES:
- Validate at system boundaries
- Sanitize all external inputs
- Use whitelist validation over blacklist
- Fail fast on invalid input
- Provide clear validation error messages
```

#### 7.2 State Management
```
STATE PRINCIPLES:
- Minimize mutable state
- Prefer immutable data structures
- Encapsulate state changes
- Make state transitions explicit
- Avoid global state when possible
```

### 8. Testing Standards

#### 8.1 Test Coverage Requirements
```
TESTING METRICS:
- Minimum 80% code coverage
- 100% coverage for critical business logic
- Test all error paths
- Include edge cases and boundary conditions
- Performance tests for critical paths
```

#### 8.2 Test Structure
```
UNIVERSAL TEST PATTERN:
1. Arrange - Set up test context
2. Act - Execute functionality
3. Assert - Verify outcomes
4. Cleanup - Reset state if needed

TEST FILE ORGANIZATION:
- Mirror source code structure
- One test file per source file
- Group related tests together
- Use descriptive test names
```

### 9. Documentation Standards

#### 9.1 Inline Documentation
```
REQUIRED DOCUMENTATION (NOT COUNTED IN LINE LIMIT):
/** 
 * Function Purpose: Clear description
 * Parameters: Type and purpose of each
 * Returns: Type and meaning
 * Throws: Possible exceptions/errors
 * Example: Usage demonstration
 * Complexity: O(n) notation if relevant
 * Side Effects: Any external changes
 */
```

#### 9.2 File Documentation
```
FILE HEADER TEMPLATE (NOT COUNTED IN LINE LIMIT):
/**
 * File: filename.ext
 * Purpose: Primary responsibility
 * Dependencies: External modules required
 * Exports: Public interfaces provided
 * Author: Creator/Maintainer
 * Created: Date
 * Modified: Last update date
 * Version: Semantic version
 * 
 * Business Context:
 * Detailed explanation of business logic
 * and domain concepts implemented
 */
```

### 10. Performance Guidelines

#### 10.1 Optimization Priorities
```
OPTIMIZATION ORDER:
1. Correctness first
2. Readability and maintainability
3. Algorithm efficiency (Big O)
4. Micro-optimizations last
5. Always measure before optimizing
```

#### 10.2 Resource Management
```
UNIVERSAL RESOURCE RULES:
- Explicitly release resources
- Use try-finally or equivalent patterns
- Implement connection pooling
- Cache expensive computations
- Lazy load when appropriate
```

### 11. Security Principles

#### 11.1 Security by Design
```
SECURITY REQUIREMENTS:
- Never hardcode credentials
- Encrypt sensitive data
- Use secure random generators
- Implement rate limiting
- Follow principle of least privilege
- Audit security-critical operations
```

#### 11.2 Data Protection
```
DATA HANDLING RULES:
- Sanitize all outputs
- Use parameterized queries
- Hash passwords with salt
- Implement secure session management
- Never log sensitive information
```

## Enforcement Rules for LLMs

### Pre-Code Generation Workflow
```
MANDATORY WORKFLOW BEFORE WRITING CODE:
1. DISCOVER: Search for ALL related existing files
2. ANALYZE: Study code patterns, conventions, and architecture
3. MAP: Document existing functionality to avoid duplication
4. PLAN: Design solution coherent with existing codebase
5. VALIDATE: Ensure no conflicts or duplications
6. IMPLEMENT: Write code following discovered patterns
7. INTEGRATE: Ensure seamless connection with existing code
```

### File Generation Strategy
```
WHEN GENERATING CODE:
1. FIRST: Check if functionality already exists
2. Search progressively from current to distant directories
3. Analyze existing similar modules for patterns
4. Check current file line count (excluding comments)
5. If approaching 300 lines, plan for splitting
6. At 500 lines, MUST create new files
7. Suggest file structure upfront for large features
8. Include comprehensive documentation without line limit concerns
```

### Error Prevention Protocol
```
BEFORE WRITING EACH FUNCTION/CLASS:
1. Search for existing implementation
2. Check for similar functionality
3. Verify naming conflicts
4. Ensure type compatibility
5. Validate interface contracts
6. Confirm no breaking changes
7. Test integration points
```

### MUST-FOLLOW Rules
1. **ALWAYS** analyze existing code BEFORE writing new code
2. **NEVER** duplicate existing functionality
3. **NEVER** write code that conflicts with existing patterns
4. **NEVER** exceed 500 lines of executable code per file
5. **NEVER** count documentation toward line limits
6. **NEVER** sacrifice documentation for line limits
7. **ALWAYS** search directories from nearest to farthest
8. **ALWAYS** maintain consistency with project conventions
9. **ALWAYS** split files by logical boundaries
10. **ALWAYS** maintain high cohesion within files
11. **ALWAYS** document all public interfaces
12. **ALWAYS** validate inputs at boundaries
13. **ALWAYS** handle errors explicitly
14. **ALWAYS** consider security implications
15. **ALWAYS** write testable code
16. **ALWAYS** prioritize correctness over everything else

### Code Review Checklist
Before presenting code, verify:
- [ ] Searched and analyzed ALL existing related files
- [ ] No duplication of existing functionality
- [ ] Consistent with project patterns and conventions
- [ ] Zero breaking changes to existing code
- [ ] Each file has <500 lines of executable code
- [ ] All public interfaces are documented
- [ ] Functions are small and focused (<20 lines)
- [ ] No code duplication (DRY principle)
- [ ] Error handling is comprehensive
- [ ] Security considerations addressed
- [ ] File has clear, single responsibility
- [ ] Dependencies are properly managed
- [ ] Code is testable and tests provided
- [ ] Documentation is complete and helpful
- [ ] Integration points verified and tested

### Multi-File Generation Guidelines
```
WHEN SPLITTING FILES:
1. First analyze existing file structure
2. Check for existing similar modules
3. Maintain consistency with project organization
4. Identify core responsibilities
5. Create interface/contract file first
6. Implement each responsibility separately
7. Create integration/coordinator file
8. Add comprehensive tests for each file
9. Document inter-file relationships
10. Verify no conflicts with existing code

EXAMPLE SPLIT:
Feature: User Management (estimated 1200 lines)
FIRST: Check if user-related files exist
THEN Generate only missing components:
- user.interface (50 lines + documentation)
- user.model (150 lines + documentation)  
- user.validator (120 lines + documentation)
- user.service (250 lines + documentation)
- user.repository (200 lines + documentation)
- user.controller (180 lines + documentation)
- user.test (300 lines + documentation)
```

### Duplication Prevention Strategy
```
ANTI-DUPLICATION PROTOCOL:
1. Before creating any utility function:
   - Search utils/, helpers/, common/ directories
   - Check for similar functionality in service layers
   - Look for existing implementations in parent modules

2. Before creating any data model:
   - Check models/, entities/, types/ directories
   - Verify no existing similar structures
   - Ensure compatibility with existing models

3. Before adding any configuration:
   - Check config/, settings/, env files
   - Verify no existing similar settings
   - Maintain consistency with configuration patterns
```

## Priority Order for Conflicting Guidelines
1. **Zero Errors and Project Coherence**
2. Security and Safety
3. Consistency with existing codebase
4. File size limits (300-500 lines)
5. Correctness and Reliability  
6. Comprehensive Documentation
7. Maintainability and Readability
8. Testability
9. Performance
10. Elegance and Brevity

## Continuous Improvement Principles
- **REFACTOR** when approaching file size limits
- **EXTRACT** common patterns into utilities
- **DOCUMENT** all architectural decisions
- **MONITOR** code metrics and complexity
- **REVIEW** and update module boundaries
- **MAINTAIN** clear separation of concerns

---

**CRITICAL INSTRUCTIONS FOR LLMs**: 
1. Always count ONLY executable code toward the 300-500 line limit
2. Documentation, comments, blank lines, and imports DO NOT count
3. Write extensive documentation without worrying about line limits
4. When generating features, proactively split into multiple files
5. Suggest file architecture before implementation for large features
6. These rules apply regardless of programming language or framework
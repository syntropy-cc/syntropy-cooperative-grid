### Test Completeness Checklist

Before considering the test suite complete:

**File Location Verification:**
- [ ] ALL test files are in `component/tests/` directory
- [ ] NO files were created in `component/src/` directory
- [ ] NO files were modified in `component/src/` directory
- [ ] Source code in `src/` remains completely unchanged

**Coverage Verification (for src/ code):**
- [ ] 100% line coverage of src/ achieved
- [ ] 100% branch coverage of src/ achieved
- [ ] 100% path coverage of src/ achieved
- [ ] All error conditions in src/ tested
- [ ] All edge cases in src/ covered
- [ ] All security vulnerabilities in src/ tested
- [ ] All performance requirements of src/ validated
- [ ] All integrations used by src/ tested
- [ ] All user journeys in src/ tested

**Test Quality Verification:**
- [ ] All tests are independent
- [ ] All tests are deterministic
- [ ] All tests execute quickly
- [ ] All tests have clear names
- [ ] All test data is isolated in tests/fixtures/
- [ ] All mocks are complete in tests/mocks/
- [ ] All helpers are implemented in tests/helpers/
- [ ] No anti-patterns present
- [ ] Quality gates passed

**Implementation Verification:**
- [ ] Tests import from src/ correctly
- [ ] Tests use relative paths properly
- [ ] No hardcoded paths to src/
- [ ] Mocks simulate src/ dependencies properly
- [ ] Fixtures cover all src/ input scenarios
- [ ] Helpers support all test scenarios# LLM Testing Rules - Professional Test Suite Implementation Guide

## Executive Summary

This document provides comprehensive, language-agnostic rules for LLMs to implement a complete, professional-grade test suite achieving 100% code coverage. These rules follow industry best practices and standards including ISO/IEC/IEEE 29119, ISTQB guidelines, and proven testing methodologies.

## CRITICAL CONSTRAINTS - READ FIRST

### File System Boundaries

**ABSOLUTE RULES:**
1. **NEVER create or modify ANY file outside the `tests/` directory**
2. **NEVER modify ANY file in the `src/` directory**
3. **NEVER add new files to the `src/` directory**
4. **The `src/` directory is READ-ONLY - treat it as immutable**
5. **ALL test code, fixtures, mocks, and helpers MUST be created inside `tests/` directory only**

### Component Structure

```
component/
├── src/          # READ-ONLY - Contains implementation code
│   └── [implementation files - DO NOT MODIFY]
└── tests/        # WRITE-ONLY - All test code goes here
    ├── unit/
    ├── integration/
    ├── e2e/
    ├── performance/
    ├── security/
    ├── fixtures/
    ├── mocks/
    └── helpers/
```

### Implementation Boundaries

- **Source Code (`src/`)**: Contains the implementation to be tested - READ ONLY
- **Test Code (`tests/`)**: Contains all testing artifacts - CREATE ALL FILES HERE
- **Independence**: Tests must NEVER modify source code to make tests pass
- **Isolation**: Tests import from `src/` but NEVER write to it
- **Black Box**: Treat `src/` as an external, immutable dependency

## Core Testing Principles

1. **Complete Coverage**: 100% code coverage is mandatory - every line, branch, and path must be tested
2. **Test Independence**: Each test must be completely isolated and idempotent
3. **Deterministic Execution**: Tests must produce consistent, reproducible results
4. **Fast Feedback**: Optimize for rapid test execution without compromising coverage
5. **Clear Intent**: Test names and structure must clearly communicate purpose and expected behavior
6. **Single Responsibility**: Each test validates exactly one behavior or requirement
7. **Test as Specification**: Tests serve as executable specifications of system behavior

## Industry Standards Applied

- **ISO/IEC/IEEE 29119**: Software Testing Standard
- **IEEE 829**: Test Documentation Standard  
- **ISTQB**: Testing Body of Knowledge
- **xUnit Patterns**: Test Architecture Patterns
- **Test Pyramid**: Mike Cohn's Test Distribution Model
- **FIRST Principles**: Fast, Independent, Repeatable, Self-validating, Timely
- **SOLID Principles**: Applied to test architecture

## Directory Structure Requirements

**IMPORTANT**: All paths below are relative to `component/tests/`. Never create files outside this directory.

```
component/
├── src/                  # [READ-ONLY] - Implementation code - DO NOT MODIFY
│   └── ...              # Existing source files - NEVER CHANGE
└── tests/               # [WRITE-ONLY] - All test files created here
    ├── unit/           # Isolated component tests
    │   ├── core/      # Core business logic tests
    │   └── utilities/ # Utility function tests
    ├── integration/   # Component interaction tests
    │   ├── api/      # External API integration tests
    │   └── database/ # Data persistence tests
    ├── e2e/          # End-to-end user journeys
    │   └── scenarios/# Complete workflow tests
    ├── performance/  # Performance validation
    ├── security/     # Security verification
    ├── fixtures/     # Test data
    │   ├── valid/   # Valid input scenarios
    │   └── invalid/ # Invalid input scenarios
    ├── mocks/       # Test doubles
    └── helpers/     # Test utilities
```

### File Creation Rules

1. **All test files**: Must be created under `component/tests/`
2. **Import statements**: Can import from `../src/` (read-only access)
3. **Test artifacts**: Fixtures, mocks, helpers stay within `tests/`
4. **No source modification**: Tests must work with existing `src/` code as-is
5. **Test independence**: Tests cannot require changes to source code

---

## Test Implementation Layers

### Layer 1: Unit Tests (`component/tests/unit/`)

#### Purpose
Validate individual components from `src/` directory in complete isolation with 100% coverage of all logic paths. Tests read from `src/` but never modify it.

#### Implementation Requirements

##### Test Coverage Targets for Source Code
- **Line Coverage**: 100% of `src/` code
- **Branch Coverage**: 100% of `src/` code
- **Path Coverage**: 100% of `src/` code
- **Condition Coverage**: 100% of `src/` code
- **Modified Condition/Decision Coverage (MC/DC)**: 100% of `src/` code

##### Test Structure Pattern
```
TEST FILE LOCATION: component/tests/unit/[test-name].test
IMPORTS FROM: ../src/[module-name]

TEST SUITE: [Component Name from src/]
  TEST CONTEXT: [Method/Function Name from src/]
    TEST GROUP: Normal Conditions
      TEST CASE: Should [expected behavior] when [valid input scenario 1]
      TEST CASE: Should [expected behavior] when [valid input scenario 2]
      TEST CASE: Should [expected behavior] with [boundary value - minimum]
      TEST CASE: Should [expected behavior] with [boundary value - maximum]
    
    TEST GROUP: Edge Cases
      TEST CASE: Should handle [edge case 1] correctly
      TEST CASE: Should handle [edge case 2] correctly
      TEST CASE: Should process [boundary condition] as expected
    
    TEST GROUP: Error Conditions
      TEST CASE: Should [error behavior] when [invalid input scenario 1]
      TEST CASE: Should [error behavior] when [invalid input scenario 2]
      TEST CASE: Should [error behavior] when [null/undefined input]
      TEST CASE: Should [error behavior] when [type mismatch]
    
    TEST GROUP: State Verification
      TEST CASE: Should maintain [invariant 1] after operation
      TEST CASE: Should maintain [invariant 2] after operation
      TEST CASE: Should transition from [state A] to [state B] when [event]
```

##### Unit Test Categories (100% Coverage Required)

###### 1. Input Space Partitioning
- **Valid Equivalence Classes**: All valid input ranges
- **Invalid Equivalence Classes**: All invalid input ranges
- **Boundary Values**: Min, Min+1, Nominal, Max-1, Max
- **Special Values**: Zero, One, Negative One, Empty, Null, Undefined

###### 2. Logic Coverage
- **Statement Coverage**: Every statement executed
- **Decision Coverage**: Both true and false outcomes
- **Condition Coverage**: Each boolean sub-expression
- **Path Coverage**: All possible execution paths
- **Loop Coverage**: Zero, One, Many, Maximum iterations

###### 3. State-Based Testing
- **State Transitions**: All valid state transitions
- **Invalid Transitions**: All invalid state attempts
- **State Invariants**: Properties that must hold
- **Initial State**: Default/constructor state
- **Terminal State**: End states and cleanup

###### 4. Exception Testing
- **Expected Exceptions**: All documented error cases
- **Exception Properties**: Error messages, codes, types
- **Exception Propagation**: Handling and re-throwing
- **Recovery Scenarios**: State after exception

###### 5. Mutation Testing Requirements
- **Arithmetic Operator Mutations**: +, -, *, /, %
- **Relational Operator Mutations**: <, >, <=, >=, ==, !=
- **Logical Operator Mutations**: &&, ||, !
- **Assignment Mutations**: =, +=, -=, *=, /=
- **Constant Mutations**: Changing literal values

##### Performance Constraints for Unit Tests
- Maximum execution time: 10ms per test
- Maximum memory usage: 10MB per test
- No I/O operations (file, network, database)
- No real time dependencies

---

### Layer 2: Integration Tests (`component/tests/integration/`)

#### Purpose
Verify correct interaction between components from `src/` and with external systems while maintaining 100% interface coverage. Tests import from `src/` but never modify source files.

#### Implementation Requirements

##### Integration Points Coverage (100% Required)
```
TEST FILE LOCATION: component/tests/integration/
IMPORTS FROM: ../src/[modules]

INTEGRATION SUITE: [Integration Name]
  CATEGORY: Component Integration
    TEST: Component A → Component B data flow (both from src/)
    TEST: Component B → Component A feedback (both from src/)
    TEST: Error propagation between components
    TEST: Transaction boundaries
    
  CATEGORY: External System Integration
    TEST: Database connectivity (using src/ database module)
    TEST: Database transactions (commit/rollback)
    TEST: External API communication (using src/ API client)
    TEST: Message queue operations (using src/ queue module)
    TEST: Cache synchronization (using src/ cache module)
    
  CATEGORY: Contract Testing
    TEST: Request schema validation (validating src/ contracts)
    TEST: Response schema validation (validating src/ contracts)
    TEST: Error response formats
    TEST: Version compatibility
```

##### Integration Test Scenarios

###### 1. API Integration Testing
- **HTTP Methods**: GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD
- **Status Codes**: 2xx, 3xx, 4xx, 5xx ranges
- **Headers**: Content-Type, Authorization, Custom headers
- **Request Formats**: JSON, XML, Form-data, Binary
- **Response Formats**: JSON, XML, Binary, Streaming
- **Pagination**: Offset, Cursor, Page-based
- **Filtering**: Field-based, Range-based, Text search
- **Sorting**: Single field, Multiple fields, Direction
- **Rate Limiting**: Request limits, Time windows, Retry-After
- **Authentication**: Basic, Bearer, OAuth, API Key
- **Versioning**: URL-based, Header-based, Query parameter

###### 2. Database Integration Testing
- **CRUD Operations**: Create, Read, Update, Delete
- **Batch Operations**: Bulk insert, Bulk update, Bulk delete
- **Transactions**: BEGIN, COMMIT, ROLLBACK, SAVEPOINT
- **Isolation Levels**: Read Uncommitted, Read Committed, Repeatable Read, Serializable
- **Constraints**: Primary Key, Foreign Key, Unique, Check, Not Null
- **Indexes**: Performance impact, Covering indexes
- **Stored Procedures**: Input/Output parameters, Return values
- **Triggers**: Before/After, Insert/Update/Delete
- **Views**: Materialized, Standard
- **Locks**: Pessimistic, Optimistic, Deadlock detection

###### 3. Message Queue Integration
- **Publishing**: Single message, Batch messages, Priority messages
- **Consuming**: Poll-based, Push-based, Acknowledgments
- **Error Handling**: Dead letter queues, Retry policies
- **Ordering**: FIFO, Priority-based, Partition-based
- **Durability**: Persistent, Non-persistent
- **Transactions**: Transactional sends/receives

##### Integration Test Data Management
- **Test Database**: Isolated test database instance
- **Data Seeding**: Consistent initial state
- **Data Cleanup**: Post-test cleanup procedures
- **Test Transactions**: Rollback after test completion
- **Fixtures**: Reusable test data sets
- **Data Factories**: Dynamic test data generation

---

### Layer 3: End-to-End Tests (`component/tests/e2e/`)

#### Purpose
Validate complete user workflows using the implementation from `src/` directory without modifying any source code.

#### Implementation Requirements

##### E2E Scenario Coverage
```
TEST FILE LOCATION: component/tests/e2e/scenarios/
SYSTEM UNDER TEST: ../src/ (entire application)

E2E SUITE: [Application Name]
  CRITICAL PATHS: (Must have 100% coverage of src/ user flows)
    SCENARIO: User Registration Flow
      STEP 1: Navigate to registration (using src/ routes)
      STEP 2: Fill registration form (testing src/ validation)
      STEP 3: Submit and verify (testing src/ business logic)
      STEP 4: Email confirmation (testing src/ email module)
      STEP 5: Account activation (testing src/ activation logic)
      
    SCENARIO: Authentication Flow
      STEP 1: Login with credentials (testing src/ auth module)
      STEP 2: Multi-factor authentication (testing src/ MFA)
      STEP 3: Session establishment (testing src/ session)
      STEP 4: Session refresh (testing src/ token refresh)
      STEP 5: Logout process (testing src/ cleanup)
      
    SCENARIO: Core Business Process
      [Test complete workflow using src/ modules]
      
  ALTERNATIVE PATHS:
    SCENARIO: Error Recovery (testing src/ error handlers)
    SCENARIO: Concurrent Users (testing src/ concurrency)
    SCENARIO: System Degradation (testing src/ fallbacks)
```

##### E2E Test Implementation Patterns

###### 1. User Journey Testing
- **Entry Points**: All possible system entry points
- **Navigation Paths**: All navigation routes
- **User Actions**: Click, Type, Select, Drag, Scroll, Upload
- **Wait Conditions**: Element visibility, Data loading, Animation completion
- **Assertions**: Visual, Functional, Data, State
- **Exit Points**: Logout, Timeout, Completion

###### 2. Cross-Browser/Platform Testing
- **Desktop Browsers**: Latest 2 versions of major browsers
- **Mobile Browsers**: iOS Safari, Chrome Android
- **Screen Resolutions**: Mobile, Tablet, Desktop, 4K
- **Orientations**: Portrait, Landscape
- **Accessibility**: Keyboard navigation, Screen readers

###### 3. Data Flow Validation
- **Input Validation**: All form fields and inputs
- **Data Persistence**: Verify storage and retrieval
- **Data Transformation**: Processing and calculations
- **Data Presentation**: Display formatting
- **Data Export**: Download and export functions

---

### Layer 4: Performance Tests (`tests/performance/`)

#### Purpose
Validate system performance characteristics meet requirements under various load conditions.

#### Implementation Requirements

##### Performance Test Types (All Required)

```
PERFORMANCE SUITE: [System Name]
  LOAD TESTING:
    TEST: Normal Load (Expected daily traffic)
    TEST: Peak Load (Expected maximum traffic)
    TEST: Sustained Load (Extended duration)
    
  STRESS TESTING:
    TEST: Breaking Point (Find system limits)
    TEST: Recovery (System recovery after stress)
    
  SPIKE TESTING:
    TEST: Sudden Load Increase
    TEST: Sudden Load Decrease
    
  VOLUME TESTING:
    TEST: Large Data Sets
    TEST: Database Performance
    
  ENDURANCE TESTING:
    TEST: Memory Leaks
    TEST: Resource Exhaustion
```

##### Performance Metrics (Must Measure All)

###### 1. Response Time Metrics
- **Minimum Response Time**
- **Average Response Time**
- **Median Response Time**
- **90th Percentile (P90)**
- **95th Percentile (P95)**
- **99th Percentile (P99)**
- **Maximum Response Time**

###### 2. Throughput Metrics
- **Requests per Second (RPS)**
- **Transactions per Second (TPS)**
- **Data Transfer Rate**
- **Concurrent Users Supported**

###### 3. Resource Utilization
- **CPU Usage**: User, System, Wait, Idle
- **Memory Usage**: Heap, Non-heap, Native
- **Disk I/O**: Read/Write operations, Queue depth
- **Network I/O**: Bandwidth, Packet loss, Latency
- **Thread/Process Count**: Active, Idle, Blocked
- **Connection Pool**: Active, Idle, Waiting

###### 4. Error Metrics
- **Error Rate**: Percentage of failed requests
- **Error Types**: Timeout, Connection, Server, Client
- **Recovery Time**: Time to recover from errors

##### Performance Test Scenarios

1. **Baseline Testing**: Single user performance
2. **Load Testing**: Expected user load
3. **Stress Testing**: Beyond capacity
4. **Spike Testing**: Sudden load changes
5. **Soak Testing**: Extended duration
6. **Scalability Testing**: Horizontal/Vertical scaling
7. **Failover Testing**: Component failure handling

---

### Layer 5: Security Tests (`tests/security/`)

#### Purpose
Validate security controls and identify vulnerabilities across all attack vectors.

#### Implementation Requirements

##### Security Test Coverage (100% Required)

```
SECURITY SUITE: [Application Name]
  AUTHENTICATION TESTING:
    TEST: Password strength requirements
    TEST: Account lockout mechanisms
    TEST: Session management
    TEST: Multi-factor authentication
    TEST: Password reset security
    TEST: Brute force protection
    
  AUTHORIZATION TESTING:
    TEST: Role-based access control
    TEST: Resource-level permissions
    TEST: Privilege escalation attempts
    TEST: Cross-tenant access prevention
    
  INPUT VALIDATION:
    TEST: SQL Injection
    TEST: NoSQL Injection
    TEST: LDAP Injection
    TEST: XPath Injection
    TEST: Command Injection
    TEST: Cross-Site Scripting (XSS)
    TEST: XML External Entity (XXE)
    TEST: Path Traversal
    
  CRYPTOGRAPHY:
    TEST: Encryption at rest
    TEST: Encryption in transit
    TEST: Key management
    TEST: Certificate validation
    TEST: Secure random generation
```

##### Security Testing Checklist

###### OWASP Top 10 Coverage
1. **Injection**: SQL, NoSQL, OS, LDAP injection tests
2. **Broken Authentication**: Session, credential, token tests
3. **Sensitive Data Exposure**: Encryption, masking, storage tests
4. **XML External Entities**: XXE prevention tests
5. **Broken Access Control**: Authorization, RBAC tests
6. **Security Misconfiguration**: Configuration, header tests
7. **Cross-Site Scripting**: Reflected, Stored, DOM XSS tests
8. **Insecure Deserialization**: Object injection tests
9. **Using Components with Known Vulnerabilities**: Dependency scans
10. **Insufficient Logging & Monitoring**: Audit, alert tests

###### Additional Security Tests
- **CSRF Protection**: Token validation
- **Clickjacking**: X-Frame-Options
- **SSL/TLS**: Protocol versions, Cipher suites
- **CORS**: Origin validation
- **Rate Limiting**: Request throttling
- **File Upload**: Type validation, Size limits, Malware scanning
- **API Security**: Authentication, Rate limiting, Input validation
- **Session Security**: Timeout, Fixation, Hijacking
- **Error Handling**: Information disclosure prevention

---

## Test Data Management (`component/tests/fixtures/`)

### Test Data Categories

**IMPORTANT**: All test data must be created in `tests/fixtures/`. Never modify data in `src/`.

#### 1. Valid Data Sets
```
component/tests/fixtures/valid/
├── minimal/          # Minimum required fields for src/ functions
├── complete/         # All fields populated for src/ functions
├── typical/          # Common use cases from src/ requirements
├── maximum/          # Maximum allowed values per src/ validation
└── variations/       # Format variations accepted by src/
```

#### 2. Invalid Data Sets
```
component/tests/fixtures/invalid/
├── missing-required/ # Missing fields that src/ requires
├── type-errors/     # Wrong types that src/ should reject
├── constraint-violations/ # Violates src/ business rules
├── malformed/       # Structural errors src/ should catch
└── malicious/       # Security payloads src/ should block
```

#### 3. Edge Case Data Sets
```
component/tests/fixtures/edge-cases/
├── boundary-values/ # Test src/ boundary handling
├── special-characters/ # Test src/ character handling
├── large-datasets/  # Test src/ volume limits
├── concurrent/      # Test src/ race conditions
└── temporal/        # Test src/ time handling
```

### Test Data Principles

1. **Deterministic**: Use seeded random generation
2. **Isolated**: Each test owns its data
3. **Realistic**: Mirror production patterns
4. **Comprehensive**: Cover all scenarios
5. **Versioned**: Track data schema changes
6. **Reusable**: Shared fixtures for common cases
7. **Documented**: Clear data intent and purpose

---

## Mock Implementation (`component/tests/mocks/`)

### Mock Types Required

**CRITICAL**: All mocks must be created in `tests/mocks/`. Mocks simulate `src/` dependencies but NEVER modify source files.

#### 1. Test Doubles Taxonomy
- **Dummy**: Objects passed to src/ functions but never used
- **Stub**: Provides predetermined responses to src/ calls
- **Spy**: Records src/ function interactions for verification
- **Mock**: Pre-programmed expectations for src/ behavior
- **Fake**: Working implementation replacing src/ dependencies

#### 2. Mock Implementation Requirements

```
FILE LOCATION: component/tests/mocks/[service-name].mock
PURPOSE: Mock external dependencies used by src/ code

MOCK: [Service Name used by src/]
  CAPABILITIES:
    - Record all calls from src/ code
    - Verify src/ call count and parameters
    - Simulate responses that src/ expects
    - Simulate errors that src/ should handle
    - Simulate timeouts that src/ should handle
    - Simulate rate limiting for src/ retry logic
    - Simulate partial failures for src/ resilience
    
  VERIFICATION METHODS:
    - Verify src/ called with specific parameters
    - Verify src/ called N times
    - Verify src/ called in specific order
    - Verify src/ did not call
    
  CONFIGURATION:
    - Response delays (test src/ timeout handling)
    - Error rates (test src/ error handling)
    - Custom responses (test src/ parsing)
    - Conditional behavior (test src/ branching)
```

### Mock Best Practices

1. **Interface Compliance**: Match interfaces that src/ expects exactly
2. **Behavior Simulation**: Realistic responses for src/ consumption
3. **State Management**: Maintain state between src/ calls
4. **Reset Capability**: Clear state between tests
5. **Error Injection**: Test src/ error handling paths
6. **Performance Characteristics**: Test src/ timeout logic
7. **Verification API**: Verify how src/ uses dependencies

---

## Test Helpers (`component/tests/helpers/`)

### Required Helper Categories

**LOCATION**: All helpers must be in `tests/helpers/`. They assist in testing `src/` code without modifying it.

#### 1. Setup Utilities
```
component/tests/helpers/setup/
- Environment configuration for testing src/
- Test database initialization for src/ DB modules
- Test server startup for src/ server code
- Service connections for src/ integrations
- Authentication setup for src/ auth testing
- Test user creation for src/ user management
```

#### 2. Assertion Helpers
```
component/tests/helpers/assertions/
- Custom matchers for src/ output validation
- Deep equality checks for src/ objects
- Schema validation for src/ data structures
- Approximate equality for src/ calculations
- Collection assertions for src/ arrays/lists
- Async assertions for src/ promises
```

#### 3. Wait/Retry Utilities
```
component/tests/helpers/timing/
- Wait for src/ async operations
- Retry src/ operations with backoff
- Timeout handling for src/ long operations
- Polling mechanisms for src/ state changes
- Event waiting for src/ event emitters
- Promise utilities for src/ async code
```

#### 4. Data Builders
```
component/tests/helpers/builders/
- Object mothers for src/ entity testing
- Test data builders for src/ inputs
- Random generators matching src/ constraints
- Factory methods for src/ object creation
- Fixture loaders for src/ data processing
- Snapshot creators for src/ state capture
```

#### 5. Cleanup Utilities
```
component/tests/helpers/cleanup/
- Database cleanup after testing src/ DB ops
- File system cleanup after src/ file ops
- Network cleanup after src/ network ops
- Process cleanup after src/ process spawning
- Memory cleanup after src/ memory ops
- State reset between src/ tests
```

---

## Test Execution Strategy

### Execution Order and Prioritization

#### Priority Levels
1. **P0 - Critical**: System cannot function without these
2. **P1 - High**: Core functionality tests
3. **P2 - Medium**: Important feature tests
4. **P3 - Low**: Nice-to-have feature tests

#### Execution Phases
```
PHASE 1: Pre-flight Checks
  - Environment validation
  - Dependency verification
  - Configuration validation

PHASE 2: Unit Tests (Parallel)
  - All unit tests run in parallel
  - Fail fast on critical failures

PHASE 3: Integration Tests (Sequential Groups)
  - Database integration tests
  - API integration tests
  - External service tests

PHASE 4: E2E Tests (Sequential)
  - Critical path scenarios
  - Alternative path scenarios

PHASE 5: Performance Tests
  - Baseline performance
  - Load testing
  - Stress testing

PHASE 6: Security Tests
  - Vulnerability scanning
  - Penetration testing
```

### Test Isolation Requirements

1. **Process Isolation**: Each test suite in separate process
2. **Data Isolation**: No shared mutable state
3. **Network Isolation**: Mock external dependencies
4. **File System Isolation**: Temporary directories
5. **Time Isolation**: Controllable time source
6. **Random Isolation**: Seeded random generators

---

## Coverage Requirements

### Mandatory Coverage Metrics (100% Required)

**TARGET**: 100% coverage of all code in `src/` directory. The `tests/` directory achieves this without modifying any source files.

```
COVERAGE REQUIREMENTS FOR src/ CODE:
  Code Coverage:
    - Line Coverage: 100% of src/
    - Branch Coverage: 100% of src/
    - Function Coverage: 100% of src/
    - Statement Coverage: 100% of src/
    
  Logic Coverage:
    - Decision Coverage: 100% of src/
    - Condition Coverage: 100% of src/
    - MC/DC Coverage: 100% of src/
    - Path Coverage: 100% of src/
    
  Data Flow Coverage:
    - All-Defs: 100% of src/
    - All-Uses: 100% of src/
    - All-P-Uses: 100% of src/
    - All-C-Uses: 100% of src/
    
  Mutation Coverage:
    - Mutation Score: 100% of src/
    - Killed Mutants: 100% of src/
```

### Coverage Scope

- **Include**: All files in `src/` directory
- **Exclude**: Only `tests/` directory files (these are tests, not code)
- **No Exceptions**: Every line in `src/` must be covered

### Coverage Validation

Tests in `tests/` directory must exercise:
1. Every line of code in `src/`
2. Every branch condition in `src/`
3. Every function/method in `src/`
4. Every error path in `src/`
5. Every edge case in `src/`

---

## Quality Gates

### Test Quality Metrics

1. **Test Effectiveness**
   - Defect detection rate: >95%
   - False positive rate: <1%
   - False negative rate: 0%

2. **Test Efficiency**
   - Unit test speed: <10ms per test
   - Integration test speed: <100ms per test
   - E2E test speed: <5s per scenario
   - Total suite execution: <10 minutes

3. **Test Maintainability**
   - Test code duplication: <5%
   - Cyclomatic complexity: <5 per test
   - Test documentation coverage: 100%

4. **Test Reliability**
   - Test flakiness: 0%
   - Test determinism: 100%
   - Platform independence: 100%

---

## Anti-Patterns to Avoid

### Test Implementation Anti-Patterns

1. **The Liar**: Tests that pass regardless of implementation
2. **The Giant**: Tests that validate too many things
3. **The Mockery**: Over-mocking leading to false confidence
4. **The Inspector**: Tests that know too much about internals
5. **The Slow Poke**: Tests that take too long to execute
6. **The Flickering**: Non-deterministic tests
7. **The Chain Gang**: Tests that depend on execution order
8. **The Free Rider**: Tests that don't add coverage
9. **The Loudmouth**: Tests with excessive logging
10. **The Secret Keeper**: Tests without clear intent

---

## LLM Implementation Instructions

### CRITICAL RULES FOR TEST GENERATION

**ABSOLUTE CONSTRAINTS:**
1. **NEVER modify any file in `component/src/` directory**
2. **NEVER create new files in `component/src/` directory**
3. **CREATE all test files ONLY in `component/tests/` directory**
4. **TREAT `src/` as completely immutable and read-only**
5. **Tests must work with existing `src/` code AS-IS**

### Test Generation Process

```
WORKING DIRECTORY: component/tests/ (ALL files created here)
SOURCE DIRECTORY: component/src/ (READ-ONLY - never modified)

STEP 1: Analysis Phase (READ from src/)
  - Analyze all code paths in src/
  - Identify all dependencies in src/
  - Map all state transitions in src/
  - List all error conditions in src/
  - Document all business rules in src/
  
  OUTPUT: Test plan (saved in tests/test-plan)

STEP 2: Test Planning
  - Create test matrix for src/ code coverage
  - Plan mocks for src/ dependencies
  - Design test data for src/ inputs
  - Define 100% coverage goals for src/
  - Establish execution order
  
  OUTPUT: Test strategy (saved in tests/test-strategy)

STEP 3: Structure Creation (CREATE in tests/)
  - Create directory structure in tests/
  - Initialize test configuration in tests/
  - Set up test utilities in tests/helpers/
  - Configure coverage tools for src/
  - Prepare CI/CD pipeline in tests/

STEP 4: Test Implementation (ALL in tests/)
  1. Generate test helpers in tests/helpers/
  2. Create mocks in tests/mocks/
  3. Prepare fixtures in tests/fixtures/
  4. Implement unit tests in tests/unit/ (100% src/ coverage)
  5. Implement integration tests in tests/integration/
  6. Implement E2E tests in tests/e2e/
  7. Implement performance tests in tests/performance/
  8. Implement security tests in tests/security/

STEP 5: Verification
  - Verify 100% coverage of src/ code
  - Verify no modifications to src/
  - Verify all tests are in tests/
  - Verify test independence
  - Verify execution speed
```

### Import Patterns for Tests

```
CORRECT IMPORT EXAMPLES:
  From test file: tests/unit/example.test
  Import source: import from '../../src/module'
  Import helper: import from '../helpers/utility'
  Import mock: import from '../mocks/service'
  Import fixture: import from '../fixtures/data'

INCORRECT (NEVER DO):
  ❌ Modify: src/module (to make test pass)
  ❌ Create: src/test-helper
  ❌ Write: any file to src/
  ❌ Delete: any file from src/
  ❌ Move: files between src/ and tests/
```

### File System Operations

```
ALLOWED OPERATIONS:
  ✅ READ from component/src/*
  ✅ WRITE to component/tests/*
  ✅ CREATE in component/tests/*
  ✅ DELETE in component/tests/* (during cleanup)

FORBIDDEN OPERATIONS:
  ❌ WRITE to component/src/*
  ❌ CREATE in component/src/*
  ❌ DELETE from component/src/*
  ❌ MODIFY component/src/*
  ❌ RENAME in component/src/*
```

### Test Completeness Checklist

Before considering the test suite complete:

- [ ] 100% line coverage achieved
- [ ] 100% branch coverage achieved
- [ ] 100% path coverage achieved
- [ ] All error conditions tested
- [ ] All edge cases covered
- [ ] All security vulnerabilities tested
- [ ] All performance requirements validated
- [ ] All integrations tested
- [ ] All user journeys tested
- [ ] All test are independent
- [ ] All tests are deterministic
- [ ] All tests execute quickly
- [ ] All tests have clear names
- [ ] All test data is isolated
- [ ] All mocks are complete
- [ ] All helpers are implemented
- [ ] No anti-patterns present
- [ ] Quality gates passed

---

## Continuous Improvement

### Test Evolution Strategy

1. **Mutation Testing**: Continuously verify test effectiveness
2. **Fault Injection**: Test system resilience
3. **Chaos Engineering**: Test production readiness
4. **Property-Based Testing**: Discover edge cases automatically
5. **Model-Based Testing**: Generate tests from specifications
6. **Contract Testing**: Verify API contracts
7. **Snapshot Testing**: Detect unexpected changes
8. **Visual Regression**: Detect UI changes

---

## Final Validation

The test suite is considered complete only when:

1. **Coverage**: 100% code coverage achieved
2. **Quality**: All quality gates passed
3. **Performance**: Execution time within limits
4. **Reliability**: Zero flaky tests
5. **Maintainability**: Low complexity, high clarity
6. **Independence**: Complete test isolation
7. **Determinism**: 100% reproducible results
8. **Completeness**: All scenarios covered

Remember: A professional test suite with 100% coverage is not optional—it's the foundation of reliable software.

---

*Version: 2.0.0 | Industry Standards Compliant | Language Agnostic | 100% Coverage Mandatory*
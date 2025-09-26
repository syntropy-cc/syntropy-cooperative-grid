# LLM Testing Rules - Professional Test Suite Implementation Guide

## Executive Summary

This document provides comprehensive, language-agnostic rules for LLMs to implement a complete, professional-grade test suite achieving 100% code coverage. These rules follow industry best practices and standards including ISO/IEC/IEEE 29119, ISTQB guidelines, and proven testing methodologies.

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

```
tests/
├── unit/                 # Isolated component tests
│   ├── core/            # Core business logic tests
│   └── utilities/       # Utility function tests
├── integration/         # Component interaction tests
│   ├── api/            # External API integration tests
│   └── database/       # Data persistence tests
├── e2e/                # End-to-end user journeys
│   └── scenarios/      # Complete workflow tests
├── performance/        # Performance validation
├── security/          # Security verification
├── fixtures/          # Test data
│   ├── valid/        # Valid input scenarios
│   └── invalid/      # Invalid input scenarios
├── mocks/            # Test doubles
└── helpers/          # Test utilities
```

---

## Test Implementation Layers

### Layer 1: Unit Tests (`tests/unit/`)

#### Purpose
Validate individual components in complete isolation with 100% coverage of all logic paths.

#### Implementation Requirements

##### Test Coverage Targets
- **Line Coverage**: 100%
- **Branch Coverage**: 100%
- **Path Coverage**: 100%
- **Condition Coverage**: 100%
- **Modified Condition/Decision Coverage (MC/DC)**: 100%

##### Test Structure Pattern
```
TEST SUITE: [Component Name]
  TEST CONTEXT: [Method/Function Name]
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

### Layer 2: Integration Tests (`tests/integration/`)

#### Purpose
Verify correct interaction between components and with external systems while maintaining 100% interface coverage.

#### Implementation Requirements

##### Integration Points Coverage (100% Required)
```
INTEGRATION SUITE: [Integration Name]
  CATEGORY: Component Integration
    TEST: Component A → Component B data flow
    TEST: Component B → Component A feedback
    TEST: Error propagation between components
    TEST: Transaction boundaries
    
  CATEGORY: External System Integration
    TEST: Database connectivity
    TEST: Database transactions (commit/rollback)
    TEST: External API communication
    TEST: Message queue operations
    TEST: Cache synchronization
    
  CATEGORY: Contract Testing
    TEST: Request schema validation
    TEST: Response schema validation
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

### Layer 3: End-to-End Tests (`tests/e2e/`)

#### Purpose
Validate complete user workflows and system behavior from entry to exit points.

#### Implementation Requirements

##### E2E Scenario Coverage
```
E2E SUITE: [Application Name]
  CRITICAL PATHS: (Must have 100% coverage)
    SCENARIO: User Registration Flow
      STEP 1: Navigate to registration
      STEP 2: Fill registration form
      STEP 3: Submit and verify
      STEP 4: Email confirmation
      STEP 5: Account activation
      
    SCENARIO: Authentication Flow
      STEP 1: Login with credentials
      STEP 2: Multi-factor authentication
      STEP 3: Session establishment
      STEP 4: Session refresh
      STEP 5: Logout process
      
    SCENARIO: Core Business Process
      [Define steps for main business workflow]
      
  ALTERNATIVE PATHS:
    SCENARIO: Error Recovery
    SCENARIO: Concurrent Users
    SCENARIO: System Degradation
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

## Test Data Management (`tests/fixtures/`)

### Test Data Categories

#### 1. Valid Data Sets
```
valid/
├── minimal/          # Minimum required fields
├── complete/         # All fields populated
├── typical/          # Common use cases
├── maximum/          # Maximum allowed values
└── variations/       # Format variations
```

#### 2. Invalid Data Sets
```
invalid/
├── missing-required/ # Missing mandatory fields
├── type-errors/     # Wrong data types
├── constraint-violations/ # Business rule violations
├── malformed/       # Structural errors
└── malicious/       # Security test payloads
```

#### 3. Edge Case Data Sets
```
edge-cases/
├── boundary-values/ # Min, Max, Zero
├── special-characters/ # Unicode, Emoji, Control
├── large-datasets/  # Volume testing
├── concurrent/      # Race condition testing
└── temporal/        # Time-based edge cases
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

## Mock Implementation (`tests/mocks/`)

### Mock Types Required

#### 1. Test Doubles Taxonomy
- **Dummy**: Objects passed but never used
- **Stub**: Provides predetermined responses
- **Spy**: Records interactions for verification
- **Mock**: Pre-programmed with expectations
- **Fake**: Working implementation for testing

#### 2. Mock Implementation Requirements

```
MOCK: [Service Name]
  CAPABILITIES:
    - Record all method calls
    - Verify call count and parameters
    - Simulate success responses
    - Simulate error conditions
    - Simulate timeouts
    - Simulate rate limiting
    - Simulate partial failures
    
  VERIFICATION METHODS:
    - Was called with specific parameters
    - Was called N times
    - Was called in specific order
    - Was not called
    
  CONFIGURATION:
    - Response delays
    - Error rates
    - Custom responses
    - Conditional behavior
```

### Mock Best Practices

1. **Interface Compliance**: Match real service interface exactly
2. **Behavior Simulation**: Realistic response patterns
3. **State Management**: Maintain state between calls
4. **Reset Capability**: Clear state between tests
5. **Error Injection**: Controllable failure modes
6. **Performance Characteristics**: Simulate latency
7. **Verification API**: Rich assertion capabilities

---

## Test Helpers (`tests/helpers/`)

### Required Helper Categories

#### 1. Setup Utilities
- Environment configuration
- Test database initialization
- Test server startup
- Service connections
- Authentication setup
- Test user creation

#### 2. Assertion Helpers
- Custom matchers
- Deep equality checks
- Schema validation
- Approximate equality
- Collection assertions
- Async assertions

#### 3. Wait/Retry Utilities
- Wait for condition
- Retry with backoff
- Timeout handling
- Polling mechanisms
- Event waiting
- Promise utilities

#### 4. Data Builders
- Object mothers
- Test data builders
- Random generators
- Factory methods
- Fixture loaders
- Snapshot creators

#### 5. Cleanup Utilities
- Database cleanup
- File system cleanup
- Network cleanup
- Process cleanup
- Memory cleanup
- State reset

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

```
COVERAGE REQUIREMENTS:
  Code Coverage:
    - Line Coverage: 100%
    - Branch Coverage: 100%
    - Function Coverage: 100%
    - Statement Coverage: 100%
    
  Logic Coverage:
    - Decision Coverage: 100%
    - Condition Coverage: 100%
    - MC/DC Coverage: 100%
    - Path Coverage: 100%
    
  Data Flow Coverage:
    - All-Defs: 100%
    - All-Uses: 100%
    - All-P-Uses: 100%
    - All-C-Uses: 100%
    
  Mutation Coverage:
    - Mutation Score: 100%
    - Killed Mutants: 100%
```

### Coverage Exclusions (Must Document)
- Generated code (with justification)
- Third-party libraries
- Unreachable code (with proof)

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

### Test Generation Process

```
STEP 1: Analysis Phase
  - Analyze all code paths
  - Identify all dependencies
  - Map all state transitions
  - List all error conditions
  - Document all business rules

STEP 2: Test Planning
  - Create test matrix (input × conditions × outputs)
  - Plan mock requirements
  - Design test data sets
  - Define coverage goals (100% mandatory)
  - Establish execution order

STEP 3: Structure Creation
  - Create directory structure
  - Initialize test configuration
  - Set up test utilities
  - Configure coverage tools
  - Prepare CI/CD pipeline

STEP 4: Test Implementation (Order is Critical)
  1. Generate test helpers and utilities
  2. Create mock implementations
  3. Prepare test fixtures
  4. Implement unit tests (100% coverage)
  5. Implement integration tests
  6. Implement E2E tests
  7. Implement performance tests
  8. Implement security tests

STEP 5: Verification
  - Verify 100% code coverage
  - Verify test independence
  - Verify execution speed
  - Verify determinism
  - Verify documentation completeness
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
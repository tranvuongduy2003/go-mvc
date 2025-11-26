# AI Coding Standards for Go MVC

## 🎯 Core Principle: Self-Documenting Code

Code should be **self-explanatory** through clear naming and structure. Comments are a **last resort**, not a first choice.

---

## 📝 Comment Policy: MINIMAL COMMENTS ONLY

### ❌ NEVER Comment For

**1. Obvious Code**
```go
// ❌ BAD - Comment states the obvious
// Create a new user
user := NewUser(email, password)

// ✅ GOOD - No comment needed, code is clear
user := NewUser(email, password)
```

**2. What Code Does**
```go
// ❌ BAD - Describing what the code does
// Loop through all users and check if they are active
for _, user := range users {
    if user.IsActive() {
        activeUsers = append(activeUsers, user)
    }
}

// ✅ GOOD - Let the code speak
for _, user := range users {
    if user.IsActive() {
        activeUsers = append(activeUsers, user)
    }
}
```

**3. Type Information**
```go
// ❌ BAD - Restating type information
// UserRepository is a repository for users
type UserRepository interface {
    Create(ctx context.Context, user *User) error
}

// ✅ GOOD - Type name is self-explanatory
type UserRepository interface {
    Create(ctx context.Context, user *User) error
}
```

**4. Function Signatures**
```go
// ❌ BAD - Duplicating function signature
// GetByID gets a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    // ...
}

// ✅ GOOD - Function name explains itself
func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    // ...
}
```

**5. Variable Declarations**
```go
// ❌ BAD - Explaining variable names
// userID is the ID of the user
var userID string

// email stores the user's email
var email string

// ✅ GOOD - Variable names are clear
var userID string
var email string
```

**6. Import Statements**
```go
// ❌ BAD - Commenting imports
// Standard library
import (
    "context"
    "time"
)

// Third-party
import (
    "github.com/gin-gonic/gin"
)

// ✅ GOOD - Imports are self-explanatory, group them logically
import (
    "context"
    "time"
    
    "github.com/gin-gonic/gin"
)
```

**7. Return Statements**
```go
// ❌ BAD - Explaining returns
// Return the user
return user

// Return error if validation fails
return nil, err

// ✅ GOOD - Return statements are clear from context
return user
return nil, err
```

**8. Error Handling**
```go
// ❌ BAD - Stating the obvious
// Check if there's an error
if err != nil {
    // Return the error
    return err
}

// ✅ GOOD - Error handling is standard Go
if err != nil {
    return err
}
```

**9. TODO Comments**
```go
// ❌ BAD - Vague TODOs
// TODO: fix this
// TODO: improve performance
// TODO: refactor

// ✅ GOOD - If you must use TODO, create a GitHub issue instead
// Then reference it in code only if absolutely necessary
```

**10. Commented-Out Code**
```go
// ❌ BAD - Dead code
// oldUser := GetOldUser(id)
// if oldUser != nil {
//     return oldUser
// }

// ✅ GOOD - Delete unused code, use version control
```

---

### ✅ ONLY Comment When

**1. Complex Business Logic That Cannot Be Expressed in Code**
```go
// Apply 3% compound interest daily with 30/360 day count convention
// as per financial regulation XYZ-2024 section 4.2
interest := principal * math.Pow(1.03, days/360.0)
```

**2. Non-Obvious Algorithms or Mathematical Formulas**
```go
// Haversine formula for calculating distance between two points on Earth
a := math.Sin(dLat/2)*math.Sin(dLat/2) +
     math.Cos(lat1Rad)*math.Cos(lat2Rad)*
     math.Sin(dLon/2)*math.Sin(dLon/2)
```

**3. Workarounds for External Libraries or Known Bugs**
```go
// WORKAROUND: GORM has a bug with preloading nested associations
// See: https://github.com/go-gorm/gorm/issues/12345
// Remove this when bug is fixed in GORM v2.0
db.Preload("User").Preload("User.Profile").Find(&orders)
```

**4. Performance-Critical Optimizations**
```go
// Pre-allocate slice with exact capacity to avoid reallocation
// Benchmarked: 40% faster than append for 10K+ items
users := make([]*User, 0, expectedCount)
```

**5. Security-Sensitive Code**
```go
// Constant-time comparison prevents timing attacks
// DO NOT replace with simple == comparison
if subtle.ConstantTimeCompare([]byte(hash1), []byte(hash2)) == 1 {
    // ...
}
```

**6. Regulatory or Compliance Requirements**
```go
// GDPR Article 17: Right to erasure
// Data must be anonymized, not just soft-deleted
user.Email = hashEmail(user.Email)
user.Name = "REDACTED"
user.DeletedAt = time.Now()
```

**7. Counterintuitive Code Required by External Constraints**
```go
// API requires sending timestamp in milliseconds, not seconds
// Converting Go's time.Now() (nanoseconds) to milliseconds
timestamp := time.Now().UnixNano() / 1e6
```

---

## 🏗️ Self-Documenting Code Guidelines

### 1. Use Descriptive Names

**Functions**
```go
// ❌ BAD
func Process(d []string) error

// ✅ GOOD
func ValidateEmailAddresses(emails []string) error
```

**Variables**
```go
// ❌ BAD
var t time.Time
var n int
var d []string

// ✅ GOOD
var expirationTime time.Time
var userCount int
var emailAddresses []string
```

**Constants**
```go
// ❌ BAD
const MAX = 100
const T = 30

// ✅ GOOD
const MaxLoginAttempts = 100
const SessionTimeoutSeconds = 30
```

### 2. Use Type System for Documentation

**Custom Types**
```go
// ❌ BAD
func TransferMoney(from, to string, amount float64) error

// ✅ GOOD
type AccountID string
type Money struct {
    Amount   int64
    Currency string
}

func TransferMoney(from, to AccountID, amount Money) error
```

**Enums**
```go
// ❌ BAD
const (
    STATUS_1 = 1
    STATUS_2 = 2
)

// ✅ GOOD
type UserStatus int

const (
    UserStatusActive UserStatus = iota
    UserStatusInactive
    UserStatusSuspended
)
```

### 3. Small, Focused Functions

```go
// ❌ BAD - Large function needs comments
func ProcessUser(user User) error {
    // Validate user
    if user.Email == "" {
        return errors.New("email required")
    }
    
    // Hash password
    hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
    if err != nil {
        return err
    }
    user.PasswordHash = string(hash)
    
    // Save to database
    if err := db.Save(&user); err != nil {
        return err
    }
    
    // Send welcome email
    if err := sendEmail(user.Email, "Welcome!"); err != nil {
        return err
    }
    
    return nil
}

// ✅ GOOD - Small functions are self-documenting
func ProcessUser(user User) error {
    if err := ValidateUser(user); err != nil {
        return err
    }
    
    if err := HashUserPassword(&user); err != nil {
        return err
    }
    
    if err := SaveUser(user); err != nil {
        return err
    }
    
    return SendWelcomeEmail(user.Email)
}
```

### 4. Guard Clauses Instead of Nested Ifs

```go
// ❌ BAD - Nested ifs need comments
func Process(user *User) error {
    if user != nil {
        if user.IsActive() {
            if user.HasPermission("write") {
                // Do something
                return DoWork(user)
            }
        }
    }
    return errors.New("invalid user")
}

// ✅ GOOD - Guard clauses are clear
func Process(user *User) error {
    if user == nil {
        return errors.New("user is nil")
    }
    
    if !user.IsActive() {
        return errors.New("user is inactive")
    }
    
    if !user.HasPermission("write") {
        return errors.New("user lacks permission")
    }
    
    return DoWork(user)
}
```

### 5. Meaningful Error Messages

```go
// ❌ BAD
return errors.New("error")

// ✅ GOOD
return errors.New("user email already exists")
return fmt.Errorf("failed to connect to database: %w", err)
```

---

## 🚫 Banned Comment Patterns

### Never Use These

1. **Separators**
```go
// ❌ BAD
// ==========================================
// User Functions
// ==========================================
```

2. **Function Headers**
```go
// ❌ BAD
// Function: GetUser
// Parameters: id string
// Returns: User, error
// Description: Gets a user by ID
func GetUser(id string) (*User, error)
```

3. **Change Log**
```go
// ❌ BAD
// 2024-01-01: Added user validation - John
// 2024-01-15: Fixed bug - Jane
// 2024-02-01: Refactored - Bob
```

4. **Author Information**
```go
// ❌ BAD
// Author: John Doe
// Created: 2024-01-01
// Modified: 2024-01-15
```

5. **File Headers**
```go
// ❌ BAD
// File: user.go
// Package: domain
// Description: User domain entity
package domain
```

6. **Section Markers**
```go
// ❌ BAD
// ========== Public Methods ==========

// ========== Private Methods ==========

// ========== Helper Functions ==========
```

---

## 📋 Code Review Checklist

### Before Submitting Code

- [ ] No comments explaining what the code does
- [ ] No comments restating type information
- [ ] No commented-out code
- [ ] No TODO comments (create GitHub issues instead)
- [ ] Function names clearly express intent
- [ ] Variable names are descriptive
- [ ] Complex logic is broken into small functions
- [ ] Only essential comments remain (algorithm, security, compliance)

### Reviewing Code

If you see a comment, ask:
1. Can this be expressed through better naming?
2. Can this be expressed through code structure?
3. Can this be expressed through types?
4. Is this truly non-obvious or just lazy naming?

**If yes to any of the above, remove the comment and improve the code.**

---

## 🎓 Examples: Before and After

### Example 1: User Creation

**❌ BEFORE (Too Many Comments)**
```go
// CreateUser creates a new user in the system
// It validates the input, hashes the password, and saves to database
func CreateUser(email, password string) (*User, error) {
    // Validate email format
    if !isValidEmail(email) {
        return nil, errors.New("invalid email")
    }
    
    // Check password strength
    if len(password) < 8 {
        return nil, errors.New("password too short")
    }
    
    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        // Return error if hashing fails
        return nil, err
    }
    
    // Create user object
    user := &User{
        Email:    email,
        Password: string(hashedPassword),
    }
    
    // Save to database
    if err := db.Save(user); err != nil {
        // Return error if save fails
        return nil, err
    }
    
    // Return the created user
    return user, nil
}
```

**✅ AFTER (Self-Documenting)**
```go
func CreateUser(email, password string) (*User, error) {
    if err := ValidateEmail(email); err != nil {
        return nil, err
    }
    
    if err := ValidatePasswordStrength(password); err != nil {
        return nil, err
    }
    
    hashedPassword, err := HashPassword(password)
    if err != nil {
        return nil, err
    }
    
    user := &User{
        Email:        email,
        PasswordHash: hashedPassword,
    }
    
    return user, SaveUser(user)
}
```

### Example 2: Complex Business Logic

**❌ BEFORE (Wrong Type of Comments)**
```go
// Calculate discount
func CalculateDiscount(amount float64, tier int) float64 {
    // Check tier level
    if tier == 1 {
        // 10% discount for tier 1
        return amount * 0.9
    } else if tier == 2 {
        // 15% discount for tier 2
        return amount * 0.85
    } else if tier == 3 {
        // 20% discount for tier 3
        return amount * 0.8
    }
    // No discount for other tiers
    return amount
}
```

**✅ AFTER (Self-Documenting with Named Constants)**
```go
type CustomerTier int

const (
    TierBasic CustomerTier = iota
    TierSilver
    TierGold
)

const (
    BasicDiscountRate  = 0.10
    SilverDiscountRate = 0.15
    GoldDiscountRate   = 0.20
)

func CalculateDiscount(amount Money, tier CustomerTier) Money {
    discountRate := GetDiscountRate(tier)
    return amount.Multiply(1 - discountRate)
}

func GetDiscountRate(tier CustomerTier) float64 {
    switch tier {
    case TierBasic:
        return BasicDiscountRate
    case TierSilver:
        return SilverDiscountRate
    case TierGold:
        return GoldDiscountRate
    default:
        return 0
    }
}
```

---

## 🤖 AI Generation Rules

When generating code:

1. **NEVER add comments** unless it meets the "ONLY Comment When" criteria
2. **Use descriptive names** that make comments unnecessary
3. **Break down complex functions** into smaller, self-explanatory ones
4. **Use type system** for documentation
5. **Prefer code clarity** over brevity
6. **If you think a comment is needed**, first try to:
   - Rename the function
   - Extract a method
   - Use a custom type
   - Create a constant

### AI Prompt Template

```
Generate Go code following these rules:
1. NO COMMENTS except for:
   - Complex algorithms that cannot be expressed in code
   - Workarounds for external library bugs
   - Security-sensitive code
   - Regulatory/compliance requirements
   
2. Use self-documenting code:
   - Descriptive function names
   - Clear variable names
   - Small, focused functions
   - Guard clauses
   - Custom types for clarity
   
3. If you're tempted to add a comment, improve the code instead
```

---

## 📊 Metrics

### Code Quality Indicators

**Good Code Has:**
- ✅ Function names that explain what they do
- ✅ Variable names that explain what they store
- ✅ Type names that explain what they represent
- ✅ < 1% of lines are comments
- ✅ No commented-out code
- ✅ No TODO comments

**Bad Code Has:**
- ❌ Comments explaining obvious code
- ❌ Comments restating type information
- ❌ Comments describing what code does
- ❌ > 5% of lines are comments
- ❌ Commented-out code
- ❌ TODO comments

---

## 🎯 Remember

> "Good code is its own best documentation. As you're about to add a comment, ask yourself, 'How can I improve the code so that this comment isn't needed?'"
> 
> — Steve McConnell, Code Complete

**The best comment is no comment. The best documentation is clear code.**

---

## 📚 Further Reading

- Clean Code by Robert C. Martin
- Code Complete by Steve McConnell
- The Pragmatic Programmer by Andrew Hunt & David Thomas
- Effective Go: https://golang.org/doc/effective_go

---

**This document is now part of the project standards. All code must comply with these rules.**

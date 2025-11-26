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

### ✅ ONLY Comment When

**1. Complex Algorithms**
```go
// ✅ GOOD - Algorithm explanation
// Boyer-Moore string search algorithm for O(n/m) performance
func BoyerMooreSearch(text, pattern string) int {
    // Implementation
}
```

**2. Non-Obvious Business Rules**
```go
// ✅ GOOD - Business rule explanation
// Discount only applies to orders over $100 placed on weekdays
// as per marketing campaign requirements (TICKET-123)
if order.Total > 100 && isWeekday(order.Date) {
    discount = order.Total * 0.1
}
```

**3. External Constraints**
```go
// ✅ GOOD - External system constraint
// API requires ISO 8601 format with milliseconds
// See: https://api.example.com/docs\#date-format
timestamp := time.Now().Format("2006-01-02T15:04:05.000Z07:00")
```

---

## 🏗️ Self-Documenting Code Guidelines

### 1. Use Descriptive Names

```go
// ❌ BAD
func p(u *U, e string) error {
    // ...
}

// ✅ GOOD
func SendVerificationEmail(user *User, emailAddress string) error {
    // ...
}
```

### 2. Functions Should Do One Thing

```go
// ❌ BAD - Function does too much
func ProcessUser(user *User) error {
    // Validate user
    // Save to database
    // Send email
    // Update cache
    // Log activity
}

// ✅ GOOD - Small, focused functions
func ValidateUser(user *User) error { }
func SaveUser(user *User) error { }
func SendWelcomeEmail(user *User) error { }
func UpdateUserCache(user *User) error { }
func LogUserActivity(user *User) error { }
```

### 3. Use Guard Clauses

```go
// ❌ BAD - Nested conditionals
func ProcessOrder(order *Order) error {
    if order != nil {
        if order.IsValid() {
            if order.HasItems() {
                // Process order
            }
        }
    }
}

// ✅ GOOD - Early returns
func ProcessOrder(order *Order) error {
    if order == nil {
        return errors.New("order is nil")
    }
    if !order.IsValid() {
        return errors.New("order is invalid")
    }
    if !order.HasItems() {
        return errors.New("order has no items")
    }
    
    // Process order
    return nil
}
```

---

## 🚫 AI Work Standards

### DO NOT Create Documentation Files After Implementation

**❌ NEVER DO:**
- Create summary markdown files after completing tasks
- Generate IMPLEMENTATION_COMPLETE.md files
- Create TASK_SUMMARY.md or similar files
- Generate architecture visualization files
- Create checklists or status files

**✅ INSTEAD DO:**
- Update existing documentation if needed
- Add to README.md if feature is significant
- Update CHANGELOG.md for version tracking
- Only create documentation files if explicitly requested by user

### Examples

**❌ BAD:**
```
After implementing feature X:
1. Create FEATURE_X_IMPLEMENTATION.md
2. Create ARCHITECTURE_UPDATE.md
3. Create CHECKLIST.md
4. Create SUMMARY.md
```

**✅ GOOD:**
```
After implementing feature X:
1. Update README.md with new feature
2. Add entry to CHANGELOG.md
3. Update relevant existing docs if needed
4. Report completion to user without creating files
```

**Rule**: If user says "don't create summary files", respect it permanently. Focus on code and essential documentation only.

---

## 📊 Code Quality Metrics

### Good Code Has:
- ✅ Function names that explain what they do
- ✅ Variable names that explain what they store
- ✅ Type names that explain what they represent
- ✅ < 1% of lines are comments
- ✅ No commented-out code
- ✅ No TODO comments

### Bad Code Has:
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

**This document is now part of the project standards. All code must comply with these rules.**

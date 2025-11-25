package valueobject

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Email represents an email value object
type Email struct {
	value string
}

// NewEmail creates a new email value object
func NewEmail(email string) (Email, error) {
	if err := validateEmail(email); err != nil {
		return Email{}, err
	}

	return Email{value: strings.ToLower(strings.TrimSpace(email))}, nil
}

// Value returns the email value
func (e Email) Value() string {
	return e.value
}

// String returns the string representation
func (e Email) String() string {
	return e.value
}

// IsEmpty checks if email is empty
func (e Email) IsEmpty() bool {
	return e.value == ""
}

// Domain returns the domain part of the email
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// LocalPart returns the local part of the email
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// validateEmail validates email format
func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	email = strings.TrimSpace(email)
	if len(email) > 254 {
		return fmt.Errorf("email is too long")
	}

	// Basic email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// Money represents a money value object
type Money struct {
	amount   int64 // stored in smallest currency unit (e.g., cents)
	currency string
}

// NewMoney creates a new money value object
func NewMoney(amount int64, currency string) (Money, error) {
	if currency == "" {
		return Money{}, fmt.Errorf("currency cannot be empty")
	}

	currency = strings.ToUpper(strings.TrimSpace(currency))
	if len(currency) != 3 {
		return Money{}, fmt.Errorf("currency must be 3 characters")
	}

	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

// Amount returns the amount in smallest currency unit
func (m Money) Amount() int64 {
	return m.amount
}

// Currency returns the currency code
func (m Money) Currency() string {
	return m.currency
}

// AmountAsFloat returns the amount as float (divided by 100 for most currencies)
func (m Money) AmountAsFloat() float64 {
	return float64(m.amount) / 100.0
}

// String returns string representation
func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.AmountAsFloat(), m.currency)
}

// Add adds two money values (must be same currency)
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot add different currencies")
	}

	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

// Subtract subtracts two money values (must be same currency)
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot subtract different currencies")
	}

	return Money{
		amount:   m.amount - other.amount,
		currency: m.currency,
	}, nil
}

// IsPositive checks if amount is positive
func (m Money) IsPositive() bool {
	return m.amount > 0
}

// IsZero checks if amount is zero
func (m Money) IsZero() bool {
	return m.amount == 0
}

// IsNegative checks if amount is negative
func (m Money) IsNegative() bool {
	return m.amount < 0
}

// PhoneNumber represents a phone number value object
type PhoneNumber struct {
	value string
}

// NewPhoneNumber creates a new phone number value object
func NewPhoneNumber(phone string) (PhoneNumber, error) {
	if err := validatePhoneNumber(phone); err != nil {
		return PhoneNumber{}, err
	}

	// Normalize phone number (remove spaces, dashes, etc.)
	normalized := normalizePhoneNumber(phone)

	return PhoneNumber{value: normalized}, nil
}

// Value returns the phone number value
func (p PhoneNumber) Value() string {
	return p.value
}

// String returns the string representation
func (p PhoneNumber) String() string {
	return p.value
}

// validatePhoneNumber validates phone number format
func validatePhoneNumber(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number cannot be empty")
	}

	// Remove all non-digit characters for validation
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	if len(digits) < 7 || len(digits) > 15 {
		return fmt.Errorf("phone number must be between 7 and 15 digits")
	}

	return nil
}

// normalizePhoneNumber normalizes phone number format
func normalizePhoneNumber(phone string) string {
	// Remove all non-digit characters except + at the beginning
	if strings.HasPrefix(phone, "+") {
		digits := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")
		return digits
	}

	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	return digits
}

// URL represents a URL value object
type URL struct {
	value string
}

// NewURL creates a new URL value object
func NewURL(url string) (URL, error) {
	if err := validateURL(url); err != nil {
		return URL{}, err
	}

	return URL{value: strings.TrimSpace(url)}, nil
}

// Value returns the URL value
func (u URL) Value() string {
	return u.value
}

// String returns the string representation
func (u URL) String() string {
	return u.value
}

// validateURL validates URL format
func validateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	url = strings.TrimSpace(url)
	if len(url) > 2048 {
		return fmt.Errorf("URL is too long")
	}

	// Basic URL validation
	urlRegex := regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`)
	if !urlRegex.MatchString(url) {
		return fmt.Errorf("invalid URL format")
	}

	return nil
}

// Address represents an address value object
type Address struct {
	street     string
	city       string
	state      string
	country    string
	postalCode string
}

// NewAddress creates a new address value object
func NewAddress(street, city, state, country, postalCode string) (Address, error) {
	street = strings.TrimSpace(street)
	city = strings.TrimSpace(city)
	state = strings.TrimSpace(state)
	country = strings.TrimSpace(country)
	postalCode = strings.TrimSpace(postalCode)

	if street == "" {
		return Address{}, fmt.Errorf("street cannot be empty")
	}
	if city == "" {
		return Address{}, fmt.Errorf("city cannot be empty")
	}
	if country == "" {
		return Address{}, fmt.Errorf("country cannot be empty")
	}

	return Address{
		street:     street,
		city:       city,
		state:      state,
		country:    country,
		postalCode: postalCode,
	}, nil
}

// Street returns the street
func (a Address) Street() string {
	return a.street
}

// City returns the city
func (a Address) City() string {
	return a.city
}

// State returns the state
func (a Address) State() string {
	return a.state
}

// Country returns the country
func (a Address) Country() string {
	return a.country
}

// PostalCode returns the postal code
func (a Address) PostalCode() string {
	return a.postalCode
}

// String returns the full address as string
func (a Address) String() string {
	parts := []string{a.street, a.city}
	if a.state != "" {
		parts = append(parts, a.state)
	}
	if a.postalCode != "" {
		parts = append(parts, a.postalCode)
	}
	parts = append(parts, a.country)

	return strings.Join(parts, ", ")
}

// AuditLog represents audit information for entities
type AuditLog struct {
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewAuditLog creates a new audit log
func NewAuditLog(userID uuid.UUID) AuditLog {
	now := time.Now()
	return AuditLog{
		CreatedBy: userID,
		UpdatedBy: userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update updates the audit log with new user and time
func (a *AuditLog) Update(userID uuid.UUID) {
	a.UpdatedBy = userID
	a.UpdatedAt = time.Now()
}

package valueobject

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {
	if err := validateEmail(email); err != nil {
		return Email{}, err
	}

	return Email{value: strings.ToLower(strings.TrimSpace(email))}, nil
}

func (e Email) Value() string {
	return e.value
}

func (e Email) String() string {
	return e.value
}

func (e Email) IsEmpty() bool {
	return e.value == ""
}

func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	email = strings.TrimSpace(email)
	if len(email) > 254 {
		return fmt.Errorf("email is too long")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

type Money struct {
	amount   int64 // stored in smallest currency unit (e.g., cents)
	currency string
}

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

func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) AmountAsFloat() float64 {
	return float64(m.amount) / 100.0
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.AmountAsFloat(), m.currency)
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot add different currencies")
	}

	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot subtract different currencies")
	}

	return Money{
		amount:   m.amount - other.amount,
		currency: m.currency,
	}, nil
}

func (m Money) IsPositive() bool {
	return m.amount > 0
}

func (m Money) IsZero() bool {
	return m.amount == 0
}

func (m Money) IsNegative() bool {
	return m.amount < 0
}

type PhoneNumber struct {
	value string
}

func NewPhoneNumber(phone string) (PhoneNumber, error) {
	if err := validatePhoneNumber(phone); err != nil {
		return PhoneNumber{}, err
	}

	normalized := normalizePhoneNumber(phone)

	return PhoneNumber{value: normalized}, nil
}

func (p PhoneNumber) Value() string {
	return p.value
}

func (p PhoneNumber) String() string {
	return p.value
}

func validatePhoneNumber(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number cannot be empty")
	}

	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	if len(digits) < 7 || len(digits) > 15 {
		return fmt.Errorf("phone number must be between 7 and 15 digits")
	}

	return nil
}

func normalizePhoneNumber(phone string) string {
	if strings.HasPrefix(phone, "+") {
		digits := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")
		return digits
	}

	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	return digits
}

type URL struct {
	value string
}

func NewURL(url string) (URL, error) {
	if err := validateURL(url); err != nil {
		return URL{}, err
	}

	return URL{value: strings.TrimSpace(url)}, nil
}

func (u URL) Value() string {
	return u.value
}

func (u URL) String() string {
	return u.value
}

func validateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	url = strings.TrimSpace(url)
	if len(url) > 2048 {
		return fmt.Errorf("URL is too long")
	}

	urlRegex := regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`)
	if !urlRegex.MatchString(url) {
		return fmt.Errorf("invalid URL format")
	}

	return nil
}

type Address struct {
	street     string
	city       string
	state      string
	country    string
	postalCode string
}

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

func (a Address) Street() string {
	return a.street
}

func (a Address) City() string {
	return a.city
}

func (a Address) State() string {
	return a.state
}

func (a Address) Country() string {
	return a.country
}

func (a Address) PostalCode() string {
	return a.postalCode
}

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

type AuditLog struct {
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewAuditLog(userID uuid.UUID) AuditLog {
	now := time.Now()
	return AuditLog{
		CreatedBy: userID,
		UpdatedBy: userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (a *AuditLog) Update(userID uuid.UUID) {
	a.UpdatedBy = userID
	a.UpdatedAt = time.Now()
}

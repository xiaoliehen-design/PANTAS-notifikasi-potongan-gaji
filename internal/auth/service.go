package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	mailaddress "net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/bcpriok/pantas/internal/config"
	"github.com/bcpriok/pantas/internal/mailer"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidCredentials  = errors.New("NIP/username atau password salah")
	ErrRateLimited         = errors.New("terlalu banyak percobaan; coba kembali beberapa menit lagi")
	ErrUnauthenticated     = errors.New("sesi tidak valid atau telah berakhir")
	ErrInvalidCSRF         = errors.New("token keamanan tidak valid")
	ErrInvalidOTP          = errors.New("kode verifikasi salah atau sudah kedaluwarsa")
	ErrContactUnavailable  = errors.New("kontak belum tersedia atau belum diverifikasi")
	ErrDeliveryUnavailable = errors.New("layanan pengiriman kode verifikasi belum dikonfigurasi")
	ErrDeliveryFailed      = errors.New("kode verifikasi belum dapat dikirim; silakan coba lagi")
)

const (
	SessionCookieName  = "pantas_session"
	CSRFCookieName     = "pantas_csrf"
	TabProofCookieName = "pantas_tab_proof"
	TabTokenHeader     = "X-PANTAS-Tab-Token"
)

type Principal struct {
	ID                 string `json:"id"`
	AccountType        string `json:"account_type"`
	NIP                string `json:"nip"`
	Username           string `json:"username,omitempty"`
	Name               string `json:"name"`
	UnitID             string `json:"unit_id"`
	UnitName           string `json:"unit_name"`
	UnitType           string `json:"unit_type"`
	PositionRole       string `json:"position_role"`
	IsAdmin            bool   `json:"is_admin"`
	MustChangePassword bool   `json:"must_change_password"`
	Email              string `json:"email,omitempty"`
	EmailVerified      bool   `json:"email_verified"`
	Phone              string `json:"phone,omitempty"`
	PhoneVerified      bool   `json:"phone_verified"`
}

func (p Principal) IsSupervisor() bool {
	return p.IsAdmin || p.PositionRole == "section_head" || p.PositionRole == "division_head" || p.PositionRole == "office_head"
}

func (p Principal) LoginIdentifier() string {
	if p.AccountType == "admin" {
		return p.Username
	}
	return p.NIP
}

type Session struct {
	ID       string
	CSRFHash []byte
}

type LoginResult struct {
	Principal    Principal
	SessionToken string
	CSRFToken    string
	TabToken     string
}

type ContactInputError struct {
	message string
}

func (e *ContactInputError) Error() string {
	return e.message
}

func IsContactInputError(err error) bool {
	var inputError *ContactInputError
	return errors.As(err, &inputError)
}

func contactInputError(message string) error {
	return &ContactInputError{message: message}
}

type Service struct {
	pool     *pgxpool.Pool
	cfg      config.Config
	delivery *mailer.Worker
}

type contextKey int

const principalKey contextKey = 1

func New(pool *pgxpool.Pool, cfg config.Config, delivery *mailer.Worker) *Service {
	return &Service{pool: pool, cfg: cfg, delivery: delivery}
}

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalKey, principal)
}

func PrincipalFrom(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalKey).(Principal)
	return principal, ok
}

func (s *Service) Login(ctx context.Context, identifier, password, ip, userAgent string) (LoginResult, error) {
	identifier = strings.ToLower(strings.TrimSpace(identifier))
	if !validLoginIdentifier(identifier) || password == "" {
		s.recordLoginAttempt(ctx, identifier, ip, false)
		return LoginResult{}, ErrInvalidCredentials
	}
	limited, err := s.loginRateLimited(ctx, identifier, ip)
	if err != nil {
		return LoginResult{}, err
	}
	if limited {
		return LoginResult{}, ErrRateLimited
	}

	var principal Principal
	var passwordHash *string
	var emailVerifiedAt, phoneVerifiedAt *time.Time
	err = s.pool.QueryRow(ctx, `
		select u.id::text, 'user', u.nip, '', u.name, u.unit_id::text, un.name, un.unit_type,
		       u.position_role, false, u.must_change_password,
		       coalesce(u.email, ''), u.email_verified_at,
		       coalesce(u.phone_e164, ''), u.phone_verified_at, u.password_hash
		from public.users u
		join public.units un on un.id = u.unit_id
		where u.nip = $1 and u.is_active and u.deleted_at is null
		union all
		select aa.account_id::text, 'admin', '', aa.username, aa.name, '', 'Administrator Sistem', 'admin',
		       'admin', true, aa.must_change_password, '', null::timestamptz,
		       '', null::timestamptz, aa.password_hash
		from public.admin_accounts aa
		where aa.username = lower($1) and aa.is_active
		limit 1`, identifier).Scan(
		&principal.ID, &principal.AccountType, &principal.NIP, &principal.Username,
		&principal.Name, &principal.UnitID, &principal.UnitName,
		&principal.UnitType, &principal.PositionRole, &principal.IsAdmin,
		&principal.MustChangePassword, &principal.Email, &emailVerifiedAt,
		&principal.Phone, &phoneVerifiedAt, &passwordHash,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		s.recordLoginAttempt(ctx, identifier, ip, false)
		return LoginResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return LoginResult{}, err
	}
	principal.EmailVerified = emailVerifiedAt != nil
	principal.PhoneVerified = phoneVerifiedAt != nil

	valid, err := s.passwordMatches(ctx, principal.LoginIdentifier(), passwordHash, password)
	if err != nil {
		return LoginResult{}, err
	}
	if !valid {
		s.recordLoginAttempt(ctx, identifier, ip, false)
		return LoginResult{}, ErrInvalidCredentials
	}

	sessionToken, sessionHash, err := randomToken(32)
	if err != nil {
		return LoginResult{}, err
	}
	csrfToken, csrfHash, err := randomToken(32)
	if err != nil {
		return LoginResult{}, err
	}
	tabToken, _, err := randomToken(32)
	if err != nil {
		return LoginResult{}, err
	}

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return LoginResult{}, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `
		insert into public.sessions (user_id, token_hash, csrf_hash, ip_address, user_agent, expires_at)
		values ($1, $2, $3, nullif($4, '')::inet, left($5, 1000), now() + $6::interval)`,
		principal.ID, sessionHash, csrfHash, ip, userAgent, s.cfg.SessionTTL.String()); err != nil {
		return LoginResult{}, err
	}
	if principal.AccountType == "admin" {
		if _, err := tx.Exec(ctx, `update public.admin_accounts set last_login_at = now() where account_id = $1`, principal.ID); err != nil {
			return LoginResult{}, err
		}
	} else {
		if _, err := tx.Exec(ctx, `update public.users set last_login_at = now() where id = $1`, principal.ID); err != nil {
			return LoginResult{}, err
		}
	}
	if _, err := tx.Exec(ctx, `
		insert into public.login_attempts (nip, ip_address, was_successful)
		values ($1, nullif($2, '')::inet, true)`, identifier, ip); err != nil {
		return LoginResult{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return LoginResult{}, err
	}
	return LoginResult{Principal: principal, SessionToken: sessionToken, CSRFToken: csrfToken, TabToken: tabToken}, nil
}

func (s *Service) AuthenticateRequest(ctx context.Context, request *http.Request) (Principal, Session, error) {
	if !s.verifyTabSession(request) {
		return Principal{}, Session{}, ErrUnauthenticated
	}
	cookie, err := request.Cookie(SessionCookieName)
	if err != nil || cookie.Value == "" {
		return Principal{}, Session{}, ErrUnauthenticated
	}
	raw, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil || len(raw) != 32 {
		return Principal{}, Session{}, ErrUnauthenticated
	}
	hash := sha256.Sum256(raw)
	var principal Principal
	var session Session
	var emailVerifiedAt, phoneVerifiedAt *time.Time
	err = s.pool.QueryRow(ctx, `
		select u.id::text, 'user', u.nip, '', u.name, u.unit_id::text, un.name, un.unit_type,
		       u.position_role, false, u.must_change_password,
		       coalesce(u.email, ''), u.email_verified_at,
		       coalesce(u.phone_e164, ''), u.phone_verified_at,
		       s.id::text, s.csrf_hash
		from public.sessions s
		join public.users u on u.id = s.user_id
		join public.units un on un.id = u.unit_id
		where s.token_hash = $1 and s.revoked_at is null and s.expires_at > now()
		  and s.last_seen_at > now() - $2::interval
		  and u.is_active and u.deleted_at is null
		union all
		select aa.account_id::text, 'admin', '', aa.username, aa.name, '', 'Administrator Sistem', 'admin',
		       'admin', true, aa.must_change_password, '', null::timestamptz,
		       '', null::timestamptz, s.id::text, s.csrf_hash
		from public.sessions s
		join public.admin_accounts aa on aa.account_id = s.user_id
		where s.token_hash = $1 and s.revoked_at is null and s.expires_at > now()
		  and s.last_seen_at > now() - $2::interval and aa.is_active
		limit 1`, hash[:], s.cfg.SessionIdleTTL.String()).Scan(
		&principal.ID, &principal.AccountType, &principal.NIP, &principal.Username,
		&principal.Name, &principal.UnitID, &principal.UnitName,
		&principal.UnitType, &principal.PositionRole, &principal.IsAdmin,
		&principal.MustChangePassword, &principal.Email, &emailVerifiedAt,
		&principal.Phone, &phoneVerifiedAt, &session.ID, &session.CSRFHash,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Principal{}, Session{}, ErrUnauthenticated
	}
	if err != nil {
		return Principal{}, Session{}, err
	}
	principal.EmailVerified = emailVerifiedAt != nil
	principal.PhoneVerified = phoneVerifiedAt != nil
	_, _ = s.pool.Exec(ctx, `
		update public.sessions set last_seen_at = now()
		where id = $1 and last_seen_at < now() - interval '5 minutes'`, session.ID)
	return principal, session, nil
}

func (s *Service) VerifyCSRF(request *http.Request, session Session) error {
	header := request.Header.Get("X-CSRF-Token")
	cookie, err := request.Cookie(CSRFCookieName)
	if err != nil || header == "" || cookie.Value == "" || subtle.ConstantTimeCompare([]byte(header), []byte(cookie.Value)) != 1 {
		return ErrInvalidCSRF
	}
	raw, err := base64.RawURLEncoding.DecodeString(header)
	if err != nil || len(raw) != 32 {
		return ErrInvalidCSRF
	}
	hash := sha256.Sum256(raw)
	if subtle.ConstantTimeCompare(hash[:], session.CSRFHash) != 1 {
		return ErrInvalidCSRF
	}
	return nil
}

func (s *Service) Logout(ctx context.Context, sessionID string) error {
	_, err := s.pool.Exec(ctx, `update public.sessions set revoked_at = now() where id = $1`, sessionID)
	return err
}

func (s *Service) ChangePassword(ctx context.Context, principal Principal, currentPassword, newPassword string) error {
	if err := ValidatePassword(newPassword, principal.LoginIdentifier()); err != nil {
		return err
	}
	valid, err := s.VerifyPassword(ctx, principal, currentPassword)
	if err != nil {
		return err
	}
	if !valid {
		return ErrInvalidCredentials
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if principal.AccountType == "admin" {
		if _, err := tx.Exec(ctx, `
			update public.admin_accounts
			set password_hash = extensions.crypt($2, extensions.gen_salt('bf', 12)), must_change_password = false
			where account_id = $1`, principal.ID, newPassword); err != nil {
			return err
		}
	} else {
		if _, err := tx.Exec(ctx, `
			update public.users
			set password_hash = extensions.crypt($2, extensions.gen_salt('bf', 12)), must_change_password = false
			where id = $1`, principal.ID, newPassword); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(ctx, `
		update public.sessions set revoked_at = now()
		where user_id = $1`, principal.ID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Service) VerifyPassword(ctx context.Context, principal Principal, password string) (bool, error) {
	var identifier string
	var passwordHash *string
	if principal.AccountType == "admin" {
		if err := s.pool.QueryRow(ctx, `select username, password_hash from public.admin_accounts where account_id = $1 and is_active`, principal.ID).Scan(&identifier, &passwordHash); err != nil {
			return false, err
		}
	} else {
		if err := s.pool.QueryRow(ctx, `select nip, password_hash from public.users where id = $1`, principal.ID).Scan(&identifier, &passwordHash); err != nil {
			return false, err
		}
	}
	return s.passwordMatches(ctx, identifier, passwordHash, password)
}

func (s *Service) RequestPasswordReset(ctx context.Context, nip, channel, ip string) error {
	nip = strings.TrimSpace(nip)
	if !validNIP(nip) || (channel != "email" && channel != "phone") {
		return nil
	}
	if s.delivery == nil || !s.delivery.ChannelConfigured(channel) {
		return ErrDeliveryUnavailable
	}
	var userID, name, destination string
	var verifiedAt *time.Time
	column := "u.email"
	verifiedColumn := "u.email_verified_at"
	if channel == "phone" {
		column = "u.phone_e164"
		verifiedColumn = "u.phone_verified_at"
	}
	query := fmt.Sprintf(`
		select u.id::text, u.name, coalesce(%s, ''), %s
		from public.users u
		where u.nip = $1 and u.is_active and u.deleted_at is null`, column, verifiedColumn)
	if err := s.pool.QueryRow(ctx, query, nip).Scan(&userID, &name, &destination, &verifiedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}
	if destination == "" || verifiedAt == nil {
		return nil
	}
	var recent int
	if err := s.pool.QueryRow(ctx, `
		select count(*) from public.recovery_otps
		where user_id = $1 and purpose = 'password_reset' and created_at > now() - interval '15 minutes'`, userID).Scan(&recent); err != nil {
		return err
	}
	if recent >= 3 {
		return nil
	}
	return s.createAndQueueOTP(ctx, userID, name, "password_reset", channel, destination, ip)
}

func (s *Service) ResetPassword(ctx context.Context, nip, channel, otp, newPassword string) error {
	if !validNIP(nip) || (channel != "email" && channel != "phone") {
		return ErrInvalidOTP
	}
	if err := ValidatePassword(newPassword, nip); err != nil {
		return err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var otpID, userID string
	var storedHash []byte
	var attempts int
	err = tx.QueryRow(ctx, `
		select ro.id::text, ro.user_id::text, ro.otp_hash, ro.attempts
		from public.recovery_otps ro
		join public.users u on u.id = ro.user_id
		where u.nip = $1 and ro.purpose = 'password_reset' and ro.channel = $2
		  and ro.consumed_at is null and ro.expires_at > now()
		order by ro.created_at desc
		limit 1 for update`, nip, channel).Scan(&otpID, &userID, &storedHash, &attempts)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInvalidOTP
	}
	if err != nil {
		return err
	}
	calculated := s.otpHash(userID, "password_reset", channel, otp)
	if attempts >= 5 || subtle.ConstantTimeCompare(calculated, storedHash) != 1 {
		_, _ = tx.Exec(ctx, `update public.recovery_otps set attempts = attempts + 1 where id = $1`, otpID)
		_ = tx.Commit(ctx)
		return ErrInvalidOTP
	}
	if _, err := tx.Exec(ctx, `
		update public.users
		set password_hash = extensions.crypt($2, extensions.gen_salt('bf', 12)), must_change_password = false
		where id = $1`, userID, newPassword); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `update public.recovery_otps set consumed_at = now() where id = $1`, otpID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `update public.sessions set revoked_at = now() where user_id = $1 and revoked_at is null`, userID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Service) StartContactChange(ctx context.Context, principal Principal, channel, destination, currentPassword string) error {
	if principal.AccountType != "user" {
		return contactInputError("kontak pemulihan hanya tersedia untuk akun pegawai")
	}
	if channel != "email" && channel != "phone" {
		return contactInputError("kanal kontak tidak valid")
	}
	valid, err := s.VerifyPassword(ctx, principal, currentPassword)
	if err != nil {
		return err
	}
	if !valid {
		return ErrInvalidCredentials
	}
	if s.delivery == nil || !s.delivery.ChannelConfigured(channel) {
		return ErrDeliveryUnavailable
	}
	if channel == "email" {
		destination, err = normalizeEmail(destination)
	} else if channel == "phone" {
		destination, err = normalizePhone(destination)
	}
	if err != nil {
		return err
	}
	var exists bool
	if channel == "email" {
		err = s.pool.QueryRow(ctx, `select exists(select 1 from public.users where lower(email) = lower($1) and id <> $2 and deleted_at is null)`, destination, principal.ID).Scan(&exists)
	} else {
		err = s.pool.QueryRow(ctx, `select exists(select 1 from public.users where phone_e164 = $1 and id <> $2 and deleted_at is null)`, destination, principal.ID).Scan(&exists)
	}
	if err != nil {
		return err
	}
	if exists {
		return contactInputError("kontak sudah digunakan akun lain")
	}

	otp, err := numericOTP()
	if err != nil {
		return err
	}
	hash := s.otpHash(principal.ID, "contact_change", channel, otp)
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `
		update public.pending_contact_changes set consumed_at = now()
		where user_id = $1 and channel = $2 and consumed_at is null`, principal.ID, channel); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		update public.notification_jobs
		set status = 'cancelled', locked_at = null, last_error = 'digantikan permintaan OTP yang lebih baru'
		where user_id = $1 and channel = $2 and template_code = 'contact_otp'
		  and status in ('pending', 'processing')`, principal.ID, channel); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		insert into public.pending_contact_changes (user_id, channel, destination, otp_hash, expires_at)
		values ($1, $2, $3, $4, now() + interval '10 minutes')`, principal.ID, channel, destination, hash); err != nil {
		return err
	}
	userID := principal.ID
	jobID, err := mailer.QueueImmediate(ctx, tx, &userID, channel, destination, "contact_otp", map[string]any{"name": principal.Name, "otp": otp})
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	if err := s.delivery.DeliverClaimed(ctx, jobID); err != nil {
		return ErrDeliveryFailed
	}
	return nil
}

func (s *Service) VerifyContactChange(ctx context.Context, principal Principal, channel, otp string) error {
	if principal.AccountType != "user" {
		return contactInputError("kontak pemulihan hanya tersedia untuk akun pegawai")
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var changeID, destination string
	var storedHash []byte
	var attempts int
	err = tx.QueryRow(ctx, `
		select id::text, destination, otp_hash, attempts
		from public.pending_contact_changes
		where user_id = $1 and channel = $2 and consumed_at is null and expires_at > now()
		order by created_at desc limit 1 for update`, principal.ID, channel).Scan(&changeID, &destination, &storedHash, &attempts)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInvalidOTP
	}
	if err != nil {
		return err
	}
	calculated := s.otpHash(principal.ID, "contact_change", channel, otp)
	if attempts >= 5 || subtle.ConstantTimeCompare(calculated, storedHash) != 1 {
		_, _ = tx.Exec(ctx, `update public.pending_contact_changes set attempts = attempts + 1 where id = $1`, changeID)
		_ = tx.Commit(ctx)
		return ErrInvalidOTP
	}
	if channel == "email" {
		_, err = tx.Exec(ctx, `update public.users set email = lower($2), email_verified_at = now() where id = $1`, principal.ID, destination)
	} else if channel == "phone" {
		_, err = tx.Exec(ctx, `update public.users set phone_e164 = $2, phone_verified_at = now() where id = $1`, principal.ID, destination)
	} else {
		return contactInputError("kanal kontak tidak valid")
	}
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `update public.pending_contact_changes set consumed_at = now() where id = $1`, changeID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Service) SetCookies(response http.ResponseWriter, result LoginResult) {
	http.SetCookie(response, &http.Cookie{
		Name: SessionCookieName, Value: result.SessionToken, Path: "/",
		HttpOnly: true, Secure: s.cfg.CookieSecure, SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(response, &http.Cookie{
		Name: CSRFCookieName, Value: result.CSRFToken, Path: "/",
		HttpOnly: false, Secure: s.cfg.CookieSecure, SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(response, &http.Cookie{
		Name: TabProofCookieName, Value: s.tabProof(result.TabToken), Path: "/",
		HttpOnly: true, Secure: s.cfg.CookieSecure, SameSite: http.SameSiteLaxMode,
	})
}

func (s *Service) ClearCookies(response http.ResponseWriter) {
	for _, name := range []string{SessionCookieName, CSRFCookieName, TabProofCookieName} {
		http.SetCookie(response, &http.Cookie{
			Name: name, Value: "", Path: "/", MaxAge: -1, Expires: time.Unix(1, 0),
			HttpOnly: name != CSRFCookieName, Secure: s.cfg.CookieSecure, SameSite: http.SameSiteLaxMode,
		})
	}
}

func (s *Service) tabProof(tabToken string) string {
	raw, err := base64.RawURLEncoding.DecodeString(tabToken)
	if err != nil || len(raw) != 32 {
		return ""
	}
	mac := hmac.New(sha256.New, []byte(s.cfg.AppSecret))
	mac.Write([]byte("pantas-tab-session|"))
	mac.Write(raw)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (s *Service) verifyTabSession(request *http.Request) bool {
	tabToken := strings.TrimSpace(request.Header.Get(TabTokenHeader))
	if tabToken == "" {
		return false
	}
	proofCookie, err := request.Cookie(TabProofCookieName)
	if err != nil || proofCookie.Value == "" {
		return false
	}
	expected := s.tabProof(tabToken)
	if expected == "" || len(expected) != len(proofCookie.Value) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(expected), []byte(proofCookie.Value)) == 1
}

func ValidatePassword(password, loginIdentifier string) error {
	if !utf8.ValidString(password) || len(password) < 10 || len(password) > 128 {
		return errors.New("password harus berisi 10–128 karakter")
	}
	if loginIdentifier != "" && strings.Contains(strings.ToLower(password), strings.ToLower(loginIdentifier)) {
		return errors.New("password baru tidak boleh memuat NIP atau username")
	}
	classes := 0
	var upper, lower, digit, symbol bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upper = true
		case unicode.IsLower(char):
			lower = true
		case unicode.IsDigit(char):
			digit = true
		default:
			symbol = true
		}
	}
	for _, present := range []bool{upper, lower, digit, symbol} {
		if present {
			classes++
		}
	}
	if classes < 3 {
		return errors.New("gunakan sedikitnya tiga jenis karakter: huruf besar, huruf kecil, angka, atau simbol")
	}
	return nil
}

func MaskDestination(channel, destination string) string {
	if channel == "email" {
		parts := strings.SplitN(destination, "@", 2)
		if len(parts) == 2 && len(parts[0]) > 1 {
			return parts[0][:1] + strings.Repeat("*", min(6, len(parts[0])-1)) + "@" + parts[1]
		}
	}
	if len(destination) > 4 {
		return strings.Repeat("*", len(destination)-4) + destination[len(destination)-4:]
	}
	return "****"
}

func (s *Service) passwordMatches(ctx context.Context, nip string, hash *string, password string) (bool, error) {
	if hash == nil || *hash == "" {
		return subtle.ConstantTimeCompare([]byte(password), []byte(nip)) == 1, nil
	}
	var valid bool
	if err := s.pool.QueryRow(ctx, `select extensions.crypt($1, $2) = $2`, password, *hash).Scan(&valid); err != nil {
		return false, err
	}
	return valid, nil
}

func (s *Service) loginRateLimited(ctx context.Context, identifier, ip string) (bool, error) {
	var identifierFailures, ipFailures int
	err := s.pool.QueryRow(ctx, `
		select
		  count(*) filter (where nip = $1 and not was_successful),
		  count(*) filter (where ip_address = nullif($2, '')::inet and not was_successful)
		from public.login_attempts
		where occurred_at > now() - interval '15 minutes'`, identifier, ip).Scan(&identifierFailures, &ipFailures)
	return identifierFailures >= 8 || ipFailures >= 30, err
}

func (s *Service) recordLoginAttempt(ctx context.Context, identifier, ip string, success bool) {
	_, _ = s.pool.Exec(ctx, `
		insert into public.login_attempts (nip, ip_address, was_successful)
		values (nullif($1, ''), nullif($2, '')::inet, $3)`, identifier, ip, success)
}

func (s *Service) createAndQueueOTP(ctx context.Context, userID, name, purpose, channel, destination, ip string) error {
	otp, err := numericOTP()
	if err != nil {
		return err
	}
	hash := s.otpHash(userID, purpose, channel, otp)
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `
		update public.recovery_otps set consumed_at = now()
		where user_id = $1 and purpose = $2 and channel = $3 and consumed_at is null`, userID, purpose, channel); err != nil {
		return err
	}
	template := "password_otp"
	if purpose != "password_reset" {
		template = "contact_otp"
	}
	if _, err := tx.Exec(ctx, `
		update public.notification_jobs
		set status = 'cancelled', locked_at = null, last_error = 'digantikan permintaan OTP yang lebih baru'
		where user_id = $1 and channel = $2 and template_code = $3
		  and status in ('pending', 'processing')`, userID, channel, template); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		insert into public.recovery_otps (user_id, purpose, channel, destination, otp_hash, requested_ip, expires_at)
		values ($1, $2, $3, $4, $5, nullif($6, '')::inet, now() + interval '10 minutes')`,
		userID, purpose, channel, destination, hash, ip); err != nil {
		return err
	}
	jobID, err := mailer.QueueImmediate(ctx, tx, &userID, channel, destination, template, map[string]any{"name": name, "otp": otp})
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	if err := s.delivery.DeliverClaimed(ctx, jobID); err != nil {
		return ErrDeliveryFailed
	}
	return nil
}

func (s *Service) otpHash(userID, purpose, channel, otp string) []byte {
	mac := hmac.New(sha256.New, []byte(s.cfg.AppSecret))
	mac.Write([]byte(userID))
	mac.Write([]byte("|" + purpose + "|" + channel + "|" + otp))
	return mac.Sum(nil)
}

func randomToken(size int) (string, []byte, error) {
	raw := make([]byte, size)
	if _, err := rand.Read(raw); err != nil {
		return "", nil, err
	}
	hash := sha256.Sum256(raw)
	return base64.RawURLEncoding.EncodeToString(raw), hash[:], nil
}

func numericOTP() (string, error) {
	raw := make([]byte, 4)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	value := uint32(raw[0])<<24 | uint32(raw[1])<<16 | uint32(raw[2])<<8 | uint32(raw[3])
	return fmt.Sprintf("%06d", value%1000000), nil
}

func validLoginIdentifier(value string) bool {
	if validNIP(value) {
		return true
	}
	if len(value) < 3 || len(value) > 64 || value[0] < 'a' || value[0] > 'z' {
		return false
	}
	for _, char := range value {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '.' || char == '_' || char == '-') {
			return false
		}
	}
	return true
}

func validNIP(nip string) bool {
	if len(nip) != 18 {
		return false
	}
	for _, char := range nip {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func normalizeEmail(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	address, err := mailaddress.ParseAddress(value)
	if err != nil || strings.ToLower(address.Address) != value || len(value) > 254 {
		return "", contactInputError("alamat email tidak valid")
	}
	return value, nil
}

var nonDigits = regexp.MustCompile(`[^0-9]+`)

func normalizePhone(value string) (string, error) {
	value = strings.TrimSpace(value)
	plus := strings.HasPrefix(value, "+")
	digits := nonDigits.ReplaceAllString(value, "")
	switch {
	case strings.HasPrefix(digits, "08"):
		digits = "62" + digits[1:]
	case strings.HasPrefix(digits, "8"):
		digits = "62" + digits
	case strings.HasPrefix(digits, "62"):
	case plus:
	default:
		return "", contactInputError("nomor HP harus menggunakan format Indonesia, misalnya 0812… atau +62812…")
	}
	if len(digits) < 9 || len(digits) > 15 {
		return "", contactInputError("nomor HP tidak valid")
	}
	return "+" + digits, nil
}

func ClientIP(request *http.Request, trustProxy bool) string {
	if trustProxy {
		if forwarded := request.Header.Get("X-Forwarded-For"); forwarded != "" {
			first := strings.TrimSpace(strings.Split(forwarded, ",")[0])
			if net.ParseIP(first) != nil {
				return first
			}
		}
	}
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err == nil && net.ParseIP(host) != nil {
		return host
	}
	if net.ParseIP(request.RemoteAddr) != nil {
		return request.RemoteAddr
	}
	return ""
}

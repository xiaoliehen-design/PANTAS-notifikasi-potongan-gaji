package httpapi

import (
	"errors"
	"net/http"

	"github.com/bcpriok/pantas/internal/auth"
)

func (a *App) login(response http.ResponseWriter, request *http.Request) {
	var input struct {
		Identifier string `json:"identifier"`
		NIP        string `json:"nip"`
		Password   string `json:"password"`
	}
	if !decodeJSON(response, request, &input) {
		return
	}
	identifier := input.Identifier
	if identifier == "" {
		identifier = input.NIP
	}
	result, err := a.auth.Login(request.Context(), identifier, input.Password, auth.ClientIP(request, a.cfg.TrustProxy), request.UserAgent())
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			writeError(response, http.StatusUnauthorized, err.Error(), "invalid_credentials")
		case errors.Is(err, auth.ErrRateLimited):
			writeError(response, http.StatusTooManyRequests, err.Error(), "rate_limited")
		default:
			a.log.Error("login", "error", err)
			writeError(response, http.StatusInternalServerError, "Login belum dapat diproses.", "internal_error")
		}
		return
	}
	a.auth.SetCookies(response, result)
	writeJSON(response, http.StatusOK, map[string]any{"user": result.Principal})
}

func (a *App) me(response http.ResponseWriter, _ *http.Request, principal auth.Principal) {
	writeJSON(response, http.StatusOK, map[string]any{"user": principal})
}

func (a *App) logout(response http.ResponseWriter, request *http.Request, _ auth.Principal) {
	_, session, err := a.auth.AuthenticateRequest(request.Context(), request)
	if err == nil {
		_ = a.auth.Logout(request.Context(), session.ID)
	}
	a.auth.ClearCookies(response)
	response.WriteHeader(http.StatusNoContent)
}

func (a *App) changePassword(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	var input struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if !decodeJSON(response, request, &input) {
		return
	}
	if err := a.auth.ChangePassword(request.Context(), principal, input.CurrentPassword, input.NewPassword); err != nil {
		status := http.StatusUnprocessableEntity
		code := "password_invalid"
		message := err.Error()
		if errors.Is(err, auth.ErrInvalidCredentials) {
			code = "current_password_invalid"
			message = "Password yang Anda masukkan salah."
		}
		writeError(response, status, message, code)
		return
	}
	a.auth.ClearCookies(response)
	writeJSON(response, http.StatusOK, map[string]any{"message": "Password berhasil diubah. Silakan login kembali."})
}

func (a *App) forgotPassword(response http.ResponseWriter, request *http.Request) {
	var input struct {
		NIP     string `json:"nip"`
		Channel string `json:"channel"`
	}
	if !decodeJSON(response, request, &input) {
		return
	}
	if err := a.auth.RequestPasswordReset(request.Context(), input.NIP, input.Channel, auth.ClientIP(request, a.cfg.TrustProxy)); err != nil {
		if errors.Is(err, auth.ErrDeliveryUnavailable) {
			writeError(response, http.StatusServiceUnavailable, "Kanal pengiriman kode belum dikonfigurasi oleh administrator.", "delivery_unavailable")
			return
		}
		a.log.Error("request password reset", "error", err)
	}
	writeJSON(response, http.StatusAccepted, map[string]any{"message": "Jika NIP dan kontak terverifikasi ditemukan, kode reset akan dikirim."})
}

func (a *App) resetPassword(response http.ResponseWriter, request *http.Request) {
	var input struct {
		NIP         string `json:"nip"`
		Channel     string `json:"channel"`
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}
	if !decodeJSON(response, request, &input) {
		return
	}
	if err := a.auth.ResetPassword(request.Context(), input.NIP, input.Channel, input.OTP, input.NewPassword); err != nil {
		status := http.StatusUnprocessableEntity
		code := "invalid_otp"
		if !errors.Is(err, auth.ErrInvalidOTP) {
			code = "password_invalid"
		}
		writeError(response, status, err.Error(), code)
		return
	}
	writeJSON(response, http.StatusOK, map[string]any{"message": "Password berhasil diatur ulang. Silakan login."})
}

func (a *App) startContactChange(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	var input struct {
		Channel         string `json:"channel"`
		Destination     string `json:"destination"`
		CurrentPassword string `json:"current_password"`
	}
	if !decodeJSON(response, request, &input) {
		return
	}
	if err := a.auth.StartContactChange(request.Context(), principal, input.Channel, input.Destination, input.CurrentPassword); err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			writeError(response, http.StatusUnprocessableEntity, "Password yang Anda masukkan salah.", "current_password_invalid")
		case errors.Is(err, auth.ErrDeliveryUnavailable):
			writeError(response, http.StatusServiceUnavailable, "Kanal pengiriman kode belum dikonfigurasi oleh administrator.", "delivery_unavailable")
		case errors.Is(err, auth.ErrDeliveryFailed):
			writeError(response, http.StatusBadGateway, err.Error(), "delivery_failed")
		case auth.IsContactInputError(err):
			writeError(response, http.StatusUnprocessableEntity, err.Error(), "contact_invalid")
		default:
			a.log.Error("start contact change", "error", err)
			writeError(response, http.StatusInternalServerError, "Kode verifikasi belum dapat dikirim.", "internal_error")
		}
		return
	}
	writeJSON(response, http.StatusAccepted, map[string]any{"message": "Kode verifikasi berhasil dikirim. Periksa kotak masuk atau pesan pada nomor HP Anda."})
}

func (a *App) verifyContactChange(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	var input struct {
		Channel string `json:"channel"`
		OTP     string `json:"otp"`
	}
	if !decodeJSON(response, request, &input) {
		return
	}
	if err := a.auth.VerifyContactChange(request.Context(), principal, input.Channel, input.OTP); err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidOTP):
			writeError(response, http.StatusUnprocessableEntity, err.Error(), "invalid_otp")
		case auth.IsContactInputError(err):
			writeError(response, http.StatusUnprocessableEntity, err.Error(), "contact_invalid")
		default:
			a.log.Error("verify contact change", "error", err)
			writeError(response, http.StatusInternalServerError, "Kontak belum dapat diverifikasi.", "internal_error")
		}
		return
	}
	writeJSON(response, http.StatusOK, map[string]any{"message": "Kontak berhasil diverifikasi."})
}

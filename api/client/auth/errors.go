package auth

import "errors"

var ErrOTPRequired error = errors.New("OTP Token is required.")
var ErrOTPInvalid error = errors.New("OTP Token is invalid.")

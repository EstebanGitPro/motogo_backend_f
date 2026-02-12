package jwt

import "github.com/EstebanGitPro/motogo-backend/core/ports/output"

// Ensure JWKSValidator implements output.JWTValidator
var _ output.JWTValidator = (*JWKSValidator)(nil)

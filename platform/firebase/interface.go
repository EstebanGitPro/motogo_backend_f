package firebase

import "github.com/EstebanGitPro/motogo-backend/core/ports/output"

// Ensure Client implements output.CustomTokenProvider
var _ output.CustomTokenProvider = (*Client)(nil)

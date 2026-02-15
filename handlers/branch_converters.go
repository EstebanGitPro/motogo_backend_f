package handlers

import "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"

// toDisplacementRanges converts []string (API input) to []domain.DisplacementRange (domain)
func toDisplacementRanges(ranges []string) []domain.DisplacementRange {
	if ranges == nil {
		return nil
	}
	result := make([]domain.DisplacementRange, len(ranges))
	for i, r := range ranges {
		result[i] = domain.DisplacementRange(r)
	}
	return result
}

// fromDisplacementRanges converts []domain.DisplacementRange (domain) to []string (API output)
func fromDisplacementRanges(ranges []domain.DisplacementRange) []string {
	if ranges == nil {
		return nil
	}
	result := make([]string, len(ranges))
	for i, r := range ranges {
		result[i] = string(r)
	}
	return result
}

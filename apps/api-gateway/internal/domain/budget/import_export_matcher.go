package budget

func MatchImportedIncome(incoming ImportIncomeRecord, existing []ImportIncomeRecord) ImportMatchResult {
	sameDate := make([]ImportIncomeRecord, 0)
	for _, candidate := range existing {
		if candidate.Date == incoming.Date {
			sameDate = append(sameDate, candidate)
		}
	}

	if len(sameDate) == 0 {
		return ImportMatchResult{Action: ImportMergeActionCreate}
	}

	for _, candidate := range sameDate {
		if candidate.Description == incoming.Description && candidate.Amount == incoming.Amount {
			id := candidate.ID
			if id == "" {
				id = candidate.Date
			}
			return ImportMatchResult{
				Action:            ImportMergeActionSkipExisting,
				MatchedExistingID: &id,
			}
		}
	}

	conflictCandidate := sameDate[0]
	conflictID := conflictCandidate.ID
	if conflictID == "" {
		conflictID = conflictCandidate.Date
	}

	return ImportMatchResult{
		Action:            ImportMergeActionConflict,
		MatchedExistingID: &conflictID,
		Conflict: &ImportConflictDetail{
			Entity:      "income",
			Date:        incoming.Date,
			ExistingID:  conflictID,
			Incoming:    incoming,
			Existing:    conflictCandidate,
			Differences: incomeDifferences(incoming, conflictCandidate),
		},
	}
}

func MatchImportedExpense(incoming ImportExpenseRecord, existing []ImportExpenseRecord) ImportMatchResult {
	sameDate := make([]ImportExpenseRecord, 0)
	for _, candidate := range existing {
		if candidate.ExpenseDate == incoming.ExpenseDate {
			sameDate = append(sameDate, candidate)
		}
	}

	if len(sameDate) == 0 {
		return ImportMatchResult{Action: ImportMergeActionCreate}
	}

	for _, candidate := range sameDate {
		if candidate.Description == incoming.Description && candidate.Amount == incoming.Amount {
			id := candidate.ID
			if id == "" {
				id = candidate.ExpenseDate
			}
			return ImportMatchResult{
				Action:            ImportMergeActionSkipExisting,
				MatchedExistingID: &id,
			}
		}
	}

	conflictCandidate := sameDate[0]
	conflictID := conflictCandidate.ID
	if conflictID == "" {
		conflictID = conflictCandidate.ExpenseDate
	}

	return ImportMatchResult{
		Action:            ImportMergeActionConflict,
		MatchedExistingID: &conflictID,
		Conflict: &ImportConflictDetail{
			Entity:      "expense",
			Date:        incoming.ExpenseDate,
			ExistingID:  conflictID,
			Incoming:    incoming,
			Existing:    conflictCandidate,
			Differences: expenseDifferences(incoming, conflictCandidate),
		},
	}
}

func incomeDifferences(incoming, existing ImportIncomeRecord) []ImportFieldDifference {
	diffs := make([]ImportFieldDifference, 0)
	if incoming.Description != existing.Description {
		diffs = append(diffs, ImportFieldDifference{Field: "description", Incoming: incoming.Description, Existing: existing.Description})
	}
	if incoming.Amount != existing.Amount {
		diffs = append(diffs, ImportFieldDifference{Field: "amount", Incoming: incoming.Amount, Existing: existing.Amount})
	}
	if incoming.Currency != existing.Currency {
		diffs = append(diffs, ImportFieldDifference{Field: "currency", Incoming: incoming.Currency, Existing: existing.Currency})
	}
	return diffs
}

func expenseDifferences(incoming, existing ImportExpenseRecord) []ImportFieldDifference {
	diffs := make([]ImportFieldDifference, 0)
	if incoming.Description != existing.Description {
		diffs = append(diffs, ImportFieldDifference{Field: "description", Incoming: incoming.Description, Existing: existing.Description})
	}
	if incoming.Amount != existing.Amount {
		diffs = append(diffs, ImportFieldDifference{Field: "amount", Incoming: incoming.Amount, Existing: existing.Amount})
	}
	if incoming.Currency != existing.Currency {
		diffs = append(diffs, ImportFieldDifference{Field: "currency", Incoming: incoming.Currency, Existing: existing.Currency})
	}
	if incoming.CategoryName != existing.CategoryName {
		diffs = append(diffs, ImportFieldDifference{Field: "category_name", Incoming: incoming.CategoryName, Existing: existing.CategoryName})
	}
	return diffs
}

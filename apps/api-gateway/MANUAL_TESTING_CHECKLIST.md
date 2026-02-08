# Manual Testing Checklist - Budget Recurring Income Features

This document provides a comprehensive testing checklist for the weekly/monthly recurring income features and expense dialog improvements.

## Prerequisites

1. API server running on `http://localhost:8080`
2. Frontend running on `http://localhost:3000`
3. Database migrations applied (005 and 006)
4. Test user account created or using demo login

---

## Test Group 1: Weekly Recurring Income

### 12.1 Test creating weekly recurring income with end_date

**Steps:**
1. Navigate to `/dashboard/budget`
2. Click "Add Income" button
3. Fill in form:
   - Amount: 5000
   - Description: "Weekly Allowance"
   - Type: "Weekly Recurring"
   - Start Date: Today's date
   - Uncheck "No end date"
   - End Date: 4 weeks from today
4. Click "Add Income"

**Expected Results:**
- ✅ Form validates successfully
- ✅ Toast shows "Income created successfully"
- ✅ Dialog closes automatically
- ✅ Income appears in "Recent Incomes" list with "weekly" badge
- ✅ Budget remaining updates to include weekly income

**Status:** ⬜ Not tested / ✅ Passed / ❌ Failed

---

### 12.2 Test creating monthly recurring income without end_date (indefinite)

**Steps:**
1. Click "Add Income" button
2. Fill in form:
   - Amount: 15000
   - Description: "Monthly Salary"
   - Type: "Monthly Recurring"
   - Start Date: First day of current month
   - Check "No end date (ongoing)"
3. Click "Add Income"

**Expected Results:**
- ✅ End date field is hidden when checkbox is checked
- ✅ Form validates successfully
- ✅ Toast shows "Income created successfully"
- ✅ Income appears with "monthly" badge
- ✅ Budget remaining includes monthly calculation
- ✅ Clicking income to edit shows "Ongoing" badge

**Status:** ⬜ Not tested / ✅ Passed / ❌ Failed

---

## Test Group 2: Budget Remaining Calculations

### 12.3 Test budget remaining calculation includes weekly/monthly income correctly

**Steps:**
1. Create a weekly recurring income: 1000 PHP starting today
2. Create a monthly recurring income: 5000 PHP starting this month
3. Select today's date in "Budget Date" filter
4. Note the budget remaining value
5. Create an expense for 500 PHP
6. Check budget remaining updated correctly

**Expected Formula:**
```
Budget Remaining = Total Income - Total Expenses

For weekly (from start_date to budget_date):
  weeks = floor(days_diff / 7) + 1
  total = amount * weeks

For monthly (from start_date to budget_date):
  months = (year_diff * 12) + month_diff + 1
  total = amount * months
```

**Expected Results:**
- ✅ Weekly income counted correctly
- ✅ Monthly income counted correctly
- ✅ Budget remaining = (weekly total + monthly total) - expenses

**Status:** ⬜ Not tested / ✅ Passed / ❌ Failed

---

### 12.4 Test budget remaining updates immediately after adding income

**Steps:**
1. Note current budget remaining value
2. Click "Add Income"
3. Add income: 1000 PHP, one-time, today
4. Submit form
5. Observe budget remaining card

**Expected Results:**
- ✅ Budget remaining updates **immediately** (no page refresh needed)
- ✅ New value = old value + 1000
- ✅ Color changes based on new status (green/red/neutral)

**Status:** ⬜ Not tested / ✅ Passed / ❌ Failed

---

### 12.5 Test budget remaining updates immediately after deleting expense

**Steps:**
1. Note current budget remaining value
2. Click any expense in "Recent Expenses"
3. Click "Delete Expense"
4. Confirm deletion
5. Observe budget remaining card

**Expected Results:**
- ✅ Budget remaining updates **immediately**
- ✅ New value = old value + deleted expense amount
- ✅ Color changes if status changed

**Status:** ⬜ Not tested / ✅ Passed / ❌ Failed

---

## Test Group 3: Expense Dialog UX

### 12.6 Test expense dialog keyboard navigation (Tab, ESC, Enter)

**Steps:**
1. Click "Add Expense"
2. Press `Tab` key repeatedly
3. Verify focus moves through all fields in order
4. Press `Shift+Tab` to go backward
5. Press `ESC` key
6. Click "Add Expense" again
7. Fill in all required fields
8. Press `Enter` in description field

**Expected Results:**
- ✅ Tab moves focus forward through: description → amount → currency → category → date → tags → notes → cancel → submit
- ✅ Shift+Tab moves focus backward
- ✅ ESC key closes dialog
- ✅ Focus is trapped within dialog (doesn't leave dialog container)
- ✅ Enter in form fields does NOT submit form (only submit button does)

**Status:** ⬜ Not tested / ✅ Passed / ❌ Failed

---

### 12.7 Test expense dialog on mobile (responsive behavior)

**Steps:**
1. Open browser DevTools
2. Toggle device toolbar (Ctrl+Shift+M)
3. Select "iPhone 12 Pro" or similar
4. Navigate to `/dashboard/budget`
5. Click "Add Expense"
6. Try scrolling within dialog
7. Fill in form and submit

**Expected Results:**
- ✅ Dialog is full-width on mobile
- ✅ Dialog content is scrollable if it exceeds viewport height
- ✅ All form fields are accessible
- ✅ Buttons are large enough to tap (min 44x44px)
- ✅ No horizontal scrolling required
- ✅ Keyboard opens without breaking layout

**Status:** ⬜ Not tested / ✅ Passed / ❌ Failed

---

## Test Group 4: Backward Compatibility

### 12.8 Test backward compatibility with existing daily recurring incomes

**Steps:**
1. Create a daily recurring income using old format (if any exist in DB)
2. View it in income list
3. Click to edit
4. Verify form loads correctly
5. Make a change and save

**Expected Results:**
- ✅ Existing daily incomes display correctly
- ✅ Edit form loads with correct data
- ✅ Can update existing daily incomes
- ✅ No errors in console
- ✅ Backend still calculates daily incomes correctly

**Status:** ⬜ Not tested / ✅ Passed / ❌ Failed / ⬜ N/A (no existing data)

---

## Test Group 5: UI/UX Validation

### 12.9 Verify color contrast of budget remaining card in different states

**Test Red Status (Over Budget):**
1. Create expenses exceeding income
2. Check Budget Remaining card in expense stats

**Expected:**
- ✅ Text color: `text-red-500` (#ef4444)
- ✅ Readable on white background
- ✅ Icon color matches text
- ✅ Description text also red

**Test Green Status (On Track):**
1. Ensure income exceeds expenses
2. Check Budget Remaining card

**Expected:**
- ✅ Text color: `text-green-900` (#14532d)
- ✅ Readable on white background
- ✅ Icon color matches text
- ✅ Description text also green

**Test Neutral Status (Break Even):**
1. Make income = expenses exactly
2. Check Budget Remaining card

**Expected:**
- ✅ Text color: `text-gray-600` (#4b5563)
- ✅ Readable on white background

**WCAG AA Compliance:**
- ✅ Contrast ratio ≥ 4.5:1 for normal text
- ✅ Use browser DevTools → Lighthouse → Accessibility to verify

**Status:** ⬜ Not tested / ✅ Passed / ❌ Failed

---

### 12.10 Test form validation (end_date >= start_date)

**Test Invalid End Date:**
1. Click "Add Income"
2. Select Type: "Weekly Recurring"
3. Start Date: Today
4. Uncheck "No end date"
5. End Date: Yesterday (before start date)
6. Try to submit

**Expected Results:**
- ✅ Red validation message appears: "End date must be after start date"
- ✅ Form does NOT submit
- ✅ Submit button remains enabled (not a browser constraint)
- ✅ Error message disappears when valid date selected

**Test Valid End Date:**
1. Change End Date to tomorrow (after start date)
2. Submit form

**Expected Results:**
- ✅ No validation error
- ✅ Form submits successfully
- ✅ Income created with correct end_date

**Status:** ⬜ Not tested / ✅ Passed / ❌ Failed

---

## Test Group 6: Deep Linking & Navigation

### Additional Test: Direct URL Navigation

**Steps:**
1. In browser, navigate to `/dashboard/budget/expenses/new`
2. Observe behavior

**Expected Results:**
- ✅ Redirects to `/dashboard/budget?openExpense=true`
- ✅ Query param removed from URL
- ✅ Expense dialog opens automatically
- ✅ No flash of content before redirect

**Status:** ⬜ Not tested / ✅ Passed / ❌ Failed

---

### Additional Test: Browser Back Button

**Steps:**
1. On `/dashboard/budget`, click "Add Expense"
2. Dialog opens
3. Press browser back button

**Expected Results:**
- ✅ Dialog closes
- ✅ Stay on `/dashboard/budget` (don't navigate away)
- ✅ Press forward button does NOT reopen dialog

**Status:** ⬜ Not tested / ✅ Passed / ❌ Failed

---

## Summary

| Test Group | Total Tests | Passed | Failed | Not Tested |
|------------|-------------|--------|--------|------------|
| Weekly Recurring Income | 2 | - | - | - |
| Budget Calculations | 3 | - | - | - |
| Expense Dialog UX | 2 | - | - | - |
| Backward Compatibility | 1 | - | - | - |
| UI/UX Validation | 2 | - | - | - |
| Deep Linking | 2 | - | - | - |
| **TOTAL** | **12** | **-** | **-** | **-** |

---

## Notes

- Mark each test with ✅ Passed, ❌ Failed, or ⬜ Not tested
- Record any bugs found in the "Issues Found" section below
- Take screenshots for failed tests if possible

## Issues Found

*(Add any issues discovered during testing here)*

1. 
2. 
3. 

---

## Test Environment

- Browser: _________________________
- OS: _______________________________
- API Version: ______________________
- Database Migration: 006 applied? ☐ Yes ☐ No
- Test Date: ________________________
- Tester: ___________________________

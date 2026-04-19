import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  fetchCategories,
  fetchCategory,
  createCategory,
  updateCategory,
  deleteCategory,
  fetchTags,
  fetchTag,
  createTag,
  updateTag,
  deleteTag,
  fetchExpenses,
  fetchExpense,
  createExpense,
  updateExpense,
  deleteExpense,
  searchExpenses,
  fetchIncomes,
  fetchIncome,
  createIncome,
  updateIncome,
  deleteIncome,
  exportBudgetJSON,
  importBudgetJSON,
  skipIncome,
  skipExpense,
  fetchRecurringIncomes,
  fetchIncomeOccurrences,
  fetchBudgetRemaining,
  fetchSummaryStats,
  fetchCategoryBreakdown,
  fetchTrends
} from '@/lib/api';
import type { ExpenseSearchParams } from '@/lib/api';
import type {
  Category,
  Tag,
  Expense,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  CreateTagRequest,
  UpdateTagRequest,
  CreateExpenseRequest,
  UpdateExpenseRequest,
  ExpenseFilters,
  Income,
  RecurringIncomeWithNextDate,
  IncomeOccurrence,
  CreateIncomeRequest,
  UpdateIncomeRequest,
  BudgetRemainingResponse,
  SummaryStats,
  CategoryBreakdown,
  TrendItem,
  PaginatedResponse,
  BudgetExportPayload,
  BudgetImportResult,
  SkipIncomeRequest,
  SkipExpenseRequest
} from '@/types/api';
import { useAuth } from '@/contexts/auth-context';

export class GuestBlockedError extends Error {
  constructor(message = 'Guest user is read-only. Create an account to save changes') {
    super(message);
    this.name = 'GuestBlockedError';
  }
}

// query keys for cache management
export const budgetKeys = {
  all: ['budget'] as const,

  categories: () => [...budgetKeys.all, 'categories'] as const,
  categoriesList: () => [...budgetKeys.categories(), 'list'] as const,
  categoryDetail: (id: string) => [...budgetKeys.categories(), 'detail', id] as const,

  priorityGroups: () => [...budgetKeys.all, 'priorityGroups'] as const,

  tags: () => [...budgetKeys.all, 'tags'] as const,
  tagsList: () => [...budgetKeys.tags(), 'list'] as const,
  tagDetail: (id: string) => [...budgetKeys.tags(), 'detail', id] as const,

  expenses: () => [...budgetKeys.all, 'expenses'] as const,
  expensesList: (filters?: ExpenseFilters) => [...budgetKeys.expenses(), 'list', filters] as const,
  expensesSearch: (query: string, filters?: Omit<ExpenseSearchParams, 'q'>) => 
    [...budgetKeys.expenses(), 'search', query, filters] as const,
  expenseDetail: (id: string) => [...budgetKeys.expenses(), 'detail', id] as const,

  incomes: () => [...budgetKeys.all, 'incomes'] as const,
  incomesList: (startDate?: string, endDate?: string, recurringType?: 'daily') => 
    [...budgetKeys.incomes(), 'list', { startDate, endDate, recurringType }] as const,
  incomeDetail: (id: string) => [...budgetKeys.incomes(), 'detail', id] as const,
  incomeOccurrences: (startDate: string, endDate: string) =>
    [...budgetKeys.incomes(), 'occurrences', { startDate, endDate }] as const,
  recurringIncomes: () => [...budgetKeys.all, 'recurringIncomes'] as const,

  budgetRemaining: (date?: string) => [...budgetKeys.all, 'budgetRemaining', date] as const,

  stats: () => [...budgetKeys.all, 'stats'] as const,
  summary: (period?: string) => [...budgetKeys.stats(), 'summary', period] as const,
  breakdown: (startDate?: string, endDate?: string) =>
    [...budgetKeys.stats(), 'breakdown', { startDate, endDate }] as const,
  trends: (granularity: 'day' | 'month', startDate?: string, endDate?: string) =>
    [...budgetKeys.stats(), 'trends', { granularity, startDate, endDate }] as const,

  // Budget Analytics
  savingsRate: (startDate?: string, endDate?: string) =>
    [...budgetKeys.stats(), 'savingsRate', { startDate, endDate }] as const,
  spendingVelocity: (startDate?: string, endDate?: string) =>
    [...budgetKeys.stats(), 'spendingVelocity', { startDate, endDate }] as const,
  upcomingBills: (days: 7 | 30) => [...budgetKeys.all, 'upcomingBills', days] as const,
  categoryBudgets: (month?: string) => [...budgetKeys.all, 'categoryBudgets', month] as const,

  // Financial Health
  healthScore: () => [...budgetKeys.stats(), 'healthScore'] as const,
  fiftyThirtyTwenty: (startDate?: string, endDate?: string) =>
    [...budgetKeys.stats(), 'fiftyThirtyTwenty', { startDate, endDate }] as const,
  weekdayPattern: (startDate?: string, endDate?: string) =>
    [...budgetKeys.stats(), 'weekdayPattern', { startDate, endDate }] as const,
  monthOverMonth: () => [...budgetKeys.stats(), 'monthOverMonth'] as const,

  // Merchant Analysis
  merchantAnalysis: (limit?: number, startDate?: string, endDate?: string) =>
    [...budgetKeys.stats(), 'merchantAnalysis', { limit, startDate, endDate }] as const,
  subscriptions: () => [...budgetKeys.all, 'subscriptions'] as const,

  // Current Total Money
  currentTotalMoney: (startDate?: string, endDate?: string) =>
    [...budgetKeys.all, 'currentTotalMoney', { startDate, endDate }] as const
};

// Categories

export function useCategories() {
  return useQuery<Category[]>({
    queryKey: budgetKeys.categoriesList(),
    queryFn: fetchCategories,
    staleTime: 300000 // 5 minutes - categories don't change often
  });
}

export function useCategory(id: string | null) {
  return useQuery<Category>({
    queryKey: budgetKeys.categoryDetail(id ?? ''),
    queryFn: () => fetchCategory(id!),
    enabled: !!id
  });
}

export function useCreateCategory() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: (data: CreateCategoryRequest) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return createCategory(data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.categories() });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to create category:", error);
      }
    }
  });
}

export function useUpdateCategory() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateCategoryRequest }) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return updateCategory(id, data);
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.categories() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.categoryDetail(variables.id) });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to update category:", error);
      }
    }
  });
}

export function useDeleteCategory() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: (id: string) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return deleteCategory(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.categories() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.expenses() });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to delete category:", error);
      }
    }
  });
}

// Tags

export function useTags() {
  return useQuery<Tag[]>({
    queryKey: budgetKeys.tagsList(),
    queryFn: fetchTags,
    staleTime: 300000 // 5 minutes
  });
}

export function useTag(id: string | null) {
  return useQuery<Tag>({
    queryKey: budgetKeys.tagDetail(id ?? ''),
    queryFn: () => fetchTag(id!),
    enabled: !!id
  });
}

export function useCreateTag() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: (data: CreateTagRequest) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return createTag(data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.tags() });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to create tag:", error);
      }
    }
  });
}

export function useUpdateTag() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTagRequest }) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return updateTag(id, data);
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.tags() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.tagDetail(variables.id) });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to update tag:", error);
      }
    }
  });
}

export function useDeleteTag() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: (id: string) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return deleteTag(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.tags() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.expenses() });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to delete tag:", error);
      }
    }
  });
}

// Expenses

export function useExpenses(filters?: ExpenseFilters, page: number = 1, limit: number = 5) {
  return useQuery<PaginatedResponse<Expense>>({
    queryKey: [...budgetKeys.expensesList(filters), page, limit],
    queryFn: () => fetchExpenses({ ...filters, page, limit }),
    staleTime: 60000,
    refetchOnMount: true
  });
}

export function useExpense(id: string | null) {
  return useQuery<Expense>({
    queryKey: budgetKeys.expenseDetail(id ?? ''),
    queryFn: () => fetchExpense(id!),
    enabled: !!id
  });
}

export function useCreateExpense() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: (data: CreateExpenseRequest) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return createExpense(data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.expenses() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.budgetRemaining() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.stats() });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to create expense:", error);
      }
    }
  });
}

export function useUpdateExpense() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateExpenseRequest }) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return updateExpense(id, data);
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.expenses() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.expenseDetail(variables.id) });
      queryClient.invalidateQueries({ queryKey: budgetKeys.budgetRemaining() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.stats() });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to update expense:", error);
      }
    }
  });
}

export function useDeleteExpense() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: (id: string) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return deleteExpense(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.expenses() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.budgetRemaining() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.stats() });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to delete expense:", error);
      }
    }
  });
}

export function useSearchExpenses(
  query: string,
  filters?: { category_id?: string; start_date?: string; end_date?: string },
  page: number = 1,
  limit: number = 5
) {
  return useQuery<PaginatedResponse<Expense>>({
    queryKey: budgetKeys.expensesSearch(query, { ...filters, page, limit }),
    queryFn: () => searchExpenses({ q: query, ...filters, page, limit }),
    enabled: query.length > 0,
    staleTime: 60000,
    refetchOnMount: true
  });
}

// Incomes

export function useIncomes(
  filters?: { start_date?: string; end_date?: string; recurring_type?: string },
  page: number = 1,
  limit: number = 5
) {
  return useQuery<PaginatedResponse<Income>>({
    queryKey: [...budgetKeys.incomes(), 'list', filters, page, limit],
    queryFn: () => fetchIncomes({ ...filters, page, limit }),
    staleTime: 60000,
    refetchOnMount: true
  });
}

export function useIncome(id: string | null) {
  return useQuery<Income>({
    queryKey: budgetKeys.incomeDetail(id ?? ''),
    queryFn: () => fetchIncome(id!),
    enabled: !!id
  });
}

export function useRecurringIncomes() {
  return useQuery<RecurringIncomeWithNextDate[]>({
    queryKey: budgetKeys.recurringIncomes(),
    queryFn: fetchRecurringIncomes,
    staleTime: 300000 // 5 minutes
  });
}

export function useIncomeOccurrences(
  startDate: string,
  endDate: string,
  page: number = 1,
  limit: number = 50,
  enabled: boolean = true
) {
  return useQuery<PaginatedResponse<IncomeOccurrence>>({
    queryKey: [...budgetKeys.incomeOccurrences(startDate, endDate), page, limit],
    queryFn: () => fetchIncomeOccurrences({ start_date: startDate, end_date: endDate, page, limit }),
    enabled: enabled && !!startDate && !!endDate,
    staleTime: 60000
  });
}

export function useCreateIncome() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: (data: CreateIncomeRequest) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return createIncome(data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.incomes() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.budgetRemaining() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.stats() });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to create income:", error);
      }
    }
  });
}

export function useUpdateIncome() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateIncomeRequest }) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return updateIncome(id, data);
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.incomes() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.incomeDetail(variables.id) });
      queryClient.invalidateQueries({ queryKey: budgetKeys.budgetRemaining() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.stats() });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to update income:", error);
      }
    }
  });
}

export function useDeleteIncome() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: (id: string) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return deleteIncome(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.incomes() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.budgetRemaining() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.stats() });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to delete income:", error);
      }
    }
  });
}

export function useExportBudgetJSON() {
  return useMutation({
    mutationFn: exportBudgetJSON
  });
}

export function useImportBudgetJSON() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: (payload: BudgetExportPayload): Promise<BudgetImportResult> => {
      if (isGuest) {
        throw new GuestBlockedError('Guest users cannot import budget data');
      }
      return importBudgetJSON(payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.all });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error('Failed to import budget JSON:', error);
      }
    }
  });
}

export function useSkipIncome() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: (data: SkipIncomeRequest) => {
      if (isGuest) {
        throw new GuestBlockedError('Guest users cannot skip incomes');
      }

      if (data.status === 'skipped') {
        throw new Error('Income occurrence is already skipped');
      }

      return skipIncome(data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.all });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error('Failed to skip income occurrence:', error);
      }
    }
  });
}

export function useSkipExpense() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: (data: SkipExpenseRequest) => {
      if (isGuest) {
        throw new GuestBlockedError('Guest users cannot skip expenses');
      }

      if (data.status === 'skipped') {
        throw new Error('Expense occurrence is already skipped');
      }

      return skipExpense(data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.all });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error('Failed to skip expense occurrence:', error);
      }
    }
  });
}

export function useBudgetRemaining(date?: string) {
  return useQuery<BudgetRemainingResponse>({
    queryKey: budgetKeys.budgetRemaining(date),
    queryFn: () => fetchBudgetRemaining(date),
    staleTime: 0
  });
}

// Statistics

export function useSummaryStats(period?: string) {
  return useQuery<SummaryStats>({
    queryKey: budgetKeys.summary(period),
    queryFn: () => fetchSummaryStats(period),
    staleTime: 60000,
    refetchOnMount: true
  });
}

export function useCategoryBreakdown(startDate?: string, endDate?: string) {
  return useQuery<CategoryBreakdown[]>({
    queryKey: budgetKeys.breakdown(startDate, endDate),
    queryFn: () => fetchCategoryBreakdown(startDate, endDate),
    staleTime: 60000
  });
}

export function useTrends(granularity: 'day' | 'month' = 'month', startDate?: string, endDate?: string) {
  return useQuery<TrendItem[]>({
    queryKey: budgetKeys.trends(granularity, startDate, endDate),
    queryFn: () => fetchTrends(granularity, startDate, endDate),
    staleTime: 60000
  });
}

// Import new API functions
import {
  fetchSavingsRate,
  fetchSpendingVelocity,
  fetchUpcomingBills,
  createCategoryBudget,
  fetchCategoryBudgets,
  updateCategoryBudget,
  deleteCategoryBudget,
  fetchPriorityGroups,
  fetchHealthScore,
  fetchFiftyThirtyTwenty,
  fetchSubscriptions,
  fetchWeekdayPattern,
  fetchMonthOverMonth,
  fetchMerchantAnalysis,
  fetchCurrentTotalMoney
} from '@/lib/api';

import type {
  PriorityGroup,
  SavingsRateResponse,
  SpendingVelocityResponse,
  UpcomingBillsResponse,
  CategoryBudgetWithVariance,
  CreateCategoryBudgetRequest,
  UpdateCategoryBudgetRequest,
  HealthScoreResponse,
  FiftyThirtyTwentyResponse,
  SubscriptionsResponse,
  WeekdayPatternResponse,
  MonthOverMonthResponse,
  MerchantAnalysisResponse,
  CurrentTotalMoneyResponse
} from '@/types/api';

// Budget Analytics Hooks

export function useSavingsRate(startDate?: string, endDate?: string) {
  return useQuery<SavingsRateResponse>({
    queryKey: budgetKeys.savingsRate(startDate, endDate),
    queryFn: () => fetchSavingsRate(startDate, endDate),
    staleTime: 300000 // 5 minutes
  });
}

export function useSpendingVelocity(startDate?: string, endDate?: string) {
  return useQuery<SpendingVelocityResponse>({
    queryKey: budgetKeys.spendingVelocity(startDate, endDate),
    queryFn: () => fetchSpendingVelocity(startDate, endDate),
    staleTime: 60000 // 1 minute
  });
}

export function useUpcomingBills(days: 7 | 30 = 30) {
  return useQuery<UpcomingBillsResponse>({
    queryKey: budgetKeys.upcomingBills(days),
    queryFn: () => fetchUpcomingBills(days),
    staleTime: 300000 // 5 minutes
  });
}

// Category Budget Hooks

export function useCategoryBudgets(month?: string) {
  return useQuery<CategoryBudgetWithVariance[]>({
    queryKey: budgetKeys.categoryBudgets(month),
    queryFn: () => fetchCategoryBudgets(month),
    staleTime: 300000 // 5 minutes
  });
}

export function useCreateCategoryBudget() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: (data: CreateCategoryBudgetRequest) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return createCategoryBudget(data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.categoryBudgets() });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to create category budget:", error);
      }
    }
  });
}

export function useUpdateCategoryBudget() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateCategoryBudgetRequest }) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return updateCategoryBudget(id, data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.categoryBudgets() });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to update category budget:", error);
      }
    }
  });
}

export function useDeleteCategoryBudget() {
  const queryClient = useQueryClient();
  const { isGuest } = useAuth();

  return useMutation({
    mutationFn: (id: string) => {
      if (isGuest) {
        throw new GuestBlockedError();
      }
      return deleteCategoryBudget(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.categoryBudgets() });
    },
    onError: (error: Error) => {
      if (!(error instanceof GuestBlockedError)) {
        console.error("Failed to delete category budget:", error);
      }
    }
  });
}

export function usePriorityGroups() {
  return useQuery<PriorityGroup[]>({
    queryKey: budgetKeys.priorityGroups(),
    queryFn: fetchPriorityGroups,
    staleTime: 300000 // 5 minutes - priority groups are static reference data
  });
}

// Financial Health Hooks

export function useHealthScore() {
  return useQuery<HealthScoreResponse>({
    queryKey: budgetKeys.healthScore(),
    queryFn: fetchHealthScore,
    staleTime: 300000
  });
}

export function useFiftyThirtyTwenty(startDate?: string, endDate?: string) {
  return useQuery<FiftyThirtyTwentyResponse>({
    queryKey: budgetKeys.fiftyThirtyTwenty(startDate, endDate),
    queryFn: () => fetchFiftyThirtyTwenty(startDate, endDate),
    staleTime: 300000
  });
}

export function useWeekdayPattern(startDate?: string, endDate?: string) {
  return useQuery<WeekdayPatternResponse>({
    queryKey: budgetKeys.weekdayPattern(startDate, endDate),
    queryFn: () => fetchWeekdayPattern(startDate, endDate),
    staleTime: 300000
  });
}

export function useMonthOverMonth() {
  return useQuery<MonthOverMonthResponse>({
    queryKey: budgetKeys.monthOverMonth(),
    queryFn: fetchMonthOverMonth,
    staleTime: 300000
  });
}

// Merchant Analysis Hooks

export function useMerchantAnalysis(limit: number = 10, startDate?: string, endDate?: string) {
  return useQuery<MerchantAnalysisResponse>({
    queryKey: budgetKeys.merchantAnalysis(limit, startDate, endDate),
    queryFn: () => fetchMerchantAnalysis(limit, startDate, endDate),
    staleTime: 300000
  });
}

export function useSubscriptions() {
  return useQuery<SubscriptionsResponse>({
    queryKey: budgetKeys.subscriptions(),
    queryFn: fetchSubscriptions,
    staleTime: 300000
  });
}

export function useCurrentTotalMoney(startDate?: string, endDate?: string) {
  return useQuery<CurrentTotalMoneyResponse>({
    queryKey: budgetKeys.currentTotalMoney(startDate, endDate),
    queryFn: () => fetchCurrentTotalMoney(startDate, endDate),
    staleTime: 300000
  });
}

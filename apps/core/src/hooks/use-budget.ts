import { useQuery, useMutation, useQueryClient, useInfiniteQuery } from '@tanstack/react-query';
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
  fetchExpensesPaginated,
  fetchIncomes,
  fetchIncome,
  createIncome,
  updateIncome,
  deleteIncome,
  fetchIncomesPaginated,
  fetchBudgetRemaining,
  fetchSummaryStats,
  fetchCategoryBreakdown,
  fetchTrends
} from '@/lib/api';
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
  CreateIncomeRequest,
  UpdateIncomeRequest,
  BudgetRemainingResponse,
  SummaryStats,
  CategoryBreakdown,
  TrendItem,
  CursorPagination,
  OffsetPagination
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

  tags: () => [...budgetKeys.all, 'tags'] as const,
  tagsList: () => [...budgetKeys.tags(), 'list'] as const,
  tagDetail: (id: string) => [...budgetKeys.tags(), 'detail', id] as const,

  expenses: () => [...budgetKeys.all, 'expenses'] as const,
  expensesList: (filters?: ExpenseFilters) => [...budgetKeys.expenses(), 'list', filters] as const,
  expenseDetail: (id: string) => [...budgetKeys.expenses(), 'detail', id] as const,

  incomes: () => [...budgetKeys.all, 'incomes'] as const,
  incomesList: (startDate?: string, endDate?: string, recurringType?: 'daily') => 
    [...budgetKeys.incomes(), 'list', { startDate, endDate, recurringType }] as const,
  incomeDetail: (id: string) => [...budgetKeys.incomes(), 'detail', id] as const,

  budgetRemaining: (date?: string) => [...budgetKeys.all, 'budgetRemaining', date] as const,

  stats: () => [...budgetKeys.all, 'stats'] as const,
  summary: (period?: string) => [...budgetKeys.stats(), 'summary', period] as const,
  breakdown: (startDate?: string, endDate?: string) =>
    [...budgetKeys.stats(), 'breakdown', { startDate, endDate }] as const,
  trends: (granularity: 'day' | 'month', startDate?: string, endDate?: string) =>
    [...budgetKeys.stats(), 'trends', { granularity, startDate, endDate }] as const
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

export function useExpenses(filters?: ExpenseFilters) {
  return useQuery<Expense[]>({
    queryKey: budgetKeys.expensesList(filters),
    queryFn: () => fetchExpenses(filters),
    staleTime: 60000,
    refetchOnMount: true
  });
}

export function useExpensesPaginated(filters?: ExpenseFilters, limit: number = 20) {
  return useInfiniteQuery({
    queryKey: [...budgetKeys.expensesList(filters), 'paginated'],
    queryFn: ({ pageParam }) => 
      fetchExpensesPaginated({ 
        cursor: pageParam, 
        limit,
        ...filters 
      }),
    getNextPageParam: (lastPage) => {
      const pagination = lastPage.pagination as CursorPagination;
      return pagination.hasMore ? pagination.nextCursor : undefined;
    },
    initialPageParam: undefined as string | undefined,
    staleTime: 60000
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

// Incomes

export function useIncomes(startDate?: string, endDate?: string, recurringType?: 'daily') {
  return useQuery<Income[]>({
    queryKey: budgetKeys.incomesList(startDate, endDate, recurringType),
    queryFn: () => fetchIncomes(startDate, endDate, recurringType),
    staleTime: 60000,
    refetchOnMount: true
  });
}

export function useIncomesPaginated(
  filters?: { start_date?: string; end_date?: string; recurring_type?: string },
  limit: number = 10
) {
  return useInfiniteQuery({
    queryKey: [...budgetKeys.incomes(), 'paginated', filters],
    queryFn: ({ pageParam = 1 }) => 
      fetchIncomesPaginated({ 
        page: pageParam as number,
        limit,
        ...filters 
      }),
    getNextPageParam: (lastPage) => {
      const pagination = lastPage.pagination as OffsetPagination;
      return pagination.hasMore ? pagination.page + 1 : undefined;
    },
    initialPageParam: 1,
    staleTime: 60000
  });
}

export function useIncome(id: string | null) {
  return useQuery<Income>({
    queryKey: budgetKeys.incomeDetail(id ?? ''),
    queryFn: () => fetchIncome(id!),
    enabled: !!id
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

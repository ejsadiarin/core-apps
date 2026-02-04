import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
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
  fetchSummaryStats,
  fetchCategoryBreakdown,
  fetchTrends,
} from "@/lib/api";
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
  SummaryStats,
  CategoryBreakdown,
  TrendItem,
} from "@/types/api";
import { useAuth } from "@/contexts/auth-context";

export class GuestBlockedError extends Error {
  constructor(message = "Guest user is read-only. Create an account to save changes") {
    super(message);
    this.name = "GuestBlockedError";
  }
}

// query keys for cache management
export const budgetKeys = {
  all: ["budget"] as const,
  
  categories: () => [...budgetKeys.all, "categories"] as const,
  categoriesList: () => [...budgetKeys.categories(), "list"] as const,
  categoryDetail: (id: string) => [...budgetKeys.categories(), "detail", id] as const,
  
  tags: () => [...budgetKeys.all, "tags"] as const,
  tagsList: () => [...budgetKeys.tags(), "list"] as const,
  tagDetail: (id: string) => [...budgetKeys.tags(), "detail", id] as const,
  
  expenses: () => [...budgetKeys.all, "expenses"] as const,
  expensesList: (filters?: ExpenseFilters) => [...budgetKeys.expenses(), "list", filters] as const,
  expenseDetail: (id: string) => [...budgetKeys.expenses(), "detail", id] as const,
  
  stats: () => [...budgetKeys.all, "stats"] as const,
  summary: (period?: string) => [...budgetKeys.stats(), "summary", period] as const,
  breakdown: (startDate?: string, endDate?: string) => 
    [...budgetKeys.stats(), "breakdown", { startDate, endDate }] as const,
  trends: (granularity: "day" | "month", startDate?: string, endDate?: string) =>
    [...budgetKeys.stats(), "trends", { granularity, startDate, endDate }] as const,
};

// Categories

export function useCategories() {
  return useQuery<Category[]>({
    queryKey: budgetKeys.categoriesList(),
    queryFn: fetchCategories,
    staleTime: 300000, // 5 minutes - categories don't change often
  });
}

export function useCategory(id: string | null) {
  return useQuery<Category>({
    queryKey: budgetKeys.categoryDetail(id ?? ""),
    queryFn: () => fetchCategory(id!),
    enabled: !!id,
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
    },
  });
}

// Tags

export function useTags() {
  return useQuery<Tag[]>({
    queryKey: budgetKeys.tagsList(),
    queryFn: fetchTags,
    staleTime: 300000, // 5 minutes
  });
}

export function useTag(id: string | null) {
  return useQuery<Tag>({
    queryKey: budgetKeys.tagDetail(id ?? ""),
    queryFn: () => fetchTag(id!),
    enabled: !!id,
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
    },
  });
}
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTagRequest }) =>
      updateTag(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.tags() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.tagDetail(variables.id) });
    },
  });
}

export function useDeleteTag() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => deleteTag(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.tags() });
    },
  });
}

// Expenses

export function useExpenses(filters?: ExpenseFilters) {
  return useQuery<Expense[]>({
    queryKey: budgetKeys.expensesList(filters),
    queryFn: () => fetchExpenses(filters),
    staleTime: 60000, // 1 minute
  });
}

export function useExpense(id: string | null) {
  return useQuery<Expense>({
    queryKey: budgetKeys.expenseDetail(id ?? ""),
    queryFn: () => fetchExpense(id!),
    enabled: !!id,
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
      queryClient.invalidateQueries({ queryKey: budgetKeys.stats() });
    },
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
      queryClient.invalidateQueries({ queryKey: budgetKeys.stats() });
    },
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
      queryClient.invalidateQueries({ queryKey: budgetKeys.stats() });
    },
  });
}

// Statistics

export function useSummaryStats(period?: string) {
  return useQuery<SummaryStats>({
    queryKey: budgetKeys.summary(period),
    queryFn: () => fetchSummaryStats(period),
    staleTime: 60000, // 1 minute
  });
}

export function useCategoryBreakdown(startDate?: string, endDate?: string) {
  return useQuery<CategoryBreakdown[]>({
    queryKey: budgetKeys.breakdown(startDate, endDate),
    queryFn: () => fetchCategoryBreakdown(startDate, endDate),
    staleTime: 60000,
  });
}

export function useTrends(
  granularity: "day" | "month" = "month",
  startDate?: string,
  endDate?: string
) {
  return useQuery<TrendItem[]>({
    queryKey: budgetKeys.trends(granularity, startDate, endDate),
    queryFn: () => fetchTrends(granularity, startDate, endDate),
    staleTime: 60000,
  });
}

export { GuestBlockedError };

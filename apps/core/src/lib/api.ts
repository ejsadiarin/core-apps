import type {
  Service,
  CreateServiceRequest,
  UpdateServiceRequest,
  ServiceHealthHistory,
  ServiceStats,
  Category,
  PriorityGroup,
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
  User,
  LoginRequest,
  RegisterRequest,
  CreateUserRequest,
  UpdateUserRequest,
  Income,
  CreateIncomeRequest,
  UpdateIncomeRequest,
  RecurringIncomeWithNextDate,
  RecurringSummary,
  BudgetRemainingResponse,
  PaginatedResponse,
  PaginationParams,
  SavingsRateResponse,
  SpendingVelocityResponse,
  UpcomingBillsResponse,
  CategoryBudget,
  CategoryBudgetWithVariance,
  CreateCategoryBudgetRequest,
  UpdateCategoryBudgetRequest,
  HealthScoreResponse,
  FiftyThirtyTwentyResponse,
  WeekdayPatternResponse,
  MonthOverMonthResponse,
  MerchantAnalysisResponse,
  SubscriptionsResponse,
  CurrentTotalMoneyResponse
} from '@/types/api';

const url = 'http://localhost:8080';

// Auth API Functions

export async function login(data: LoginRequest): Promise<User> {
  const res = await fetch(url + '/api/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    const error = await res.json().catch(() => ({ message: 'Login failed' }));
    throw new Error(error.message || 'Login failed');
  }
  return res.json();
}

export async function register(data: RegisterRequest): Promise<User> {
  const res = await fetch(url + '/api/auth/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    // TODO: detect duplicate email error and send "Email already exists"
    const error = await res.json().catch(() => ({ message: 'Registration failed' }));
    throw new Error(error.message || 'Registration failed');
  }
  return res.json();
}

export async function logout(): Promise<void> {
  const res = await fetch(url + '/api/auth/logout', {
    method: 'POST',
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error('Logout failed');
  }
}

export async function fetchCurrentUser(): Promise<User> {
  const res = await fetch(url + '/api/auth/me', {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error('Not authenticated');
  }
  return res.json();
}

export async function demoLogin(): Promise<User> {
  const res = await fetch(url + '/api/auth/demo', {
    method: 'POST',
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error('Demo login failed');
  }
  return res.json();
}

// User Management API Functions (Admin)

export async function fetchUsers(): Promise<User[]> {
  const res = await fetch(url + '/api/users', {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching users: ${res.status}`);
  }
  return res.json();
}

export async function fetchUser(id: string): Promise<User> {
  const res = await fetch(url + `/api/users/${id}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching user: ${res.status}`);
  }
  return res.json();
}

export async function createUser(data: CreateUserRequest): Promise<User> {
  const res = await fetch(url + '/api/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    const error = await res.json().catch(() => ({ message: 'Create user failed' }));
    throw new Error(error.message || 'Create user failed');
  }
  return res.json();
}

export async function updateUser(id: string, data: UpdateUserRequest): Promise<User> {
  const res = await fetch(url + `/api/users/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    const error = await res.json().catch(() => ({ message: 'Update user failed' }));
    throw new Error(error.message || 'Update user failed');
  }
  return res.json();
}

export async function deleteUser(id: string): Promise<void> {
  const res = await fetch(url + `/api/users/${id}`, {
    method: 'DELETE',
    credentials: 'include'
  });
  if (!res.ok) {
    const error = await res.json().catch(() => ({ message: 'Delete user failed' }));
    throw new Error(error.message || 'Delete user failed');
  }
}

export async function fetchSystemStatus() {
  const res = await fetch(url + '/api/system/stats');
  if (!res.ok) {
    throw new Error(`Error fetching ${res.status}`);
  }
  return res.json();
}

// Budget API Functions

// Categories

export async function fetchCategories(): Promise<Category[]> {
  const res = await fetch(url + '/api/budget/categories', {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching categories: ${res.status}`);
  }
  return res.json();
}

export async function fetchCategory(id: string): Promise<Category> {
  const res = await fetch(url + `/api/budget/categories/${id}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching category: ${res.status}`);
  }
  return res.json();
}

export async function createCategory(data: CreateCategoryRequest): Promise<Category> {
  const res = await fetch(url + '/api/budget/categories', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(`Error creating category: ${res.status}`);
  }
  return res.json();
}

export async function updateCategory(id: string, data: UpdateCategoryRequest): Promise<Category> {
  const res = await fetch(url + `/api/budget/categories/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(`Error updating category: ${res.status}`);
  }
  return res.json();
}

export async function deleteCategory(id: string): Promise<void> {
  const res = await fetch(url + `/api/budget/categories/${id}`, {
    method: 'DELETE',
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error deleting category: ${res.status}`);
  }
}

// Priority Groups

export async function fetchPriorityGroups(): Promise<PriorityGroup[]> {
  const res = await fetch(url + '/api/budget/priority-groups', {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching priority groups: ${res.status}`);
  }
  return res.json();
}

// Tags

export async function fetchTags(): Promise<Tag[]> {
  const res = await fetch(url + '/api/budget/tags', {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching tags: ${res.status}`);
  }
  return res.json();
}

export async function fetchTag(id: string): Promise<Tag> {
  const res = await fetch(url + `/api/budget/tags/${id}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching tag: ${res.status}`);
  }
  return res.json();
}

export async function createTag(data: CreateTagRequest): Promise<Tag> {
  const res = await fetch(url + '/api/budget/tags', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    if (res.status === 403) {
      throw new Error(`Error creating tag: ${res.status}`);
    }
    throw new Error(`Error creating tag: ${res.status}`);
  }
  return res.json();
}

export async function updateTag(id: string, data: UpdateTagRequest): Promise<Tag> {
  const res = await fetch(url + `/api/budget/tags/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(`Error updating tag: ${res.status}`);
  }
  return res.json();
}

export async function deleteTag(id: string): Promise<void> {
  const res = await fetch(url + `/api/budget/tags/${id}`, {
    method: 'DELETE',
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error deleting tag: ${res.status}`);
  }
}

// Expenses

export async function fetchExpenses(
  params: PaginationParams & ExpenseFilters
): Promise<PaginatedResponse<Expense>> {
  const searchParams = new URLSearchParams();
  if (params.page) searchParams.append('page', params.page.toString());
  if (params.limit) searchParams.append('limit', params.limit.toString());
  if (params.start_date) searchParams.append('start_date', params.start_date);
  if (params.end_date) searchParams.append('end_date', params.end_date);
  if (params.category_id) searchParams.append('category_id', params.category_id);

  const res = await fetch(url + `/api/budget/expenses?${searchParams.toString()}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching expenses: ${res.status}`);
  }
  return res.json();
}

export async function fetchExpense(id: string): Promise<Expense> {
  const res = await fetch(url + `/api/budget/expenses/${id}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching expense: ${res.status}`);
  }
  return res.json();
}

export async function createExpense(data: CreateExpenseRequest): Promise<Expense> {
  const res = await fetch(url + '/api/budget/expenses', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(`Error creating expense: ${res.status}`);
  }
  return res.json();
}

export async function updateExpense(id: string, data: UpdateExpenseRequest): Promise<Expense> {
  const res = await fetch(url + `/api/budget/expenses/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(`Error updating expense: ${res.status}`);
  }
  return res.json();
}

export async function deleteExpense(id: string): Promise<void> {
  const res = await fetch(url + `/api/budget/expenses/${id}`, {
    method: 'DELETE',
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error deleting expense: ${res.status}`);
  }
}

export interface ExpenseSearchParams {
  q: string;
  page?: number;
  limit?: number;
  category_id?: string;
  start_date?: string;
  end_date?: string;
}

export async function searchExpenses(
  params: ExpenseSearchParams
): Promise<PaginatedResponse<Expense>> {
  const searchParams = new URLSearchParams();
  searchParams.append('q', params.q);
  if (params.page) searchParams.append('page', params.page.toString());
  if (params.limit) searchParams.append('limit', params.limit.toString());
  if (params.start_date) searchParams.append('start_date', params.start_date);
  if (params.end_date) searchParams.append('end_date', params.end_date);
  if (params.category_id) searchParams.append('category_id', params.category_id);

  const res = await fetch(url + `/api/budget/expenses/search?${searchParams.toString()}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error searching expenses: ${res.status}`);
  }
  return res.json();
}

// Statistics

export async function fetchSummaryStats(period?: string): Promise<SummaryStats> {
  const params = period ? `?period=${period}` : '';
  const res = await fetch(url + `/api/budget/stats/summary${params}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching summary stats: ${res.status}`);
  }
  return res.json();
}

export async function fetchCategoryBreakdown(startDate?: string, endDate?: string): Promise<CategoryBreakdown[]> {
  const params = new URLSearchParams();
  if (startDate) params.append('start_date', startDate);
  if (endDate) params.append('end_date', endDate);

  const queryString = params.toString();
  const endpoint = queryString
    ? `/api/budget/stats/category-breakdown?${queryString}`
    : '/api/budget/stats/category-breakdown';

  const res = await fetch(url + endpoint, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching category breakdown: ${res.status}`);
  }
  return res.json();
}

export async function fetchTrends(
  granularity: 'day' | 'month' = 'month',
  startDate?: string,
  endDate?: string
): Promise<TrendItem[]> {
  const params = new URLSearchParams();
  params.append('granularity', granularity);
  if (startDate) params.append('start_date', startDate);
  if (endDate) params.append('end_date', endDate);

  const res = await fetch(url + `/api/budget/stats/trends?${params.toString()}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching trends: ${res.status}`);
  }
  return res.json();
}

// Service Monitoring API Functions

export async function fetchServices(): Promise<Service[]> {
  const res = await fetch(url + '/api/services/list');
  if (!res.ok) {
    throw new Error(`Error fetching services: ${res.status}`);
  }
  return res.json();
}

export async function fetchService(id: string): Promise<Service> {
  const res = await fetch(url + `/api/services/${id}`);
  if (!res.ok) {
    throw new Error(`Error fetching service: ${res.status}`);
  }
  return res.json();
}

export async function createService(data: CreateServiceRequest): Promise<Service> {
  const res = await fetch(url + '/api/services', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(`Error creating service: ${res.status}`);
  }
  return res.json();
}

export async function updateService(id: string, data: UpdateServiceRequest): Promise<Service> {
  const res = await fetch(url + `/api/services/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(`Error updating service: ${res.status}`);
  }
  return res.json();
}

export async function deleteService(id: string): Promise<void> {
  const res = await fetch(url + `/api/services/${id}`, {
    method: 'DELETE'
  });
  if (!res.ok) {
    throw new Error(`Error deleting service: ${res.status}`);
  }
}

export async function triggerHealthCheck(id: string): Promise<void> {
  const res = await fetch(url + `/api/services/${id}/check`, {
    method: 'POST'
  });
  if (!res.ok) {
    throw new Error(`Error triggering health check: ${res.status}`);
  }
}

export async function fetchServiceHistory(id: string): Promise<ServiceHealthHistory[]> {
  const res = await fetch(url + `/api/services/${id}/history`);
  if (!res.ok) {
    throw new Error(`Error fetching service history: ${res.status}`);
  }
  return res.json();
}

export async function fetchServiceStats(id: string): Promise<ServiceStats> {
  const res = await fetch(url + `/api/services/${id}/stats`);
  if (!res.ok) {
    throw new Error(`Error fetching service stats: ${res.status}`);
  }
  return res.json();
}

export async function fetchAllServicesStats(): Promise<{
  total_services: number;
  total_checks: number;
  successful_checks: number;
  overall_uptime: number;
  avg_response_time: number;
}> {
  const res = await fetch(url + '/api/services/stats/all');
  if (!res.ok) {
    throw new Error(`Error fetching services stats: ${res.status}`);
  }
  return res.json();
}

// Incomes

export async function fetchIncomes(
  params: PaginationParams & { 
    start_date?: string; 
    end_date?: string; 
    recurring_type?: string;
  }
): Promise<PaginatedResponse<Income>> {
  const searchParams = new URLSearchParams();
  if (params.page) searchParams.append('page', params.page.toString());
  if (params.limit) searchParams.append('limit', params.limit.toString());
  if (params.start_date) searchParams.append('start_date', params.start_date);
  if (params.end_date) searchParams.append('end_date', params.end_date);
  if (params.recurring_type) searchParams.append('recurring_type', params.recurring_type);

  const res = await fetch(url + `/api/budget/incomes?${searchParams.toString()}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching incomes: ${res.status}`);
  }
  return res.json();
}

export async function fetchIncome(id: string): Promise<Income> {
  const res = await fetch(url + `/api/budget/incomes/${id}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching income: ${res.status}`);
  }
  return res.json();
}

export async function createIncome(data: CreateIncomeRequest): Promise<Income> {
  const res = await fetch(url + '/api/budget/incomes', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(`Error creating income: ${res.status}`);
  }
  return res.json();
}

export async function updateIncome(id: string, data: UpdateIncomeRequest): Promise<Income> {
  const res = await fetch(url + `/api/budget/incomes/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(`Error updating income: ${res.status}`);
  }
  return res.json();
}

export async function deleteIncome(id: string): Promise<void> {
  const res = await fetch(url + `/api/budget/incomes/${id}`, {
    method: 'DELETE',
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error deleting income: ${res.status}`);
  }
}

export async function fetchRecurringIncomes(): Promise<RecurringIncomeWithNextDate[]> {
  const res = await fetch(url + '/api/budget/recurring-incomes', {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching recurring incomes: ${res.status}`);
  }
  return res.json();
}

export async function checkSkippedIncome(date: string): Promise<boolean> {
  const res = await fetch(url + `/api/budget/incomes/check-skipped?date=${date}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error checking skipped income: ${res.status}`);
  }
  return res.json();
}

export async function checkSkippedExpense(date: string): Promise<boolean> {
  const res = await fetch(url + `/api/budget/expenses/check-skipped?date=${date}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error checking skipped expense: ${res.status}`);
  }
  return res.json();
}

export async function cancelRecurringExpense(id: string): Promise<Expense> {
  const today = new Date().toISOString().split('T')[0];
  return updateExpense(id, { end_date: today });
}

export async function cancelRecurringIncome(id: string): Promise<Income> {
  const today = new Date().toISOString().split('T')[0];
  return updateIncome(id, { end_date: today });
}

export async function fetchBudgetRemaining(date?: string): Promise<BudgetRemainingResponse> {
  const params = date ? `?date=${date}` : '';
  const res = await fetch(url + `/api/budget/remaining${params}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching budget remaining: ${res.status}`);
  }
  return res.json();
}

// Budget Analytics API

export async function fetchSavingsRate(startDate?: string, endDate?: string): Promise<SavingsRateResponse> {
  const params = new URLSearchParams();
  if (startDate) params.append('start_date', startDate);
  if (endDate) params.append('end_date', endDate);

  const queryString = params.toString();
  const endpoint = queryString
    ? `/api/budget/stats/savings-rate?${queryString}`
    : '/api/budget/stats/savings-rate';

  const res = await fetch(url + endpoint, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching savings rate: ${res.status}`);
  }
  return res.json();
}

export async function fetchSpendingVelocity(startDate?: string, endDate?: string): Promise<SpendingVelocityResponse> {
  const params = new URLSearchParams();
  if (startDate) params.append('start_date', startDate);
  if (endDate) params.append('end_date', endDate);

  const queryString = params.toString();
  const endpoint = queryString
    ? `/api/budget/velocity?${queryString}`
    : '/api/budget/velocity';

  const res = await fetch(url + endpoint, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching spending velocity: ${res.status}`);
  }
  return res.json();
}

export async function fetchUpcomingBills(days: 7 | 30 = 30): Promise<UpcomingBillsResponse> {
  const res = await fetch(url + `/api/budget/forecast/upcoming?days=${days}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching upcoming bills: ${res.status}`);
  }
  return res.json();
}

// Category Budget API

export async function createCategoryBudget(data: CreateCategoryBudgetRequest): Promise<CategoryBudget> {
  const res = await fetch(url + '/api/budget/category-budgets', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(`Error creating category budget: ${res.status}`);
  }
  return res.json();
}

export async function fetchCategoryBudgets(month?: string): Promise<CategoryBudgetWithVariance[]> {
  const params = month ? `?month=${month}` : '';
  const res = await fetch(url + `/api/budget/category-budgets${params}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching category budgets: ${res.status}`);
  }
  return res.json();
}

export async function updateCategoryBudget(id: string, data: UpdateCategoryBudgetRequest): Promise<CategoryBudget> {
  const res = await fetch(url + `/api/budget/category-budgets/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(`Error updating category budget: ${res.status}`);
  }
  return res.json();
}

export async function deleteCategoryBudget(id: string): Promise<void> {
  const res = await fetch(url + `/api/budget/category-budgets/${id}`, {
    method: 'DELETE',
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error deleting category budget: ${res.status}`);
  }
}

// Financial Health API

export async function fetchHealthScore(): Promise<HealthScoreResponse> {
  const res = await fetch(url + '/api/budget/stats/health-score', {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching health score: ${res.status}`);
  }
  return res.json();
}

export async function fetchFiftyThirtyTwenty(startDate?: string, endDate?: string): Promise<FiftyThirtyTwentyResponse> {
  const params = new URLSearchParams();
  if (startDate) params.append('start_date', startDate);
  if (endDate) params.append('end_date', endDate);

  const queryString = params.toString();
  const endpoint = queryString
    ? `/api/budget/analysis/503020?${queryString}`
    : '/api/budget/analysis/503020';

  const res = await fetch(url + endpoint, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching 50/30/20 analysis: ${res.status}`);
  }
  return res.json();
}

export async function fetchWeekdayPattern(startDate?: string, endDate?: string): Promise<WeekdayPatternResponse> {
  const params = new URLSearchParams();
  if (startDate) params.append('start_date', startDate);
  if (endDate) params.append('end_date', endDate);

  const queryString = params.toString();
  const endpoint = queryString
    ? `/api/budget/analysis/weekday-pattern?${queryString}`
    : '/api/budget/analysis/weekday-pattern';

  const res = await fetch(url + endpoint, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching weekday pattern: ${res.status}`);
  }
  return res.json();
}

export async function fetchMonthOverMonth(): Promise<MonthOverMonthResponse> {
  const res = await fetch(url + '/api/budget/trends/month-over-month', {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching month-over-month trends: ${res.status}`);
  }
  return res.json();
}

export async function fetchMerchantAnalysis(limit: number = 10, startDate?: string, endDate?: string): Promise<MerchantAnalysisResponse> {
  const params = new URLSearchParams();
  params.append('limit', limit.toString());
  if (startDate) params.append('start_date', startDate);
  if (endDate) params.append('end_date', endDate);

  const res = await fetch(url + `/api/budget/analysis/merchants?${params.toString()}`, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching merchant analysis: ${res.status}`);
  }
  return res.json();
}

export async function fetchSubscriptions(): Promise<SubscriptionsResponse> {
  const res = await fetch(url + '/api/budget/subscriptions', {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching subscriptions: ${res.status}`);
  }
  return res.json();
}

export async function fetchCurrentTotalMoney(startDate?: string, endDate?: string): Promise<CurrentTotalMoneyResponse> {
  const params = new URLSearchParams();
  if (startDate) params.append('start_date', startDate);
  if (endDate) params.append('end_date', endDate);

  const queryString = params.toString();
  const endpoint = queryString
    ? `/api/budget/current-total-money?${queryString}`
    : '/api/budget/current-total-money';

  const res = await fetch(url + endpoint, {
    credentials: 'include'
  });
  if (!res.ok) {
    throw new Error(`Error fetching current total money: ${res.status}`);
  }
  return res.json();
}

import type {
  Service,
  CreateServiceRequest,
  UpdateServiceRequest,
  ServiceHealthHistory,
  ServiceStats,
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
  TrendItem
} from '@/types/api';

const url = 'http://localhost:8080';

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
  const res = await fetch(url + '/api/budget/categories');
  if (!res.ok) {
    throw new Error(`Error fetching categories: ${res.status}`);
  }
  return res.json();
}

export async function fetchCategory(id: string): Promise<Category> {
  const res = await fetch(url + `/api/budget/categories/${id}`);
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
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(`Error updating category: ${res.status}`);
  }
  return res.json();
}

export async function deleteCategory(id: string): Promise<void> {
  const res = await fetch(url + `/api/budget/categories/${id}`, {
    method: 'DELETE'
  });
  if (!res.ok) {
    throw new Error(`Error deleting category: ${res.status}`);
  }
}

// Tags

export async function fetchTags(): Promise<Tag[]> {
  const res = await fetch(url + '/api/budget/tags');
  if (!res.ok) {
    throw new Error(`Error fetching tags: ${res.status}`);
  }
  return res.json();
}

export async function fetchTag(id: string): Promise<Tag> {
  const res = await fetch(url + `/api/budget/tags/${id}`);
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
    body: JSON.stringify(data)
  });
  if (!res.ok) {
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
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(`Error updating tag: ${res.status}`);
  }
  return res.json();
}

export async function deleteTag(id: string): Promise<void> {
  const res = await fetch(url + `/api/budget/tags/${id}`, {
    method: 'DELETE'
  });
  if (!res.ok) {
    throw new Error(`Error deleting tag: ${res.status}`);
  }
}

// Expenses

export async function fetchExpenses(filters?: ExpenseFilters): Promise<Expense[]> {
  const params = new URLSearchParams();
  if (filters?.start_date) params.append('start_date', filters.start_date);
  if (filters?.end_date) params.append('end_date', filters.end_date);
  if (filters?.category_id) params.append('category_id', filters.category_id);

  const queryString = params.toString();
  const endpoint = queryString ? `/api/budget/expenses?${queryString}` : '/api/budget/expenses';

  const res = await fetch(url + endpoint);
  if (!res.ok) {
    throw new Error(`Error fetching expenses: ${res.status}`);
  }
  return res.json();
}

export async function fetchExpense(id: string): Promise<Expense> {
  const res = await fetch(url + `/api/budget/expenses/${id}`);
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
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(`Error updating expense: ${res.status}`);
  }
  return res.json();
}

export async function deleteExpense(id: string): Promise<void> {
  const res = await fetch(url + `/api/budget/expenses/${id}`, {
    method: 'DELETE'
  });
  if (!res.ok) {
    throw new Error(`Error deleting expense: ${res.status}`);
  }
}

// Statistics

export async function fetchSummaryStats(period?: string): Promise<SummaryStats> {
  const params = period ? `?period=${period}` : '';
  const res = await fetch(url + `/api/budget/stats/summary${params}`);
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

  const res = await fetch(url + endpoint);
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

  const res = await fetch(url + `/api/budget/stats/trends?${params.toString()}`);
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

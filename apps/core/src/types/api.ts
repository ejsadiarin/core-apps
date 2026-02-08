export interface SystemStats {
    cpu: number;
    memory: number;
    storage: number;
    temperature: number;
    uptime: string;
    network: {
        up: string;
        down: string;
    };
}

export interface ServiceStatus {
    name: string;
    status: "online" | "offline" | "degraded" | "maintenance" | "unknown";
    type: string;
}

export interface Service {
    id: string;
    name: string;
    url: string;
    icon?: string;
    description?: string;
    service_type?: string;
    health_check_interval: number;
    health_check_method: string;
    expected_status_codes: number[];
    timeout: number;
    created_at: string;
    updated_at: string;
    is_active: boolean;
    current_status?:
        | "online"
        | "offline"
        | "degraded"
        | "maintenance"
        | "unknown";
    last_check?: string;
    response_time?: number;
}

export interface ServiceHealthHistory {
    id: string;
    service_id: string;
    status: string;
    response_time?: number;
    status_code?: number;
    error_message?: string;
    checked_at: string;
}

export interface ServiceStats {
    service_id: string;
    uptime_24h: number;
    uptime_7d: number;
    uptime_30d: number;
    avg_response_time: number;
    total_checks: number;
    successful_checks: number;
}

export interface CreateServiceRequest {
    name: string;
    url: string;
    icon?: string;
    description?: string;
    service_type?: string;
    health_check_interval?: number;
    health_check_method?: string;
    expected_status_codes?: number[];
    timeout?: number;
}

export interface UpdateServiceRequest {
    name?: string;
    url?: string;
    icon?: string;
    description?: string;
    service_type?: string;
    health_check_interval?: number;
    health_check_method?: string;
    expected_status_codes?: number[];
    timeout?: number;
    is_active?: boolean;
}

// Budget Types

export interface Category {
    id: string;
    name: string;
    color?: string;
    icon?: string;
}

export interface Tag {
    id: string;
    name: string;
    color?: string;
}

export interface Expense {
    id: string;
    description: string;
    amount: number;
    currency: string;
    category?: Category;
    expense_date: string;
    notes?: string;
    tags?: Tag[];
    created_at: string;
    updated_at: string;
}

export interface CreateCategoryRequest {
    name: string;
    color?: string;
    icon?: string;
}

export interface UpdateCategoryRequest {
    name?: string;
    color?: string;
    icon?: string;
}

export interface CreateTagRequest {
    name: string;
    color?: string;
}

export interface UpdateTagRequest {
    name?: string;
    color?: string;
}

export interface CreateExpenseRequest {
    description: string;
    amount: number;
    currency?: string;
    category_id?: string;
    expense_date: string;
    notes?: string;
    tag_ids?: string[];
}

export interface UpdateExpenseRequest {
    description?: string;
    amount?: number;
    currency?: string;
    category_id?: string;
    expense_date?: string;
    notes?: string;
    tag_ids?: string[];
}

export interface ExpenseFilters {
    start_date?: string;
    end_date?: string;
    category_id?: string;
}

export interface Income {
    id: string;
    amount: number;
    currency: string;
    date: string;
    description?: string;
    recurring_type?: "daily" | "weekly" | "monthly" | null;
    start_date?: string;
    end_date?: string;
    created_at: string;
    updated_at: string;
}

export interface CreateIncomeRequest {
    amount: number;
    currency?: string;
    date: string;
    description?: string;
    recurring_type?: "daily" | "weekly" | "monthly" | null;
    start_date?: string;
    end_date?: string;
}

export interface UpdateIncomeRequest {
    amount?: number;
    currency?: string;
    date?: string;
    description?: string;
    recurring_type?: "daily" | "weekly" | "monthly" | null;
    start_date?: string;
    end_date?: string;
}

export interface BudgetRemainingResponse {
    budget_remaining: number;
    budget_remaining_status: "green" | "red" | "neutral";
}

export interface SummaryStats {
    total_spent: number;
    transaction_count: number;
    period: string;
    budget_remaining?: number;
    budget_remaining_status?: "green" | "red" | "neutral";
}

export interface CategoryBreakdown {
    category_id: string;
    category_name: string;
    color?: string;
    total_amount: number;
    count: number;
    percentage: number;
}

export interface TrendItem {
    date: string;
    total_amount: number;
    count: number;
}

// Auth Types

export type UserRole = "guest" | "user" | "admin";

export interface User {
    id: string;
    email: string;
    role: UserRole;
    created_at: string;
}

export interface LoginRequest {
    email: string;
    password: string;
    remember_me?: boolean;
}

export interface RegisterRequest {
    email: string;
    password: string;
}

export interface CreateUserRequest {
    email: string;
    password: string;
    role: UserRole;
}

export interface UpdateUserRequest {
    email?: string;
    password?: string;
    role?: UserRole;
}

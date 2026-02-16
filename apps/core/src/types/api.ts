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
    category_type?: "need" | "want" | "savings";
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
    recurring_type?: "daily" | "weekly" | "monthly" | "yearly" | null;
    start_date?: string;
    end_date?: string;
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
    recurring_type?: "daily" | "weekly" | "monthly" | "yearly" | null;
    start_date?: string;
    end_date?: string;
}

export interface UpdateExpenseRequest {
    description?: string;
    amount?: number;
    currency?: string;
    category_id?: string;
    expense_date?: string;
    notes?: string;
    tag_ids?: string[];
    recurring_type?: "daily" | "weekly" | "monthly" | "yearly" | null;
    start_date?: string;
    end_date?: string;
}

export interface ExpenseFilters {
    start_date?: string;
    end_date?: string;
    category_id?: string;
    recurring_type?: string;
}

// Pagination Types

export interface OffsetPagination {
    total: number;
    page: number;
    limit: number;
    totalPages: number;
    hasMore: boolean;
}

export interface PaginatedResponse<T> {
    data: T[];
    pagination: OffsetPagination;
}

export interface PaginationParams {
    page?: number;
    limit?: number;
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

// Budget Analytics Types

export interface SavingsRateResponse {
    income: number;
    expenses: number;
    savings: number;
    savings_rate: number;
    status: "excellent" | "good" | "fair" | "poor" | "negative";
    period: string;
}

export interface SpendingVelocityResponse {
    amount_spent: number;
    days_elapsed: number;
    days_in_month: number;
    projected_spend: number;
    total_budget: number;
    status: "on_track" | "warning" | "at_risk" | "over_pace" | "unknown";
}

export interface UpcomingBill {
    id: string;
    description: string;
    amount: number;
    currency: string;
    recurring_type: string;
    due_date: string;
    category_id?: string;
}

export interface UpcomingBillsResponse {
    bills: UpcomingBill[];
    total_amount: number;
    period: number;
}

// Category Budget Types

export interface CategoryBudget {
    id: string;
    category_id: string;
    month: string;
    budget_amount: number;
    created_at: string;
    updated_at: string;
}

export interface CategoryBudgetWithVariance {
    category_id: string;
    category_name: string;
    category_color?: string;
    category_type?: "need" | "want" | "savings";
    budget_amount: number;
    spent_amount: number;
    variance?: number;
    percentage: number;
    transaction_count: number;
}

export interface CreateCategoryBudgetRequest {
    category_id: string;
    month: string;
    budget_amount: number;
}

export interface UpdateCategoryBudgetRequest {
    budget_amount?: number;
}

export interface UpdateCategoryTypeRequest {
    category_type?: "need" | "want" | "savings";
}

// Financial Health Types

export interface CategoryTypeSpending {
    category_type: "need" | "want" | "savings";
    total_amount: number;
    transaction_count: number;
}

export interface FiftyThirtyTwentyData {
    needs_percentage: number;
    wants_percentage: number;
    savings_percentage: number;
    needs_amount: number;
    wants_amount: number;
    savings_amount: number;
}

export interface DayOfWeekSpending {
    day_of_week: number;
    total_amount: number;
    transaction_count: number;
    avg_amount: number;
}

export interface MerchantSpending {
    description: string;
    total_amount: number;
    transaction_count: number;
    avg_amount: number;
}

// Financial Health Response Types

export interface HealthScoreResponse {
    score: number;
    status: string;
    savings_rate: number;
    debt_to_income: number;
    emergency_fund_months: number;
    recommendations: string[];
}

export interface FiftyThirtyTwentyItem {
    category: string;
    amount: number;
    target_percentage: number;
    actual_percentage: number;
    status: string;
}

export interface FiftyThirtyTwentyResponse {
    needs: FiftyThirtyTwentyItem;
    wants: FiftyThirtyTwentyItem;
    savings: FiftyThirtyTwentyItem;
    total_income: number;
}

export interface WeekdaySpendingItem {
    day: string;
    total_amount: number;
    count: number;
    average_amount: number;
}

export interface WeekdayPatternResponse {
    weekdays: WeekdaySpendingItem[];
    highest_spending_day: string;
    lowest_spending_day: string;
}

export interface MonthOverMonthItem {
    month: string;
    income: number;
    expenses: number;
    savings: number;
    savings_rate: number;
    expense_change_percent: number;
}

export interface MonthOverMonthResponse {
    trends: MonthOverMonthItem[];
    average_savings_rate: number;
}

export interface MerchantItem {
    name: string;
    total_spent: number;
    count: number;
    average_amount: number;
    percentage: number;
    category_name?: string;
}

export interface MerchantAnalysisResponse {
    merchants: MerchantItem[];
    total_spent: number;
    unique_merchant_count: number;
}

export interface SubscriptionItem {
    id: string;
    description: string;
    amount: number;
    currency: string;
    recurring_type: string;
    category_name?: string;
    next_due_date: string;
}

export interface SubscriptionsResponse {
    subscriptions: SubscriptionItem[];
    total_monthly: number;
    count: number;
}

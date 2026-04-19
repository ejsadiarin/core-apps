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

export interface PriorityGroup {
    id: string;
    name: string;
    slug: string;
    display_order: number;
}

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
    priority_group?: PriorityGroup;
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
    priority_group_id?: string;
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
    priority_group_id?: string;
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
    status?: string;
    source_rule_id?: string;
    exclude_from_calculations?: boolean;
    created_at: string;
    updated_at: string;
}

export interface RecurringIncomeWithNextDate {
    id: string;
    amount: number;
    currency: string;
    date: string;
    description?: string;
    recurring_type?: "daily" | "weekly" | "monthly" | null;
    start_date?: string;
    end_date?: string;
    next_occurrence: string;
    monthly_equivalent: number;
}

export interface RecurringSummary {
    total_recurring_income: number;
    total_recurring_expenses: number;
    net_recurring_cash_flow: number;
    recurring_income_count: number;
    recurring_expense_count: number;
}

export interface CreateIncomeRequest {
    amount: number;
    currency?: string;
    date: string;
    description?: string;
    notes?: string;
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
    tracking_period_start?: string;
    date_range?: DateRangeMetadata;
}

export interface SpendingVelocityResponse {
    amount_spent: number;
    days_elapsed: number;
    days_in_month: number;
    projected_spend: number;
    total_budget: number;
    status: "on_track" | "warning" | "at_risk" | "over_pace" | "unknown";
    date_range?: DateRangeMetadata;
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

// Financial Health Types

// Financial Health Response Types

export interface FactorScoreBreakdown {
    savings_rate: number;
    debt_to_income: number;
    emergency_fund: number;
}

export interface HealthScoreResponse {
    score: number;
    status: string;
    savings_rate: number;
    debt_to_income: number;
    emergency_fund_months: number;
    recommendations: string[];
    factor_scores?: FactorScoreBreakdown;
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
    investments: FiftyThirtyTwentyItem;
    total_income: number;
    unclassified_count: number;
    unclassified_amount: number;
    tracking_period_start?: string;
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
    date_range?: DateRangeMetadata;
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

export interface DateRangeMetadata {
    start: string;
    end: string;
    source: string;
}

export interface CurrentTotalMoneyResponse {
    current_total: number;
    money_baseline: number;
    income_since_start: number;
    expenses_since_start: number;
    net_change: number;
    tracking_start_date: string;
    date_range?: DateRangeMetadata;
}

export interface IncomeOccurrence {
    id: string;
    source_income_id: string;
    amount: number;
    currency: string;
    date: string;
    description?: string;
    recurring_type?: 'daily' | 'weekly' | 'monthly' | null;
    is_virtual: boolean;
    is_skipped: boolean;
}

// Budget Import / Export Types

export interface ImportExportMetadata {
    schema_version: string;
    exported_at: string;
    source: string;
}

export interface ImportIncomeRecord {
    id?: string;
    amount: number;
    currency: string;
    date: string;
    description?: string;
    recurring_type?: "daily" | "weekly" | "monthly" | null;
    start_date?: string;
    end_date?: string;
}

export interface ImportExpenseRecord {
    id?: string;
    description: string;
    amount: number;
    currency: string;
    expense_date: string;
    category_name?: string;
    notes?: string;
    recurring_type?: "daily" | "weekly" | "monthly" | "yearly" | null;
    start_date?: string;
    end_date?: string;
    priority_group_id?: string;
    is_debt: boolean;
}

export interface BudgetExportPayload {
    metadata: ImportExportMetadata;
    incomes: ImportIncomeRecord[];
    expenses: ImportExpenseRecord[];
}

export type ImportMergeAction = "CREATE" | "SKIP_EXISTING" | "CONFLICT";

export interface ImportFieldDifference {
    field: string;
    incoming: unknown;
    existing: unknown;
}

export interface ImportConflictDetail {
    entity: string;
    date: string;
    existing_id?: string;
    incoming: unknown;
    existing: unknown;
    differences: ImportFieldDifference[];
}

export interface BudgetImportResultMetadata {
    imported_at: string;
    dry_run: boolean;
}

export interface BudgetImportSummary {
    incomes_created: number;
    incomes_skipped: number;
    incomes_conflicts: number;
    expenses_created: number;
    expenses_skipped: number;
    expenses_conflicts: number;
}

export interface ImportIncomeResult {
    input_index: number;
    action: ImportMergeAction;
    existing_id?: string;
    conflict?: ImportConflictDetail;
}

export interface ImportExpenseResult {
    input_index: number;
    action: ImportMergeAction;
    existing_id?: string;
    conflict?: ImportConflictDetail;
}

export interface BudgetImportResult {
    metadata: BudgetImportResultMetadata;
    summary: BudgetImportSummary;
    incomes: ImportIncomeResult[];
    expenses: ImportExpenseResult[];
    conflicts?: ImportConflictDetail[];
}

export interface SkipIncomeRequest {
    date: string;
    source_rule_id: string;
}

export interface SkipExpenseRequest {
    expense_date: string;
    source_rule_id: string;
}

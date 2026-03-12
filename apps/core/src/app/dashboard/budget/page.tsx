"use client";

import { motion } from "motion/react";
import { useState, useEffect, useMemo } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { ExpenseStats } from "@/components/budget/expense-stats";
import { IncomeForm } from "@/components/budget/income-form";
import { ExpenseFormDialog } from "@/components/budget/expense-form-dialog";
import {
  SavingsRateCard,
  SpendingVelocityCard,
  UpcomingBillsCard,
  RecurringSummaryCard,
  RecurringIncomeList,
  BudgetVarianceTable,
  CurrentTotalMoneyCard,
  PeriodPresetFilter
} from "@/components/budget";
import type { PeriodPresetFilterValue } from "@/components/budget";
import {
  useSummaryStats,
  useExpenses,
  useIncomeOccurrences,
  useCreateIncome,
  useUpdateIncome,
  useDeleteIncome,
  useUpdateExpense,
  useDeleteExpense,
  useCreateExpense,
  useBudgetRemaining,
} from "@/hooks/use-budget";
import { GuestBlockedError } from "@/hooks/use-budget";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Plus, ArrowRight } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import Link from "next/link";
import { formatDistanceToNow } from "date-fns";
import { useAuth } from "@/contexts/auth-context";
import { useToast } from "@/components/ui/toast";
import { ExpenseDetailDialog } from "@/components/budget/expense-detail-dialog";
import type { Expense, Income, CreateIncomeRequest, UpdateIncomeRequest, CreateExpenseRequest, UpdateExpenseRequest } from "@/types/api";

export default function BudgetDashboard() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [period] = useState<string>("month");
  const [selectedDate, setSelectedDate] = useState<string>(new Date().toISOString().split('T')[0]);
  
  // Date range for analytics cards
  const [dateRange, setDateRange] = useState<PeriodPresetFilterValue>({
    startDate: '',
    endDate: '',
    preset: null
  });
  
  const { data: summaryStats, isLoading: statsLoading } = useSummaryStats(period);
  const { data: expensesData, isLoading: expensesLoading } = useExpenses(undefined, 1, 5);
  const { data: budgetRemainingData } = useBudgetRemaining(selectedDate);

  // fetch recent income occurrences (last 7 days) including virtual recurring entries
  const recentOccurrenceDates = useMemo(() => {
    const end = new Date();
    const start = new Date();
    start.setDate(start.getDate() - 6);
    return {
      startDate: start.toISOString().split('T')[0],
      endDate: end.toISOString().split('T')[0]
    };
  }, []);
  const { data: recentOccurrences, isLoading: occurrencesLoading } = useIncomeOccurrences(
    recentOccurrenceDates.startDate,
    recentOccurrenceDates.endDate,
    1,
    10
  );
  const { isGuest } = useAuth();
  const { showToast } = useToast();

  const [viewingExpense, setViewingExpense] = useState<Expense | null>(null);
  const [showIncomeForm, setShowIncomeForm] = useState(false);
  const [showExpenseDialog, setShowExpenseDialog] = useState(false);
  const [editingIncome, setEditingIncome] = useState<Income | null>(null);

  // check for query param to open expense dialog (deep linking support)
  useEffect(() => {
    const openExpense = searchParams.get("openExpense");
    if (openExpense === "true") {
      setShowExpenseDialog(true);
      // remove query param from URL
      router.replace("/dashboard/budget", { scroll: false });
    }
  }, [searchParams, router]);

  const updateExpense = useUpdateExpense();
  const deleteExpense = useDeleteExpense();
  const createIncome = useCreateIncome();
  const createExpense = useCreateExpense();
  const updateIncome = useUpdateIncome();
  const deleteIncome = useDeleteIncome();

  const displayExpenses = expensesData?.data || [];

  const handleActionClick = () => {
    if (isGuest) {
      showToast("Guest user is read-only. Create an account to save changes", "warning");
    }
  };

  const handleEditExpense = async (data: any) => {
    try {
      if (!viewingExpense) return;
      await updateExpense.mutateAsync({ id: viewingExpense.id, data });
      showToast("Expense updated successfully", "success");
      setViewingExpense(null);
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, "warning");
      } else {
        showToast("Failed to update expense", "error");
      }
    }
  };

  const handleDeleteExpense = async (id: string) => {
    try {
      await deleteExpense.mutateAsync(id);
      showToast("Expense deleted successfully", "success");
      setViewingExpense(null);
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, "warning");
      } else {
        showToast("Failed to delete expense", "error");
      }
    }
  };

  const handleCreateIncome = async (data: CreateIncomeRequest | UpdateIncomeRequest) => {
    try {
      await createIncome.mutateAsync(data as CreateIncomeRequest);
      showToast("Income created successfully", "success");
      setShowIncomeForm(false);
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, "warning");
      } else {
        showToast("Failed to create income", "error");
      }
    }
  };

  const handleUpdateIncome = async (data: UpdateIncomeRequest) => {
    try {
      if (!editingIncome) return;
      await updateIncome.mutateAsync({ id: editingIncome.id, data });
      showToast("Income updated successfully", "success");
      setEditingIncome(null);
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, "warning");
      } else {
        showToast("Failed to update income", "error");
      }
    }
  };

  const handleDeleteIncome = async (income: Income) => {
    try {
      await deleteIncome.mutateAsync(income.id);
      showToast("Income deleted successfully", "success");
      setEditingIncome(null);
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, "warning");
      } else {
        showToast("Failed to delete income", "error");
      }
    }
  };

  const handleCreateExpense = async (data: CreateExpenseRequest | UpdateExpenseRequest) => {
    try {
      await createExpense.mutateAsync(data as CreateExpenseRequest);
      showToast("Expense created successfully", "success");
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, "warning");
      } else {
        showToast("Failed to create expense", "error");
      }
    }
  };

  return (
    <div className="px-4 md:px-6 py-6">
      {/* Header */}
      <motion.div
        className="mb-8 flex items-center justify-between"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <div>
          <div className="flex items-center gap-2 mb-2">
            <h1 className="text-primary">BUDGET TRACKER</h1>
            {isGuest && (
              <Badge variant="outline" className="bg-accent/10 text-accent border-accent/30">
                Demo Mode
              </Badge>
            )}
          </div>
          <p className="text-sm text-muted-foreground">
            {isGuest
              ? "Viewing sample data. Register to create your own budget."
              : "Track and manage your expenses and income"}
          </p>
        </div>
        <div className="flex gap-2">
          <Button onClick={() => { handleActionClick(); if (!isGuest) setShowIncomeForm(true); }}>
            <Plus className="mr-2 h-4 w-4" />
            Add Income
          </Button>
          <Button onClick={() => { handleActionClick(); if (!isGuest) setShowExpenseDialog(true); }}>
            <Plus className="mr-2 h-4 w-4" />
            Add Expense
          </Button>
        </div>
      </motion.div>

      {/* Date Filter & Budget Remaining */}
      <motion.div
        className="mb-4 flex items-center justify-between"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.1 }}
      >
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <Label htmlFor="date-filter" className="text-sm font-medium">Budget Date:</Label>
            <Input
              id="date-filter"
              type="date"
              value={selectedDate}
              onChange={(e) => setSelectedDate(e.target.value)}
              className="w-auto"
            />
          </div>
          {budgetRemainingData && (
            <div className={`text-sm font-semibold ${
              budgetRemainingData.budget_remaining_status === 'red' ? 'text-red-500' :
              budgetRemainingData.budget_remaining_status === 'green' ? 'text-green-600' :
              'text-gray-600'
            }`}>
              Budget Remaining: ₱{budgetRemainingData.budget_remaining.toFixed(2)}
              {budgetRemainingData.budget_remaining_status === 'red' && ' ⚠️ Over Budget'}
              {budgetRemainingData.budget_remaining_status === 'green' && ' ✓ On Track'}
            </div>
          )}
        </div>
      </motion.div>

      {/* Summary Statistics */}
      <motion.div
        className="mb-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.15 }}
      >
        <ExpenseStats stats={summaryStats} isLoading={statsLoading} />
      </motion.div>

      {/* Date Range Picker for Analytics */}
      <motion.div
        className="mb-4"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.18 }}
      >
        <PeriodPresetFilter value={dateRange} onChange={setDateRange} />
      </motion.div>

      {/* Budget Analytics Cards */}
      <motion.div
        className="mb-8 grid gap-4 md:grid-cols-2 lg:grid-cols-4"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.2 }}
      >
        <CurrentTotalMoneyCard startDate={dateRange.startDate || undefined} endDate={dateRange.endDate || undefined} />
        <SavingsRateCard startDate={dateRange.startDate || undefined} endDate={dateRange.endDate || undefined} />
        <SpendingVelocityCard startDate={dateRange.startDate || undefined} endDate={dateRange.endDate || undefined} />
        <RecurringSummaryCard />
      </motion.div>

      {/* Budget Variance Table */}
      <motion.div
        className="mb-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.25 }}
      >
        <BudgetVarianceTable />
      </motion.div>

      {/* Recurring Income and Bills */}
      <motion.div
        className="mb-8 grid gap-6 lg:grid-cols-2"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.3 }}
      >
        <RecurringIncomeList />
        <UpcomingBillsCard days={7} />
      </motion.div>

      {/* Recent Transactions */}
      <motion.div
        className="grid gap-6 lg:grid-cols-3"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.25 }}
      >
        {/* Recent Incomes */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Recent Incomes</CardTitle>
                <CardDescription>Last 7 days (including recurring)</CardDescription>
              </div>
              <Link href="/dashboard/budget/incomes">
                <Button variant="ghost" size="sm">
                  View All
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </Link>
            </div>
          </CardHeader>
          <CardContent>
            {occurrencesLoading ? (
              <div className="space-y-4">
                {[...Array(3)].map((_, i) => (
                  <div key={i} className="flex items-center space-x-4 animate-pulse">
                    <div className="h-10 w-10 bg-muted rounded-full" />
                    <div className="flex-1 space-y-2">
                      <div className="h-4 w-3/4 bg-muted rounded" />
                      <div className="h-3 w-1/2 bg-muted rounded" />
                    </div>
                    <div className="h-6 w-20 bg-muted rounded" />
                  </div>
                ))}
              </div>
            ) : !recentOccurrences?.data?.length ? (
              <div className="text-center py-8 text-muted-foreground">
                <p>No incomes in the last 7 days</p>
                <Button
                  variant="outline"
                  size="sm"
                  className="mt-4"
                  onClick={() => { handleActionClick(); if (!isGuest) setShowIncomeForm(true); }}
                >
                  <Plus className="mr-2 h-4 w-4" />
                  Add your first income
                </Button>
              </div>
            ) : (
              <div className="space-y-4">
                {recentOccurrences.data.map((occ) => (
                  <div
                    key={occ.id}
                    className={`flex items-center justify-between p-3 rounded-lg border border-border/50 hover:border-border transition-colors ${occ.is_skipped ? 'bg-red-500/5' : ''} ${occ.is_virtual ? 'cursor-default' : 'cursor-pointer hover:bg-accent/50'}`}
                    onClick={() => {
                      if (!occ.is_virtual) {
                        // for real entries, allow editing
                        const incomeForEdit: Income = {
                          id: occ.source_income_id,
                          amount: occ.amount,
                          currency: occ.currency,
                          date: occ.date,
                          description: occ.description,
                          recurring_type: occ.recurring_type,
                          created_at: '',
                          updated_at: ''
                        };
                        setEditingIncome(incomeForEdit);
                      }
                    }}
                  >
                    <div className="flex-1 min-w-0">
                      <p className={`font-medium truncate ${occ.is_skipped ? 'line-through text-muted-foreground' : ''}`}>
                        {occ.description || (occ.is_skipped ? "Skipped Income" : "Income")}
                        {occ.is_virtual && (
                          <Badge variant="secondary" className="ml-2 text-xs">
                            {occ.recurring_type || 'recurring'}
                          </Badge>
                        )}
                        {occ.is_skipped && (
                          <Badge variant="destructive" className="ml-2 text-xs">
                            Skipped
                          </Badge>
                        )}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        {formatDistanceToNow(new Date(occ.date), { addSuffix: true })}
                      </p>
                    </div>
                    <div className="text-right ml-4">
                      <p className={`font-semibold shrink-0 ${occ.is_skipped ? 'text-red-600' : 'text-green-600'}`}>
                        {occ.is_skipped ? '' : '+'}{occ.currency} {Math.abs(occ.amount).toFixed(2)}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Recent Expenses */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Recent Expenses</CardTitle>
                <CardDescription>Your latest transactions</CardDescription>
              </div>
              <Link href="/dashboard/budget/expenses">
                <Button variant="ghost" size="sm">
                  View All
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </Link>
            </div>
          </CardHeader>
          <CardContent>
            {expensesLoading ? (
              <div className="space-y-4">
                {[...Array(3)].map((_, i) => (
                  <div key={i} className="flex items-center space-x-4 animate-pulse">
                    <div className="h-10 w-10 bg-muted rounded-full" />
                    <div className="flex-1 space-y-2">
                      <div className="h-4 w-3/4 bg-muted rounded" />
                      <div className="h-3 w-1/2 bg-muted rounded" />
                    </div>
                    <div className="h-6 w-20 bg-muted rounded" />
                  </div>
                ))}
              </div>
            ) : displayExpenses.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                {isGuest ? (
                  <div className="space-y-2">
                    <p>No sample expenses available.</p>
                    <p className="text-xs">Register to create your own budget and track expenses.</p>
                  </div>
                ) : (
                  <>
                    <p>No expenses yet</p>
                    <Button 
                      variant="outline" 
                      size="sm" 
                      className="mt-4"
                      onClick={() => { handleActionClick(); if (!isGuest) setShowExpenseDialog(true); }}
                    >
                      <Plus className="mr-2 h-4 w-4" />
                      Add your first expense
                    </Button>
                  </>
                )}
              </div>
            ) : (
              <div className="space-y-4">
                {displayExpenses.map((expense) => (
                  <div
                    key={expense.id}
                    className="flex items-center justify-between p-3 rounded-lg border border-border/50 hover:border-border transition-colors cursor-pointer hover:bg-accent/50"
                    onClick={() => setViewingExpense(expense)}
                  >
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 flex-wrap">
                        <p className="font-medium truncate">{expense.description}</p>
                        {expense.category && (
                          <span
                            className="text-xs px-2 py-0.5 rounded-full truncate max-w-[100px]"
                            style={{
                              backgroundColor: expense.category.color
                                ? `${expense.category.color}20`
                                : undefined,
                              color: expense.category.color || "inherit",
                            }}
                            title={expense.category.name}
                          >
                            {expense.category.name}
                          </span>
                        )}
                        {expense.tags && expense.tags.length > 0 && (
                          <div className="flex gap-1">
                            {expense.tags.slice(0, 2).map((tag) => (
                              <Badge
                                key={tag.id}
                                variant="secondary"
                                className="text-xs py-0 h-5"
                                style={{
                                  backgroundColor: tag.color ? `${tag.color}15` : undefined,
                                  color: tag.color || "inherit",
                                }}
                              >
                                {tag.name}
                              </Badge>
                            ))}
                            {expense.tags.length > 2 && (
                              <Badge variant="outline" className="text-xs py-0 h-5">
                                +{expense.tags.length - 2}
                              </Badge>
                            )}
                          </div>
                        )}
                      </div>
                      <p className="text-xs text-muted-foreground">
                        {formatDistanceToNow(new Date(expense.expense_date), { addSuffix: true })}
                      </p>
                    </div>
                    <div className="text-right ml-4">
                      <p className="font-semibold shrink-0 text-red-600">
                        -{expense.currency} {expense.amount.toFixed(2)}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </motion.div>

      {/* Income Form Dialog */}
      <Dialog open={showIncomeForm} onOpenChange={setShowIncomeForm}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>Add New Income</DialogTitle>
          </DialogHeader>
          <IncomeForm
            onSubmit={handleCreateIncome}
            onCancel={() => setShowIncomeForm(false)}
            isLoading={createIncome.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Income Dialog */}
      <Dialog open={!!editingIncome} onOpenChange={() => setEditingIncome(null)}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>Edit Income</DialogTitle>
          </DialogHeader>
          {editingIncome && (
            <IncomeForm
              initialData={{
                id: editingIncome.id,
                amount: editingIncome.amount,
                currency: editingIncome.currency,
                date: editingIncome.date,
                description: editingIncome.description,
                recurring_type: editingIncome.recurring_type,
                start_date: editingIncome.start_date,
                end_date: editingIncome.end_date,
              }}
              onSubmit={handleUpdateIncome}
              onCancel={() => setEditingIncome(null)}
              isLoading={updateIncome.isPending}
            />
          )}
          {editingIncome && (
            <div className="mt-4 pt-4 border-t">
              <Button
                variant="destructive"
                className="w-full"
                onClick={() => handleDeleteIncome(editingIncome)}
                disabled={deleteIncome.isPending}
              >
                Delete Income
              </Button>
            </div>
          )}
        </DialogContent>
      </Dialog>

      <ExpenseDetailDialog
        expense={viewingExpense}
        open={!!viewingExpense}
        onOpenChange={(open) => {
          if (!open) {
            setViewingExpense(null);
          }
        }}
        onEdit={handleEditExpense}
        onDelete={handleDeleteExpense}
        showToast={showToast}
        isGuest={isGuest}
      />

      {/* Expense Form Dialog */}
      <ExpenseFormDialog
        open={showExpenseDialog}
        onOpenChange={setShowExpenseDialog}
        onSubmit={handleCreateExpense}
      />
    </div>
  );
}

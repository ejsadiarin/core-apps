"use client";

import { motion } from "motion/react";
import { useState } from "react";
import { ExpenseStats } from "@/components/budget/expense-stats";
import { IncomeForm } from "@/components/budget/income-form";
import { ExpenseFormDialog } from "@/components/budget/expense-form-dialog";
import {
  useSummaryStats,
  useExpenses,
  useIncomes,
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
import { Plus, ArrowRight, Receipt, Settings, Wallet } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import Link from "next/link";
import { formatDistanceToNow } from "date-fns";
import { useAuth } from "@/contexts/auth-context";
import { useToast } from "@/components/ui/toast";
import { ExpenseDetailDialog } from "@/components/budget/expense-detail-dialog";
import { formatAmount } from "@/lib/utils";
import type { Expense, Income, CreateIncomeRequest, UpdateIncomeRequest } from "@/types/api";

export default function BudgetDashboard() {
  const [period] = useState<string>("month");
  const [selectedDate, setSelectedDate] = useState<string>(new Date().toISOString().split('T')[0]);
  const { data: summaryStats, isLoading: statsLoading } = useSummaryStats(period);
  const { data: recentExpenses, isLoading: expensesLoading } = useExpenses();
  const { data: recentIncomes, isLoading: incomesLoading } = useIncomes();
  const { data: budgetRemainingData } = useBudgetRemaining(selectedDate);
  const { isGuest } = useAuth();
  const { showToast } = useToast();

  const [viewingExpense, setViewingExpense] = useState<Expense | null>(null);
  const [showIncomeForm, setShowIncomeForm] = useState(false);
  const [showExpenseDialog, setShowExpenseDialog] = useState(false);
  const [editingIncome, setEditingIncome] = useState<Income | null>(null);

  const updateExpense = useUpdateExpense();
  const deleteExpense = useDeleteExpense();
  const createIncome = useCreateIncome();
  const createExpense = useCreateExpense();
  const updateIncome = useUpdateIncome();
  const deleteIncome = useDeleteIncome();

  const displayExpenses = recentExpenses?.slice(0, 5) || [];
  const displayIncomes = recentIncomes?.slice(0, 5) || [];

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

  const handleCreateIncome = async (data: CreateIncomeRequest) => {
    try {
      await createIncome.mutateAsync(data);
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

  const handleCreateExpense = async (data: CreateExpenseRequest) => {
    try {
      await createExpense.mutateAsync(data);
      showToast("Expense created successfully", "success");
      setShowExpenseDialog(false);
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, "warning");
      } else {
        showToast("Failed to create expense", "error");
      }
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
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
              budgetRemainingData.budget_remaining_status === 'green' ? 'text-green-900' :
              'text-gray-600'
            }`}>
              Budget Remaining: {formatAmount(budgetRemainingData.budget_remaining, budgetRemainingData.currency || 'PHP')}
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
                <CardDescription>Your latest income entries</CardDescription>
              </div>
              <Wallet className="h-5 w-5 text-muted-foreground" />
            </div>
          </CardHeader>
          <CardContent>
            {incomesLoading ? (
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
            ) : displayIncomes.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                <p>No incomes yet</p>
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
                {displayIncomes.map((income) => (
                  <div
                    key={income.id}
                    className="flex items-center justify-between p-3 rounded-lg border border-border/50 hover:border-border transition-colors cursor-pointer hover:bg-accent/50"
                    onClick={() => setEditingIncome(income)}
                  >
                    <div className="flex-1 min-w-0">
                      <p className="font-medium truncate">
                        {income.description || "Income"}
                        {income.recurring_type && (
                          <Badge variant="secondary" className="ml-2 text-xs">
                            {income.recurring_type}
                          </Badge>
                        )}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        {formatDistanceToNow(new Date(income.date), { addSuffix: true })}
                      </p>
                    </div>
                    <div className="text-right ml-4">
                      <p className="font-semibold shrink-0 text-green-600">
                        +{income.currency} {income.amount.toFixed(2)}
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

      {/* Quick Actions */}
      <motion.div
        className="mt-6 grid gap-6 lg:grid-cols-3"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.3 }}
      >
        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
            <CardDescription>Manage your budget</CardDescription>
          </CardHeader>
          <CardContent className="space-y-2">
            <Button
              variant="outline"
              className="w-full justify-start"
              onClick={() => { handleActionClick(); if (!isGuest) setShowIncomeForm(true); }}
            >
              <Plus className="mr-2 h-4 w-4" />
              Add New Income
            </Button>
            <Button
              variant="outline"
              className="w-full justify-start"
              onClick={() => { handleActionClick(); if (!isGuest) setShowExpenseDialog(true); }}
            >
              <Plus className="mr-2 h-4 w-4" />
              Add New Expense
            </Button>
            <Link href="/dashboard/budget/expenses" className="block" onClick={handleActionClick}>
              <Button variant="outline" className="w-full justify-start">
                <Receipt className="mr-2 h-4 w-4" />
                View All Expenses
              </Button>
            </Link>
            <Link href="/dashboard/budget/settings" className="block" onClick={handleActionClick}>
              <Button variant="outline" className="w-full justify-start">
                <Settings className="mr-2 h-4 w-4" />
                Manage Categories & Tags
              </Button>
            </Link>
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
        isLoading={createExpense.isPending}
      />
    </div>
  );
}

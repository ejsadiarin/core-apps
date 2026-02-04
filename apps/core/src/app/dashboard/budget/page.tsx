"use client";

import { motion } from "motion/react";
import { useState } from "react";
import { ExpenseStats } from "@/components/budget/expense-stats";
import { useSummaryStats, useExpenses, useUpdateExpense, useDeleteExpense } from "@/hooks/use-budget";
import { GuestBlockedError } from "@/hooks/use-budget";
import { Button } from "@/components/ui/button";
import { Plus, ArrowRight, Receipt, Settings } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import Link from "next/link";
import { formatDistanceToNow } from "date-fns";
import { useAuth } from "@/contexts/auth-context";
import { useToast } from "@/components/ui/toast";
import { ExpenseDetailDialog } from "@/components/budget/expense-detail-dialog";
import type { Expense } from "@/types/api";

export default function BudgetDashboard() {
  const [period] = useState<string>("month");
  const { data: summaryStats, isLoading: statsLoading } = useSummaryStats(period);
  const { data: recentExpenses, isLoading: expensesLoading } = useExpenses();
  const { isGuest } = useAuth();
  const { showToast } = useToast();
  const [viewingExpense, setViewingExpense] = useState<Expense | null>(null);
  const updateExpense = useUpdateExpense();
  const deleteExpense = useDeleteExpense();

  const displayExpenses = recentExpenses?.slice(0, 5) || [];

  const handleActionClick = () => {
    if (isGuest) {
      showToast("Guest user is read-only. Create an account to save changes", "warning");
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteExpense.mutateAsync(id);
      showToast("Expense deleted successfully", "success");
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, "warning");
      } else {
        showToast("Failed to delete expense", "error");
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
              : "Track and manage your expenses"}
          </p>
        </div>
        <Link href="/dashboard/budget/expenses/new" onClick={handleActionClick}>
          <Button>
            <Plus className="mr-2 h-4 w-4" />
            Add Expense
          </Button>
        </Link>
      </motion.div>

      {/* Summary Statistics */}
      <motion.div
        className="mb-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.1 }}
      >
        <ExpenseStats stats={summaryStats} isLoading={statsLoading} />
      </motion.div>

      {/* Recent Expenses */}
      <motion.div
        className="grid gap-6 lg:grid-cols-2"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.2 }}
      >
        <Card>
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
                    <Link href="/dashboard/budget/expenses/new">
                      <Button variant="outline" size="sm" className="mt-4">
                        <Plus className="mr-2 h-4 w-4" />
                        Add your first expense
                      </Button>
                    </Link>
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
                        {formatDistanceToNow(new Date(expense.expense_date), {
                          addSuffix: true,
                        })}
                      </p>
                    </div>
                    <div className="text-right ml-4">
                      <p className="font-semibold shrink-0">
                        {expense.currency} {expense.amount.toFixed(2)}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Quick Actions */}
        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
            <CardDescription>Manage your budget</CardDescription>
          </CardHeader>
          <CardContent className="space-y-2">
            <Link href="/dashboard/budget/expenses/new" className="block" onClick={handleActionClick}>
              <Button variant="outline" className="w-full justify-start">
                <Plus className="mr-2 h-4 w-4" />
                Add New Expense
              </Button>
            </Link>
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

      <ExpenseDetailDialog
        expense={viewingExpense}
        open={!!viewingExpense}
        onOpenChange={(open) => !open && setViewingExpense(null)}
        onDelete={handleDelete}
        showToast={showToast}
        isGuest={isGuest}
      />
    </div>
  );
}

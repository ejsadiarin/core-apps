"use client";

import { motion } from "motion/react";
import { useState } from "react";
import { useExpenses, useCategories, useDeleteExpense, useUpdateExpense, GuestBlockedError } from "@/hooks/use-budget";
import { ExpenseCard } from "@/components/budget/expense-card";
import { EditExpenseDialog } from "@/components/budget/expense-edit-dialog";
import { ExpenseDetailDialog } from "@/components/budget/expense-detail-dialog";
import { Button } from "@/components/ui/button";
import { Plus, Filter, ArrowLeft, EyeOff } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import type { ExpenseFilters, Expense } from "@/types/api";
import { useAuth } from "@/contexts/auth-context";
import { useToast } from "@/components/ui/toast";
import Link from "next/link";

export default function ExpensesPage() {
  const [filters, setFilters] = useState<ExpenseFilters>({});
  const [searchTerm, setSearchTerm] = useState("");
  const [editingExpense, setEditingExpense] = useState<Expense | null>(null);
  const [viewingExpense, setViewingExpense] = useState<Expense | null>(null);

  const { data: expenses, isLoading } = useExpenses(filters);
  const { data: categories } = useCategories();
  const deleteExpense = useDeleteExpense();
  const updateExpense = useUpdateExpense();
  const { isGuest } = useAuth();
  const { showToast } = useToast();

  const filteredExpenses = expenses?.filter((expense) =>
    expense.description.toLowerCase().includes(searchTerm.toLowerCase())
  );

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

  const handleUpdate = async (data: any) => {
    try {
      await updateExpense.mutateAsync({ id: editingExpense!.id, data });
      showToast("Expense updated successfully", "success");
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, "warning");
      } else {
        showToast("Failed to update expense", "error");
      }
    }
  };

  const handleView = (expense: Expense) => {
    setViewingExpense(expense);
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
          <Link href="/dashboard/budget">
            <Button variant="ghost" size="sm" className="mb-2 pl-0">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Budget
            </Button>
          </Link>
          <div className="flex items-center gap-2 mb-2">
            <h1 className="text-primary">EXPENSE LIST</h1>
            {isGuest && (
              <Badge variant="outline" className="bg-accent/10 text-accent border-accent/30">
                Demo Mode
              </Badge>
            )}
          </div>
          <p className="text-sm text-muted-foreground">
            {isGuest
              ? "Viewing sample expenses"
              : "View and manage all your expenses"}
          </p>
        </div>
        <Link href="/dashboard/budget/expenses/new">
          <Button>
            <Plus className="mr-2 h-4 w-4" />
            Add Expense
          </Button>
        </Link>
      </motion.div>

      {/* Filters */}
      <motion.div
        className="mb-6 flex flex-col sm:flex-row gap-4"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.1 }}
      >
        <Input
          placeholder="Search expenses..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="sm:max-w-xs"
        />
        
        <Select
          value={filters.category_id || "all"}
          onValueChange={(value) =>
            setFilters((prev) => ({
              ...prev,
              category_id: value === "all" ? undefined : value,
            }))
          }
        >
          <SelectTrigger className="sm:max-w-xs">
            <Filter className="mr-2 h-4 w-4" />
            <SelectValue placeholder="All Categories" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Categories</SelectItem>
            {categories?.map((category) => (
              <SelectItem key={category.id} value={category.id}>
                {category.icon && <span className="mr-2">{category.icon}</span>}
                {category.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <div className="flex gap-2">
          <Input
            type="date"
            value={filters.start_date || ""}
            onChange={(e) =>
              setFilters((prev) => ({ ...prev, start_date: e.target.value }))
            }
            className="sm:max-w-[150px]"
          />
          <Input
            type="date"
            value={filters.end_date || ""}
            onChange={(e) =>
              setFilters((prev) => ({ ...prev, end_date: e.target.value }))
            }
            className="sm:max-w-[150px]"
          />
        </div>
      </motion.div>

      {/* Expense List */}
      <motion.div
        className="space-y-4"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.2 }}
      >
        {isLoading ? (
          <div className="space-y-4">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="h-32 bg-muted animate-pulse rounded-lg" />
            ))}
          </div>
        ) : filteredExpenses && filteredExpenses.length > 0 ? (
          filteredExpenses.map((expense) => (
            <ExpenseCard
              key={expense.id}
              expense={expense}
              onView={handleView}
              onEdit={(exp) => setEditingExpense(exp)}
              onDelete={handleDelete}
              showToast={showToast}
            />
          ))
        ) : (
          <div className="text-center py-12 text-muted-foreground">
            <div className="space-y-2">
              <EyeOff className="w-8 h-8 mx-auto opacity-50" />
              <p>No expenses found.</p>
              {isGuest ? (
                <p className="text-xs">Sign in or create an account to manage your own expenses.</p>
              ) : (
                <Link href="/dashboard/budget/expenses/new">
                  <Button variant="outline" size="sm" className="mt-2">
                    <Plus className="mr-2 h-4 w-4" />
                    Add your first expense
                  </Button>
                </Link>
              )}
            </div>
          </div>
        )}
      </motion.div>

      {/* Edit Dialog */}
      <EditExpenseDialog
        expense={editingExpense}
        open={!!editingExpense}
        onOpenChange={(open) => !open && setEditingExpense(null)}
        onSubmit={handleUpdate}
        isLoading={updateExpense.isPending}
      />

      {/* Detail Dialog */}
      <ExpenseDetailDialog
        expense={viewingExpense}
        open={!!viewingExpense}
        onOpenChange={(open) => !open && setViewingExpense(null)}
        onEdit={(exp) => {
          setViewingExpense(null);
          setEditingExpense(exp);
        }}
        onDelete={handleDelete}
        showToast={showToast}
      />
    </div>
  );
}

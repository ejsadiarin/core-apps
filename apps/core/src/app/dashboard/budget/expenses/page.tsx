"use client";

import { motion } from "motion/react";
import { useState } from "react";
import { useExpenses, useCategories, useDeleteExpense } from "@/hooks/use-budget";
import { ExpenseCard } from "@/components/budget/expense-card";
import { Button } from "@/components/ui/button";
import { Plus, Filter } from "lucide-react";
import Link from "next/link";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import type { ExpenseFilters } from "@/types/api";

export default function ExpensesPage() {
  const [filters, setFilters] = useState<ExpenseFilters>({});
  const [searchTerm, setSearchTerm] = useState("");
  
  const { data: expenses, isLoading } = useExpenses(filters);
  const { data: categories } = useCategories();
  const deleteExpense = useDeleteExpense();

  // filter by search term
  const filteredExpenses = expenses?.filter((expense) =>
    expense.description.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleDelete = async (id: string) => {
    if (confirm("Are you sure you want to delete this expense?")) {
      await deleteExpense.mutateAsync(id);
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
          <h1 className="text-primary mb-2">EXPENSE LIST</h1>
          <p className="text-sm text-muted-foreground">
            View and manage all your expenses
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
              onDelete={handleDelete}
            />
          ))
        ) : (
          <div className="text-center py-12 text-muted-foreground">
            <p className="mb-4">No expenses found</p>
            <Link href="/dashboard/budget/expenses/new">
              <Button variant="outline">
                <Plus className="mr-2 h-4 w-4" />
                Add your first expense
              </Button>
            </Link>
          </div>
        )}
      </motion.div>
    </div>
  );
}

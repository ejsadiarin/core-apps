"use client";

import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Trash2, Edit } from "lucide-react";
import type { Expense } from "@/types/api";
import { format } from "date-fns";

interface ExpenseCardProps {
  expense: Expense;
  onEdit?: (expense: Expense) => void;
  onDelete?: (id: string) => void;
  disabled?: boolean;
  showToast?: (message: string, type?: "info" | "warning" | "error" | "success") => void;
}

export function ExpenseCard({ expense, onEdit, onDelete, disabled, showToast }: ExpenseCardProps) {
  const handleEditClick = () => {
    if (disabled && showToast) {
      showToast("Guest user is read-only. Create an account to save changes", "warning");
      return;
    }
    onEdit?.(expense);
  };

  const handleDeleteClick = () => {
    if (disabled && showToast) {
      showToast("Guest user is read-only. Create an account to save changes", "warning");
      return;
    }
    if (onDelete) {
      if (confirm("Are you sure you want to delete this expense?")) {
        onDelete(expense.id);
      }
    }
  };

  return (
    <Card className="hover:border-primary/50 transition-colors opacity-90">
      <CardContent className="p-4">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <h3 className="font-semibold">{expense.description}</h3>
              {expense.category && (
                <Badge
                  variant="outline"
                  style={{
                    backgroundColor: expense.category.color
                      ? `${expense.category.color}20`
                      : undefined,
                    borderColor: expense.category.color || undefined,
                    color: expense.category.color || "inherit",
                  }}
                >
                  {expense.category.icon && <span className="mr-1">{expense.category.icon}</span>}
                  {expense.category.name}
                </Badge>
              )}
            </div>
            
            {expense.notes && (
              <p className="text-sm text-muted-foreground mb-2">{expense.notes}</p>
            )}
            
            <div className="flex items-center gap-2 flex-wrap">
              {expense.tags && expense.tags.length > 0 && (
                <div className="flex gap-1 flex-wrap">
                  {expense.tags.map((tag) => (
                    <Badge
                      key={tag.id}
                      variant="secondary"
                      className="text-xs"
                      style={{
                        backgroundColor: tag.color ? `${tag.color}15` : undefined,
                        color: tag.color || "inherit",
                      }}
                    >
                      {tag.name}
                    </Badge>
                  ))}
                </div>
              )}
            </div>
            
            <p className="text-xs text-muted-foreground mt-2">
              {format(new Date(expense.expense_date), "MMM dd, yyyy")}
            </p>
          </div>
          
          <div className="flex flex-col items-end gap-2">
            <div className="text-xl font-bold">
              {expense.currency} {expense.amount.toFixed(2)}
            </div>
            <div className="flex gap-1">
              <Button
                variant="ghost"
                size="sm"
                onClick={handleEditClick}
                disabled={disabled}
                className={disabled ? "opacity-50" : ""}
              >
                <Edit className="h-4 w-4" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleDeleteClick}
                disabled={disabled}
                className={disabled ? "opacity-50 text-muted-foreground" : "text-destructive hover:text-destructive"}
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

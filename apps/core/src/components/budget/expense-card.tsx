"use client";

import { useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Trash2, Edit } from "lucide-react";
import type { Expense } from "@/types/api";
import { format } from "date-fns";

interface ExpenseCardProps {
  expense: Expense;
  onView?: (expense: Expense) => void;
  onEdit?: (expense: Expense) => void;
  onDelete?: (id: string) => void;
  disabled?: boolean;
  showToast?: (message: string, type?: "info" | "warning" | "error" | "success") => void;
}

export function ExpenseCard({ expense, onView, onEdit, onDelete, disabled, showToast }: ExpenseCardProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  const handleCardClick = () => {
    onView?.(expense);
  };

  const handleEditClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (disabled && showToast) {
      showToast("Guest user is read-only. Create an account to save changes", "warning");
      return;
    }
    onEdit?.(expense);
  };

  const handleDeleteClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (disabled && showToast) {
      showToast("Guest user is read-only. Create an account to save changes", "warning");
      return;
    }
    setDeleteDialogOpen(true);
  };

  const handleConfirmDelete = () => {
    onDelete?.(expense.id);
    setDeleteDialogOpen(false);
  };

  return (
    <>
      <Card
        className="hover:border-primary/50 transition-colors cursor-pointer"
        onClick={handleCardClick}
      >
        <CardContent className="p-4">
          <div className="flex items-start justify-between">
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-2 flex-wrap">
                <h3 className="font-semibold">{expense.description}</h3>
                {expense.category && (
                  <Badge
                    variant="outline"
                    className="text-xs"
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
                <p className="text-sm text-muted-foreground mb-2 line-clamp-2">{expense.notes}</p>
              )}

              <div className="flex items-center gap-2 flex-wrap">
                {expense.tags && expense.tags.length > 0 && (
                  <div className="flex gap-1 flex-wrap">
                    {expense.tags.slice(0, 3).map((tag) => (
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
                    {expense.tags.length > 3 && (
                      <Badge variant="outline" className="text-xs">
                        +{expense.tags.length - 3}
                      </Badge>
                    )}
                  </div>
                )}
              </div>

              <p className="text-xs text-muted-foreground mt-2">
                {format(new Date(expense.expense_date), "MMM dd, yyyy")}
              </p>
            </div>

            <div className="flex flex-col items-end gap-2">
              <div className="text-xl font-bold text-red-500">
                -{expense.currency} {expense.amount.toFixed(2)}
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

      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Expense</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete this expense? This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleConfirmDelete} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}

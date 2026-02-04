"use client";

import { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@/components/ui/dialog";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Edit, Trash2, Calendar, Tag as TagIcon } from "lucide-react";
import type { Expense } from "@/types/api";
import { format } from "date-fns";
import { EditExpenseDialog } from "./expense-edit-dialog";
import { GuestBlockedError } from "@/hooks/use-budget";

interface ExpenseDetailDialogProps {
  expense: Expense | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onEdit?: (expense: Expense) => void;
  onDelete?: (id: string) => void;
  showToast?: (message: string, type?: "info" | "warning" | "error" | "success") => void;
}

export function ExpenseDetailDialog({
  expense,
  open,
  onOpenChange,
  onEdit,
  onDelete,
  showToast,
}: ExpenseDetailDialogProps) {
  const [editDialogOpen, setEditDialogOpen] = useState(false);

  if (!expense) return null;

  const handleEditClick = () => {
    onEdit?.(expense);
    setEditDialogOpen(true);
  };

  const handleDeleteClick = () => {
    if (showToast) {
      showToast("Guest user is read-only. Create an account to save changes", "warning");
      return;
    }
    if (onDelete) {
      if (confirm("Are you sure you want to delete this expense?")) {
        onDelete(expense.id);
        onOpenChange(false);
      }
    }
  };

  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="sm:max-w-[450px]">
          <DialogHeader>
            <DialogTitle>Expense Details</DialogTitle>
            <DialogDescription>
              View the full details of this expense.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 mt-4">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-lg font-semibold">{expense.description}</h3>
                <p className="text-2xl font-bold text-primary">
                  {expense.currency} {expense.amount.toFixed(2)}
                </p>
              </div>
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
              <Card>
                <CardContent className="p-3">
                  <p className="text-sm text-muted-foreground">{expense.notes}</p>
                </CardContent>
              </Card>
            )}

            <div className="flex flex-wrap gap-2">
              {expense.tags && expense.tags.length > 0 && (
                <div className="flex flex-wrap gap-1">
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
                      <TagIcon className="mr-1 h-3 w-3" />
                      {tag.name}
                    </Badge>
                  ))}
                </div>
              )}
            </div>

            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Calendar className="h-4 w-4" />
              <span>{format(new Date(expense.expense_date), "MMMM dd, yyyy")}</span>
            </div>

            <div className="flex gap-2 pt-4">
              <Button
                variant="outline"
                size="sm"
                onClick={handleEditClick}
                className="flex-1"
              >
                <Edit className="mr-2 h-4 w-4" />
                Edit
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={handleDeleteClick}
                className="flex-1 text-destructive hover:text-destructive"
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Delete
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>

      <EditExpenseDialog
        expense={expense}
        open={editDialogOpen}
        onOpenChange={setEditDialogOpen}
        onSubmit={async (data) => {
          await onEdit?.(expense);
        }}
        isLoading={false}
      />
    </>
  );
}

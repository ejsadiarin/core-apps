"use client";

import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@/components/ui/dialog";
import { ExpenseForm } from "./expense-form";
import type { Expense, CreateExpenseRequest, UpdateExpenseRequest } from "@/types/api";

interface EditExpenseDialogProps {
  expense: Expense | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (data: CreateExpenseRequest | UpdateExpenseRequest) => Promise<void>;
  isLoading?: boolean;
}

export function EditExpenseDialog({ expense, open, onOpenChange, onSubmit, isLoading }: EditExpenseDialogProps) {
  const handleSubmit = async (data: CreateExpenseRequest | UpdateExpenseRequest) => {
    await onSubmit(data);
    onOpenChange(false);
  };

  if (!expense) return null;

  const initialData: UpdateExpenseRequest & { id: string } = {
    id: expense.id,
    description: expense.description,
    amount: expense.amount,
    currency: expense.currency,
    category_id: expense.category?.id,
    expense_date: expense.expense_date,
    notes: expense.notes,
    tag_ids: expense.tags?.map((t) => t.id) || [],
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Edit Expense</DialogTitle>
          <DialogDescription>
            Update the expense details below.
          </DialogDescription>
        </DialogHeader>
        <ExpenseForm
          initialData={initialData}
          onSubmit={handleSubmit}
          onCancel={() => onOpenChange(false)}
          isLoading={isLoading}
        />
      </DialogContent>
    </Dialog>
  );
}

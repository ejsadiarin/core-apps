"use client";

import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { ExpenseForm } from "./expense-form";
import type { CreateExpenseRequest, UpdateExpenseRequest } from "@/types/api";

interface ExpenseFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (data: CreateExpenseRequest | UpdateExpenseRequest) => Promise<void>;
  isLoading?: boolean;
}

export function ExpenseFormDialog({
  open,
  onOpenChange,
  onSubmit,
  isLoading,
}: ExpenseFormDialogProps) {
  const handleSubmit = async (data: CreateExpenseRequest | UpdateExpenseRequest) => {
    await onSubmit(data);
    // close dialog on success
    onOpenChange(false);
  };

  const handleCancel = () => {
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Add New Expense</DialogTitle>
        </DialogHeader>
        <ExpenseForm
          onSubmit={handleSubmit}
          onCancel={handleCancel}
          isLoading={isLoading}
        />
      </DialogContent>
    </Dialog>
  );
}

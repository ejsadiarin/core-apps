"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
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
  const router = useRouter();

  // handle browser back button by adding history state
  useEffect(() => {
    if (open) {
      // push a new history state when dialog opens
      window.history.pushState({ dialogOpen: true }, "");

      const handlePopState = (event: PopStateEvent) => {
        // close dialog when back button is pressed
        onOpenChange(false);
      };

      window.addEventListener("popstate", handlePopState);

      return () => {
        window.removeEventListener("popstate", handlePopState);
      };
    }
  }, [open, onOpenChange]);

  const handleSubmit = async (data: CreateExpenseRequest | UpdateExpenseRequest) => {
    await onSubmit(data);
    // close dialog on success
    onOpenChange(false);
  };

  const handleCancel = () => {
    onOpenChange(false);
    // go back in history if we pushed a state
    if (window.history.state?.dialogOpen) {
      router.back();
    }
  };

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen && window.history.state?.dialogOpen) {
      // go back in history when closing
      router.back();
    } else {
      onOpenChange(newOpen);
    }
  };

  // handle ESC key
  const handleEscapeKeyDown = (e: KeyboardEvent) => {
    if (e.key === "Escape") {
      handleCancel();
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent 
        className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto"
        onEscapeKeyDown={handleEscapeKeyDown as any}
      >
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

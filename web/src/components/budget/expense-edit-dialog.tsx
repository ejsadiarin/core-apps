
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@/components/ui/dialog";
import { ExpenseForm } from "./expense-form";
import type { Expense, CreateExpenseRequest, UpdateExpenseRequest } from "@/types/api";

interface EditExpenseDialogProps {
  expense: Expense | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (data: CreateExpenseRequest | UpdateExpenseRequest) => void | Promise<void>;
}

export function EditExpenseDialog({ expense, open, onOpenChange, onSubmit }: EditExpenseDialogProps) {
  const handleSubmit = async (data: CreateExpenseRequest | UpdateExpenseRequest) => {
    onOpenChange(false);
    await onSubmit(data);
  };

  if (!expense) return null;

  const initialData: UpdateExpenseRequest & { id: string } = {
    id: expense.id,
    description: expense.description,
    amount: expense.amount,
    currency: expense.currency,
    category_id: expense.category_id,
    expense_date: expense.expense_date,
    notes: expense.notes,
    recurring_type: expense.recurring_type,
    priority: expense.priority,
    status: expense.status,
    is_debt: expense.is_debt,
    start_date: expense.start_date,
    end_date: expense.end_date,
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Edit Expense</DialogTitle>
          <DialogDescription>Update the details of this expense.</DialogDescription>
        </DialogHeader>
        <ExpenseForm initialData={initialData} onSubmit={handleSubmit} />
      </DialogContent>
    </Dialog>
  );
}

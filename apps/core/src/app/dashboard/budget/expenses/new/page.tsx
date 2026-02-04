"use client";

import { motion } from "motion/react";
import { useRouter } from "next/navigation";
import { ExpenseForm } from "@/components/budget/expense-form";
import { useCreateExpense, GuestBlockedError } from "@/hooks/use-budget";
import { useToast } from "@/components/ui/toast";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import type { CreateExpenseRequest, UpdateExpenseRequest } from "@/types/api";

export default function NewExpensePage() {
  const router = useRouter();
  const createExpense = useCreateExpense();
  const { showToast } = useToast();

  const handleSubmit = async (data: CreateExpenseRequest | UpdateExpenseRequest) => {
    try {
      await createExpense.mutateAsync(data as CreateExpenseRequest);
      router.push("/dashboard/budget");
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, "warning");
      }
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <motion.div
        className="mb-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <Link href="/dashboard/budget">
          <Button variant="ghost" size="sm" className="mb-4">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Budget
          </Button>
        </Link>
        <h1 className="text-primary mb-2">ADD EXPENSE</h1>
        <p className="text-sm text-muted-foreground">
          Create a new expense entry
        </p>
      </motion.div>

      {/* Form */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.1 }}
      >
        <Card className="max-w-2xl mx-auto">
          <CardHeader>
            <CardTitle>Expense Details</CardTitle>
            <CardDescription>
              Fill in the information about your expense
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ExpenseForm
              onSubmit={handleSubmit}
              onCancel={() => router.back()}
              isLoading={createExpense.isPending}
            />
          </CardContent>
        </Card>
      </motion.div>
    </div>
  );
}

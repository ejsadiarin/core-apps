
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useRecurringIncomes, useSubscriptions } from '@/hooks/use-budget';
import { Repeat, ArrowDownRight, ArrowUpRight, DollarSign } from 'lucide-react';
import { cn } from '@/lib/utils';

interface RecurringSummaryCardProps {
  className?: string;
}

export function RecurringSummaryCard({ className }: RecurringSummaryCardProps) {
  const { data: recurringIncomes, isLoading: incomesLoading } = useRecurringIncomes();
  const { data: subscriptionsData, isLoading: subscriptionsLoading } = useSubscriptions();

  const isLoading = incomesLoading || subscriptionsLoading;

  if (isLoading) {
    return (
      <Card className={className}>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium">Recurring Cash Flow</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-20 flex items-center justify-center">
            <div className="h-8 w-8 rounded-full border-2 border-primary border-t-transparent animate-spin" />
          </div>
        </CardContent>
      </Card>
    );
  }

  const totalRecurringIncome = recurringIncomes?.reduce(
    (sum, inc) => sum + inc.monthly_equivalent,
    0
  ) || 0;

  const totalRecurringExpenses = subscriptionsData?.total_monthly || 0;
  const netRecurringCashFlow = totalRecurringIncome - totalRecurringExpenses;

  const recurringIncomeCount = recurringIncomes?.length || 0;
  const recurringExpenseCount = subscriptionsData?.count || 0;

  return (
    <Card className={className}>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">Recurring Cash Flow</CardTitle>
        <CardDescription className="text-xs">Monthly equivalents</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <div className="p-2 rounded-full bg-green-500/10">
                <ArrowDownRight className="h-4 w-4 text-green-500" />
              </div>
              <div>
                <div className="text-xs text-muted-foreground">Income</div>
                <div className="text-lg font-semibold text-green-600">
                  ₱{totalRecurringIncome.toLocaleString()}
                </div>
              </div>
            </div>
            <div className="text-xs text-muted-foreground">
              {recurringIncomeCount} sources
            </div>
          </div>

          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <div className="p-2 rounded-full bg-red-500/10">
                <ArrowUpRight className="h-4 w-4 text-red-500" />
              </div>
              <div>
                <div className="text-xs text-muted-foreground">Expenses</div>
                <div className="text-lg font-semibold text-red-600">
                  ₱{totalRecurringExpenses.toLocaleString()}
                </div>
              </div>
            </div>
            <div className="text-xs text-muted-foreground">
              {recurringExpenseCount} subscriptions
            </div>
          </div>

          <div className="pt-3 border-t">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className={cn(
                  "p-2 rounded-full",
                  netRecurringCashFlow >= 0 ? "bg-green-500/10" : "bg-red-500/10"
                )}>
                  <Repeat className={cn(
                    "h-4 w-4",
                    netRecurringCashFlow >= 0 ? "text-green-500" : "text-red-500"
                  )} />
                </div>
                <div>
                  <div className="text-xs text-muted-foreground">Net Cash Flow</div>
                  <div className={cn(
                    "text-xl font-bold",
                    netRecurringCashFlow >= 0 ? "text-green-600" : "text-red-600"
                  )}>
                    {netRecurringCashFlow >= 0 ? '+' : ''}₱{netRecurringCashFlow.toLocaleString()}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useMonthOverMonth } from '@/hooks/use-budget';
import { TrendingUp, TrendingDown, Minus } from 'lucide-react';
import { cn } from '@/lib/utils';

interface SpendingTrendCardProps {
  className?: string;
}

export function SpendingTrendCard({ className }: SpendingTrendCardProps) {
  const { data, isLoading, error } = useMonthOverMonth();

  if (isLoading) {
    return (
      <Card className={className}>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium">Spending Trends</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-16 flex items-center justify-center">
            <div className="h-8 w-8 rounded-full border-2 border-primary border-t-transparent animate-spin" />
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error || !data) {
    return (
      <Card className={className}>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium">Spending Trends</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-muted-foreground">Failed to load trends</div>
        </CardContent>
      </Card>
    );
  }

  const trends = data.trends || [];
  const latest = trends[trends.length - 1];

  const expenseChange = latest?.expense_change_percent || 0;
  const changeIcon = expenseChange > 0 ? TrendingUp : expenseChange < 0 ? TrendingDown : Minus;
  const ChangeIcon = changeIcon;

  // for expenses, decreasing is good
  const changeColor = expenseChange > 5
    ? 'text-red-500'
    : expenseChange < -5
      ? 'text-green-500'
      : 'text-yellow-500';

  const maxExpense = Math.max(...trends.map((t) => t.expenses), 1);

  return (
    <Card className={className}>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-sm font-medium text-muted-foreground">Spending Trends</CardTitle>
            <CardDescription className="text-xs">
              Avg savings rate: {data.average_savings_rate.toFixed(1)}%
            </CardDescription>
          </div>
          {latest && (
            <div className="flex items-center gap-1">
              <ChangeIcon className={cn('h-4 w-4', changeColor)} />
              <span className={cn('text-sm font-medium', changeColor)}>
                {expenseChange > 0 ? '+' : ''}{expenseChange.toFixed(1)}%
              </span>
            </div>
          )}
        </div>
      </CardHeader>
      <CardContent>
        {trends.length === 0 ? (
          <div className="text-sm text-muted-foreground text-center py-4">
            Not enough data for trends
          </div>
        ) : (
          <>
            {/* mini bar chart of monthly expenses */}
            <div className="flex items-end gap-1 h-20 mb-3">
              {trends.slice(-6).map((t) => {
                const heightPct = (t.expenses / maxExpense) * 100;
                return (
                  <div key={t.month} className="flex-1 flex flex-col items-center gap-0.5">
                    <div className="w-full relative" style={{ height: '60px' }}>
                      <div
                        className="absolute bottom-0 w-full rounded-t bg-primary/60 transition-all"
                        style={{ height: `${Math.max(heightPct, 4)}%` }}
                      />
                    </div>
                    <div className="text-[9px] text-muted-foreground truncate w-full text-center">
                      {t.month.slice(5, 7)}/{t.month.slice(2, 4)}
                    </div>
                  </div>
                );
              })}
            </div>

            {/* current month details */}
            {latest && (
              <div className="space-y-1 pt-3 border-t border-border">
                <div className="flex justify-between text-xs">
                  <span className="text-muted-foreground">Income</span>
                  <span className="font-medium">₱{latest.income.toLocaleString()}</span>
                </div>
                <div className="flex justify-between text-xs">
                  <span className="text-muted-foreground">Expenses</span>
                  <span className="font-medium">₱{latest.expenses.toLocaleString()}</span>
                </div>
                <div className="flex justify-between text-xs">
                  <span className="text-muted-foreground">Savings Rate</span>
                  <span className={cn(
                    'font-medium',
                    latest.savings_rate >= 20 ? 'text-green-500' : latest.savings_rate >= 10 ? 'text-yellow-500' : 'text-red-500'
                  )}>
                    {latest.savings_rate.toFixed(1)}%
                  </span>
                </div>
              </div>
            )}
          </>
        )}
      </CardContent>
    </Card>
  );
}

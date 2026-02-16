'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useWeekdayPattern } from '@/hooks/use-budget';
import { cn } from '@/lib/utils';

interface WeekdaySpendingChartProps {
  startDate?: string;
  endDate?: string;
  className?: string;
}

export function WeekdaySpendingChart({ startDate, endDate, className }: WeekdaySpendingChartProps) {
  const { data, isLoading, error } = useWeekdayPattern(startDate, endDate);

  if (isLoading) {
    return (
      <Card className={className}>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium">Spending by Day</CardTitle>
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
          <CardTitle className="text-sm font-medium">Spending by Day</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-muted-foreground">Failed to load weekday data</div>
        </CardContent>
      </Card>
    );
  }

  const maxAmount = Math.max(...data.weekdays.map((d) => d.total_amount), 1);

  return (
    <Card className={className}>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">Spending by Day</CardTitle>
        <CardDescription className="text-xs">
          Highest: {data.highest_spending_day} · Lowest: {data.lowest_spending_day}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex items-end gap-2 h-32">
          {data.weekdays.map((day) => {
            const heightPct = (day.total_amount / maxAmount) * 100;
            const isHighest = day.day === data.highest_spending_day;
            const isLowest = day.day === data.lowest_spending_day;

            return (
              <div key={day.day} className="flex-1 flex flex-col items-center gap-1">
                <div className="text-[10px] text-muted-foreground font-medium">
                  ₱{Math.round(day.total_amount).toLocaleString()}
                </div>
                <div className="w-full relative" style={{ height: '80px' }}>
                  <div
                    className={cn(
                      'absolute bottom-0 w-full rounded-t transition-all',
                      isHighest ? 'bg-red-500/80' : isLowest ? 'bg-green-500/80' : 'bg-primary/60'
                    )}
                    style={{ height: `${Math.max(heightPct, 4)}%` }}
                  />
                </div>
                <div className="text-[10px] text-muted-foreground">{day.day.slice(0, 3)}</div>
              </div>
            );
          })}
        </div>

        <div className="mt-4 pt-3 border-t border-border grid grid-cols-2 gap-2">
          <div className="text-xs text-muted-foreground">
            Avg highest: <span className="font-medium text-foreground">
              ₱{data.weekdays.find((d) => d.day === data.highest_spending_day)?.average_amount.toFixed(0) || '0'}
            </span>
          </div>
          <div className="text-xs text-muted-foreground">
            Avg lowest: <span className="font-medium text-foreground">
              ₱{data.weekdays.find((d) => d.day === data.lowest_spending_day)?.average_amount.toFixed(0) || '0'}
            </span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

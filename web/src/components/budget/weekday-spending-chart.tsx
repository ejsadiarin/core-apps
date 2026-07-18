'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { useTrends } from '@/hooks/use-budget';
import { cn } from '@/lib/utils';

interface WeekdaySpendingChartProps {
  startDate?: string;
  endDate?: string;
  className?: string;
}

export function WeekdaySpendingChart({ startDate, endDate, className }: WeekdaySpendingChartProps) {
  const { data, isLoading, error } = useTrends('day', startDate, endDate);

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

  const maxAmount = Math.max(...data.map((d) => d.total_amount), 1);
  const minAmountItem = data.reduce((min, item) => item.total_amount < min.total_amount ? item : min, data[0]);
  const maxAmountItem = data.reduce((max, item) => item.total_amount > max.total_amount ? item : max, data[0]);

  // Helper to format date for display
  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  };

  // Helper to abbreviate amounts for compact display
  const formatAmount = (amount: number) => {
    if (amount >= 10000) {
      return `₱${(amount / 1000).toFixed(1)}k`;
    }
    return `₱${Math.round(amount).toLocaleString()}`;
  };

  return (
    <Card className={className}>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">Spending by Day</CardTitle>
        <CardDescription className="text-xs">
          {startDate && endDate ? `${formatDate(startDate)} - ${formatDate(endDate)}` : 'This Month'}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto overflow-y-hidden">
          <div className="flex items-end gap-1 h-32 min-w-0" style={{ minWidth: `${Math.max(data.length * 36, 100)}px` }}>
            {data.map((day) => {
              const heightPct = (day.total_amount / maxAmount) * 100;
              const isHighest = day.date === maxAmountItem?.date;
              const isLowest = day.date === minAmountItem?.date;

              return (
                <TooltipProvider>
                <Tooltip key={day.date}>
                  <TooltipTrigger asChild>
                    <div className="flex flex-col items-center gap-1 min-w-[32px] flex-shrink-0">
                      <div className="text-[10px] text-muted-foreground font-medium overflow-hidden text-ellipsis whitespace-nowrap max-w-[40px]">
                        {formatAmount(day.total_amount)}
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
                      <div className="text-[10px] text-muted-foreground">{new Date(day.date).getDate()}</div>
                    </div>
                  </TooltipTrigger>
                  <TooltipContent>
                    <p>{formatDate(day.date)}: ₱{day.total_amount.toLocaleString()}</p>
                  </TooltipContent>
                </Tooltip>
                </TooltipProvider>
              );
            })}
          </div>
        </div>

        <div className="mt-4 pt-3 border-t border-border grid grid-cols-2 gap-2">
          <div className="text-xs text-muted-foreground">
            Highest: <span className="font-medium text-foreground">
              ₱{maxAmountItem?.total_amount.toLocaleString()}
            </span>
          </div>
          <div className="text-xs text-muted-foreground">
            Lowest: <span className="font-medium text-foreground">
              ₱{minAmountItem?.total_amount.toLocaleString()}
            </span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

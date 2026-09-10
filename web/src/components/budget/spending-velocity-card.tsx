
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useSpendingVelocity } from '@/hooks/use-budget';
import { Gauge } from 'lucide-react';
import { cn, safeFixed } from '@/lib/utils';

interface SpendingVelocityCardProps {
  startDate?: string;
  endDate?: string;
  className?: string;
}

export function SpendingVelocityCard({ startDate, endDate, className }: SpendingVelocityCardProps) {
  const { data, isLoading, error } = useSpendingVelocity(startDate, endDate);

  if (isLoading) {
    return (
      <Card className={className}>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium">Spending Velocity</CardTitle>
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
          <CardTitle className="text-sm font-medium">Spending Velocity</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-muted-foreground">Failed to load spending velocity</div>
        </CardContent>
      </Card>
    );
  }

  const avgMonthly = Number(data.avg_monthly_spending) || 0;
  const monthsData = data.months_with_data ?? 0;

  return (
    <Card className={className}>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">Avg Monthly Spending</CardTitle>
        <CardDescription className="text-xs">
          Based on {monthsData} month{monthsData !== 1 ? 's' : ''} of data
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between">
          <div>
            <div className="text-3xl font-bold text-foreground">
              ₱{safeFixed(avgMonthly)}
            </div>
            <div className="text-xs text-muted-foreground mt-1">
              per month average
            </div>
          </div>
          <div className={cn('p-3 rounded-full bg-muted')}>
            <Gauge className={cn('h-5 w-5 text-muted-foreground')} />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

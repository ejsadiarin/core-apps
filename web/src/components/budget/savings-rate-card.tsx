
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useSavingsRate } from '@/hooks/use-budget';
import { TrendingUp, TrendingDown, Minus } from 'lucide-react';
import { cn, safeFixed } from '@/lib/utils';

interface SavingsRateCardProps {
  startDate?: string;
  endDate?: string;
  className?: string;
}

export function SavingsRateCard({ startDate, endDate, className }: SavingsRateCardProps) {
  const { data, isLoading, error } = useSavingsRate(startDate, endDate);

  if (isLoading) {
    return (
      <Card className={className}>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium">Savings Rate</CardTitle>
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
          <CardTitle className="text-sm font-medium">Savings Rate</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-muted-foreground">Failed to load savings rate</div>
        </CardContent>
      </Card>
    );
  }

  const rate = Number(data.savings_rate) || 0;
  const income = Number(data.total_incomes) || 0;
  const expenses = Number(data.total_expenses) || 0;
  const savings = income - expenses;

  let status: 'excellent' | 'good' | 'fair' | 'poor' | 'negative';
  if (rate >= 20) status = 'excellent';
  else if (rate >= 15) status = 'good';
  else if (rate >= 10) status = 'fair';
  else if (rate >= 0) status = 'poor';
  else status = 'negative';

  const statusConfig = {
    excellent: { color: 'text-green-500', bgColor: 'bg-green-500/10', icon: TrendingUp },
    good: { color: 'text-emerald-500', bgColor: 'bg-emerald-500/10', icon: TrendingUp },
    fair: { color: 'text-yellow-500', bgColor: 'bg-yellow-500/10', icon: Minus },
    poor: { color: 'text-orange-500', bgColor: 'bg-orange-500/10', icon: TrendingDown },
    negative: { color: 'text-red-500', bgColor: 'bg-red-500/10', icon: TrendingDown }
  };

  const config = statusConfig[status];
  const Icon = config.icon;

  return (
    <Card className={className}>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">Savings Rate</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between">
          <div>
            <div className={cn('text-3xl font-bold', config.color)}>
              {safeFixed(rate, 1)}%
            </div>
            <div className="text-xs text-muted-foreground mt-1 capitalize">
              {status.replace('_', ' ')}
            </div>
          </div>
          <div className={cn('p-3 rounded-full', config.bgColor)}>
            <Icon className={cn('h-5 w-5', config.color)} />
          </div>
        </div>

        <div className="mt-4 space-y-1">
          <div className="flex justify-between text-xs">
            <span className="text-muted-foreground">Income</span>
            <span className="font-medium">₱{income.toLocaleString()}</span>
          </div>
          <div className="flex justify-between text-xs">
            <span className="text-muted-foreground">Expenses</span>
            <span className="font-medium">₱{expenses.toLocaleString()}</span>
          </div>
          <div className="flex justify-between text-xs">
            <span className="text-muted-foreground">Savings</span>
            <span className={cn('font-medium', savings >= 0 ? 'text-green-500' : 'text-red-500')}>
              ₱{savings.toLocaleString()}
            </span>
          </div>
        </div>

        <div className="mt-4 pt-3 border-t border-border">
          <div className="text-xs text-muted-foreground">
            Target: <span className="font-medium text-foreground">15-20%</span>
          </div>
          <div className="w-full bg-muted rounded-full h-2 mt-2">
            <div
              className={cn('h-2 rounded-full transition-all', config.color.replace('text-', 'bg-'))}
              style={{ width: `${Math.min(Math.max(rate, 0), 100)}%` }}
            />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

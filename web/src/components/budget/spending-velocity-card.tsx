
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { useSpendingVelocity } from '@/hooks/use-budget';
import { Gauge, AlertTriangle, CheckCircle, XCircle } from 'lucide-react';
import { cn } from '@/lib/utils';

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

  const statusConfig = {
    on_track: {
      color: 'text-green-500',
      bgColor: 'bg-green-500',
      lightBgColor: 'bg-green-500/10',
      icon: CheckCircle,
      label: 'On Track',
      message: 'You are spending within budget'
    },
    warning: {
      color: 'text-yellow-500',
      bgColor: 'bg-yellow-500',
      lightBgColor: 'bg-yellow-500/10',
      icon: AlertTriangle,
      label: 'Warning',
      message: 'You are approaching your budget limit'
    },
    at_risk: {
      color: 'text-orange-500',
      bgColor: 'bg-orange-500',
      lightBgColor: 'bg-orange-500/10',
      icon: AlertTriangle,
      label: 'At Risk',
      message: 'You may exceed your budget this month'
    },
    over_pace: {
      color: 'text-red-500',
      bgColor: 'bg-red-500',
      lightBgColor: 'bg-red-500/10',
      icon: XCircle,
      label: 'Over Pace',
      message: 'You are spending faster than planned'
    },
    unknown: {
      color: 'text-gray-500',
      bgColor: 'bg-gray-500',
      lightBgColor: 'bg-gray-500/10',
      icon: Gauge,
      label: 'Unknown',
      message: 'Set budgets to track your spending velocity'
    }
  };

  const config = statusConfig[data.status];
  const Icon = config.icon;

  const percentage = data.total_budget > 0
    ? (data.projected_spend / data.total_budget) * 100
    : 0;

  return (
    <Card className={className}>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">Spending Velocity</CardTitle>
        <CardDescription className="text-xs">
          Day {data.days_elapsed} of {data.days_in_month}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between">
          <div>
            <div className={cn('text-3xl font-bold', config.color)}>
              ₱{Math.round(data.projected_spend).toLocaleString()}
            </div>
            <div className="text-xs text-muted-foreground mt-1">
              Projected spend
            </div>
          </div>
          <div className={cn('p-3 rounded-full', config.lightBgColor)}>
            <Icon className={cn('h-5 w-5', config.color)} />
          </div>
        </div>

        {data.total_budget > 0 && (
          <div className="mt-4">
            <div className="flex justify-between text-xs mb-1">
              <span className="text-muted-foreground">Budget Usage</span>
              <span className={cn('font-medium', config.color)}>{percentage.toFixed(0)}%</span>
            </div>
            <Progress
              value={Math.min(percentage, 100)}
              className="h-2"
            />
            <div className="flex justify-between text-xs mt-1 text-muted-foreground">
              <span>₱{data.amount_spent.toLocaleString()} spent</span>
              <span>₱{data.total_budget.toLocaleString()} budget</span>
            </div>
          </div>
        )}

        <div className={cn('mt-4 p-3 rounded-lg text-xs', config.lightBgColor)}>
          <div className={cn('font-medium mb-1', config.color)}>{config.label}</div>
          <div className="text-muted-foreground">{config.message}</div>
        </div>
      </CardContent>
    </Card>
  );
}

'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useFiftyThirtyTwenty } from '@/hooks/use-budget';
import { cn } from '@/lib/utils';

interface FiftyThirtyTwentyChartProps {
  startDate?: string;
  endDate?: string;
  className?: string;
}

export function FiftyThirtyTwentyChart({ startDate, endDate, className }: FiftyThirtyTwentyChartProps) {
  const { data, isLoading, error } = useFiftyThirtyTwenty(startDate, endDate);

  if (isLoading) {
    return (
      <Card className={className}>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium">50/30/20 Budget</CardTitle>
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
          <CardTitle className="text-sm font-medium">50/30/20 Budget</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-muted-foreground">Failed to load analysis</div>
        </CardContent>
      </Card>
    );
  }

  const categories = [
    {
      label: 'Needs',
      data: data.needs,
      color: 'bg-blue-500',
      textColor: 'text-blue-500',
      target: 50
    },
    {
      label: 'Wants',
      data: data.wants,
      color: 'bg-purple-500',
      textColor: 'text-purple-500',
      target: 30
    },
    {
      label: 'Savings',
      data: data.investments,
      color: 'bg-green-500',
      textColor: 'text-green-500',
      target: 20
    }
  ];

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'on_target': return 'text-green-500';
      case 'over': return 'text-red-500';
      case 'under': return 'text-yellow-500';
      default: return 'text-muted-foreground';
    }
  };

  return (
    <Card className={className}>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">50/30/20 Budget</CardTitle>
        <CardDescription className="text-xs">
          Income: ₱{data.total_income.toLocaleString()}
        </CardDescription>
      </CardHeader>
      <CardContent>
        {/* stacked bar visualization */}
        <div className="mb-4">
          <div className="flex rounded-full overflow-hidden h-4">
            {categories.map((c) => (
              <div
                key={c.label}
                className={cn(c.color, 'transition-all')}
                style={{ width: `${c.data.actual_percentage || 0}%` }}
                title={`${c.label}: ${c.data.actual_percentage.toFixed(1)}%`}
              />
            ))}
          </div>
          {/* target markers */}
          <div className="relative h-2 mt-1">
            <div className="absolute left-[50%] w-px h-2 bg-muted-foreground/50" title="50% target" />
            <div className="absolute left-[80%] w-px h-2 bg-muted-foreground/50" title="80% target" />
          </div>
        </div>

        <div className="space-y-3">
          {categories.map((c) => (
            <div key={c.label} className="space-y-1">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className={cn('w-3 h-3 rounded-full', c.color)} />
                  <span className="text-sm font-medium">{c.label}</span>
                </div>
                <span className={cn('text-xs font-medium', getStatusColor(c.data.status))}>
                  {c.data.actual_percentage.toFixed(1)}% / {c.data.target_percentage}%
                </span>
              </div>
              <div className="flex justify-between text-xs text-muted-foreground">
                <span>₱{c.data.amount.toLocaleString()}</span>
                <span className="capitalize">{c.data.status.replace('_', ' ')}</span>
              </div>
            </div>
          ))}

          {data.unclassified_count > 0 && (
            <div className="space-y-1 pt-2 border-t">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 rounded-full bg-muted-foreground/30" />
                  <span className="text-sm font-medium text-muted-foreground">Unclassified</span>
                </div>
                <span className="text-xs font-medium text-muted-foreground">
                  {data.unclassified_count} expense{data.unclassified_count !== 1 ? 's' : ''}
                </span>
              </div>
              <div className="flex justify-between text-xs text-muted-foreground">
                <span>₱{data.unclassified_amount.toLocaleString()}</span>
                <span>Not assigned to a priority group</span>
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

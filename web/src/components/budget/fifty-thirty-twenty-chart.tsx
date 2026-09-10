
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useFiftyThirtyTwenty } from '@/hooks/use-budget';
import { safeFixed, cn } from '@/lib/utils';

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

  const totalIncome = Number(data.total_incomes) || 0;
  const needs = Number(data.needs) || 0;
  const wants = Number(data.wants) || 0;
  const savings = Number(data.savings) || 0;
  const needsPct = Number(data.needs_pct) || 0;
  const wantsPct = Number(data.wants_pct) || 0;
  const savingsPct = Number(data.savings_pct) || 0;

  const categories = [
    { label: 'Needs', amount: needs, pct: needsPct, target: 50, color: 'bg-blue-500', textColor: 'text-blue-500' },
    { label: 'Wants', amount: wants, pct: wantsPct, target: 30, color: 'bg-purple-500', textColor: 'text-purple-500' },
    { label: 'Savings', amount: savings, pct: savingsPct, target: 20, color: 'bg-green-500', textColor: 'text-green-500' },
  ];

  return (
    <Card className={className}>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">50/30/20 Budget</CardTitle>
        <CardDescription className="text-xs">
          Income: ₱{totalIncome.toLocaleString()}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="mb-4">
          <div className="flex rounded-full overflow-hidden h-4">
            {categories.map((c) => (
              <div
                key={c.label}
                className={cn(c.color, 'transition-all')}
                style={{ width: `${c.pct}%` }}
                title={`${c.label}: ${safeFixed(c.pct, 1)}%`}
              />
            ))}
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
                <span className="text-xs font-medium">
                  {safeFixed(c.pct, 1)}% / {c.target}%
                </span>
              </div>
              <div className="flex justify-between text-xs text-muted-foreground">
                <span>₱{c.amount.toLocaleString()}</span>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

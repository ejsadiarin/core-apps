
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useMonthOverMonth } from '@/hooks/use-budget';

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

  const trends = Array.isArray(data) ? data : [];
  const maxAmount = Math.max(...trends.map((t) => Number(t.total_amount) || 0), 1);

  return (
    <Card className={className}>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">Spending Trends</CardTitle>
        <CardDescription className="text-xs">
          {trends.length} month{trends.length !== 1 ? 's' : ''} of data
        </CardDescription>
      </CardHeader>
      <CardContent>
        {trends.length === 0 ? (
          <div className="text-sm text-muted-foreground text-center py-4">
            Not enough data for trends
          </div>
        ) : (
          <div className="flex items-end gap-1 h-20">
            {trends.slice(-6).map((t) => {
              const amount = Number(t.total_amount) || 0;
              const heightPct = (amount / maxAmount) * 100;
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
        )}
      </CardContent>
    </Card>
  );
}

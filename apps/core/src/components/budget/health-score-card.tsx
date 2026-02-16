'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useHealthScore } from '@/hooks/use-budget';
import { Heart, ShieldCheck, AlertTriangle, XCircle } from 'lucide-react';
import { cn } from '@/lib/utils';

interface HealthScoreCardProps {
  className?: string;
}

export function HealthScoreCard({ className }: HealthScoreCardProps) {
  const { data, isLoading, error } = useHealthScore();

  if (isLoading) {
    return (
      <Card className={className}>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium">Financial Health</CardTitle>
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
          <CardTitle className="text-sm font-medium">Financial Health</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-muted-foreground">Failed to load health score</div>
        </CardContent>
      </Card>
    );
  }

  const statusConfig: Record<string, { color: string; bgColor: string; icon: typeof Heart }> = {
    excellent: { color: 'text-green-500', bgColor: 'bg-green-500/10', icon: ShieldCheck },
    good: { color: 'text-emerald-500', bgColor: 'bg-emerald-500/10', icon: Heart },
    fair: { color: 'text-yellow-500', bgColor: 'bg-yellow-500/10', icon: AlertTriangle },
    poor: { color: 'text-red-500', bgColor: 'bg-red-500/10', icon: XCircle }
  };

  const config = statusConfig[data.status] || statusConfig.fair;
  const Icon = config.icon;

  // score gauge color based on value
  const getScoreColor = (score: number) => {
    if (score >= 80) return 'bg-green-500';
    if (score >= 60) return 'bg-emerald-500';
    if (score >= 40) return 'bg-yellow-500';
    return 'bg-red-500';
  };

  return (
    <Card className={className}>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-sm font-medium text-muted-foreground">Financial Health</CardTitle>
            <CardDescription className="text-xs capitalize">{data.status}</CardDescription>
          </div>
          <div className={cn('p-3 rounded-full', config.bgColor)}>
            <Icon className={cn('h-5 w-5', config.color)} />
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex items-center gap-4 mb-4">
          <div className={cn('text-4xl font-bold', config.color)}>{data.score}</div>
          <div className="text-xs text-muted-foreground">/ 100</div>
        </div>

        <div className="w-full bg-muted rounded-full h-3 mb-4">
          <div
            className={cn('h-3 rounded-full transition-all', getScoreColor(data.score))}
            style={{ width: `${Math.min(data.score, 100)}%` }}
          />
        </div>

        <div className="space-y-2">
          <div className="flex justify-between text-xs">
            <span className="text-muted-foreground">Savings Rate</span>
            <span className="font-medium">{data.savings_rate.toFixed(1)}%</span>
          </div>
          <div className="flex justify-between text-xs">
            <span className="text-muted-foreground">Emergency Fund</span>
            <span className="font-medium">{data.emergency_fund_months.toFixed(1)} months</span>
          </div>
        </div>

        {data.recommendations.length > 0 && (
          <div className="mt-4 pt-3 border-t border-border">
            <div className="text-xs font-medium mb-2">Recommendations</div>
            <div className="space-y-1">
              {data.recommendations.slice(0, 3).map((rec, i) => (
                <div key={i} className="text-xs text-muted-foreground flex items-start gap-1">
                  <span className="shrink-0 mt-0.5">•</span>
                  <span>{rec}</span>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

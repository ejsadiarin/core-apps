
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { useCurrentTotalMoney } from '@/hooks/use-budget';
import { TrendingUp, TrendingDown, Wallet, Info } from 'lucide-react';
import { cn } from '@/lib/utils';

interface CurrentTotalMoneyCardProps {
  startDate?: string;
  endDate?: string;
  className?: string;
}

export function CurrentTotalMoneyCard({ startDate, endDate, className }: CurrentTotalMoneyCardProps) {
  const { data, isLoading, error } = useCurrentTotalMoney(startDate, endDate);

  if (isLoading) {
    return (
      <Card className={className}>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium">Total Money</CardTitle>
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
          <CardTitle className="text-sm font-medium">Total Money</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-muted-foreground">Failed to load total money</div>
        </CardContent>
      </Card>
    );
  }

  const isPositiveChange = data.net_change >= 0;
  const Icon = isPositiveChange ? TrendingUp : TrendingDown;

  return (
      <TooltipProvider>
     
    <Card className={className}>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">Total Money</CardTitle>
        <CardDescription className="text-xs">
          Since {new Date(data.tracking_start_date).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between">
          <div>
            <div className="text-3xl font-bold">
              ₱{data.current_total.toLocaleString()}
            </div>
            <div className={cn('text-xs mt-1 flex items-center gap-1', isPositiveChange ? 'text-green-500' : 'text-red-500')}>
              <Icon className="h-3 w-3" />
              <span>₱{Math.abs(data.net_change).toLocaleString()}</span>
            </div>
          </div>
          <div className="p-3 rounded-full bg-primary/10">
            <Wallet className="h-5 w-5 text-primary" />
          </div>
        </div>

        <div className="mt-4 space-y-1">
          <div className="flex justify-between text-xs">
            <span className="text-muted-foreground flex items-center gap-1">
              Starting Baseline
              <Tooltip>
                <TooltipTrigger asChild>
                  <Info className="h-3 w-3 text-muted-foreground/50 cursor-help" />
                </TooltipTrigger>
                <TooltipContent side="top" className="max-w-[200px]">
                  <p>Your total money at the tracking start date, calculated backwards from your current balance.</p>
                </TooltipContent>
              </Tooltip>
            </span>
            <span className="font-medium">₱{data.money_baseline.toLocaleString()}</span>
          </div>
          <div className="flex justify-between text-xs">
            <span className="text-muted-foreground">Income</span>
            <span className="font-medium text-green-500">+₱{data.income_since_start.toLocaleString()}</span>
          </div>
          <div className="flex justify-between text-xs">
            <span className="text-muted-foreground">Expenses</span>
            <span className="font-medium text-red-500">-₱{data.expenses_since_start.toLocaleString()}</span>
          </div>
        </div>

        <div className="mt-4 pt-3 border-t border-border">
          <div className="text-xs text-muted-foreground">
            Net Change: <span className={cn('font-medium', isPositiveChange ? 'text-green-500' : 'text-red-500')}>
              {isPositiveChange ? '+' : ''}₱{data.net_change.toLocaleString()}
            </span>
          </div>
        </div>
      </CardContent>
    </Card>
 </TooltipProvider>
  );
}

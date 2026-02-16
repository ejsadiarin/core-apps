'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useMerchantAnalysis } from '@/hooks/use-budget';
import { Store } from 'lucide-react';
import { cn } from '@/lib/utils';

interface TopMerchantsTableProps {
  limit?: number;
  startDate?: string;
  endDate?: string;
  className?: string;
}

export function TopMerchantsTable({ limit = 10, startDate, endDate, className }: TopMerchantsTableProps) {
  const { data, isLoading, error } = useMerchantAnalysis(limit, startDate, endDate);

  if (isLoading) {
    return (
      <Card className={className}>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium">Top Merchants</CardTitle>
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
          <CardTitle className="text-sm font-medium">Top Merchants</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-muted-foreground">Failed to load merchant data</div>
        </CardContent>
      </Card>
    );
  }

  const merchants = data.merchants || [];

  return (
    <Card className={className}>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-sm font-medium text-muted-foreground">Top Merchants</CardTitle>
            <CardDescription className="text-xs">
              {data.unique_merchant_count} unique merchants · ${data.total_spent.toLocaleString()} total
            </CardDescription>
          </div>
          <Store className="h-4 w-4 text-muted-foreground" />
        </div>
      </CardHeader>
      <CardContent>
        {merchants.length === 0 ? (
          <div className="text-sm text-muted-foreground text-center py-4">
            No merchant data available
          </div>
        ) : (
          <div className="space-y-3">
            {merchants.map((merchant, i) => (
              <div key={merchant.name} className="space-y-1">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <span className="text-xs text-muted-foreground w-5 text-right">{i + 1}.</span>
                    <span className="text-sm font-medium truncate max-w-[200px]">{merchant.name}</span>
                  </div>
                  <span className="text-sm font-medium">${merchant.total_spent.toLocaleString()}</span>
                </div>
                <div className="flex items-center gap-2 ml-7">
                  <div className="flex-1 bg-muted rounded-full h-1.5">
                    <div
                      className="h-1.5 rounded-full bg-primary/70 transition-all"
                      style={{ width: `${merchant.percentage}%` }}
                    />
                  </div>
                  <span className="text-[10px] text-muted-foreground w-10 text-right">
                    {merchant.percentage.toFixed(1)}%
                  </span>
                </div>
                <div className="flex justify-between text-[10px] text-muted-foreground ml-7">
                  <span>{merchant.count} transactions</span>
                  <span>avg ${merchant.average_amount.toFixed(0)}</span>
                </div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

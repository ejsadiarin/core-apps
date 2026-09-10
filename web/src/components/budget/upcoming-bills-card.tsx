
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { useUpcomingBills } from '@/hooks/use-budget';
import { Calendar, Clock, DollarSign } from 'lucide-react';
import { cn, safeFormat } from '@/lib/utils';
import { format, parseISO, isWithinInterval, addDays } from 'date-fns';

interface UpcomingBillsCardProps {
  days?: 7 | 30;
  className?: string;
}

export function UpcomingBillsCard({ days = 30, className }: UpcomingBillsCardProps) {
  const { data, isLoading, error } = useUpcomingBills();

  if (isLoading) {
    return (
      <Card className={className}>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium">Upcoming Bills</CardTitle>
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
          <CardTitle className="text-sm font-medium">Upcoming Bills</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-muted-foreground">Failed to load upcoming bills</div>
        </CardContent>
      </Card>
    );
  }

  const today = new Date();
  const next7Days = addDays(today, 7);

  const getDueStatus = (dueDate: string) => {
    const date = parseISO(dueDate);
    if (isWithinInterval(date, { start: today, end: next7Days })) {
      return { label: 'Soon', color: 'bg-yellow-500/10 text-yellow-500 border-yellow-500/20' };
    }
    return { label: 'Upcoming', color: 'bg-blue-500/10 text-blue-500 border-blue-500/20' };
  };

  const displayedBills = (data.bills || []).slice(0, 5);

  return (
    <Card className={className}>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-sm font-medium text-muted-foreground">Upcoming Bills</CardTitle>
            <CardDescription className="text-xs">
              Next {days} days
            </CardDescription>
          </div>
          <div className="text-right">
            <div className="text-2xl font-bold">₱{Math.round(data.total_amount).toLocaleString()}</div>
            <div className="text-xs text-muted-foreground">Total due</div>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        {(data.bills || []).length === 0 ? (
          <div className="text-sm text-muted-foreground text-center py-4">
            No upcoming bills found
          </div>
        ) : (
          <div className="space-y-3">
            {displayedBills.map((bill) => {
              const status = getDueStatus(bill.due_date);
              return (
                <div
                  key={`${bill.id}-${bill.due_date}`}
                  className="flex items-center justify-between p-2 rounded-lg border border-border/50 hover:bg-muted/50 transition-colors"
                >
                  <div className="flex items-center gap-3">
                    <div className="p-2 rounded-full bg-muted">
                      <Calendar className="h-4 w-4 text-muted-foreground" />
                    </div>
                    <div>
                      <div className="text-sm font-medium">{bill.description}</div>
                      <div className="text-xs text-muted-foreground flex items-center gap-1">
                        <Clock className="h-3 w-3" />
                        Due {safeFormat(bill.due_date, (d) => format(d, 'MMM d'))}
                      </div>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-sm font-medium">₱{bill.amount.toLocaleString()}</div>
                    <Badge variant="outline" className={cn('text-xs mt-1', status.color)}>
                      {status.label}
                    </Badge>
                  </div>
                </div>
              );
            })}

            {(data.bills || []).length > 5 && (
              <div className="text-center text-xs text-muted-foreground pt-2">
                +{(data.bills || []).length - 5} more bills
              </div>
            )}
          </div>
        )}

        <div className="mt-4 pt-3 border-t border-border">
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            <DollarSign className="h-3 w-3" />
            <span>{(data.bills || []).length} recurring payments</span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

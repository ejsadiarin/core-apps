
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { Badge } from '@/components/ui/badge';
import { useCategoryBudgets, useDeleteCategoryBudget } from '@/hooks/use-budget';
import { CategoryBudgetForm } from './category-budget-form';
import { format } from 'date-fns';
import { AlertTriangle, CheckCircle, TrendingDown, Trash2 } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';

interface BudgetVarianceTableProps {
  month?: string;
  className?: string;
}

export function BudgetVarianceTable({ month, className }: BudgetVarianceTableProps) {
  const { data: budgets, isLoading, error } = useCategoryBudgets(month);
  const deleteBudget = useDeleteCategoryBudget();

  const currentMonth = month || format(new Date(), 'yyyy-MM-dd');

  if (isLoading) {
    return (
      <Card className={className}>
        <CardHeader>
          <CardTitle className="text-sm font-medium">Budget vs Actual</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-32 flex items-center justify-center">
            <div className="h-8 w-8 rounded-full border-2 border-primary border-t-transparent animate-spin" />
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className={className}>
        <CardHeader>
          <CardTitle className="text-sm font-medium">Budget vs Actual</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-muted-foreground">Failed to load budget data</div>
        </CardContent>
      </Card>
    );
  }

  const getStatusConfig = (percentage: number) => {
    if (percentage >= 100) {
      return {
        icon: AlertTriangle,
        color: 'text-red-500',
        bgColor: 'bg-red-500',
        label: 'Over Budget',
        badgeColor: 'bg-red-500/10 text-red-500 border-red-500/20'
      };
    }
    if (percentage >= 80) {
      return {
        icon: TrendingDown,
        color: 'text-yellow-500',
        bgColor: 'bg-yellow-500',
        label: 'Warning',
        badgeColor: 'bg-yellow-500/10 text-yellow-500 border-yellow-500/20'
      };
    }
    return {
      icon: CheckCircle,
      color: 'text-green-500',
      bgColor: 'bg-green-500',
      label: 'On Track',
      badgeColor: 'bg-green-500/10 text-green-500 border-green-500/20'
    };
  };

  const categoriesWithBudgets = budgets?.filter(b => b.budget_amount > 0) || [];
  const categoriesWithoutBudgets = budgets?.filter(b => b.budget_amount === 0) || [];

  return (
    <Card className={className}>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-sm font-medium text-muted-foreground">Budget vs Actual</CardTitle>
            <CardDescription className="text-xs">
              {format(new Date(currentMonth), 'MMMM yyyy')}
            </CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        {categoriesWithBudgets.length === 0 ? (
          <div className="text-center py-8">
            <div className="text-sm text-muted-foreground mb-4">
              No budgets set for this month
            </div>
            <p className="text-xs text-muted-foreground">
              Set budgets for your categories to track spending
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {categoriesWithBudgets.map((budget) => {
              const percentage = budget.percentage;
              const status = getStatusConfig(percentage);
              const Icon = status.icon;

              return (
                <div
                  key={budget.category_id}
                  className="space-y-2 p-3 rounded-lg border border-border/50 hover:bg-muted/30 transition-colors"
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      {budget.category_color && (
                        <div
                          className="w-3 h-3 rounded-full"
                          style={{ backgroundColor: budget.category_color }}
                        />
                      )}
                      <span className="font-medium text-sm">{budget.category_name}</span>
                      <Badge variant="outline" className={cn('text-xs', status.badgeColor)}>
                        {status.label}
                      </Badge>
                    </div>
                    <div className="flex items-center gap-2">
                      <CategoryBudgetForm
                        categoryBudget={budget}
                        trigger={
                          <Button variant="ghost" size="icon" className="h-8 w-8">
                            <Icon className={cn('h-4 w-4', status.color)} />
                          </Button>
                        }
                      />
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 text-muted-foreground hover:text-red-500"
                        onClick={() => deleteBudget.mutate(budget.category_id)}
                        disabled={deleteBudget.isPending}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>

                  <div className="space-y-1">
                    <Progress
                      value={Math.min(percentage, 100)}
                      className="h-2"
                    />
                    <div className="flex justify-between text-xs text-muted-foreground">
                      <span>
                        ₱{budget.spent_amount.toLocaleString()} of ₱{budget.budget_amount.toLocaleString()}
                      </span>
                      <span className={status.color}>{percentage.toFixed(0)}%</span>
                    </div>
                  </div>

                  {budget.variance !== undefined && (
                    <div className={cn(
                      'text-xs',
                      budget.variance >= 0 ? 'text-green-500' : 'text-red-500'
                    )}>
                      {budget.variance >= 0 ? '+' : ''}₱{budget.variance.toLocaleString()} remaining
                    </div>
                  )}
                </div>
              );
            })}

            {categoriesWithoutBudgets.length > 0 && (
              <div className="pt-4 border-t border-border">
                <div className="text-xs text-muted-foreground mb-2">
                  Categories without budgets ({categoriesWithoutBudgets.length})
                </div>
                <div className="flex flex-wrap gap-2">
                  {categoriesWithoutBudgets.slice(0, 5).map((category) => (
                    <CategoryBudgetForm
                      key={category.category_id}
                      categoryId={category.category_id}
                      month={currentMonth}
                      trigger={
                        <Badge
                          variant="outline"
                          className="cursor-pointer hover:bg-muted"
                        >
                          {category.category_name}
                        </Badge>
                      }
                    />
                  ))}
                  {categoriesWithoutBudgets.length > 5 && (
                    <Badge variant="outline">
                      +{categoriesWithoutBudgets.length - 5} more
                    </Badge>
                  )}
                </div>
              </div>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

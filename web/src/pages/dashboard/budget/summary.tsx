import { useState, useMemo } from 'react';
import { motion } from 'motion/react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { PeriodPresetFilter } from '@/components/budget';
import type { PeriodPresetFilterValue } from '@/components/budget';
import {
  useSavingsRate,
  useIncomeOccurrences,
  useExpenses,
} from '@/hooks/use-budget';
import { TrendingUp, TrendingDown, PiggyBank, Percent } from 'lucide-react';

function computeThisMonthDates(): PeriodPresetFilterValue {
  const today = new Date();
  const start = new Date(today.getFullYear(), today.getMonth(), 1);
  return {
    startDate: start.toISOString().split('T')[0],
    endDate: today.toISOString().split('T')[0],
    preset: 'this-month'
  };
}

export default function BudgetSummaryPage() {
  const [dateRange, setDateRange] = useState<PeriodPresetFilterValue>(computeThisMonthDates);

  const hasDateRange = !!dateRange.startDate && !!dateRange.endDate;

  const { data: savingsData, isLoading: savingsLoading } = useSavingsRate(
    dateRange.startDate || undefined,
    dateRange.endDate || undefined
  );

  const { data: occurrencesData, isLoading: occurrencesLoading } = useIncomeOccurrences(
    dateRange.startDate,
    dateRange.endDate,
    1,
    200,
    hasDateRange
  );

  const { data: expensesData, isLoading: expensesLoading } = useExpenses(
    {
      start_date: dateRange.startDate || undefined,
      end_date: dateRange.endDate || undefined
    },
    1,
    200
  );

  const totalIncome = savingsData?.income ?? 0;
  const totalExpenses = savingsData?.expenses ?? 0;
  const netSavings = savingsData?.savings ?? 0;
  const savingsRate = savingsData?.savings_rate ?? 0;

  const incomeEntries = occurrencesData?.data ?? [];
  const expenseEntries = expensesData?.data ?? [];

  const savingsStatusColor = useMemo(() => {
    if (!savingsData) return 'text-muted-foreground';
    switch (savingsData.status) {
      case 'excellent': return 'text-green-600';
      case 'good': return 'text-green-500';
      case 'fair': return 'text-yellow-500';
      case 'poor': return 'text-orange-500';
      case 'negative': return 'text-red-500';
      default: return 'text-muted-foreground';
    }
  }, [savingsData]);

  return (
    <div className="px-4 md:px-6 py-6">
      {/* Header */}
      <motion.div
        className="mb-6"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <h1 className="text-primary mb-1">BUDGET SUMMARY</h1>
        <p className="text-sm text-muted-foreground">
          Financial overview for the selected period
        </p>
      </motion.div>

      {/* Period Filter */}
      <motion.div
        className="mb-6"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.1 }}
      >
        <PeriodPresetFilter value={dateRange} onChange={setDateRange} />
      </motion.div>

      {!hasDateRange ? (
        <motion.div
          className="text-center py-16 text-muted-foreground"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
        >
          <p>Select a period above to view your financial summary.</p>
        </motion.div>
      ) : (
        <>
          {/* Summary Cards */}
          <motion.div
            className="mb-8 grid gap-4 md:grid-cols-2 lg:grid-cols-4"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.15 }}
          >
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Income</CardTitle>
                <TrendingUp className="h-4 w-4 text-green-500" />
              </CardHeader>
              <CardContent>
                {savingsLoading ? (
                  <div className="h-8 w-32 bg-muted rounded animate-pulse" />
                ) : (
                  <div className="text-2xl font-bold text-green-600">
                    +PHP {totalIncome.toFixed(2)}
                  </div>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Expenses</CardTitle>
                <TrendingDown className="h-4 w-4 text-red-500" />
              </CardHeader>
              <CardContent>
                {savingsLoading ? (
                  <div className="h-8 w-32 bg-muted rounded animate-pulse" />
                ) : (
                  <div className="text-2xl font-bold text-red-600">
                    -PHP {totalExpenses.toFixed(2)}
                  </div>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Net Savings</CardTitle>
                <PiggyBank className="h-4 w-4 text-blue-500" />
              </CardHeader>
              <CardContent>
                {savingsLoading ? (
                  <div className="h-8 w-32 bg-muted rounded animate-pulse" />
                ) : (
                  <div className={`text-2xl font-bold ${netSavings >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                    {netSavings >= 0 ? '+' : ''}PHP {netSavings.toFixed(2)}
                  </div>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Savings Rate</CardTitle>
                <Percent className="h-4 w-4 text-purple-500" />
              </CardHeader>
              <CardContent>
                {savingsLoading ? (
                  <div className="h-8 w-32 bg-muted rounded animate-pulse" />
                ) : (
                  <div className="flex items-baseline gap-2">
                    <span className={`text-2xl font-bold ${savingsStatusColor}`}>
                      {savingsRate.toFixed(1)}%
                    </span>
                    {savingsData?.status && (
                      <Badge
                        variant="secondary"
                        className="text-xs capitalize"
                      >
                        {savingsData.status}
                      </Badge>
                    )}
                  </div>
                )}
              </CardContent>
            </Card>
          </motion.div>

          {/* Itemized Lists */}
          <motion.div
            className="grid gap-6 lg:grid-cols-2"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
          >
            {/* Income List */}
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <CardTitle>Income Entries</CardTitle>
                  <Badge variant="outline" className="text-xs">
                    {incomeEntries.length} entries
                  </Badge>
                </div>
              </CardHeader>
              <CardContent>
                {occurrencesLoading ? (
                  <div className="space-y-3">
                    {[...Array(3)].map((_, i) => (
                      <div key={i} className="flex items-center space-x-4 animate-pulse">
                        <div className="flex-1 space-y-2">
                          <div className="h-4 w-3/4 bg-muted rounded" />
                          <div className="h-3 w-1/2 bg-muted rounded" />
                        </div>
                        <div className="h-5 w-20 bg-muted rounded" />
                      </div>
                    ))}
                  </div>
                ) : incomeEntries.length === 0 ? (
                  <p className="text-center py-8 text-muted-foreground">
                    No income entries for this period
                  </p>
                ) : (
                  <div className="space-y-2 max-h-[400px] overflow-y-auto">
                    {incomeEntries.map((occ) => (
                      <div
                        key={occ.id}
                        className={`flex items-center justify-between p-2.5 rounded-lg border ${
                          occ.is_virtual ? 'border-dashed border-border/70' : 'border-border/50'
                        } ${occ.is_skipped ? 'bg-red-500/5' : ''}`}
                      >
                        <div className="flex-1 min-w-0">
                          <p className={`text-sm font-medium truncate ${occ.is_skipped ? 'line-through text-muted-foreground' : ''}`}>
                            {occ.description || 'Income'}
                            {occ.is_virtual && (
                              <Badge variant="secondary" className="ml-2 text-xs">
                                {occ.recurring_type || 'recurring'}
                              </Badge>
                            )}
                            {occ.is_skipped && (
                              <Badge variant="destructive" className="ml-2 text-xs">
                                Skipped
                              </Badge>
                            )}
                          </p>
                          <p className="text-xs text-muted-foreground">
                            {occ.date}
                          </p>
                        </div>
                        <div className="text-right ml-3">
                          <p className={`text-sm font-semibold ${occ.is_skipped ? 'text-red-600' : 'text-green-600'}`}>
                            {occ.is_skipped ? '' : '+'}PHP {Math.abs(occ.amount).toFixed(2)}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Expense List */}
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <CardTitle>Expense Entries</CardTitle>
                  <Badge variant="outline" className="text-xs">
                    {expenseEntries.length} entries
                  </Badge>
                </div>
              </CardHeader>
              <CardContent>
                {expensesLoading ? (
                  <div className="space-y-3">
                    {[...Array(3)].map((_, i) => (
                      <div key={i} className="flex items-center space-x-4 animate-pulse">
                        <div className="flex-1 space-y-2">
                          <div className="h-4 w-3/4 bg-muted rounded" />
                          <div className="h-3 w-1/2 bg-muted rounded" />
                        </div>
                        <div className="h-5 w-20 bg-muted rounded" />
                      </div>
                    ))}
                  </div>
                ) : expenseEntries.length === 0 ? (
                  <p className="text-center py-8 text-muted-foreground">
                    No expenses for this period
                  </p>
                ) : (
                  <div className="space-y-2 max-h-[400px] overflow-y-auto">
                    {expenseEntries.map((expense) => (
                      <div
                        key={expense.id}
                        className="flex items-center justify-between p-2.5 rounded-lg border border-border/50"
                      >
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2 flex-wrap">
                            <p className="text-sm font-medium truncate">{expense.description}</p>
                            <span className="text-xs px-1.5 py-0.5 rounded-full capitalize">
                              {expense.priority}
                            </span>
                          </div>
                          <p className="text-xs text-muted-foreground">
                            {expense.expense_date}
                          </p>
                        </div>
                        <div className="text-right ml-3">
                          <p className="text-sm font-semibold text-red-600">
                            -{expense.currency} {expense.amount.toFixed(2)}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </motion.div>
        </>
      )}
    </div>
  );
}

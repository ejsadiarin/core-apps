'use client';

import { Suspense, useState, useEffect } from 'react';
import { useRouter, useSearchParams, usePathname } from 'next/navigation';
import { motion } from 'motion/react';
import {
  HealthScoreCard,
  FiftyThirtyTwentyChart,
  WeekdaySpendingChart,
  SpendingTrendCard,
  SavingsRateCard,
  SpendingVelocityCard
} from '@/components/budget';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { ArrowLeft } from 'lucide-react';
import Link from 'next/link';
import { cn } from '@/lib/utils';

function HealthDashboardContent() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const queryStartDate = searchParams.get('start_date');
  const queryEndDate = searchParams.get('end_date');

  const [dateRangeStart, setDateRangeStart] = useState<string>('');
  const [dateRangeEnd, setDateRangeEnd] = useState<string>('');

  useEffect(() => {
    // Set default boundaries on mount if not provided in URL
    if (!queryStartDate && !queryEndDate) {
      const today = new Date().toISOString().split('T')[0];
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setDateRangeStart('2026-01-15');
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setDateRangeEnd(today);
    } else {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setDateRangeStart(queryStartDate || '');
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setDateRangeEnd(queryEndDate || '');
    }
  }, [queryStartDate, queryEndDate]);

  const handleApplyFilter = () => {
    const params = new URLSearchParams(searchParams.toString());
    
    if (dateRangeStart) params.set('start_date', dateRangeStart);
    else params.delete('start_date');

    if (dateRangeEnd) params.set('end_date', dateRangeEnd);
    else params.delete('end_date');

    router.push(`${pathname}?${params.toString()}`);
  };

  const handleClearFilter = () => {
    setDateRangeStart('');
    setDateRangeEnd('');
    router.push(pathname);
  };

  const currentStartDate = queryStartDate || undefined;
  const currentEndDate = queryEndDate || undefined;

  return (
    <div className="px-4 md:px-6 py-6">
      <motion.div
        className="mb-8 flex flex-col md:flex-row md:items-center justify-between gap-4"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <div>
          <div className="flex items-center gap-3 mb-2">
            <Link href="/dashboard/budget">
              <Button variant="ghost" size="sm">
                <ArrowLeft className="h-4 w-4" />
              </Button>
            </Link>
            <h1 className="text-primary">FINANCIAL HEALTH</h1>
          </div>
          <p className="text-sm text-muted-foreground ml-11">
            Analyze your spending patterns and financial wellness
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2 text-sm bg-muted/30 p-2 rounded-lg">
          <span className="text-muted-foreground font-medium mr-1">Date Range:</span>
          <Input
            type="date"
            value={dateRangeStart}
            onChange={(e) => setDateRangeStart(e.target.value)}
            className="w-36 h-8 bg-background"
            placeholder="Start date"
          />
          <span className="text-muted-foreground">to</span>
          <Input
            type="date"
            value={dateRangeEnd}
            onChange={(e) => setDateRangeEnd(e.target.value)}
            className="w-36 h-8 bg-background"
            placeholder="End date"
          />
          <Button
            size="sm"
            onClick={handleApplyFilter}
            className="h-8 ml-1"
          >
            Apply
          </Button>
          {(queryStartDate || queryEndDate) && (
            <Button
              variant="ghost"
              size="sm"
              onClick={handleClearFilter}
              className="h-8 px-2"
            >
              Clear
            </Button>
          )}
        </div>
      </motion.div>

      {/* Financial Insights - Key Metrics */}
      <motion.div
        className="mb-6 grid gap-4 md:grid-cols-3"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.1 }}
      >
        <SavingsRateCard startDate={currentStartDate} endDate={currentEndDate} />
        <SpendingVelocityCard startDate={currentStartDate} endDate={currentEndDate} />
        <HealthScoreCard />
      </motion.div>

      {/* Spending Analysis - Charts */}
      <motion.div
        className="mb-6 grid gap-4 md:grid-cols-2"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.2 }}
      >
        <WeekdaySpendingChart 
          startDate={currentStartDate} 
          endDate={currentEndDate} 
        />
        <SpendingTrendCard />
      </motion.div>
      
       {/* Budget Rule */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.3 }}
      >
        <FiftyThirtyTwentyChart />
      </motion.div>
    </div>
  );
}

export default function HealthDashboardPage() {
  return (
    <Suspense fallback={
      <div className="flex h-[50vh] items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    }>
      <HealthDashboardContent />
    </Suspense>
  );
}

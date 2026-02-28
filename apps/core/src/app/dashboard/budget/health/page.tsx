'use client';

import { useState } from 'react';
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
import { ArrowLeft } from 'lucide-react';
import Link from 'next/link';
import { cn } from '@/lib/utils';

type Period = 'week' | 'month';

function getDatesForPeriod(period: Period): { startDate: string; endDate: string } {
  const now = new Date();
  const year = now.getFullYear();
  const month = now.getMonth();
  const day = now.getDate();
  const currentDay = now.getDay(); // 0 = Sunday

  if (period === 'week') {
    // Get start of week (Sunday)
    const startOfWeek = new Date(now);
    startOfWeek.setDate(day - currentDay);
    // Get end of week (Saturday)
    const endOfWeek = new Date(now);
    endOfWeek.setDate(day + (6 - currentDay));

    // If current date is before start of week (shouldn't happen with above logic but just in case)
    if (now < startOfWeek) {
        startOfWeek.setDate(startOfWeek.getDate() - 7);
        endOfWeek.setDate(endOfWeek.getDate() - 7);
    }

    return {
      startDate: startOfWeek.toISOString().split('T')[0],
      endDate: endOfWeek.toISOString().split('T')[0]
    };
  } else {
    // This month
    const startOfMonth = new Date(year, month, 1);
    const endOfMonth = new Date(year, month + 1, 0);
    return {
      startDate: startOfMonth.toISOString().split('T')[0],
      endDate: endOfMonth.toISOString().split('T')[0]
    };
  }
}

export default function HealthDashboardPage() {
  const [period, setPeriod] = useState<Period>('month');
  
  // Calculate dates based on period
  const { startDate, endDate } = getDatesForPeriod(period);

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

        <div className="flex items-center gap-2 bg-muted/50 p-1 rounded-lg">
          <Button
            variant={period === 'week' ? 'default' : 'ghost'}
            size="sm"
            onClick={() => setPeriod('week')}
            className={cn("h-8 px-3", period === 'week' && "shadow-sm")}
          >
            Week
          </Button>
          <Button
            variant={period === 'month' ? 'default' : 'ghost'}
            size="sm"
            onClick={() => setPeriod('month')}
            className={cn("h-8 px-3", period === 'month' && "shadow-sm")}
          >
            Month
          </Button>
        </div>
      </motion.div>

      {/* Financial Insights - Key Metrics */}
      <motion.div
        className="mb-6 grid gap-4 md:grid-cols-3"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.1 }}
      >
        <SavingsRateCard startDate={startDate} endDate={endDate} />
        <SpendingVelocityCard />
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
          startDate={startDate} 
          endDate={endDate} 
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

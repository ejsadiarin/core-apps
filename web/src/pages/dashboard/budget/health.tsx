import { useState, useEffect, useCallback } from 'react';
import { useNavigate, useSearchParams, useLocation } from 'react-router-dom';
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
import { Link } from 'react-router-dom';

export default function HealthDashboardContent() {
    const navigate = useNavigate();
    const { pathname } = useLocation();
    const [searchParams] = useSearchParams();

    const queryStartDate = searchParams.get('start_date');
    const queryEndDate = searchParams.get('end_date');

    const [dateRangeStart, setDateRangeStart] = useState<string>('');
    const [dateRangeEnd, setDateRangeEnd] = useState<string>('');

    useEffect(() => {
        if (!queryStartDate && !queryEndDate) {
            const today = new Date().toISOString().split('T')[0];
            setDateRangeStart('2026-01-15');
            setDateRangeEnd(today);
        } else {
            setDateRangeStart(queryStartDate || '');
            setDateRangeEnd(queryEndDate || '');
        }
    }, [queryStartDate, queryEndDate]);

    const updateUrlParams = useCallback((start: string, end: string) => {
        const params = new URLSearchParams();
        if (start) params.set('start_date', start);
        if (end) params.set('end_date', end);

        if (params.toString()) {
            navigate(`${pathname}?${params.toString()}`);
        } else {
            navigate(pathname);
        }
    }, [pathname, navigate]);

    const handleStartDateChange = (value: string) => {
        setDateRangeStart(value);
        updateUrlParams(value, dateRangeEnd);
    };

    const handleEndDateChange = (value: string) => {
        setDateRangeEnd(value);
        updateUrlParams(dateRangeStart, value);
    };

    const handleClearFilter = () => {
        setDateRangeStart('');
        setDateRangeEnd('');
        navigate(pathname);
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
                <Link to="/dashboard/budget">
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

        <div className="flex flex-col sm:flex-row sm:flex-wrap items-start sm:items-center gap-2 text-sm bg-muted/30 p-2 rounded-lg">
            <span className="text-muted-foreground font-medium">Date Range:</span>
      <div className="flex items-center gap-2 w-full sm:w-auto">
        <Input
          type="date"
          value={dateRangeStart}
          onChange={(e) => handleStartDateChange(e.target.value)}
          className="w-full sm:w-40 h-9 bg-background pr-8"
          placeholder="Start"
        />
        <span className="text-muted-foreground shrink-0">to</span>
        <Input
          type="date"
          value={dateRangeEnd}
          onChange={(e) => handleEndDateChange(e.target.value)}
          className="w-full sm:w-40 h-9 bg-background pr-8"
          placeholder="End"
        />
      </div>
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
            className="mb-6 grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3"
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
            className="mb-6 grid gap-4 grid-cols-1 lg:grid-cols-2"
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

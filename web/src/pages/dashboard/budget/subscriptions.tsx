import { motion } from 'motion/react';
import {
  RecurringSummaryCard,
  RecurringExpensesList,
  RecurringIncomesList,
  TopMerchantsTable
} from '@/components/budget';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';
import { Link } from 'react-router-dom';

export default function SubscriptionsPage() {
  return (
    <div className="px-4 md:px-6 py-6">
      <motion.div
        className="mb-8 flex items-center justify-between"
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
            <h1 className="text-primary">RECURRING TRACKER</h1>
          </div>
          <p className="text-sm text-muted-foreground ml-11">
            Manage recurring expenses and incomes, skip or cancel occurrences
          </p>
        </div>
      </motion.div>

      {/* summary card */}
      <motion.div
        className="mb-6 grid gap-4 md:grid-cols-3"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.1 }}
      >
        <RecurringSummaryCard />
      </motion.div>

      {/* recurring items tabs */}
      <motion.div
        className="mb-6"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.2 }}
      >
        <Tabs defaultValue="expenses">
          <TabsList>
            <TabsTrigger value="expenses">Expenses</TabsTrigger>
            <TabsTrigger value="incomes">Incomes</TabsTrigger>
          </TabsList>
          <TabsContent value="expenses">
            <RecurringExpensesList />
          </TabsContent>
          <TabsContent value="incomes">
            <RecurringIncomesList />
          </TabsContent>
        </Tabs>
      </motion.div>

      {/* merchants table */}
      <motion.div
        className="mb-6"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.3 }}
      >
        <TopMerchantsTable />
      </motion.div>
    </div>
  );
}

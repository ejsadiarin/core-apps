'use client';

import { motion } from 'motion/react';
import {
  SubscriptionList,
  SubscriptionTotalCard,
  TopMerchantsTable
} from '@/components/budget';
import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';
import Link from 'next/link';

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
            <Link href="/dashboard/budget">
              <Button variant="ghost" size="sm">
                <ArrowLeft className="h-4 w-4" />
              </Button>
            </Link>
            <h1 className="text-primary">SUBSCRIPTIONS & MERCHANTS</h1>
          </div>
          <p className="text-sm text-muted-foreground ml-11">
            Track recurring expenses and top spending merchants
          </p>
        </div>
      </motion.div>

      {/* subscription total card */}
      <motion.div
        className="mb-6 grid gap-4 md:grid-cols-3"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.1 }}
      >
        <SubscriptionTotalCard />
      </motion.div>

      {/* subscriptions list + merchants */}
      <motion.div
        className="mb-6 grid gap-4 lg:grid-cols-2"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.2 }}
      >
        <SubscriptionList />
        <TopMerchantsTable />
      </motion.div>
    </div>
  );
}

"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DollarSign, TrendingUp, Receipt, Calendar } from "lucide-react";
import type { SummaryStats } from "@/types/api";

interface ExpenseStatsProps {
  stats: SummaryStats | undefined;
  isLoading?: boolean;
}

export function ExpenseStats({ stats, isLoading }: ExpenseStatsProps) {
  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {[...Array(4)].map((_, i) => (
          <Card key={i} className="animate-pulse">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <div className="h-4 w-24 bg-muted rounded" />
              <div className="h-4 w-4 bg-muted rounded" />
            </CardHeader>
            <CardContent>
              <div className="h-8 w-32 bg-muted rounded mb-1" />
              <div className="h-3 w-20 bg-muted rounded" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (!stats) {
    return null;
  }

  const statCards = [
    {
      title: "Total Spent",
      value: `₱${stats.total_spent.toFixed(2)}`,
      description: stats.period,
      icon: DollarSign,
      iconColor: "text-green-500",
    },
    {
      title: "Transactions",
      value: stats.transaction_count.toString(),
      description: `${stats.period} period`,
      icon: Receipt,
      iconColor: "text-blue-500",
    },
    {
      title: "Avg per Transaction",
      value: stats.transaction_count > 0
        ? `₱${(stats.total_spent / stats.transaction_count).toFixed(2)}`
        : "₱0.00",
      description: "average amount",
      icon: TrendingUp,
      iconColor: "text-purple-500",
    },
    {
      title: "Period",
      value: stats.period,
      description: "current view",
      icon: Calendar,
      iconColor: "text-orange-500",
    },
  ];

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {statCards.map((stat, index) => {
        const Icon = stat.icon;
        return (
          <Card key={index}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{stat.title}</CardTitle>
              <Icon className={`h-4 w-4 ${stat.iconColor}`} />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <p className="text-xs text-muted-foreground">{stat.description}</p>
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}

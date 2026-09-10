
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DollarSign, TrendingUp, Receipt, Wallet } from "lucide-react";
import { safeFixed } from "@/lib/utils";
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

  const budgetRemaining = stats.remaining ?? 0;
  const transactionCount = stats.transaction_count ?? 0;
  const budgetStatus = stats.remaining > 0 ? "green" : stats.remaining < 0 ? "red" : "neutral";
  const statusColors = {
    green: "text-green-600",
    red: "text-red-500",
    neutral: "text-gray-600",
  };
  // const bgColors = {
  //   green: "bg-green-50 border-green-200",
  //   red: "bg-red-50 border-red-200",
  //   neutral: "",
  // };

  const statCards = [
    {
      title: "Total Spent",
      value: `₱${safeFixed(stats.total_expense)}`,
      description: stats.period,
      icon: DollarSign,
      iconColor: "text-green-500",
    },
    {
      title: "Transactions",
      value: String(transactionCount),
      description: `${stats.period} period`,
      icon: Receipt,
      iconColor: "text-blue-500",
    },
    {
      title: "Avg per Transaction",
      value: transactionCount > 0
        ? `₱${safeFixed(stats.total_expense / transactionCount)}`
        : "₱0.00",
      description: "average amount",
      icon: TrendingUp,
      iconColor: "text-purple-500",
    },
    {
      title: "Budget Remaining",
      value: `₱${safeFixed(budgetRemaining)}`,
      description: budgetStatus === "green" ? "on track" : budgetStatus === "red" ? "over budget" : "break even",
      icon: Wallet,
      iconColor: statusColors[budgetStatus as keyof typeof statusColors],
    },
  ];

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {statCards.map((stat, index) => {
        const Icon = stat.icon;
        const isBudgetCard = stat.title === "Budget Remaining";
        const valueColor = isBudgetCard ? stat.iconColor : "";
        return (
          <Card key={index}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{stat.title}</CardTitle>
              <Icon className={`h-4 w-4 ${stat.iconColor}`} />
            </CardHeader>
            <CardContent>
              <div className={`text-2xl font-bold ${valueColor}`}>{stat.value}</div>
              <p className={`text-xs ${isBudgetCard ? stat.iconColor : "text-muted-foreground"}`}>
                {stat.description}
              </p>
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}

"use client";

import { motion } from "motion/react";
import { SystemStatus } from "@/components/dashboard/system-status";
import { QuickActions } from "@/components/dashboard/quick-actions";
import { ServiceGrid } from "@/components/dashboard/service-grid";
import { useAllServicesStats } from "@/hooks/use-services";

export default function Dashboard() {
    const { data: overallStats } = useAllServicesStats();

    return (
        <div className="container mx-auto px-4 py-8">
            {/* Welcome Banner */}
            <motion.div
                className="mb-8 p-6 rounded-lg border border-primary/30 bg-gradient-to-r from-primary/10 via-accent/10 to-secondary/10 backdrop-blur-sm"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5 }}
            >
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-primary mb-2">DASHBOARD ONLINE</h1>
                        <p className="text-sm text-muted-foreground">
                            All critical systems are operational. Last sync:{" "}
                            <span className="text-accent">just now</span>
                        </p>
                    </div>
                    <div className="hidden md:flex items-center gap-4 text-xs">
                        {overallStats && (
                            <>
                                <div className="flex flex-col items-end">
                                    <span className="text-muted-foreground">Services</span>
                                    <span className="text-card-foreground">
                                        {overallStats.total_services} monitored
                                    </span>
                                </div>
                                <div className="w-px h-8 bg-border" />
                                <div className="flex flex-col items-end">
                                    <span className="text-muted-foreground">Uptime (24h)</span>
                                    <span className="text-card-foreground">
                                        {overallStats.successful_checks > 0
                                            ? `${((overallStats.successful_checks / overallStats.total_checks) * 100).toFixed(1)}%`
                                            : "N/A"}
                                    </span>
                                </div>
                                <div className="w-px h-8 bg-border" />
                                <div className="flex flex-col items-end">
                                    <span className="text-muted-foreground">Avg Response</span>
                                    <span className="text-card-foreground">
                                        {overallStats.avg_response_time
                                            ? `${Math.round(overallStats.avg_response_time)}ms`
                                            : "N/A"}
                                    </span>
                                </div>
                            </>
                        )}
                    </div>
                </div>
            </motion.div>

            {/* Quick Actions */}
            <div className="mb-8">
                <QuickActions />
            </div>

            {/* System Status */}
            <div className="mb-8">
                <SystemStatus />
            </div>

            {/* Services Grid */}
            <ServiceGrid />
        </div>
    );
}

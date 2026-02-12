"use client";

import { DashboardShell } from "@/components/layout/dashboard-shell";
import { useAuth } from "@/contexts/auth-context";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/contexts/query-context";

function ProtectedLayout({ children }: { children: React.ReactNode }) {
    const { isAuthenticated, isLoading } = useAuth();
    const router = useRouter();

    useEffect(() => {
        if (!isLoading && !isAuthenticated) {
            router.replace('/login');
        }
    }, [isAuthenticated, isLoading, router]);

    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-background">
                <div className="flex flex-col items-center gap-3">
                    <div className="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin" />
                    <span className="text-[10px] uppercase tracking-[0.3em] text-muted-foreground">
                        Initializing...
                    </span>
                </div>
            </div>
        )
    }

    if (!isAuthenticated) {
        return null;
    }

    return <DashboardShell>{children}</DashboardShell>
}


export default function DashboardLayout({
    children,
}: {
    children: React.ReactNode
}) {
    return (
        <QueryClientProvider client={queryClient}>
            <ProtectedLayout>
                {children}
            </ProtectedLayout>
        </QueryClientProvider>
    )
}

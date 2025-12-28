"use client";

import { AuthProvider } from "@/contexts/auth-context";
import { MainLayout } from "@/components/layout/main-layout";
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

    if (isLoading || !isAuthenticated) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
            </div>
        )
    }

    return <MainLayout>{children}</MainLayout>
}


export default function DashboardLayout({
    children,
}: {
    children: React.ReactNode
}) {
    return (
        <QueryClientProvider client={queryClient}>
            <AuthProvider>
                <ProtectedLayout>
                    {children}
                </ProtectedLayout>
            </AuthProvider>
        </QueryClientProvider>
    )
}

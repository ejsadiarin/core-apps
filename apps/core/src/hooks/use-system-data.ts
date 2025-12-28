import { useQuery } from "@tanstack/react-query";
import { fetchServicesStats, fetchSystemStatus } from "@/lib/api";
import { SystemStats, ServiceStatus } from "@/types/api";

export function useSystemStatus() {
    return useQuery<SystemStats>({
        queryKey: ['system-stats'],
        queryFn: fetchSystemStatus,
        refetchInterval: 2000
    })
}

export function useServicesStatus() {
    return useQuery<ServiceStatus[]>({
        queryKey: ['services-stats'],
        queryFn: fetchServicesStats,
        refetchInterval: 5000
    })
}

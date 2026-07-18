import { useQuery } from '@tanstack/react-query';
import { fetchServices, fetchSystemStatus } from '@/lib/api';
import { SystemStats, ServiceStatus } from '@/types/api';

export function useSystemStatus() {
  return useQuery<SystemStats>({
    queryKey: ['system-stats'],
    queryFn: fetchSystemStatus,
    refetchInterval: 2000
  });
}

export function useServicesStatus() {
  return useQuery<ServiceStatus[]>({
    queryKey: ['services-status'],
    queryFn: async (): Promise<ServiceStatus[]> => {
      const services = await fetchServices();
      return services
        .filter((s) => s.current_status != null)
        .map((service) => ({
          name: service.name,
          status: service.current_status as 'online' | 'offline' | 'degraded' | 'maintenance' | 'unknown',
          type: service.service_type || 'unknown'
        }));
    },
    refetchInterval: 30000
  });
}

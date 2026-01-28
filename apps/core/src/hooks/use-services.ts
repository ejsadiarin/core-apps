import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  fetchServices,
  fetchService,
  createService,
  updateService,
  deleteService,
  triggerHealthCheck,
  fetchServiceHistory,
  fetchServiceStats,
  fetchAllServicesStats,
} from "@/lib/api";
import type {
  Service,
  CreateServiceRequest,
  UpdateServiceRequest,
  ServiceHealthHistory,
  ServiceStats,
} from "@/types/api";

// Query keys for cache management
export const serviceKeys = {
  all: ["services"] as const,
  lists: () => [...serviceKeys.all, "list"] as const,
  list: () => [...serviceKeys.lists()] as const,
  details: () => [...serviceKeys.all, "detail"] as const,
  detail: (id: string) => [...serviceKeys.details(), id] as const,
  history: (id: string) => [...serviceKeys.all, "history", id] as const,
  stats: (id: string) => [...serviceKeys.all, "stats", id] as const,
  allStats: () => [...serviceKeys.all, "stats", "all"] as const,
};

/**
 * Hook to fetch all services with current health status
 * Auto-refreshes every 30 seconds
 */
export function useServices() {
  return useQuery<Service[]>({
    queryKey: serviceKeys.list(),
    queryFn: fetchServices,
    refetchInterval: 30000, // 30 seconds
    staleTime: 10000, // Consider data stale after 10 seconds
  });
}

/**
 * Hook to fetch a single service by ID
 */
export function useService(id: string | null) {
  return useQuery<Service>({
    queryKey: serviceKeys.detail(id ?? ""),
    queryFn: () => fetchService(id!),
    enabled: !!id, // Only fetch when id is provided
  });
}

/**
 * Hook to fetch service health history
 * Only fetches when enabled (e.g., when modal is open)
 */
export function useServiceHistory(id: string | null, enabled: boolean = true) {
  return useQuery<ServiceHealthHistory[]>({
    queryKey: serviceKeys.history(id ?? ""),
    queryFn: () => fetchServiceHistory(id!),
    enabled: !!id && enabled,
    staleTime: 30000, // History doesn't change that often
  });
}

/**
 * Hook to fetch service statistics (uptime percentages)
 * Only fetches when enabled (e.g., when modal is open)
 */
export function useServiceStats(id: string | null, enabled: boolean = true) {
  return useQuery<ServiceStats>({
    queryKey: serviceKeys.stats(id ?? ""),
    queryFn: () => fetchServiceStats(id!),
    enabled: !!id && enabled,
    staleTime: 60000, // Stats can be cached for longer
  });
}

/**
 * Hook to fetch overall services statistics
 */
export function useAllServicesStats() {
  return useQuery({
    queryKey: serviceKeys.allStats(),
    queryFn: fetchAllServicesStats,
    refetchInterval: 60000, // 1 minute
    staleTime: 30000,
  });
}

/**
 * Hook for creating a new service
 */
export function useCreateService() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateServiceRequest) => createService(data),
    onSuccess: () => {
      // Invalidate services list to trigger refetch
      queryClient.invalidateQueries({ queryKey: serviceKeys.lists() });
    },
  });
}

/**
 * Hook for updating a service
 */
export function useUpdateService() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateServiceRequest }) =>
      updateService(id, data),
    onSuccess: (_, variables) => {
      // Invalidate both the list and the specific service
      queryClient.invalidateQueries({ queryKey: serviceKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: serviceKeys.detail(variables.id),
      });
    },
  });
}

/**
 * Hook for deleting a service
 */
export function useDeleteService() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => deleteService(id),
    onSuccess: () => {
      // Invalidate services list
      queryClient.invalidateQueries({ queryKey: serviceKeys.lists() });
    },
  });
}

/**
 * Hook for triggering a manual health check
 */
export function useTriggerHealthCheck() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => triggerHealthCheck(id),
    onSuccess: (_, id) => {
      // Invalidate the service and its history
      queryClient.invalidateQueries({ queryKey: serviceKeys.lists() });
      queryClient.invalidateQueries({ queryKey: serviceKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: serviceKeys.history(id) });
      queryClient.invalidateQueries({ queryKey: serviceKeys.stats(id) });
    },
  });
}

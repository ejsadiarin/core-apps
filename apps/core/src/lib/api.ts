import type { 
  Service, 
  CreateServiceRequest, 
  UpdateServiceRequest, 
  ServiceHealthHistory, 
  ServiceStats 
} from "@/types/api";

const url = "http://localhost:8080";

export async function fetchSystemStatus() {
  const res = await fetch(url + "/api/system/stats");
  if (!res.ok) {
    throw new Error(`Error fetching ${res.status}`);
  }
  return res.json();
}

export async function fetchServicesStats() {
  const res = await fetch(url + "/api/services");
  if (!res.ok) {
    throw new Error(`Error fetching ${res.status}`);
  }
  return res.json();
}

// Service Monitoring API Functions

export async function fetchServices(): Promise<Service[]> {
  const res = await fetch(url + "/api/services/list");
  if (!res.ok) {
    throw new Error(`Error fetching services: ${res.status}`);
  }
  return res.json();
}

export async function fetchService(id: string): Promise<Service> {
  const res = await fetch(url + `/api/services/${id}`);
  if (!res.ok) {
    throw new Error(`Error fetching service: ${res.status}`);
  }
  return res.json();
}

export async function createService(data: CreateServiceRequest): Promise<Service> {
  const res = await fetch(url + "/api/services", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    throw new Error(`Error creating service: ${res.status}`);
  }
  return res.json();
}

export async function updateService(id: string, data: UpdateServiceRequest): Promise<Service> {
  const res = await fetch(url + `/api/services/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    throw new Error(`Error updating service: ${res.status}`);
  }
  return res.json();
}

export async function deleteService(id: string): Promise<void> {
  const res = await fetch(url + `/api/services/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) {
    throw new Error(`Error deleting service: ${res.status}`);
  }
}

export async function triggerHealthCheck(id: string): Promise<void> {
  const res = await fetch(url + `/api/services/${id}/check`, {
    method: "POST",
  });
  if (!res.ok) {
    throw new Error(`Error triggering health check: ${res.status}`);
  }
}

export async function fetchServiceHistory(id: string): Promise<ServiceHealthHistory[]> {
  const res = await fetch(url + `/api/services/${id}/history`);
  if (!res.ok) {
    throw new Error(`Error fetching service history: ${res.status}`);
  }
  return res.json();
}

export async function fetchServiceStats(id: string): Promise<ServiceStats> {
  const res = await fetch(url + `/api/services/${id}/stats`);
  if (!res.ok) {
    throw new Error(`Error fetching service stats: ${res.status}`);
  }
  return res.json();
}

export async function fetchAllServicesStats(): Promise<{
  total_services: number;
  total_checks: number;
  successful_checks: number;
  overall_uptime: number;
  avg_response_time: number;
}> {
  const res = await fetch(url + "/api/services/stats/all");
  if (!res.ok) {
    throw new Error(`Error fetching services stats: ${res.status}`);
  }
  return res.json();
}

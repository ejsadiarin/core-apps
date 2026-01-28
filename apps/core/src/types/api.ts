export interface SystemStats {
    cpu: number;
    memory: number;
    storage: number;
    temperature: number;
    uptime: string;
    network: {
        up: string;
        down: string;
    };
}

export interface ServiceStatus {
    name: string;
    status: "online" | "offline" | "maintenance";
    type: string;
}

export interface Service {
    id: string;
    name: string;
    url: string;
    icon?: string;
    description?: string;
    service_type?: string;
    health_check_interval: number;
    health_check_method: string;
    expected_status_codes: number[];
    timeout: number;
    created_at: string;
    updated_at: string;
    is_active: boolean;
    current_status?:
        | "online"
        | "offline"
        | "degraded"
        | "maintenance"
        | "unknown";
    last_check?: string;
    response_time?: number;
}

export interface ServiceHealthHistory {
    id: string;
    service_id: string;
    status: string;
    response_time?: number;
    status_code?: number;
    error_message?: string;
    checked_at: string;
}

export interface ServiceStats {
    service_id: string;
    uptime_24h: number;
    uptime_7d: number;
    uptime_30d: number;
    avg_response_time: number;
    total_checks: number;
    successful_checks: number;
}

export interface CreateServiceRequest {
    name: string;
    url: string;
    icon?: string;
    description?: string;
    service_type?: string;
    health_check_interval?: number;
    health_check_method?: string;
    expected_status_codes?: number[];
    timeout?: number;
}

export interface UpdateServiceRequest {
    name?: string;
    url?: string;
    icon?: string;
    description?: string;
    service_type?: string;
    health_check_interval?: number;
    health_check_method?: string;
    expected_status_codes?: number[];
    timeout?: number;
    is_active?: boolean;
}


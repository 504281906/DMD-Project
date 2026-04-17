const BASE_URL = 'http://localhost:8080';

interface RequestOptions {
  method?: string;
  body?: unknown;
  headers?: Record<string, string>;
}

async function request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
  const { method = 'GET', body, headers = {} } = options;

  const config: RequestInit = {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...headers,
    },
  };

  if (body) {
    config.body = JSON.stringify(body);
  }

  const response = await fetch(`${BASE_URL}${endpoint}`, config);

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return response.json();
}

// Dashboard APIs
export interface DashboardStats {
  totalDevices: number;
  totalApps: number;
  onlineDevices: number;
  pendingAlerts: number;
}

export interface PlatformDistribution {
  platform: string;
  count: number;
}

export interface RecentDevice {
  id: string;
  name: string;
  platform: 'iOS' | 'Android';
  status: 'online' | 'offline';
  registeredAt: string;
}

export async function getDashboardStats(): Promise<DashboardStats> {
  return request<DashboardStats>('/api/dashboard/stats');
}

export async function getPlatformDistribution(): Promise<PlatformDistribution[]> {
  return request<PlatformDistribution[]>('/api/dashboard/platforms');
}

export async function getRecentDevices(): Promise<RecentDevice[]> {
  return request<RecentDevice[]>('/api/dashboard/recent-devices');
}

// Device APIs
export interface Device {
  id: string;
  name: string;
  platform: 'iOS' | 'Android';
  status: 'online' | 'offline';
  registeredAt: string;
  model?: string;
  osVersion?: string;
}

export interface DeviceListResponse {
  items: Device[];
  total: number;
  page: number;
  pageSize: number;
}

export async function getDevices(params: { page?: number; pageSize?: number; search?: string }): Promise<DeviceListResponse> {
  const query = new URLSearchParams();
  if (params.page) query.set('page', String(params.page));
  if (params.pageSize) query.set('pageSize', String(params.pageSize));
  if (params.search) query.set('search', params.search);
  return request<DeviceListResponse>(`/api/devices?${query}`);
}

export async function lockDevice(id: string): Promise<void> {
  await request<void>(`/api/devices/${id}/lock`, { method: 'POST' });
}

export async function wipeDevice(id: string): Promise<void> {
  await request<void>(`/api/devices/${id}/wipe`, { method: 'POST' });
}

export async function getDeviceDetail(id: string): Promise<Device> {
  return request<Device>(`/api/devices/${id}`);
}

// App APIs
export interface App {
  id: string;
  name: string;
  bundleId: string;
  platform: 'iOS' | 'Android';
  version: string;
  status: 'active' | 'inactive';
  installCount: number;
  iconUrl?: string;
}

export interface AppListResponse {
  items: App[];
  total: number;
  page: number;
  pageSize: number;
}

export async function getApps(params: { page?: number; pageSize?: number; platform?: string }): Promise<AppListResponse> {
  const query = new URLSearchParams();
  if (params.page) query.set('page', String(params.page));
  if (params.pageSize) query.set('pageSize', String(params.pageSize));
  if (params.platform) query.set('platform', params.platform);
  return request<AppListResponse>(`/api/apps?${query}`);
}

export async function uploadApp(data: FormData): Promise<App> {
  const response = await fetch(`${BASE_URL}/api/apps/upload`, {
    method: 'POST',
    body: data,
  });
  if (!response.ok) {
    throw new Error(`Upload failed: ${response.status}`);
  }
  return response.json();
}

export async function updateApp(id: string, data: Partial<App>): Promise<App> {
  return request<App>(`/api/apps/${id}`, { method: 'PUT', body: data });
}

export async function deleteApp(id: string): Promise<void> {
  await request<void>(`/api/apps/${id}`, { method: 'DELETE' });
}

// Policy APIs
export interface Policy {
  id: string;
  name: string;
  type: 'password' | 'encryption' | 'compliance';
  status: 'active' | 'inactive';
  deviceCount: number;
  description?: string;
}

export interface PolicyListResponse {
  items: Policy[];
  total: number;
  page: number;
  pageSize: number;
}

export async function getPolicies(params?: { page?: number; pageSize?: number }): Promise<PolicyListResponse> {
  const query = new URLSearchParams();
  if (params?.page) query.set('page', String(params.page));
  if (params?.pageSize) query.set('pageSize', String(params.pageSize));
  return request<PolicyListResponse>(`/api/policies?${query}`);
}

export async function createPolicy(data: Omit<Policy, 'id'>): Promise<Policy> {
  return request<Policy>('/api/policies', { method: 'POST', body: data });
}

export async function updatePolicy(id: string, data: Partial<Policy>): Promise<Policy> {
  return request<Policy>(`/api/policies/${id}`, { method: 'PUT', body: data });
}

export async function deletePolicy(id: string): Promise<void> {
  await request<void>(`/api/policies/${id}`, { method: 'DELETE' });
}

export async function assignPolicyToDevices(policyId: string, deviceIds: string[]): Promise<void> {
  await request<void>(`/api/policies/${policyId}/assign`, { method: 'POST', body: { deviceIds } });
}

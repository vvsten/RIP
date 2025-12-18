/**
 * Типы для работы с заявками (логистическими запросами)
 */

import { TransportService } from './TransportService';

export interface LogisticRequestService {
  id: number;
  logistic_request_id: number;
  transport_service_id: number;
  quantity: number;
  comment?: string;
  sort_order: number;
  service?: TransportService;
}

export interface LogisticRequest {
  id: number;
  session_id?: string;
  is_draft: boolean;
  from_city?: string;
  to_city?: string;
  weight: number;
  length: number;
  width: number;
  height: number;
  services: LogisticRequestService[];
  total_cost: number;
  total_days: number;
  status: 'draft' | 'formed' | 'completed' | 'rejected' | 'deleted';
  creator_id: number;
  moderator_id?: number;
  created_at: string;
  formed_at?: string;
  completed_at?: string;
  updated_at: string;
}

export interface OrdersListResponse {
  status: string;
  orders: LogisticRequest[];
}

export interface OrderResponse {
  status: string;
  order: LogisticRequest;
}


import { useEffect, useState, useRef, useCallback } from 'react';
import type { SystemMetrics, SiteSettings } from '../types';

interface NetworkSpeed {
  rx_sec: number;
  tx_sec: number;
}

export interface ServerConfig {
  id: string;
  name: string;
  type: 'real' | 'local';
  location?: string;
  provider?: string;
  tag?: string;
  version?: string;
}

export interface ServerState {
  config: ServerConfig;
  metrics: SystemMetrics | null;
  speed: NetworkSpeed;
  isConnected: boolean;
  error: string | null;
}

// Loading state for initial data fetch
export type LoadingState = 'idle' | 'loading' | 'ready' | 'error';

// Message from dashboard WebSocket
interface DashboardMessage {
  type: string;
  servers: ServerMetricsUpdate[];
  site_settings?: SiteSettings;
}

interface ServerMetricsUpdate {
  server_id: string;
  server_name: string;
  location: string;
  provider: string;
  tag?: string;
  version?: string;
  online: boolean;
  metrics: SystemMetrics | null;
}

const defaultSiteSettings: SiteSettings = {
  site_name: 'vStats Dashboard',
  site_description: 'Real-time Server Monitoring',
  social_links: []
};

export function useServerManager() {
  const [servers, setServers] = useState<ServerState[]>([]);
  const [siteSettings, setSiteSettings] = useState<SiteSettings>(defaultSiteSettings);
  const [loadingState, setLoadingState] = useState<LoadingState>('loading');
  const [isInitialLoad, setIsInitialLoad] = useState(true);
  
  const lastMetricsMap = useRef<Map<string, { metrics: SystemMetrics, time: number }>>(new Map());
  const wsRef = useRef<WebSocket | null>(null);
  const initialDataReceived = useRef(false);

  // Connect to dashboard WebSocket - all data (local + remote) comes through here
  useEffect(() => {
    const connect = () => {
      try {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;

        const ws = new WebSocket(wsUrl);
        wsRef.current = ws;

        ws.onopen = () => {
          console.log('[Dashboard] WebSocket connected');
        };

        ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data) as DashboardMessage;
            
            // Update site settings if provided
            if (data.site_settings) {
              setSiteSettings(data.site_settings);
            }
            
            if (data.type === 'metrics' && data.servers) {
              // Mark initial data as received
              if (!initialDataReceived.current) {
                initialDataReceived.current = true;
                setLoadingState('ready');
                setIsInitialLoad(false);
              }
              const now = Date.now();
              
              // All servers (local + remote) come through WebSocket
              const allServers: ServerState[] = data.servers.map(serverUpdate => {
                const lastData = lastMetricsMap.current.get(serverUpdate.server_id);
                
                let newSpeed = { rx_sec: 0, tx_sec: 0 };

                // Use pre-calculated speeds from server if available
                if (serverUpdate.metrics?.network.rx_speed !== undefined && 
                    serverUpdate.metrics?.network.tx_speed !== undefined) {
                  newSpeed = {
                    rx_sec: serverUpdate.metrics.network.rx_speed,
                    tx_sec: serverUpdate.metrics.network.tx_speed
                  };
                } else if (lastData && serverUpdate.metrics) {
                  // Fallback: calculate from totals difference
                  const timeDiff = (now - lastData.time) / 1000;
                  if (timeDiff > 0) {
                    const rxDiff = serverUpdate.metrics.network.total_rx - lastData.metrics.network.total_rx;
                    const txDiff = serverUpdate.metrics.network.total_tx - lastData.metrics.network.total_tx;
                    newSpeed = {
                      rx_sec: Math.max(0, rxDiff / timeDiff),
                      tx_sec: Math.max(0, txDiff / timeDiff)
                    };
                  }
                }

                if (serverUpdate.metrics) {
                  lastMetricsMap.current.set(serverUpdate.server_id, { 
                    metrics: serverUpdate.metrics, 
                    time: now 
                  });
                }

                // Determine if this is the local server
                const isLocal = serverUpdate.server_id === 'local';

                return {
                  config: {
                    id: serverUpdate.server_id,
                    name: serverUpdate.server_name,
                    type: isLocal ? 'local' as const : 'real' as const,
                    location: serverUpdate.location,
                    provider: serverUpdate.provider,
                    tag: serverUpdate.tag,
                    version: serverUpdate.version || serverUpdate.metrics?.version,
                  },
                  metrics: serverUpdate.metrics,
                  speed: newSpeed,
                  isConnected: serverUpdate.online,
                  error: null
                };
              });

              setServers(allServers);
            }
          } catch (e) {
            console.error('[Dashboard] Parse error', e);
          }
        };

        ws.onclose = () => {
          console.log('[Dashboard] WebSocket disconnected, reconnecting...');
          wsRef.current = null;
          // Don't reset loading state on reconnect if we already have data
          if (!initialDataReceived.current) {
            setLoadingState('loading');
          }
          setTimeout(connect, 3000);
        };

        ws.onerror = (err) => {
          console.error('[Dashboard] WebSocket error', err);
        };
      } catch (e) {
        console.error('[Dashboard] Connection error', e);
        setTimeout(connect, 3000);
      }
    };

    connect();

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

  // Get a server by ID (with cached lookup)
  const getServerById = useCallback((id: string): ServerState | undefined => {
    return servers.find(s => s.config.id === id);
  }, [servers]);

  return { servers, siteSettings, loadingState, isInitialLoad, getServerById };
}

export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
}

export function formatSpeed(bytesPerSec: number): string {
  return `${formatBytes(bytesPerSec)}/s`;
}

export function formatUptime(seconds: number): string {
  const days = Math.floor(seconds / 86400);
  if (days > 0) return `${days}d`;
  const hours = Math.floor(seconds / 3600);
  if (hours > 0) return `${hours}h`;
  return `${Math.floor(seconds / 60)}m`;
}

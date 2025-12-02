use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

use crate::config::SiteSettings;

// ============================================================================
// Auth Types
// ============================================================================

#[derive(Debug, Serialize, Deserialize)]
pub struct Claims {
    pub sub: String,
    pub exp: i64,
}

#[derive(Debug, Deserialize)]
pub struct LoginRequest {
    pub password: String,
}

#[derive(Debug, Serialize)]
pub struct LoginResponse {
    pub token: String,
    pub expires_at: DateTime<Utc>,
}

#[derive(Debug, Deserialize)]
pub struct ChangePasswordRequest {
    pub current_password: String,
    pub new_password: String,
}

#[derive(Debug, Deserialize)]
pub struct AddServerRequest {
    pub name: String,
    #[serde(default)]
    pub url: String,
    #[serde(default)]
    pub location: String,
    #[serde(default)]
    pub provider: String,
    #[serde(default)]
    pub tag: String,
}

#[derive(Debug, Deserialize)]
pub struct UpdateServerRequest {
    pub name: Option<String>,
    #[serde(default)]
    pub location: Option<String>,
    #[serde(default)]
    pub provider: Option<String>,
    #[serde(default)]
    pub tag: Option<String>,
}

// ============================================================================
// Agent Registration Types
// ============================================================================

#[derive(Debug, Deserialize)]
pub struct AgentRegisterRequest {
    pub name: String,
    #[serde(default)]
    pub location: String,
    #[serde(default)]
    pub provider: String,
}

#[derive(Debug, Serialize)]
pub struct AgentRegisterResponse {
    pub id: String,
    pub token: String,
}

// ============================================================================
// Historical Data Types
// ============================================================================

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HistoryPoint {
    pub timestamp: String,
    pub cpu: f32,
    pub memory: f32,
    pub disk: f32,
    pub net_rx: i64,
    pub net_tx: i64,
    #[serde(default)]
    pub ping_ms: Option<f64>,
}

#[derive(Debug, Deserialize)]
pub struct HistoryQuery {
    #[serde(default = "default_range")]
    pub range: String,
}

fn default_range() -> String {
    "24h".to_string()
}

#[derive(Debug, Serialize)]
pub struct HistoryResponse {
    pub server_id: String,
    pub range: String,
    pub data: Vec<HistoryPoint>,
    #[serde(default)]
    pub ping_targets: Vec<PingHistoryTarget>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct PingHistoryTarget {
    pub name: String,
    pub host: String,
    pub data: Vec<PingHistoryPoint>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PingHistoryPoint {
    pub timestamp: String,
    pub latency_ms: Option<f64>,
    pub status: String,
}

// ============================================================================
// System Metrics Types
// ============================================================================

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SystemMetrics {
    pub timestamp: DateTime<Utc>,
    pub hostname: String,
    pub os: OsInfo,
    pub cpu: CpuMetrics,
    pub memory: MemoryMetrics,
    pub disks: Vec<DiskMetrics>,
    pub network: NetworkMetrics,
    pub uptime: u64,
    pub load_average: LoadAverage,
    #[serde(default)]
    pub ping: Option<PingMetrics>,
    #[serde(default)]
    pub version: Option<String>,
    #[serde(default)]
    pub ip_addresses: Option<Vec<String>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OsInfo {
    pub name: String,
    pub version: String,
    pub kernel: String,
    pub arch: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CpuMetrics {
    pub brand: String,
    pub cores: usize,
    pub usage: f32,
    pub frequency: u64,
    pub per_core: Vec<f32>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MemoryMetrics {
    pub total: u64,
    pub used: u64,
    pub available: u64,
    pub swap_total: u64,
    pub swap_used: u64,
    pub usage_percent: f32,
    #[serde(default)]
    pub modules: Vec<MemoryModule>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MemoryModule {
    #[serde(default)]
    pub slot: Option<String>,
    pub size: u64,
    #[serde(default)]
    pub mem_type: Option<String>,
    #[serde(default)]
    pub speed: Option<u32>,
    #[serde(default)]
    pub manufacturer: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiskMetrics {
    pub name: String,
    #[serde(default)]
    pub model: Option<String>,
    #[serde(default)]
    pub serial: Option<String>,
    pub total: u64,
    #[serde(default)]
    pub disk_type: Option<String>,
    #[serde(default)]
    pub mount_points: Vec<String>,
    pub usage_percent: f32,
    pub used: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NetworkMetrics {
    pub interfaces: Vec<NetworkInterface>,
    pub total_rx: u64,
    pub total_tx: u64,
    #[serde(default)]
    pub rx_speed: u64,
    #[serde(default)]
    pub tx_speed: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NetworkInterface {
    pub name: String,
    #[serde(default)]
    pub mac: Option<String>,
    #[serde(default)]
    pub speed: Option<u32>,
    pub rx_bytes: u64,
    pub tx_bytes: u64,
    pub rx_packets: u64,
    pub tx_packets: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LoadAverage {
    pub one: f64,
    pub five: f64,
    pub fifteen: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PingMetrics {
    pub targets: Vec<PingTarget>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PingTarget {
    pub name: String,
    pub host: String,
    pub latency_ms: Option<f64>,
    pub packet_loss: f64,
    pub status: String,
}

// ============================================================================
// Dashboard/WebSocket Message Types
// ============================================================================

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AgentMetricsData {
    pub server_id: String,
    pub metrics: SystemMetrics,
    pub last_updated: DateTime<Utc>,
}

/// Full state message - sent on initial connection
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DashboardMessage {
    #[serde(rename = "type")]
    pub msg_type: String,
    pub servers: Vec<ServerMetricsUpdate>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub site_settings: Option<SiteSettings>,
}

/// Full server info with all fields
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ServerMetricsUpdate {
    pub server_id: String,
    pub server_name: String,
    pub location: String,
    pub provider: String,
    #[serde(default)]
    pub tag: String,
    #[serde(default)]
    pub version: String,
    #[serde(default)]
    pub ip: String,
    pub online: bool,
    pub metrics: Option<SystemMetrics>,
}

/// Compact delta message - only contains changed metrics
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DeltaMessage {
    #[serde(rename = "type")]
    pub msg_type: String,  // "delta"
    /// Timestamp for this update
    pub ts: i64,
    /// Compact server updates (only changed data)
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub d: Vec<CompactServerUpdate>,
}

/// Compact server update with minimal fields
#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct CompactServerUpdate {
    /// Server ID
    pub id: String,
    /// Online status (only if changed)
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub on: Option<bool>,
    /// Compact metrics (only changed values)
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub m: Option<CompactMetrics>,
}

/// Compact metrics with only essential real-time values
#[derive(Debug, Clone, Serialize, Deserialize, Default, PartialEq)]
pub struct CompactMetrics {
    /// CPU usage percentage (0-100)
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub c: Option<u8>,
    /// Memory usage percentage (0-100)
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub m: Option<u8>,
    /// Disk usage percentage (0-100)
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub d: Option<u8>,
    /// Network RX speed (bytes/s)
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub rx: Option<u64>,
    /// Network TX speed (bytes/s)
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub tx: Option<u64>,
    /// Uptime in seconds
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub up: Option<u64>,
}

impl CompactMetrics {
    pub fn from_system_metrics(m: &SystemMetrics) -> Self {
        Self {
            c: Some(m.cpu.usage.round() as u8),
            m: Some(m.memory.usage_percent.round() as u8),
            d: m.disks.first().map(|d| d.usage_percent.round() as u8),
            rx: Some(m.network.rx_speed),
            tx: Some(m.network.tx_speed),
            up: Some(m.uptime),
        }
    }
    
    /// Check if metrics have changed significantly
    pub fn has_changed(&self, other: &CompactMetrics) -> bool {
        self.c != other.c || self.m != other.m || self.d != other.d ||
        self.rx != other.rx || self.tx != other.tx
    }
    
    /// Get only the changed fields compared to previous
    pub fn diff(&self, prev: &CompactMetrics) -> CompactMetrics {
        CompactMetrics {
            c: if self.c != prev.c { self.c } else { None },
            m: if self.m != prev.m { self.m } else { None },
            d: if self.d != prev.d { self.d } else { None },
            rx: if self.rx != prev.rx { self.rx } else { None },
            tx: if self.tx != prev.tx { self.tx } else { None },
            up: None, // Don't send uptime in delta, it always changes
        }
    }
    
    pub fn is_empty(&self) -> bool {
        self.c.is_none() && self.m.is_none() && self.d.is_none() &&
        self.rx.is_none() && self.tx.is_none() && self.up.is_none()
    }
}

/// Tracks the last sent state for delta calculation
#[derive(Debug, Clone, Default)]
pub struct LastSentState {
    pub servers: std::collections::HashMap<String, (bool, CompactMetrics)>,
}

#[derive(Debug, Deserialize)]
pub struct AgentMessage {
    #[serde(rename = "type")]
    pub msg_type: String,
    pub server_id: Option<String>,
    pub token: Option<String>,
    #[serde(default)]
    pub version: Option<String>,
    pub metrics: Option<SystemMetrics>,
}

// ============================================================================
// Installation Script Types
// ============================================================================

#[derive(Debug, Serialize)]
pub struct InstallCommand {
    pub command: String,
    pub script_url: String,
}

// ============================================================================
// Update Agent Types
// ============================================================================

#[derive(Debug, Deserialize)]
pub struct UpdateAgentRequest {
    #[serde(default)]
    pub download_url: Option<String>,
}

#[derive(Debug, Serialize)]
pub struct UpdateAgentResponse {
    pub success: bool,
    pub message: String,
}

/// Command message sent to agent
#[derive(Debug, Serialize)]
pub struct AgentCommand {
    #[serde(rename = "type")]
    pub cmd_type: String,
    pub command: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub download_url: Option<String>,
}


use chrono::Utc;
use sysinfo::{Disks, Networks, System};

use crate::types::{
    CpuMetrics, DiskMetrics, LoadAverage, MemoryMetrics, NetworkInterface, NetworkMetrics, OsInfo,
    SystemMetrics,
};

pub fn collect_metrics(sys: &mut System, disks: &Disks, networks: &Networks) -> SystemMetrics {
    sys.refresh_all();

    let cpu_usage: Vec<f32> = sys.cpus().iter().map(|cpu| cpu.cpu_usage()).collect();
    let avg_cpu = if cpu_usage.is_empty() {
        0.0
    } else {
        cpu_usage.iter().sum::<f32>() / cpu_usage.len() as f32
    };

    let cpu_brand = sys
        .cpus()
        .first()
        .map(|c| c.brand().to_string())
        .unwrap_or_else(|| "Unknown".to_string());

    let cpu_freq = sys.cpus().first().map(|c| c.frequency()).unwrap_or(0);

    // Simplified disk collection for local server (show main disk)
    let disk_metrics: Vec<DiskMetrics> = disks
        .list()
        .iter()
        .filter(|disk| {
            let mount = disk.mount_point().to_string_lossy();
            mount == "/" || mount.starts_with("C:")
        })
        .map(|disk| {
            let total = disk.total_space();
            let available = disk.available_space();
            let used = total.saturating_sub(available);
            let usage = if total > 0 {
                (used as f32 / total as f32) * 100.0
            } else {
                0.0
            };

            DiskMetrics {
                name: disk.name().to_string_lossy().to_string(),
                model: None,
                serial: None,
                total,
                disk_type: Some("SSD".to_string()),
                mount_points: vec![disk.mount_point().to_string_lossy().to_string()],
                usage_percent: usage,
                used,
            }
        })
        .collect();

    let mut total_rx = 0u64;
    let mut total_tx = 0u64;
    // Filter to physical interfaces only
    let interfaces: Vec<NetworkInterface> = networks
        .list()
        .iter()
        .filter(|(name, _)| {
            let n = name.to_lowercase();
            n != "lo" && n != "lo0" && 
            !n.starts_with("veth") && !n.starts_with("docker") && 
            !n.starts_with("br-") && !n.starts_with("virbr") &&
            !n.starts_with("utun") && !n.starts_with("awdl") && !n.starts_with("llw")
        })
        .map(|(name, data)| {
            total_rx += data.total_received();
            total_tx += data.total_transmitted();
            NetworkInterface {
                name: name.clone(),
                mac: None,
                speed: None,
                rx_bytes: data.total_received(),
                tx_bytes: data.total_transmitted(),
                rx_packets: data.total_packets_received(),
                tx_packets: data.total_packets_transmitted(),
            }
        })
        .collect();

    let total_mem = sys.total_memory();
    let used_mem = sys.used_memory();
    let mem_usage = if total_mem > 0 {
        (used_mem as f32 / total_mem as f32) * 100.0
    } else {
        0.0
    };

    let load_avg = System::load_average();

    SystemMetrics {
        timestamp: Utc::now(),
        hostname: System::host_name().unwrap_or_else(|| "Unknown".to_string()),
        os: OsInfo {
            name: System::name().unwrap_or_else(|| "Unknown".to_string()),
            version: System::os_version().unwrap_or_else(|| "Unknown".to_string()),
            kernel: System::kernel_version().unwrap_or_else(|| "Unknown".to_string()),
            arch: std::env::consts::ARCH.to_string(),
        },
        cpu: CpuMetrics {
            brand: cpu_brand,
            cores: sys.cpus().len(),
            usage: avg_cpu,
            frequency: cpu_freq,
            per_core: cpu_usage,
        },
        memory: MemoryMetrics {
            total: total_mem,
            used: used_mem,
            available: sys.available_memory(),
            swap_total: sys.total_swap(),
            swap_used: sys.used_swap(),
            usage_percent: mem_usage,
            modules: Vec::new(), // Detailed module info not collected for local server
        },
        disks: disk_metrics,
        network: NetworkMetrics {
            interfaces,
            total_rx,
            total_tx,
            rx_speed: 0,
            tx_speed: 0,
        },
        uptime: System::uptime(),
        load_average: LoadAverage {
            one: load_avg.one,
            five: load_avg.five,
            fifteen: load_avg.fifteen,
        },
        ping: None,
        version: None,
        ip_addresses: None,
    }
}


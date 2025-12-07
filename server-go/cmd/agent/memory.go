package main

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/mem"
)

// collectMemoryModules collects detailed memory module information
func collectMemoryModules() []MemoryModule {
	var modules []MemoryModule

	switch runtime.GOOS {
	case "linux":
		// Use dmidecode (requires root)
		cmd := exec.Command("dmidecode", "-t", "memory")
		output, err := cmd.Output()
		if err == nil {
			scanner := bufio.NewScanner(strings.NewReader(string(output)))
			var currentModule *MemoryModule
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if strings.HasPrefix(line, "Memory Device") {
					if currentModule != nil && currentModule.Size > 0 {
						modules = append(modules, *currentModule)
					}
					currentModule = &MemoryModule{}
				} else if currentModule != nil {
					if strings.HasPrefix(line, "Size:") {
						val := strings.TrimSpace(strings.TrimPrefix(line, "Size:"))
						if val != "No Module Installed" {
							parts := strings.Fields(val)
							if len(parts) >= 2 {
								if size, err := strconv.ParseUint(parts[0], 10, 64); err == nil {
									switch strings.ToUpper(parts[1]) {
									case "GB":
										currentModule.Size = size * 1024 * 1024 * 1024
									case "MB":
										currentModule.Size = size * 1024 * 1024
									case "KB":
										currentModule.Size = size * 1024
									default:
										currentModule.Size = size
									}
								}
							}
						}
					} else if strings.HasPrefix(line, "Type:") {
						val := strings.TrimSpace(strings.TrimPrefix(line, "Type:"))
						if val != "Unknown" && val != "" {
							currentModule.MemType = val
						}
					} else if strings.HasPrefix(line, "Speed:") {
						val := strings.TrimSpace(strings.TrimPrefix(line, "Speed:"))
						parts := strings.Fields(val)
						if len(parts) > 0 {
							if speed, err := strconv.ParseUint(parts[0], 10, 32); err == nil {
								currentModule.Speed = uint32(speed)
							}
						}
					} else if strings.HasPrefix(line, "Locator:") {
						val := strings.TrimSpace(strings.TrimPrefix(line, "Locator:"))
						if val != "" {
							currentModule.Slot = val
						}
					} else if strings.HasPrefix(line, "Manufacturer:") {
						val := strings.TrimSpace(strings.TrimPrefix(line, "Manufacturer:"))
						if val != "Unknown" && val != "" && val != "Not Specified" {
							currentModule.Manufacturer = val
						}
					}
				}
			}
			if currentModule != nil && currentModule.Size > 0 {
				modules = append(modules, *currentModule)
			}
		}
	case "darwin":
		// Use system_profiler
		cmd := exec.Command("system_profiler", "SPMemoryDataType", "-json")
		output, err := cmd.Output()
		if err == nil {
			var data map[string]interface{}
			if json.Unmarshal(output, &data) == nil {
				if memoryData, ok := data["SPMemoryDataType"].([]interface{}); ok {
					for _, item := range memoryData {
						if itemMap, ok := item.(map[string]interface{}); ok {
							if items, ok := itemMap["_items"].([]interface{}); ok {
								for _, moduleItem := range items {
									if module, ok := moduleItem.(map[string]interface{}); ok {
										sizeStr, _ := module["dimm_size"].(string)
										var size uint64
										if strings.Contains(sizeStr, "GB") {
											sizeStr = strings.ReplaceAll(sizeStr, " GB", "")
											sizeStr = strings.TrimSpace(sizeStr)
											if s, err := strconv.ParseUint(sizeStr, 10, 64); err == nil {
												size = s * 1024 * 1024 * 1024
											}
										}
										if size > 0 {
											memModule := MemoryModule{
												Size: size,
											}
											if name, ok := module["_name"].(string); ok {
												memModule.Slot = name
											}
											if dimmType, ok := module["dimm_type"].(string); ok {
												memModule.MemType = dimmType
											}
											if dimmSpeed, ok := module["dimm_speed"].(string); ok {
												parts := strings.Fields(dimmSpeed)
												if len(parts) > 0 {
													if s, err := strconv.ParseUint(parts[0], 10, 32); err == nil {
														memModule.Speed = uint32(s)
													}
												}
											}
											if manufacturer, ok := module["dimm_manufacturer"].(string); ok {
												memModule.Manufacturer = manufacturer
											}
											modules = append(modules, memModule)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	case "windows":
		// Use WMIC
		cmd := exec.Command("wmic", "memorychip", "get", "Capacity,Speed,MemoryType,Manufacturer,DeviceLocator", "/format:csv")
		output, err := cmd.Output()
		if err == nil {
			scanner := bufio.NewScanner(strings.NewReader(string(output)))
			firstLine := true
			for scanner.Scan() {
				if firstLine {
					firstLine = false
					continue
				}
				line := scanner.Text()
				parts := strings.Split(line, ",")
				if len(parts) >= 5 {
					if size, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64); err == nil && size > 0 {
						memModule := MemoryModule{
							Size: size,
						}
						if slot := strings.TrimSpace(parts[2]); slot != "" {
							memModule.Slot = slot
						}
						if memTypeCode, err := strconv.ParseUint(strings.TrimSpace(parts[3]), 10, 32); err == nil {
							switch memTypeCode {
							case 20:
								memModule.MemType = "DDR"
							case 21:
								memModule.MemType = "DDR2"
							case 24:
								memModule.MemType = "DDR3"
							case 26:
								memModule.MemType = "DDR4"
							case 34:
								memModule.MemType = "DDR5"
							}
						}
						if speed, err := strconv.ParseUint(strings.TrimSpace(parts[4]), 10, 32); err == nil {
							memModule.Speed = uint32(speed)
						}
						if manufacturer := strings.TrimSpace(parts[5]); manufacturer != "" {
							memModule.Manufacturer = manufacturer
						}
						modules = append(modules, memModule)
					}
				}
			}
		}
	}

	return modules
}

// SwapInfo contains swap memory information
type SwapInfo struct {
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	Free         uint64  `json:"free"`
	UsagePercent float32 `json:"usage_percent"`
}

// collectSwapInfo collects swap memory information
func collectSwapInfo() SwapInfo {
	swapInfo, err := mem.SwapMemory()
	if err != nil {
		return SwapInfo{}
	}

	return SwapInfo{
		Total:        swapInfo.Total,
		Used:         swapInfo.Used,
		Free:         swapInfo.Free,
		UsagePercent: float32(swapInfo.UsedPercent),
	}
}

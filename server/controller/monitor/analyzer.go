/*
 * Copyright (c) 2022 Yunshan Networks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package monitor

import (
	"sort"
	"time"

	"github.com/deckarep/golang-set"

	"github.com/deepflowys/deepflow/server/controller/common"
	"github.com/deepflowys/deepflow/server/controller/db/mysql"
	"github.com/deepflowys/deepflow/server/controller/monitor/config"
)

type AnalyzerCheck struct {
	cfg                   config.MonitorConfig
	ch                    chan string
	normalAnalyzerDict    map[string]*dfHostCheck
	exceptionAnalyzerDict map[string]*dfHostCheck
}

func NewAnalyzerCheck(cfg config.MonitorConfig) *AnalyzerCheck {
	return &AnalyzerCheck{
		cfg:                   cfg,
		ch:                    make(chan string, cfg.HealthCheckHandleChannelLen),
		normalAnalyzerDict:    make(map[string]*dfHostCheck),
		exceptionAnalyzerDict: make(map[string]*dfHostCheck),
	}
}

func (c *AnalyzerCheck) Start() {
	go func() {
		for range time.Tick(time.Duration(c.cfg.HealthCheckInterval) * time.Second) {
			// 数据节点健康检查
			c.healthCheck()
			// 检查没有分配数据节点的采集器，并进行分配
			c.vtapAnalyzerCheck()
		}
	}()

	// 根据ch信息，针对部分采集器分配/重新分配数据节点
	go func() {
		for {
			excludeIPs := <-c.ch
			c.vtapAnalyzerAlloc(excludeIPs)
		}
	}()
}

func (c *AnalyzerCheck) healthCheck() {
	var analyzers []mysql.Analyzer
	var exceptionIPs []string

	log.Info("analyzer health check start")

	mysql.Db.Not("state = ?", common.HOST_STATE_MAINTENANCE).Find(&analyzers)
	for _, analyzer := range analyzers {
		checkIP := analyzer.IP
		if analyzer.NATIPEnabled != 0 {
			checkIP = analyzer.NATIP
		}

		// 检查逻辑同控制器
		active := isActive(common.HEALTH_CHECK_URL, checkIP, c.cfg.HealthCheckPort)
		if analyzer.State == common.HOST_STATE_COMPLETE {
			if active {
				if _, ok := c.normalAnalyzerDict[analyzer.IP]; ok {
					delete(c.normalAnalyzerDict, analyzer.IP)
				}
				if _, ok := c.exceptionAnalyzerDict[analyzer.IP]; ok {
					delete(c.exceptionAnalyzerDict, analyzer.IP)
				}
			} else {
				if _, ok := c.exceptionAnalyzerDict[analyzer.IP]; ok {
					if c.exceptionAnalyzerDict[analyzer.IP].duration() >= int64(3*common.HEALTH_CHECK_INTERVAL.Seconds()) {
						delete(c.exceptionAnalyzerDict, analyzer.IP)
						mysql.Db.Model(&analyzer).Update("state", common.HOST_STATE_EXCEPTION)
						exceptionIPs = append(exceptionIPs, analyzer.IP)
						log.Infof("set analyzer (%s) state to exception", analyzer.IP)
						// 根据exceptionIP，重新分配对应采集器的数据节点
						c.TriggerReallocAnalyzer(analyzer.IP)
					}
				} else {
					c.exceptionAnalyzerDict[analyzer.IP] = newDFHostCheck()
				}
			}
		} else {
			if active {
				if _, ok := c.normalAnalyzerDict[analyzer.IP]; ok {
					if c.normalAnalyzerDict[analyzer.IP].duration() >= int64(3*common.HEALTH_CHECK_INTERVAL.Seconds()) {
						delete(c.normalAnalyzerDict, analyzer.IP)
						mysql.Db.Model(&analyzer).Update("state", common.HOST_STATE_COMPLETE)
						log.Infof("set analyzer (%s) state to normal", analyzer.IP)
					}
				} else {
					c.normalAnalyzerDict[analyzer.IP] = newDFHostCheck()
				}
			} else {
				if _, ok := c.normalAnalyzerDict[analyzer.IP]; ok {
					delete(c.normalAnalyzerDict, analyzer.IP)
				}
				if _, ok := c.exceptionAnalyzerDict[analyzer.IP]; ok {
					delete(c.exceptionAnalyzerDict, analyzer.IP)
				}
			}
		}
	}
	log.Info("analyzer health check end")
}

func (c *AnalyzerCheck) TriggerReallocAnalyzer(analyzerIP string) {
	c.ch <- analyzerIP
}

func (c *AnalyzerCheck) vtapAnalyzerCheck() {
	var vtaps []mysql.VTap
	var noAnalyzerVtapCount int64

	log.Info("vtap analyzer check start")

	mysql.Db.Find(&vtaps)
	for _, vtap := range vtaps {
		if vtap.AnalyzerIP == "" {
			noAnalyzerVtapCount += 1
		} else if vtap.Exceptions&common.VTAP_EXCEPTION_ALLOC_ANALYZER_FAILED != 0 {
			// 检查是否存在已分配数据节点，但异常未清除的采集器
			exceptions := vtap.Exceptions ^ common.VTAP_EXCEPTION_ALLOC_ANALYZER_FAILED
			mysql.Db.Model(vtap).Update("exceptions", exceptions)
		}
	}
	// 如果存在没有数据节点的采集器，触发数据节点重新分配
	if noAnalyzerVtapCount > 0 {
		c.TriggerReallocAnalyzer("")
	}
	log.Info("vtap analyzer check end")
}

func (c *AnalyzerCheck) vtapAnalyzerAlloc(excludeIP string) {
	var vtaps []mysql.VTap
	var analyzers []mysql.Analyzer
	var azs []mysql.AZ
	var azAnalyzerConns []mysql.AZAnalyzerConnection

	log.Info("vtap analyzer alloc start")

	mysql.Db.Find(&vtaps)
	mysql.Db.Where("state = ?", common.HOST_STATE_COMPLETE).Find(&analyzers)

	// 获取待分配采集器对应的可用区信息
	// 获取数据节点当前已分配的采集器个数
	azToNoAnalyzerVTaps := make(map[string][]*mysql.VTap)
	analyzerIPToUsedVTapNum := make(map[string]int)
	azLcuuids := mapset.NewSet()
	for i, vtap := range vtaps {
		if vtap.AnalyzerIP != "" && vtap.AnalyzerIP != excludeIP {
			analyzerIPToUsedVTapNum[vtap.AnalyzerIP] += 1
			continue
		}
		azToNoAnalyzerVTaps[vtap.AZ] = append(azToNoAnalyzerVTaps[vtap.AZ], &vtaps[i])
		azLcuuids.Add(vtap.AZ)
	}
	// 获取数据节点的剩余采集器个数
	analyzerIPToAvailableVTapNum := make(map[string]int)
	for _, analyzer := range analyzers {
		analyzerIPToAvailableVTapNum[analyzer.IP] -= analyzer.VTapMax
		if usedVTapNum, ok := analyzerIPToUsedVTapNum[analyzer.IP]; ok {
			analyzerIPToAvailableVTapNum[analyzer.IP] = analyzer.VTapMax - usedVTapNum
		}
	}

	// 根据可用区查询region信息
	mysql.Db.Where("lcuuid IN (?)", azLcuuids.ToSlice()).Find(&azs)
	regionToAZLcuuids := make(map[string][]string)
	regionLcuuids := mapset.NewSet()
	for _, az := range azs {
		regionToAZLcuuids[az.Region] = append(regionToAZLcuuids[az.Region], az.Lcuuid)
		regionLcuuids.Add(az.Region)
	}

	// 获取可用区中的数据节点IP
	mysql.Db.Where("region IN (?)", regionLcuuids.ToSlice()).Find(&azAnalyzerConns)
	azToAnalyzerIPs := make(map[string][]string)
	for _, conn := range azAnalyzerConns {
		if conn.AZ == "ALL" {
			if azLcuuids, ok := regionToAZLcuuids[conn.Region]; ok {
				for _, azLcuuid := range azLcuuids {
					azToAnalyzerIPs[azLcuuid] = append(azToAnalyzerIPs[azLcuuid], conn.AnalyzerIP)
				}
			}
		} else {
			azToAnalyzerIPs[conn.AZ] = append(azToAnalyzerIPs[conn.AZ], conn.AnalyzerIP)
		}
	}

	// 遍历待分配采集器，分配数据节点IP
	for az, noAnalyzerVtaps := range azToNoAnalyzerVTaps {
		// 获取可分配的数据节点列表
		analyzerAvailableVTapNum := []common.KVPair{}
		if analyzerIPs, ok := azToAnalyzerIPs[az]; ok {
			for _, analyzerIP := range analyzerIPs {
				if availableVTapNum, ok := analyzerIPToAvailableVTapNum[analyzerIP]; ok {
					analyzerAvailableVTapNum = append(
						analyzerAvailableVTapNum,
						common.KVPair{Key: analyzerIP, Value: availableVTapNum},
					)
				}
			}
		}

		for _, vtap := range noAnalyzerVtaps {
			// 分配数据节点失败，更新异常错误码
			if len(analyzerAvailableVTapNum) == 0 {
				log.Warningf("no available analyzer for vtap (%s)", vtap.Name)
				exceptions := vtap.Exceptions | common.VTAP_EXCEPTION_ALLOC_ANALYZER_FAILED
				mysql.Db.Model(vtap).Update("exceptions", exceptions)
				continue
			}
			sort.Slice(analyzerAvailableVTapNum, func(i, j int) bool {
				return analyzerAvailableVTapNum[i].Value > analyzerAvailableVTapNum[j].Value
			})
			analyzerAvailableVTapNum[0].Value -= 1
			analyzerIPToAvailableVTapNum[analyzerAvailableVTapNum[0].Key] -= 1

			// 分配数据节点成功，更新数据节点IP + 清空数据节点分配失败的错误码
			log.Infof("alloc analyzer (%s) for vtap (%s)", analyzerAvailableVTapNum[0].Key, vtap.Name)
			mysql.Db.Model(vtap).Update("analyzer_ip", analyzerAvailableVTapNum[0].Key)
			if vtap.Exceptions&common.VTAP_EXCEPTION_ALLOC_ANALYZER_FAILED != 0 {
				exceptions := vtap.Exceptions ^ common.VTAP_EXCEPTION_ALLOC_ANALYZER_FAILED
				mysql.Db.Model(vtap).Update("exceptions", exceptions)
			}
		}
	}
	log.Info("vtap analyzer alloc end")
}

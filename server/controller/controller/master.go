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

package controller

import (
	"context"
	"os"
	"time"

	"github.com/deepflowys/deepflow/server/controller/config"
	"github.com/deepflowys/deepflow/server/controller/db/mysql/migrator"
	"github.com/deepflowys/deepflow/server/controller/election"
	"github.com/deepflowys/deepflow/server/controller/monitor"
	"github.com/deepflowys/deepflow/server/controller/monitor/license"
	"github.com/deepflowys/deepflow/server/controller/recorder"
	recorderdb "github.com/deepflowys/deepflow/server/controller/recorder/db"
	"github.com/deepflowys/deepflow/server/controller/service"
	"github.com/deepflowys/deepflow/server/controller/tagrecorder"
)

func IsMasterRegion(cfg *config.ControllerConfig) bool {
	if cfg.TrisolarisCfg.NodeType == "master" {
		return true
	}
	return false
}

// try to check until success
func IsMasterController(cfg *config.ControllerConfig) bool {
	if IsMasterRegion(cfg) {
		for range time.Tick(time.Second * 5) {
			isMasterController, err := election.IsMasterController()
			if err == nil {
				if isMasterController {
					return true
				} else {
					return false
				}
			} else {
				log.Errorf("check whether I am master controller failed: %s", err.Error())
			}
		}
	}
	return false
}

// migrate db by master region master controller
func migrateMySQL(cfg *config.ControllerConfig) {
	ok := migrator.MigrateMySQL(cfg.MySqlCfg)
	if !ok {
		log.Error("migrate mysql failed")
		time.Sleep(time.Second)
		os.Exit(0)
	}
}

func checkAndStartMasterFunctions(
	cfg *config.ControllerConfig, ctx context.Context, tr *tagrecorder.TagRecorder,
	controllerCheck *monitor.ControllerCheck, analyzerCheck *monitor.AnalyzerCheck,
) {

	// 定时检查当前是否为master controller
	// 仅master controller才启动以下goroutine
	// - tagrecorder
	// - 控制器和数据节点检查
	// - license分配和检查
	// - resource id manager
	// - clean deleted resources

	// 从区域控制器无需判断是否为master controller
	if !IsMasterRegion(cfg) {
		return
	}

	vtapCheck := monitor.NewVTapCheck(cfg.MonitorCfg, ctx)
	vtapLicenseAllocation := license.NewVTapLicenseAllocation(cfg.MonitorCfg, ctx)
	softDeletedResourceCleaner := recorder.NewSoftDeletedResourceCleaner(&cfg.ManagerCfg.TaskCfg.RecorderCfg, ctx)
	domainChecker := service.NewDomainCheck(ctx)

	masterController := ""
	thisIsMasterController := false
	for range time.Tick(time.Minute) {
		newThisIsMasterController, newMasterController, err := election.IsMasterControllerAndReturnIP()
		if err != nil {
			continue
		}
		if masterController != newMasterController {
			if newThisIsMasterController {
				thisIsMasterController = true
				log.Infof("I am the master controller now, previous master controller is %s", masterController)

				migrateMySQL(cfg)

				if _, enabled := os.LookupEnv("FEATURE_FLAG_ALLOCATE_ID"); enabled {
					// 启动资源ID管理器
					err := recorderdb.IDMNG.Start()
					if err != nil {
						log.Error("resource id mananger start failed")
						time.Sleep(time.Second)
						os.Exit(0)
					}
				}

				// 启动tagrecorder
				tr.Start()

				// 控制器检查
				controllerCheck.Start()

				// 数据节点检查
				analyzerCheck.Start()

				// vtap check
				vtapCheck.Start()

				// license分配和检查
				vtapLicenseAllocation.Start()

				// 启动软删除数据清理
				softDeletedResourceCleaner.Start()

				if _, enabled := os.LookupEnv("FEATURE_FLAG_CHECK_DOMAIN_CONTROLLER"); enabled {
					// 自动切换domain控制器
					domainChecker.Start()
				}
			} else if thisIsMasterController {
				thisIsMasterController = false
				log.Infof("I am not the master controller anymore, new master controller is %s", newMasterController)

				// stop tagrecorder
				tr.Stop()

				// stop controller check
				controllerCheck.Stop()

				// stop analyzer check
				analyzerCheck.Stop()

				// stop vtap check
				vtapCheck.Stop()

				// stop vtap license allocation and check
				vtapLicenseAllocation.Stop()

				softDeletedResourceCleaner.Stop()

				if _, enabled := os.LookupEnv("FEATURE_FLAG_CHECK_DOMAIN_CONTROLLER"); enabled {
					domainChecker.Stop()
				}

				if _, enabled := os.LookupEnv("FEATURE_FLAG_ALLOCATE_ID"); enabled {
					recorderdb.IDMNG.Stop()
				}
			} else {
				log.Infof(
					"current master controller is %s, previous master controller is %s",
					newMasterController, masterController,
				)
			}
		}
		masterController = newMasterController
	}
}

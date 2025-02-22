// Copyright 2021 Chaos Mesh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chaos-mesh/chaos-mesh/pkg/log"
	"github.com/chaos-mesh/chaos-mesh/pkg/time"
	"github.com/chaos-mesh/chaos-mesh/pkg/time/utils"
	"github.com/chaos-mesh/chaos-mesh/pkg/version"
)

var (
	pid           int
	secDelta      int64
	nsecDelta     int64
	printVersion  bool
	clockIdsSlice string
)

func initFlag() {
	flag.IntVar(&pid, "pid", 0, "pid of target program")
	flag.Int64Var(&secDelta, "sec_delta", 0, "delta time of sec field")
	flag.Int64Var(&nsecDelta, "nsec_delta", 0, "delta time of nsec field")
	flag.StringVar(&clockIdsSlice, "clk_ids", "CLOCK_REALTIME", "all affected clock ids split with \",\"")
	flag.BoolVar(&printVersion, "version", false, "print version information and exit")

	flag.Parse()
}

func main() {
	initFlag()

	version.PrintVersionInfo("watchmaker")

	if printVersion {
		os.Exit(0)
	}

	logger, err := log.NewDefaultZapLogger()
	if err != nil {
		panic(fmt.Sprintf("error while creating zap logger: %v", err))
	}

	clkIds := strings.Split(clockIdsSlice, ",")
	mask, err := utils.EncodeClkIds(clkIds)
	if err != nil {
		logger.Error(err, "error while converting clock ids to mask")
		os.Exit(1)
	}
	logger.Info("get clock ids mask", "mask", mask)

	err = time.ModifyTime(pid, secDelta, nsecDelta, mask, logger.WithName("time"))

	if err != nil {
		logger.Error(err, "error while modifying time", "pid", pid, "secDelta", secDelta, "nsecDelta", nsecDelta, "mask", mask)
	}
}

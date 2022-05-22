package cluster

import (
	"fmt"

	"github.com/ysicing/ergo/common"
	qcexec "github.com/ysicing/ergo/pkg/util/exec"
	"github.com/ysicing/ergo/pkg/util/log"
)

func (p *Cluster) SystemInit() (err error) {
	initShell := fmt.Sprintf("%s/hack/scripts/system-init.sh", common.GetDefaultDataDir())
	log.Flog.Debugf("gen init shell: %v", initShell)
	if err := qcexec.RunCmd("/bin/bash", initShell); err != nil {
		return err
	}
	return nil
}

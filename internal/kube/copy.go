// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package kube

import (
	"context"
	"fmt"
	"io"
	"os"
)

const (
	defaultReadFromByteCmd = "tail -c+%d %s"
	defaultMaxTries        = 5
)

// CopyFromPod is to copy srcFile in a given pod to local destFile with defaultMaxTries.
func (c *Client) CopyFromPod(ctx context.Context, namespace, pod, container string, srcFile, destFile string) error {
	pipe := newPipe(&CopyOptions{
		MaxTries: defaultMaxTries,
		ReadFunc: readFromPod(ctx, c, namespace, pod, container, srcFile),
	})

	outFile, err := os.OpenFile(destFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, pipe); err != nil {
		return err
	}
	return nil
}

func readFromPod(ctx context.Context, client *Client, namespace, pod, container, srcFile string) ReadFunc {
	return func(offset uint64, writer io.Writer) error {
		command := []string{"sh", "-c", fmt.Sprintf(defaultReadFromByteCmd, offset, srcFile)}
		return client.execInPodWithWriters(ctx, ExecParameters{
			Namespace: namespace,
			Pod:       pod,
			Container: container,
			Command:   command,
		}, writer, writer)
	}
}

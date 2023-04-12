# Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.

FROM ysicing/alpine

COPY dist/ergo_linux_amd64 /usr/local/bin/ergo

COPY hack/docker/entrypoint.sh /entrypoint.sh

RUN chmod +x /usr/local/bin/ergo /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD ["ergo", "-h"]
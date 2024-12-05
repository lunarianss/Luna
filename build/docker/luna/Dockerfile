# Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


FROM BASE_IMAGE
LABEL maintainer="cyan0908@163.com"

WORKDIR /opt/Luna

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone && \
      mkdir -p /var/log/Luna && \
      mkdir -p /etc/Luna

COPY luna /opt/Luna/bin/

COPY luna.yaml /etc/Luna/

ENTRYPOINT ["/opt/Luna/bin/luna"]
CMD ["-c", "/etc/Luna/luna.yaml"]

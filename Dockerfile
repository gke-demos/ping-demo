#Copyright 2025 Google LLC.
#SPDX-License-Identifier: Apache-2.0

FROM golang AS builder
WORKDIR /src
ADD . /src
RUN CGO_ENABLED=0 go build -o /ping-demo

FROM gcr.io/distroless/static
COPY --from=builder /ping-demo ping-demo
CMD ["/ping-demo"]
FROM bz-findata-base AS builder

# Add timezone data
RUN apk add --no-cache tzdata

# Set the target architecture
ENV GOARCH=amd64

# Copy the entire project
COPY . .

# Build the liquidator app
RUN go build -o liquidator ./cmd/liquidator

# Create final image
FROM alpine:latest
COPY --from=builder /app/liquidator /usr/local/bin/

# Add timezone data to the final image
RUN apk add --no-cache tzdata

CMD ["liquidator"]
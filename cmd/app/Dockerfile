# Stage 1
# Use the base image
FROM bz-findata-base AS builder

# Set the target architecture
ENV GOARCH=amd64

# Copy the entire project
COPY . .

# Build the main app
RUN go build -o main ./cmd/app

# Create final image
FROM alpine:latest
COPY --from=builder /app/main /usr/local/bin/
CMD ["main"]
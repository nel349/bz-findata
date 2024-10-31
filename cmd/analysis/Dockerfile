# Use the base image
FROM bz-findata-base AS builder

# Add timezone data
RUN apk add --no-cache tzdata

# Set the target architecture
ENV GOARCH=amd64

# Copy the entire project
COPY . .

# Build the analysis app
RUN go build -o analysis ./cmd/analysis

# Create final image
FROM alpine:latest
COPY --from=builder /app/analysis /usr/local/bin/

# Add timezone data to the final image
RUN apk add --no-cache tzdata

CMD ["analysis"]
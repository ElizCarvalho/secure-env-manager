FROM alpine:latest

# Create non-root user for security
RUN adduser -D envkeeper

# Set working directory
WORKDIR /app

# Copy pre-built binary from GoReleaser
COPY secure-env .

# Set ownership to non-root user
RUN chown -R envkeeper:envkeeper /app

# Switch to non-root user
USER envkeeper

# Configure container startup command
ENTRYPOINT ["/app/secure-env"] 
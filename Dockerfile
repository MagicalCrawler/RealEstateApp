# Base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the project files into the container
COPY . .

# Install Go modules and dependencies
RUN go mod tidy \
    && go get -u github.com/playwright-community/playwright-go

# Install Playwright browsers and OS dependencies
RUN go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps

# Expose the application port (if needed)
EXPOSE 8080

# Command to run the application (update this based on your project)
CMD ["go", "run", "main.go"]

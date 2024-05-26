FROM alpine:latest

# Create a directory in the container for the application
RUN mkdir /app
WORKDIR /app

# Copy the binary into the container
COPY ./build/adminServiceApp /app

# Copy the config.yaml file into the container
COPY ./configs/config.yaml /app/configs/config.yaml

# Copy the config.yaml file into the container
COPY ./web/templates /app/web/templates

# Expose the port on which the application will run
EXPOSE 9002

# Command to run the application
CMD ["./adminServiceApp"]

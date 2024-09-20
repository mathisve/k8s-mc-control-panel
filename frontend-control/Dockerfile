# Step 1: Build the React app
FROM node:18-alpine AS build

# Set working directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json to install dependencies
COPY package.json package-lock.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application files
COPY . .

# Build the React app for production
RUN npm run build

# Step 2: Serve the build using NGINX
FROM nginx:alpine

# Copy the React app build to NGINX's default static directory
COPY --from=build /app/build /usr/share/nginx/html

# Copy custom NGINX config file (optional)
COPY ./nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80 to be accessible externally
EXPOSE 80

# Start NGINX server
CMD ["nginx", "-g", "daemon off;"]
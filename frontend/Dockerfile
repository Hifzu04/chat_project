FROM node:20-alpine AS builder

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy application files
COPY . .

# Pass in the API URL for build time
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

# Build the Vite application
RUN npm run build

# Use Nginx to serve the built files
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose the frontend port
EXPOSE 80

# Run Nginx
CMD ["nginx", "-g", "daemon off;"]


#Stage 1 (Builder): Uses a Node.js image to install NPM packages and run npm run build (Vite's build command). 
#This generates static HTML, CSS, and JS files into a dist/ folder.
#Stage 2 (Production): Uses Nginx (a fast web server) to serve those static files to users.

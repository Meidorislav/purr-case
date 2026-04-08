FROM node:20-alpine AS base
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .

FROM base AS dev
EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host"]

FROM base AS build
ARG VITE_XSOLLA_CLIENT_ID
ARG VITE_XSOLLA_LOGIN_ID
ARG VITE_XSOLLA_PROJECT_ID
ARG VITE_XSOLLA_RETURN_URL
ENV VITE_XSOLLA_CLIENT_ID=$VITE_XSOLLA_CLIENT_ID
ENV VITE_XSOLLA_LOGIN_ID=$VITE_XSOLLA_LOGIN_ID
ENV VITE_XSOLLA_PROJECT_ID=$VITE_XSOLLA_PROJECT_ID
ENV VITE_XSOLLA_RETURN_URL=$VITE_XSOLLA_RETURN_URL
RUN npm run build

FROM nginx:alpine AS production
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

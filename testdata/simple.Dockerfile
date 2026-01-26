FROM node:18-alpine
WORKDIR /app
ENV NODE_ENV=production
EXPOSE 3000
CMD ["npm", "start"]

FROM nginx:1.27.5-alpine3.21

# 将构建好的 React 应用静态文件复制到 Nginx 的服务目录
COPY dist /usr/share/nginx/html

# 将自定义 Nginx 配置文件复制到容器中
COPY nginx.conf /etc/nginx/nginx.conf

# 暴露 80 端口
EXPOSE 80

# 启动 Nginx
CMD ["nginx", "-g", "daemon off;"]
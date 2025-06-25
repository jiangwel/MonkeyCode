# baizhiyun-app-airport

一个基于React 19和TypeScript的前端应用，使用Vite构建，Material UI作为UI框架。

## 技术栈

- React 19
- TypeScript
- Vite
- Material UI (@mui) + Emotion
- 状态管理: ahooks
- 路由: react-router-dom v7
- 表单: react-hook-form
- 图表: echarts
- 动画: lottie-react
- API请求: axios
- 日期处理: dayjs
- Markdown处理: react-markdown

## 主要功能

1. **仪表盘**: 应用主界面，展示关键指标和数据概览
2. **聊天功能**: 用户交流界面
3. **完成/结果**: 展示操作结果或完成状态
4. **模型管理**: 模型相关操作和展示
5. **用户管理**: 用户信息管理
6. **管理员功能**: 系统管理后台
7. **认证和登录**: 用户认证流程
8. **邀请系统**: 用户邀请功能

## 项目结构

```
src/
├── api/          # API请求相关
├── assets/       # 静态资源
├── components/   # 公共组件
├── layouts/      # 布局组件
├── pages/        # 页面组件
│   ├── admin/    # 管理员相关
│   ├── auth/     # 认证相关
│   ├── chat/     # 聊天功能
│   ├── completion/ # 完成/结果页面
│   ├── dashboard/  # 仪表盘
│   ├── invite/    # 邀请相关
│   ├── login/     # 登录页面
│   ├── model/     # 模型相关
│   └── user/      # 用户相关
├── router.tsx    # 路由配置
├── theme.ts      # 主题配置
├── utils/        # 工具函数
└── main.tsx      # 应用入口
```

## 开发环境配置

1. 确保已安装Node.js (>=18.0.0) 和 pnpm
2. 克隆项目仓库
3. 安装依赖: `pnpm install`
4. 启动开发服务器: `pnpm dev`
5. 打开浏览器访问: `http://localhost:5173`

## 常用命令

- `pnpm dev`: 启动开发服务器
- `pnpm build`: 生产环境构建
- `pnpm preview`: 预览生产构建
- `pnpm lint`: 运行代码检查

export default [
  {
    path: '/user',
    layout: false,
    routes: [
      { path: '/user/login', layout: false, name: '登录', component: './user/login' },
      { path: '/user', redirect: '/user/login' },
      {
        name: '注册结果',
        icon: 'smile',
        path: '/user/register-result',
        component: './user/register-result',
      },
      { name: '注册', icon: 'smile', path: '/user/register', component: './user/register' },
      { component: '404', path: '/user/*' },
    ],
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    icon: 'dashboard',
    routes: [
      { path: '/dashboard', redirect: '/dashboard/analysis' },
      {
        name: '分析页',
        icon: 'smile',
        path: '/dashboard/analysis',
        component: './dashboard/analysis',
      },
      {
        name: '监控页',
        icon: 'smile',
        path: '/dashboard/monitor',
        component: './dashboard/monitor',
      },
      {
        name: '工作台',
        icon: 'smile',
        path: '/dashboard/workplace',
        component: './dashboard/workplace',
      },
    ],
  },
  {
    path: '/form',
    icon: 'form',
    name: '表单页',
    routes: [
      { path: '/form', redirect: '/form/basic-form' },
      { name: '基础表单', icon: 'smile', path: '/form/basic-form', component: './form/basic-form' },
      { name: '分步表单', icon: 'smile', path: '/form/step-form', component: './form/step-form' },
      {
        name: '高级表单',
        icon: 'smile',
        path: '/form/advanced-form',
        component: './form/advanced-form',
      },
    ],
  },
  {
    path: '/list',
    icon: 'table',
    name: '列表页',
    routes: [
      {
        path: '/list/search',
        name: '搜索列表',
        component: './list/search',
        routes: [
          { path: '/list/search', redirect: '/list/search/articles' },
          {
            name: '搜索列表（文章）',
            icon: 'smile',
            path: '/list/search/articles',
            component: './list/search/articles',
          },
          {
            name: '搜索列表（项目）',
            icon: 'smile',
            path: '/list/search/projects',
            component: './list/search/projects',
          },
          {
            name: '搜索列表（应用）',
            icon: 'smile',
            path: '/list/search/applications',
            component: './list/search/applications',
          },
        ],
      },
      { path: '/list', redirect: '/list/table-list' },
      { name: '查询表格', icon: 'smile', path: '/list/table-list', component: './table-list' },
      { name: '标准列表', icon: 'smile', path: '/list/basic-list', component: './list/basic-list' },
      { name: '卡片列表', icon: 'smile', path: '/list/card-list', component: './list/card-list' },
    ],
  },
  {
    path: '/profile',
    name: '详情页',
    icon: 'profile',
    routes: [
      { path: '/profile', redirect: '/profile/basic' },
      { name: '基础详情页', icon: 'smile', path: '/profile/basic', component: './profile/basic' },
      {
        name: '高级详情页',
        icon: 'smile',
        path: '/profile/advanced',
        component: './profile/advanced',
      },
    ],
  },
  {
    name: '结果页',
    icon: 'CheckCircleOutlined',
    path: '/result',
    routes: [
      { path: '/result', redirect: '/result/success' },
      { name: '成功页', icon: 'smile', path: '/result/success', component: './result/success' },
      { name: '失败页', icon: 'smile', path: '/result/fail', component: './result/fail' },
    ],
  },
  {
    name: '异常页',
    icon: 'warning',
    path: '/exception',
    routes: [
      { path: '/exception', redirect: '/exception/403' },
      { name: '403', icon: 'smile', path: '/exception/403', component: './exception/403' },
      { name: '404', icon: 'smile', path: '/exception/404', component: './exception/404' },
      { name: '500', icon: 'smile', path: '/exception/500', component: './exception/500' },
    ],
  },
  {
    name: '个人页',
    icon: 'user',
    path: '/account',
    routes: [
      { path: '/account', redirect: '/account/center' },
      { name: '个人中心', icon: 'smile', path: '/account/center', component: './account/center' },
      {
        name: '个人设置',
        icon: 'smile',
        path: '/account/settings',
        component: './account/settings',
      },
    ],
  },
  { path: '/', redirect: '/dashboard/analysis' },
  { component: '404', path: '/*' },
];

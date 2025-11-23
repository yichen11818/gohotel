export default [
  {
    path: '/user',
    layout: false,
    routes: [{ name: '登录', path: '/user/login', component: './user/login' }],
  },
  { path: '/welcome', name: '欢迎', icon: 'smile', component: './Welcome' },
  {
    path: '/admin',
    name: '管理页',
    icon: 'crown',
    access: 'canAdmin',
    routes: [
      { path: '/admin', redirect: '/admin/sub-page' },
      { path: '/admin/sub-page', name: '二级管理页', component: './Admin' },
    ],
  },
  { name: '用户管理', icon: 'table', path: '/user-manage', component: './user-manage' },
  {
    path: '/hotel',
    name: '酒店管理',
    icon: 'dashboard',
    routes: [
      { path: '/hotel', redirect: '/hotel/room-manage' },
      { name: '房间管理', icon: 'table', path: '/hotel/room-manage', component: './hotel/room-manage' },
      { name: '酒店前台', icon: 'table', path: '/hotel/hotel-front', component: './hotel/hotel-front' },
      { name: '订单管理', icon: 'table', path: '/hotel/order-manage', component: './hotel/order-manage' },
      { name: '折扣管理', icon: 'table', path: '/hotel/discount-manage', component: './hotel/discount-manage' },
      { name: '活动管理', icon: 'table', path: '/hotel/activity-manage', component: './hotel/activity-manage' },
    ],
  },
  { name: '人员管理', icon: 'table', path: '/staff-manage', component: './staff-manage' },
  { name: '系统设置', icon: 'table', path: '/system-setting', component: './system-setting' },
  { path: '/', redirect: '/welcome' },
  { path: '*', layout: false, component: './404' },
];

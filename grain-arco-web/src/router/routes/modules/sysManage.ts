import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const SysManage: AppRouteRecordRaw = {
  path: '/sysManage',
  name: 'sysManage',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.sysManage',
    requiresAuth: true,
    icon: 'icon-mind-mapping',
    order: 0,
  },
  children: [
    {
      path: 'sysApi',
      name: 'sysApi',
      component: () => import('@/views/sysManage/sysApi/index.vue'),
      meta: {
        locale: 'menu.sysApi',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'sysRole',
      name: 'sysRole',
      component: () => import('@/views/sysManage/sysRole/index.vue'),
      meta: {
        locale: 'menu.sysRole',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'sysMenu',
      name: 'sysMenu',
      component: () => import('@/views/sysManage/sysMenu/index.vue'),
      meta: {
        locale: 'menu.sysMenu',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'sysUser',
      name: 'sysUser',
      component: () => import('@/views/sysManage/sysUser/index.vue'),
      meta: {
        locale: 'menu.sysUser',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'sysOrganize',
      name: 'sysOrganize',
      component: () => import('@/views/sysManage/sysOrganize/index.vue'),
      meta: {
        locale: 'menu.organize',
        requiresAuth: true,
        roles: ['*'],
      },
    },

    {
      path: 'generateCode',
      name: 'generateCode',
      component: () => import('@/views/codeAssistant/index.vue'),
      meta: {
        locale: 'menu.generateCode',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'attachments',
      name: 'attachments',
      component: () => import('@/views/sysManage/sysFile/index.vue'),
      meta: {
        locale: 'menu.sysFile',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'sysLog',
      name: 'sysLog',
      component: () => import('@/views/sysManage/sysLog/index.vue'),
      meta: {
        locale: 'menu.sysLog',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
};

export default SysManage;

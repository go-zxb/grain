import Mock from 'mockjs';
import setupMock, {
  successResponseWrap,
  failResponseWrap,
} from '@/utils/setup-mock';

import { MockParams } from '@/types/mock';
import { isLogin } from '@/utils/auth';

setupMock({
  mock: false,
  setup() {
    // Mock.XHR.prototype.withCredentials = true;
    // 用户信息
    Mock.mock(new RegExp('/api/sysUser/info'), () => {
      if (isLogin()) {
        const role = window.localStorage.getItem('userRole') || 'admin';
        return successResponseWrap({
          name: '张漳',
          avatar:
            '//p3-armor.byteimg.com/tos-cn-i-49unhts6dw/dfdba5317c0c20ce20e64fac803d52bc.svg~tplv-49unhts6dw-image.image',
          email: 'wangliqun@email.com',
          job: 'frontend',
          jobName: '前端艺术家',
          organization: 'Frontend',
          organizationName: '前端',
          location: 'beijing',
          locationName: '北京',
          introduction: '人潇洒，性温存',
          personalWebsite: 'https://www.arco.design',
          phone: '150****0000',
          registrationDate: '2013-05-10 12:10:00',
          accountId: '15012312300',
          certification: 1,
          role,
        });
      }
      return failResponseWrap(null, '未登录', 50008);
    });

    // 登录
    Mock.mock(new RegExp('/api/sysUser/login'), (params: MockParams) => {
      const { username, password } = JSON.parse(params.body);
      if (!username) {
        return failResponseWrap(null, '用户名不能为空', 50000);
      }
      if (!password) {
        return failResponseWrap(null, '密码不能为空', 50000);
      }
      if (username === 'admin' && password === 'admin') {
        window.localStorage.setItem('userRole', 'admin');
        return successResponseWrap({
          token: '12345',
        });
      }
      if (username === 'user' && password === 'user') {
        window.localStorage.setItem('userRole', 'user');
        return successResponseWrap({
          token: '54321',
        });
      }
      return failResponseWrap(null, '账号或者密码错误', 50000);
    });

    // 登出
    Mock.mock(new RegExp('/api/sysUser/logout'), () => {
      return successResponseWrap(null);
    });

    // 用户的服务端菜单
    Mock.mock(new RegExp('/api/v1/sysMenu/userMenu'), () => {
      const menuList = [
        {
          path: '/sysManage',
          name: 'sysManage',
          meta: {
            locale: 'menu.sysManage',
            requiresAuth: true,
            icon: 'icon-command',
            order: 2,
          },
          children: [
            {
              path: 'api',
              name: 'api',
              meta: {
                locale: 'sysManage.api',
                requiresAuth: true,
                icon: '',
                order: 2,
                roles: ['2023'],
              },
            },
            {
              path: 'role',
              name: 'role',
              meta: {
                locale: 'sysManage.role',
                requiresAuth: true,
                icon: '',
                order: 1,
                roles: ['2023'],
              },
            },
          ],
        },
      ];
      return successResponseWrap(menuList);
    });
  },
});

import localeMessageBox from '@/components/message-box/locale/zh-CN';
import localeLogin from '@/views/login/locale/zh-CN';

import localeWorkplace from '@/views/dashboard/workplace/locale/zh-CN';

import localeMonitor from '@/views/dashboard/monitor/locale/zh-CN';

import api from '@/views/sysManage/sysApi/locale/zh-CN';

import role from '@/views/sysManage/sysRole/locale/zh-CN';

import menu from '@/views/sysManage/sysMenu/locale/zh-CN';

import user from '@/views/sysManage/sysUser/locale/zh-CN';

import log from '@/views/sysManage/sysLog/locale/zh-CN';

import organize from '@/views/sysManage/sysOrganize/locale/zh-CN';

import project from '@/views/codeAssistant/locale/zh-CN';
import model from '@/views/codeAssistant/component/model/zh-CN';
import Field from '@/views/codeAssistant/component/field/zh-CN';

import sysFile from '@/views/sysManage/sysFile/locale/zh-CN';

import localeSettings from './zh-CN/settings';

export default {
  'menu.dashboard': '仪表盘',
  'menu.server.dashboard': '仪表盘-服务端',
  'menu.server.workplace': '工作台-服务端',
  'menu.user': '个人中心',
  'menu.faq': '常见问题',
  'navbar.docs': '文档中心',
  'navbar.action.locale': '切换为中文',
  'menu.sysManage': '系统管理',
  'menu.attachments': '附件管理',
  'menu.codeAssistant': '代码助手',
  'menu.codeAssistantField': '数据字段',
  'menu.codeAssistantModel': '项目模块',
  'menu.arcoWebsite': 'Arco官方文档',
  ...localeMessageBox,
  ...localeSettings,
  ...localeWorkplace,
  ...localeMonitor,
  ...localeLogin,
  ...project,
  ...sysFile,
  ...organize,
  ...model,
  ...Field,
  ...api,
  ...role,
  ...menu,
  ...user,
  ...log,
};

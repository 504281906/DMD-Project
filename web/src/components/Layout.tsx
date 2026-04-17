import React from 'react';
import { Layout, Menu } from 'antd';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  DashboardOutlined,
  TabletOutlined,
  AppstoreOutlined,
  SafetyOutlined,
} from '@ant-design/icons';
import Header from './Header';

const { Content, Sider } = Layout;

interface AppLayoutProps {
  children: React.ReactNode;
}

const menuItems = [
  {
    key: '/dashboard',
    icon: <DashboardOutlined />,
    label: '仪表盘',
  },
  {
    key: '/devices',
    icon: <TabletOutlined />,
    label: '设备列表',
  },
  {
    key: '/apps',
    icon: <AppstoreOutlined />,
    label: '应用列表',
  },
  {
    key: '/policies',
    icon: <SafetyOutlined />,
    label: '策略列表',
  },
];

const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  const navigate = useNavigate();
  const location = useLocation();

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header />
      <Layout>
        <Sider width={200} style={{ background: '#fff' }}>
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
            style={{ height: '100%', borderRight: 0 }}
          />
        </Sider>
        <Layout style={{ padding: '24px' }}>
          <Content
            style={{
              background: '#fff',
              padding: 24,
              margin: 0,
              minHeight: 280,
              borderRadius: 8,
            }}
          >
            {children}
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
};

export default AppLayout;

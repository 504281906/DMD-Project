import React from 'react';
import { Layout, Typography, Avatar, Dropdown, Space } from 'antd';
import { UserOutlined, SettingOutlined, LogoutOutlined } from '@ant-design/icons';
import type { MenuProps } from 'antd';

const { Header: AntHeader } = Layout;
const { Title } = Typography;

const items: MenuProps['items'] = [
  {
    key: 'settings',
    icon: <SettingOutlined />,
    label: '设置',
  },
  {
    key: 'logout',
    icon: <LogoutOutlined />,
    label: '退出登录',
  },
];

const Header: React.FC = () => {
  return (
    <AntHeader
      style={{
        background: '#001529',
        padding: '0 24px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
      }}
    >
      <Title level={4} style={{ color: '#fff', margin: 0 }}>
        OpenMDM 管理后台
      </Title>
      <Dropdown menu={{ items }} placement="bottomRight">
        <Space style={{ cursor: 'pointer' }}>
          <Avatar icon={<UserOutlined />} />
          <span style={{ color: '#fff' }}>管理员</span>
        </Space>
      </Dropdown>
    </AntHeader>
  );
};

export default Header;

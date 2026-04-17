import React, { useEffect, useState } from 'react';
import { Row, Col, Card, Statistic, Table, Pie, Tag, Typography } from 'antd';
import {
  TabletOutlined,
  AppstoreOutlined,
  WifiOutlined,
  AlertOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import {
  getDashboardStats,
  getPlatformDistribution,
  getRecentDevices,
  DashboardStats,
  PlatformDistribution,
  RecentDevice,
} from '../services/api';

const { Title } = Typography;

const Dashboard: React.FC = () => {
  const [stats, setStats] = useState<DashboardStats>({
    totalDevices: 0,
    totalApps: 0,
    onlineDevices: 0,
    pendingAlerts: 0,
  });
  const [platformData, setPlatformData] = useState<PlatformDistribution[]>([]);
  const [recentDevices, setRecentDevices] = useState<RecentDevice[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [statsData, platformData, recentData] = await Promise.all([
          getDashboardStats(),
          getPlatformDistribution(),
          getRecentDevices(),
        ]);
        setStats(statsData);
        setPlatformData(platformData);
        setRecentDevices(recentData);
      } catch {
        // Use mock data when API is not available
        setStats({
          totalDevices: 128,
          totalApps: 45,
          onlineDevices: 86,
          pendingAlerts: 7,
        });
        setPlatformData([
          { platform: 'iOS', count: 72 },
          { platform: 'Android', count: 56 },
        ]);
        setRecentDevices([
          { id: '1', name: 'iPhone 15 Pro', platform: 'iOS', status: 'online', registeredAt: '2026-04-17 10:30' },
          { id: '2', name: 'Samsung Galaxy S24', platform: 'Android', status: 'online', registeredAt: '2026-04-17 09:15' },
          { id: '3', name: 'iPad Pro', platform: 'iOS', status: 'offline', registeredAt: '2026-04-16 18:00' },
          { id: '4', name: 'Xiaomi 14', platform: 'Android', status: 'online', registeredAt: '2026-04-16 15:45' },
          { id: '5', name: 'MacBook Air', platform: 'iOS', status: 'online', registeredAt: '2026-04-16 11:20' },
        ]);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const columns: ColumnsType<RecentDevice> = [
    {
      title: '设备名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '平台',
      dataIndex: 'platform',
      key: 'platform',
      render: (platform: string) => (
        <Tag color={platform === 'iOS' ? 'blue' : 'green'}>{platform}</Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'online' ? 'success' : 'default'}>
          {status === 'online' ? '在线' : '离线'}
        </Tag>
      ),
    },
    {
      title: '注册时间',
      dataIndex: 'registeredAt',
      key: 'registeredAt',
    },
  ];

  const pieData = platformData.map((item) => ({
    name: item.platform,
    value: item.count,
  }));

  return (
    <div>
      <Title level={4}>仪表盘</Title>

      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="设备总数"
              value={stats.totalDevices}
              prefix={<TabletOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="应用总数"
              value={stats.totalApps}
              prefix={<AppstoreOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="在线设备"
              value={stats.onlineDevices}
              prefix={<WifiOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="待处理告警"
              value={stats.pendingAlerts}
              prefix={<AlertOutlined />}
              valueStyle={{ color: '#cf1322' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={16}>
        <Col span={12}>
          <Card title="设备平台分布" loading={loading}>
            <Pie
              data={pieData}
              height={250}
              radius={['40%', '70%']}
              label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
              legend={{ position: 'bottom' }}
            />
          </Card>
        </Col>
        <Col span={12}>
          <Card title="最近注册设备" loading={loading}>
            <Table
              columns={columns}
              dataSource={recentDevices}
              rowKey="id"
              pagination={false}
              size="small"
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;

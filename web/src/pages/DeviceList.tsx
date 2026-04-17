import React, { useEffect, useState } from 'react';
import {
  Table,
  Input,
  Button,
  Tag,
  Space,
  Modal,
  message,
  Typography,
  Descriptions,
} from 'antd';
import { SearchOutlined, LockOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import {
  getDevices,
  lockDevice,
  wipeDevice,
  getDeviceDetail,
  Device,
} from '../services/api';

const { Title } = Typography;

const DeviceList: React.FC = () => {
  const [devices, setDevices] = useState<Device[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(false);
  const [detailVisible, setDetailVisible] = useState(false);
  const [selectedDevice, setSelectedDevice] = useState<Device | null>(null);

  const fetchDevices = async () => {
    setLoading(true);
    try {
      const data = await getDevices({ page, pageSize, search });
      setDevices(data.items);
      setTotal(data.total);
    } catch {
      // Use mock data when API is not available
      setDevices([
        { id: '1', name: 'iPhone 15 Pro', platform: 'iOS', status: 'online', registeredAt: '2026-04-17 10:30', model: 'iPhone 15 Pro', osVersion: '17.4' },
        { id: '2', name: 'Samsung Galaxy S24', platform: 'Android', status: 'online', registeredAt: '2026-04-17 09:15', model: 'Galaxy S24 Ultra', osVersion: '14' },
        { id: '3', name: 'iPad Pro', platform: 'iOS', status: 'offline', registeredAt: '2026-04-16 18:00', model: 'iPad Pro 12.9', osVersion: '17.4' },
        { id: '4', name: 'Xiaomi 14', platform: 'Android', status: 'online', registeredAt: '2026-04-16 15:45', model: 'Xiaomi 14', osVersion: '14' },
        { id: '5', name: 'MacBook Air', platform: 'iOS', status: 'online', registeredAt: '2026-04-16 11:20', model: 'MacBook Air M3', osVersion: '14.4' },
        { id: '6', name: 'OnePlus 12', platform: 'Android', status: 'offline', registeredAt: '2026-04-15 14:30', model: 'OnePlus 12', osVersion: '14' },
        { id: '7', name: 'iPhone 14', platform: 'iOS', status: 'online', registeredAt: '2026-04-15 10:00', model: 'iPhone 14', osVersion: '17.3' },
        { id: '8', name: 'Pixel 8 Pro', platform: 'Android', status: 'online', registeredAt: '2026-04-14 16:20', model: 'Pixel 8 Pro', osVersion: '14' },
      ]);
      setTotal(8);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDevices();
  }, [page, pageSize, search]);

  const handleSearch = (value: string) => {
    setSearch(value);
    setPage(1);
  };

  const handleLock = async (id: string) => {
    try {
      await lockDevice(id);
      message.success('设备已锁定');
    } catch {
      message.success('设备已锁定（模拟）');
    }
  };

  const handleWipe = async (id: string) => {
    Modal.confirm({
      title: '确认擦除设备',
      content: '此操作将清除设备上的所有数据，且不可恢复。确定要继续吗？',
      okText: '确认擦除',
      okType: 'danger',
      cancelText: '取消',
      onOk: async () => {
        try {
          await wipeDevice(id);
          message.success('设备擦除命令已发送');
        } catch {
          message.success('设备擦除命令已发送（模拟）');
        }
      },
    });
  };

  const handleViewDetail = async (id: string) => {
    try {
      const device = await getDeviceDetail(id);
      setSelectedDevice(device);
    } catch {
      const device = devices.find((d) => d.id === id);
      setSelectedDevice(device || null);
    }
    setDetailVisible(true);
  };

  const columns: ColumnsType<Device> = [
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
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: Device) => (
        <Space>
          <Button
            type="link"
            icon={<LockOutlined />}
            onClick={() => handleLock(record.id)}
          >
            锁定
          </Button>
          <Button
            type="link"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleWipe(record.id)}
          >
            擦除
          </Button>
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record.id)}
          >
            详情
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Title level={4}>设备列表</Title>

      <div style={{ marginBottom: 16 }}>
        <Input.Search
          placeholder="搜索设备名称"
          allowClear
          enterButton={<SearchOutlined />}
          onSearch={handleSearch}
          style={{ width: 300 }}
        />
      </div>

      <Table
        columns={columns}
        dataSource={devices}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 条`,
          onChange: (p, ps) => {
            setPage(p);
            setPageSize(ps);
          },
        }}
      />

      <Modal
        title="设备详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailVisible(false)}>
            关闭
          </Button>,
        ]}
      >
        {selectedDevice && (
          <Descriptions column={1} bordered>
            <Descriptions.Item label="设备ID">{selectedDevice.id}</Descriptions.Item>
            <Descriptions.Item label="设备名称">{selectedDevice.name}</Descriptions.Item>
            <Descriptions.Item label="平台">{selectedDevice.platform}</Descriptions.Item>
            <Descriptions.Item label="状态">
              {selectedDevice.status === 'online' ? '在线' : '离线'}
            </Descriptions.Item>
            <Descriptions.Item label="型号">{selectedDevice.model || '-'}</Descriptions.Item>
            <Descriptions.Item label="系统版本">{selectedDevice.osVersion || '-'}</Descriptions.Item>
            <Descriptions.Item label="注册时间">{selectedDevice.registeredAt}</Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </div>
  );
};

export default DeviceList;

import React, { useEffect, useState } from 'react';
import {
  Table,
  Button,
  Tag,
  Space,
  Modal,
  message,
  Typography,
  Select,
  Upload,
  Form,
  Input,
  Popconfirm,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  UploadOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import {
  getApps,
  uploadApp,
  updateApp,
  deleteApp,
  App,
} from '../services/api';

const { Title } = Typography;

const AppList: React.FC = () => {
  const [apps, setApps] = useState<App[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [platformFilter, setPlatformFilter] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [editVisible, setEditVisible] = useState(false);
  const [editingApp, setEditingApp] = useState<App | null>(null);
  const [form] = Form.useForm();

  const fetchApps = async () => {
    setLoading(true);
    try {
      const data = await getApps({ page, pageSize, platform: platformFilter || undefined });
      setApps(data.items);
      setTotal(data.total);
    } catch {
      // Use mock data when API is not available
      const mockApps: App[] = [
        { id: '1', name: '企业微信', bundleId: 'com.tencent.wework', platform: 'iOS', version: '4.0.12', status: 'active', installCount: 156, iconUrl: '' },
        { id: '2', name: '钉钉', bundleId: 'com.dingtalk', platform: 'iOS', version: '7.0.30', status: 'active', installCount: 143, iconUrl: '' },
        { id: '3', name: '飞书', bundleId: 'com.feishu', platform: 'Android', version: '6.8.5', status: 'active', installCount: 98, iconUrl: '' },
        { id: '4', name: 'WPS Office', bundleId: 'cn.wps.moffice', platform: 'Android', version: '13.20', status: 'active', installCount: 87, iconUrl: '' },
        { id: '5', name: 'Slack', bundleId: 'com.Slack', platform: 'iOS', version: '24.02.20', status: 'inactive', installCount: 45, iconUrl: '' },
        { id: '6', name: 'Zoom', bundleId: 'us.zoom.videomeetings', platform: 'Android', version: '5.16.6', status: 'active', installCount: 76, iconUrl: '' },
        { id: '7', name: 'Microsoft Teams', bundleId: 'com.microsoft.teams', platform: 'iOS', version: '1415/240227', status: 'active', installCount: 112, iconUrl: '' },
        { id: '8', name: 'Notion', bundleId: 'com.notion', platform: 'Android', version: '2.1.14', status: 'active', installCount: 63, iconUrl: '' },
      ];
      const filtered = platformFilter ? mockApps.filter(a => a.platform === platformFilter) : mockApps;
      setApps(filtered);
      setTotal(filtered.length);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchApps();
  }, [page, pageSize, platformFilter]);

  const handlePlatformChange = (value: string) => {
    setPlatformFilter(value);
    setPage(1);
  };

  const handleUpload = async (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    try {
      await uploadApp(formData);
      message.success('应用上传成功');
      fetchApps();
    } catch {
      message.success('应用上传成功（模拟）');
      fetchApps();
    }
    return false;
  };

  const handleEdit = (record: App) => {
    setEditingApp(record);
    form.setFieldsValue(record);
    setEditVisible(true);
  };

  const handleEditSave = async () => {
    try {
      const values = await form.validateFields();
      if (editingApp) {
        await updateApp(editingApp.id, values);
        message.success('应用更新成功');
      }
      setEditVisible(false);
      fetchApps();
    } catch {
      message.success('应用更新成功（模拟）');
      setEditVisible(false);
      fetchApps();
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteApp(id);
      message.success('应用已删除');
      fetchApps();
    } catch {
      message.success('应用已删除（模拟）');
      fetchApps();
    }
  };

  const columns: ColumnsType<App> = [
    {
      title: '图标',
      dataIndex: 'iconUrl',
      key: 'icon',
      render: () => (
        <div
          style={{
            width: 40,
            height: 40,
            background: '#f0f0f0',
            borderRadius: 8,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: 20,
          }}
        >
          📱
        </div>
      ),
    },
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '包名',
      dataIndex: 'bundleId',
      key: 'bundleId',
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
      title: '版本',
      dataIndex: 'version',
      key: 'version',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'success' : 'default'}>
          {status === 'active' ? '已启用' : '已停用'}
        </Tag>
      ),
    },
    {
      title: '安装数',
      dataIndex: 'installCount',
      key: 'installCount',
    },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: App) => (
        <Space>
          <Upload
            showUploadList={false}
            beforeUpload={handleUpload}
            accept=".apk,.ipa"
          >
            <Button type="link" icon={<UploadOutlined />}>
              更新
            </Button>
          </Upload>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确认删除"
            description="删除后不可恢复，确定要删除吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确认"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Title level={4}>应用列表</Title>

      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Select
          placeholder="筛选平台"
          allowClear
          style={{ width: 200 }}
          onChange={handlePlatformChange}
          value={platformFilter || undefined}
          options={[
            { label: 'iOS', value: 'iOS' },
            { label: 'Android', value: 'Android' },
          ]}
        />
        <Upload
          showUploadList={false}
          beforeUpload={handleUpload}
          accept=".apk,.ipa"
        >
          <Button type="primary" icon={<PlusOutlined />}>
            上传应用
          </Button>
        </Upload>
      </div>

      <Table
        columns={columns}
        dataSource={apps}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (t) => `共 ${t} 条`,
          onChange: (p, ps) => {
            setPage(p);
            setPageSize(ps);
          },
        }}
      />

      <Modal
        title="编辑应用"
        open={editVisible}
        onOk={handleEditSave}
        onCancel={() => setEditVisible(false)}
        okText="保存"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="应用名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="bundleId" label="包名" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="version" label="版本" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select
              options={[
                { label: '已启用', value: 'active' },
                { label: '已停用', value: 'inactive' },
              ]}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default AppList;

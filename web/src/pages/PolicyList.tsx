import React, { useEffect, useState } from 'react';
import {
  Table,
  Button,
  Tag,
  Space,
  Modal,
  message,
  Typography,
  Form,
  Input,
  Select,
  Popconfirm,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SafetyOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import {
  getPolicies,
  createPolicy,
  updatePolicy,
  deletePolicy,
  assignPolicyToDevices,
  Policy,
} from '../services/api';

const { Title } = Typography;

const PolicyList: React.FC = () => {
  const [policies, setPolicies] = useState<Policy[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [loading, setLoading] = useState(false);
  const [createVisible, setCreateVisible] = useState(false);
  const [editVisible, setEditVisible] = useState(false);
  const [assignVisible, setAssignVisible] = useState(false);
  const [editingPolicy, setEditingPolicy] = useState<Policy | null>(null);
  const [assigningPolicy, setAssigningPolicy] = useState<Policy | null>(null);
  const [form] = Form.useForm();
  const [assignForm] = Form.useForm();

  const fetchPolicies = async () => {
    setLoading(true);
    try {
      const data = await getPolicies({ page, pageSize });
      setPolicies(data.items);
      setTotal(data.total);
    } catch {
      // Use mock data when API is not available
      const mockPolicies: Policy[] = [
        { id: '1', name: '密码复杂度策略', type: 'password', status: 'active', deviceCount: 45, description: '要求密码包含大小写字母、数字和特殊字符，长度不少于8位' },
        { id: '2', name: '设备加密策略', type: 'encryption', status: 'active', deviceCount: 38, description: '强制启用设备加密，保护数据安全' },
        { id: '3', name: '合规检查策略', type: 'compliance', status: 'active', deviceCount: 52, description: '定期检查设备是否越狱/ROOT' },
        { id: '4', name: '弱密码检测', type: 'password', status: 'inactive', deviceCount: 0, description: '检测并禁止使用弱密码' },
        { id: '5', name: '数据加密传输', type: 'encryption', status: 'active', deviceCount: 60, description: '要求所有数据传输使用TLS加密' },
        { id: '6', name: '应用白名单', type: 'compliance', status: 'inactive', deviceCount: 0, description: '仅允许安装白名单中的应用' },
      ];
      setPolicies(mockPolicies);
      setTotal(mockPolicies.length);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPolicies();
  }, [page, pageSize]);

  const getTypeTag = (type: string) => {
    const config: Record<string, { color: string; label: string }> = {
      password: { color: 'blue', label: '密码' },
      encryption: { color: 'green', label: '加密' },
      compliance: { color: 'orange', label: '合规' },
    };
    const { color, label } = config[type] || { color: 'default', label: type };
    return <Tag color={color}>{label}</Tag>;
  };

  const handleCreate = () => {
    form.resetFields();
    setCreateVisible(true);
  };

  const handleCreateSave = async () => {
    try {
      const values = await form.validateFields();
      await createPolicy(values);
      message.success('策略创建成功');
      setCreateVisible(false);
      fetchPolicies();
    } catch {
      message.success('策略创建成功（模拟）');
      setCreateVisible(false);
      fetchPolicies();
    }
  };

  const handleEdit = (record: Policy) => {
    setEditingPolicy(record);
    form.setFieldsValue(record);
    setEditVisible(true);
  };

  const handleEditSave = async () => {
    try {
      const values = await form.validateFields();
      if (editingPolicy) {
        await updatePolicy(editingPolicy.id, values);
        message.success('策略更新成功');
      }
      setEditVisible(false);
      fetchPolicies();
    } catch {
      message.success('策略更新成功（模拟）');
      setEditVisible(false);
      fetchPolicies();
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deletePolicy(id);
      message.success('策略已删除');
      fetchPolicies();
    } catch {
      message.success('策略已删除（模拟）');
      fetchPolicies();
    }
  };

  const handleAssign = (record: Policy) => {
    setAssigningPolicy(record);
    assignForm.resetFields();
    setAssignVisible(true);
  };

  const handleAssignSave = async () => {
    try {
      const values = await assignForm.validateFields();
      if (assigningPolicy && values.deviceIds) {
        await assignPolicyToDevices(assigningPolicy.id, values.deviceIds);
        message.success('策略分配成功');
      }
      setAssignVisible(false);
      fetchPolicies();
    } catch {
      message.success('策略分配成功（模拟）');
      setAssignVisible(false);
      fetchPolicies();
    }
  };

  const columns: ColumnsType<Policy> = [
    {
      title: '策略名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => getTypeTag(type),
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
      title: '关联设备数',
      dataIndex: 'deviceCount',
      key: 'deviceCount',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: Policy) => (
        <Space>
          <Button
            type="link"
            icon={<SafetyOutlined />}
            onClick={() => handleAssign(record)}
          >
            分配
          </Button>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确认删除"
            description="确定要删除此策略吗？"
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
      <Title level={4}>策略列表</Title>

      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          创建策略
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={policies}
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
        title="创建策略"
        open={createVisible}
        onOk={handleCreateSave}
        onCancel={() => setCreateVisible(false)}
        okText="创建"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="策略名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="type" label="策略类型" rules={[{ required: true }]}>
            <Select
              options={[
                { label: '密码策略', value: 'password' },
                { label: '加密策略', value: 'encryption' },
                { label: '合规策略', value: 'compliance' },
              ]}
            />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue="active">
            <Select
              options={[
                { label: '已启用', value: 'active' },
                { label: '已停用', value: 'inactive' },
              ]}
            />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="编辑策略"
        open={editVisible}
        onOk={handleEditSave}
        onCancel={() => setEditVisible(false)}
        okText="保存"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="策略名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="type" label="策略类型" rules={[{ required: true }]}>
            <Select
              options={[
                { label: '密码策略', value: 'password' },
                { label: '加密策略', value: 'encryption' },
                { label: '合规策略', value: 'compliance' },
              ]}
            />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select
              options={[
                { label: '已启用', value: 'active' },
                { label: '已停用', value: 'inactive' },
              ]}
            />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="分配策略到设备"
        open={assignVisible}
        onOk={handleAssignSave}
        onCancel={() => setAssignVisible(false)}
        okText="分配"
        cancelText="取消"
      >
        <p>将策略 <strong>{assigningPolicy?.name}</strong> 分配给以下设备：</p>
        <Form form={assignForm} layout="vertical">
          <Form.Item
            name="deviceIds"
            label="选择设备"
            rules={[{ required: true, message: '请选择至少一个设备' }]}
          >
            <Select
              mode="multiple"
              placeholder="搜索并选择设备"
              options={[
                { label: 'iPhone 15 Pro (device-001)', value: 'device-001' },
                { label: 'Samsung Galaxy S24 (device-002)', value: 'device-002' },
                { label: 'iPad Pro (device-003)', value: 'device-003' },
                { label: 'Xiaomi 14 (device-004)', value: 'device-004' },
                { label: 'MacBook Air (device-005)', value: 'device-005' },
              ]}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default PolicyList;

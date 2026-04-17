import React from 'react';
import { Routes, Route } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import AppLayout from './components/Layout';
import Dashboard from './pages/Dashboard';
import DeviceList from './pages/DeviceList';
import AppList from './pages/AppList';
import PolicyList from './pages/PolicyList';

const App: React.FC = () => {
  return (
    <ConfigProvider locale={zhCN}>
      <AppLayout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/devices" element={<DeviceList />} />
          <Route path="/apps" element={<AppList />} />
          <Route path="/policies" element={<PolicyList />} />
        </Routes>
      </AppLayout>
    </ConfigProvider>
  );
};

export default App;

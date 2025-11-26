import { PageContainer } from '@ant-design/pro-components';
import { Card } from 'antd';
import React from 'react';
import CheckInForm from '../components/CheckInForm';

const CheckInPage: React.FC = () => {
  return (
    <PageContainer title="办理入住">
      <Card>
        <CheckInForm />
      </Card>
    </PageContainer>
  );
};

export default CheckInPage;
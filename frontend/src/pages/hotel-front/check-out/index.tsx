import { PageContainer } from '@ant-design/pro-components';
import { Card } from 'antd';
import React from 'react';
import CheckOutForm from '../components/CheckOutForm';

const CheckOutPage: React.FC = () => {
  return (
    <PageContainer title="办理退房">
      <Card>
        <CheckOutForm />
      </Card>
    </PageContainer>
  );
};

export default CheckOutPage;
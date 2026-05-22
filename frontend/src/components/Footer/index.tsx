import { DefaultFooter } from '@ant-design/pro-components';
import React from 'react';

const Footer: React.FC = () => {
  const year = new Date().getFullYear();

  return (
    <DefaultFooter
      style={{
        background: 'none',
      }}
      copyright={`${year} GoHotel 酒店运营中台`}
      links={[
        {
          key: 'support',
          title: 'GoHotel 运营手册',
          href: 'mailto:ops@gohotel.com',
          blankTarget: false,
        },
        {
          key: 'room-inventory',
          title: '房态库存矩阵',
          href: '/room-manage/inventory',
          blankTarget: false,
        },
        {
          key: 'order-manage',
          title: '订单与账务中心',
          href: '/order-manage',
          blankTarget: false,
        },
      ]}
    />
  );
};

export default Footer;

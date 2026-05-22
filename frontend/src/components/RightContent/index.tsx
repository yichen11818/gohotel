import { CustomerServiceOutlined } from '@ant-design/icons';

export type SiderTheme = 'light' | 'dark';

export const SelectLang: React.FC = () => null;

export const SupportLink: React.FC = () => {
  return (
    <a
      href="mailto:ops@gohotel.com?subject=GoHotel%20运营支持"
      rel="noreferrer"
      style={{
        display: 'inline-flex',
        padding: '4px',
        fontSize: '18px',
        color: 'inherit',
      }}
    >
      <CustomerServiceOutlined />
    </a>
  );
};

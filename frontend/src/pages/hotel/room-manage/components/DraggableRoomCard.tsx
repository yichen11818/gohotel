import React from 'react';
import { Card, Tag, Button, Space, Popconfirm } from 'antd';
import { EditOutlined, DeleteOutlined, UserOutlined, ExpandOutlined } from '@ant-design/icons';
import { useDrag } from 'react-dnd';

const ItemTypes = {
  ROOM_CARD: 'room_card',
};

interface DraggableRoomCardProps {
  room: API.Room;
  left: number;
  top: number;
  onEdit: () => void;
  onDelete: (id: number) => void;
  onDrop: (id: number, left: number, top: number) => void;
}

const DraggableRoomCard: React.FC<DraggableRoomCardProps> = ({
  room,
  left,
  top,
  onEdit,
  onDelete,
  onDrop,
}) => {
  const [{ isDragging }, drag] = useDrag(
    () => ({
      type: ItemTypes.ROOM_CARD,
      item: { id: room.id, left, top },
      collect: (monitor) => ({
        isDragging: monitor.isDragging(),
      }),
      end: (item, monitor) => {
        const delta = monitor.getDifferenceFromInitialOffset();
        if (delta) {
          const newLeft = Math.round(item.left + delta.x);
          const newTop = Math.round(item.top + delta.y);
          onDrop(item.id, newLeft, newTop);
        }
      },
    }),
    [room.id, left, top, onDrop],
  );

  const getStatusColor = (status?: string) => {
    switch (status) {
      case 'available':
        return 'success';
      case 'occupied':
        return 'error';
      case 'maintenance':
        return 'warning';
      default:
        return 'default';
    }
  };

  const getStatusText = (status?: string) => {
    switch (status) {
      case 'available':
        return '可用';
      case 'occupied':
        return '占用';
      case 'maintenance':
        return '维护中';
      default:
        return '未知';
    }
  };

  return (
    <div
      ref={drag}
      style={{
        position: 'absolute',
        left,
        top,
        opacity: isDragging ? 0.5 : 1,
        cursor: 'move',
        width: 280,
      }}
    >
      <Card
        hoverable
        size="small"
        title={
          <Space>
            <span style={{ fontSize: '16px', fontWeight: 'bold' }}>
              房间 {room.room_number}
            </span>
            <Tag color={getStatusColor(room.status)}>{getStatusText(room.status)}</Tag>
          </Space>
        }
        extra={
          <Space>
            <Button type="text" size="small" icon={<EditOutlined />} onClick={onEdit} />
            <Popconfirm
              title="确定要删除这个房间吗？"
              onConfirm={() => onDelete(room.id)}
              okText="确定"
              cancelText="取消"
            >
              <Button type="text" size="small" danger icon={<DeleteOutlined />} />
            </Popconfirm>
          </Space>
        }
        style={{
          boxShadow: isDragging ? '0 8px 16px rgba(0,0,0,0.2)' : '0 2px 8px rgba(0,0,0,0.1)',
        }}
      >
        <div style={{ fontSize: '13px' }}>
          <div style={{ marginBottom: 8 }}>
            <Tag color="blue">{room.room_type}</Tag>
            <Tag>{room.floor}楼</Tag>
          </div>

          <div style={{ marginBottom: 4 }}>
            <span style={{ fontWeight: 'bold', color: '#ff4d4f', fontSize: '18px' }}>
              ¥{room.price}
            </span>
            <span style={{ color: '#999', fontSize: '12px' }}>/晚</span>
            {room.original_price && room.original_price > room.price && (
              <span
                style={{
                  marginLeft: 8,
                  textDecoration: 'line-through',
                  color: '#999',
                  fontSize: '12px',
                }}
              >
                ¥{room.original_price}
              </span>
            )}
          </div>

          <div style={{ color: '#666', lineHeight: '22px' }}>
            <div>
              <UserOutlined /> 可住 {room.capacity} 人
            </div>
            {room.area && (
              <div>
                <ExpandOutlined /> {room.area} m²
              </div>
            )}
            {room.bed_type && <div>{room.bed_type}</div>}
          </div>

          {room.description && (
            <div
              style={{
                marginTop: 8,
                paddingTop: 8,
                borderTop: '1px solid #f0f0f0',
                color: '#999',
                fontSize: '12px',
              }}
            >
              {room.description}
            </div>
          )}
        </div>
      </Card>
    </div>
  );
};

export default DraggableRoomCard;

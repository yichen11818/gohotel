import { PageContainer } from '@ant-design/pro-components';
import { Card, Button, message, Spin, theme } from 'antd';
import React, { useState, useEffect } from 'react';
import { DndProvider } from 'react-dnd';
import { HTML5Backend } from 'react-dnd-html5-backend';
import DraggableRoomCard from './components/DraggableRoomCard';
import { useRequest } from '@umijs/max';
import { getRooms } from '@/services/api/fangjian';
import RoomFormModal from './components/RoomFormModal';

interface RoomPosition {
  id: number;
  left: number;
  top: number;
}

const RoomManage: React.FC = () => {
  const { token } = theme.useToken();
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [editingRoom, setEditingRoom] = useState<API.Room | null>(null);
  const [roomPositions, setRoomPositions] = useState<RoomPosition[]>([]);
  const [selectedFloor, setSelectedFloor] = useState<number>(1);

  // 使用 useRequest 获取房间数据
  const { data: rooms, loading, run: reloadRooms } = useRequest(() => getRooms({ page_size: 100 }), {
    formatResult: (res) => (Array.isArray(res) ? res : res.data || []),
  });

  // 提取楼层数据
  const floors = rooms ? [...new Set(rooms.map((r) => r.floor))].sort((a, b) => a - b) : [];

  // 根据楼层筛选房间
  const filteredRooms = rooms?.filter(
    (room) => room.floor === selectedFloor,
  );

  // 当房间数据加载完成后,初始化房间位置
  useEffect(() => {
    if (rooms && rooms.length > 0) {
      // 检查是否已经有保存的位置数据(从 localStorage 或后端获取)
      const savedPositions = localStorage.getItem('roomPositions');
      
      if (savedPositions) {
        try {
          const parsed = JSON.parse(savedPositions);
          // 确保所有房间都有位置数据
          const positionsMap = new Map(parsed.map((p: RoomPosition) => [p.id, p]));
          const newPositions: RoomPosition[] = [];
          
          rooms.forEach((room, index) => {
            if (positionsMap.has(room.id)) {
              newPositions.push(positionsMap.get(room.id)!);
            } else {
              // 新房间,生成默认位置
              newPositions.push(generateDefaultPosition(room.id, index));
            }
          });
          
          setRoomPositions(newPositions);
        } catch (error) {
          console.error('Failed to parse saved positions:', error);
          initializeDefaultPositions();
        }
      } else {
        initializeDefaultPositions();
      }
    }
  }, [rooms]);

  // 生成默认位置(网格布局)
  const generateDefaultPosition = (id: number, index: number): RoomPosition => {
    const cardWidth = 300; // 卡片宽度 + 间距
    const cardHeight = 240; // 卡片高度 + 间距
    const columns = 4; // 每行显示4个卡片
    
    const row = Math.floor(index / columns);
    const col = index % columns;
    
    return {
      id,
      left: col * cardWidth + 20,
      top: row * cardHeight + 20,
    };
  };

  // 初始化所有房间的默认位置
  const initializeDefaultPositions = () => {
    if (!rooms) return;
    
    const positions = rooms.map((room, index) => generateDefaultPosition(room.id, index));
    setRoomPositions(positions);
  };

  // 处理拖拽结束,更新房间位置
  const handleDrop = (id: number, left: number, top: number) => {
    setRoomPositions((prevPositions) => {
      const newPositions = prevPositions.map((pos) =>
        pos.id === id ? { ...pos, left, top } : pos,
      );
      
      // 保存到 localStorage
      localStorage.setItem('roomPositions', JSON.stringify(newPositions));
      
      // TODO: 这里可以调用 API 保存到后端
      // await updateRoomPosition({ id, x: left, y: top });
      
      return newPositions;
    });
  };

  const handleOpenModal = (room: API.Room | null) => {
    setEditingRoom(room);
    setIsModalVisible(true);
  };

  const handleCancelModal = () => {
    setIsModalVisible(false);
    setEditingRoom(null);
  };

  const handleSuccess = () => {
    setIsModalVisible(false);
    setEditingRoom(null);
    reloadRooms(); // 重新加载数据
  };

  const handleDelete = async (id: number) => {
    try {
      // 假设你有一个 deleteRoom 的 API 方法
      // await deleteRoom({ id });
      message.success(`房间 ${id} 删除成功(模拟)`);
      
      // 从位置数据中移除
      setRoomPositions((prev) => prev.filter((pos) => pos.id !== id));
      
      reloadRooms(); // 重新加载数据
    } catch (error) {
      message.error('删除失败');
    }
  };

  // 重置布局
  const handleResetLayout = () => {
    initializeDefaultPositions();
    localStorage.removeItem('roomPositions');
    message.success('布局已重置');
  };

  return (
    <PageContainer
      title="房间可视化管理"
      extra={[
        <Button key="reset" onClick={handleResetLayout}>
          重置布局
        </Button>,
        <Button key="create" type="primary" onClick={() => handleOpenModal(null)}>
          新建房间
        </Button>,
      ]}
    >
      <DndProvider backend={HTML5Backend}>
        <Card bodyStyle={{ padding: 0, position: 'relative' }}>
          <div
            style={{
              position: 'absolute',
              left: 24,
              top: 24,
              zIndex: 10,
              backgroundColor: token.colorBgContainer,
              padding: '12px 8px',
              borderRadius: token.borderRadiusLG,
              boxShadow: token.boxShadowSecondary,
              border: `1px solid ${token.colorBorderSecondary}`,
              minWidth: 80,
            }}
          >
            <div
              style={{
                fontSize: 14,
                fontWeight: 'bold',
                textAlign: 'center',
                marginBottom: 12,
                color: token.colorPrimary,
              }}
            >
              楼层
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
              {floors.map((floor) => (
                <Button
                  key={floor}
                  type={selectedFloor === floor ? 'primary' : 'default'}
                  onClick={() => setSelectedFloor(floor)}
                  style={{ width: '100%' }}
                >
                  {floor}楼
                </Button>
              ))}
            </div>
          </div>
          {loading ? (
            <div style={{ textAlign: 'center', padding: '50px 0' }}>
              <Spin size="large" />
            </div>
          ) : (
            <div
              style={{
                position: 'relative',
                width: '100%',
                minHeight: '800px',
                height: 'calc(100vh - 200px)',
                backgroundColor: token.colorBgLayout,
                backgroundImage:
                  `linear-gradient(${token.colorSplit} 1px, transparent 1px), linear-gradient(90deg, ${token.colorSplit} 1px, transparent 1px)`,
                backgroundSize: '20px 20px',
                overflow: 'auto',
              }}
            >
              {filteredRooms?.map((room) => {
                const position = roomPositions.find((pos) => pos.id === room.id);
                if (!position) return null;

                return (
                  <DraggableRoomCard
                    key={room.id}
                    room={room}
                    left={position.left}
                    top={position.top}
                    onEdit={() => handleOpenModal(room)}
                    onDelete={handleDelete}
                    onDrop={handleDrop}
                  />
                );
              })}
            </div>
          )}
        </Card>
      </DndProvider>

      <RoomFormModal
        visible={isModalVisible}
        room={editingRoom}
        onCancel={handleCancelModal}
        onSuccess={handleSuccess}
      />
    </PageContainer>
  );
};

export default RoomManage;

import React, { useState, useCallback, useRef, useEffect } from 'react';
import { Card, Button, Space, Popconfirm } from 'antd';
import { DeleteOutlined, GatewayOutlined, FormOutlined } from '@ant-design/icons';
import { useDrag } from 'react-dnd';
import { ItemTypes } from './DraggableFacilityCard';
import Iconfont from '@/components/Iconfont';

interface DraggableRoomCardProps {
  room: API.Room;
  left: number;
  top: number;
  width: number;
  height: number;
  onEdit?: (room: API.Room) => void;
  onDelete: (id: number) => void;
  onDrop: (id: number, left: number, top: number) => void;
  onResizeComplete?: (id: number, newWidth: number, newHeight: number, roomType?: string) => void;
}

// 网格大小（与背景网格一致）
const GRID_SIZE = 20;
const MIN_SIZE = 60;
const MAX_SIZE = 400;

// 对齐到网格的函数
const snapToGrid = (value: number): number => {
  return Math.round(value / GRID_SIZE) * GRID_SIZE;
};

const DraggableRoomCard: React.FC<DraggableRoomCardProps> = ({
  room,
  left,
  top,
  width,
  height,
  onEdit,
  onDelete,
  onDrop,
  onResizeComplete,
}) => {
  const [isResizing, setIsResizing] = useState(false);
  const [resizeWidth, setResizeWidth] = useState(width);
  const [resizeHeight, setResizeHeight] = useState(height);
  const [isDraggingResize, setIsDraggingResize] = useState(false);
  const startPosRef = useRef({ x: 0, y: 0, width: 0, height: 0 });

  // 同步外部尺寸变化
  useEffect(() => {
    if (!isResizing) {
      setResizeWidth(width);
      setResizeHeight(height);
    }
  }, [width, height, isResizing]);

  const [{ isDragging }, drag, preview] = useDrag(
    () => ({
      type: ItemTypes.ROOM_CARD,
      item: { id: room.id, left, top, room, width, height },
      canDrag: !isResizing,
      collect: (monitor) => ({
        isDragging: monitor.isDragging(),
      }),
      end: (item, monitor) => {
        const delta = monitor.getDifferenceFromInitialOffset();
        if (delta && item.id) {
          const newLeft = item.left + delta.x;
          const newTop = item.top + delta.y;
          const snappedLeft = snapToGrid(newLeft);
          const snappedTop = snapToGrid(newTop);
          onDrop(item.id, snappedLeft, snappedTop);
        }
      },
    }),
    [room, left, top, width, height, onDrop, isResizing],
  );

  // 使用空图片作为拖动预览
  useEffect(() => {
    preview(new Image(), { captureDraggingState: true });
  }, [preview]);

  // 开始调整大小拖拽
  const handleResizeMouseDown = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    e.preventDefault();
    setIsDraggingResize(true);
    startPosRef.current = {
      x: e.clientX,
      y: e.clientY,
      width: resizeWidth,
      height: resizeHeight,
    };
  }, [resizeWidth, resizeHeight]);

  // 处理拖拽移动和结束
  useEffect(() => {
    if (!isDraggingResize) return;

    const handleMouseMove = (e: MouseEvent) => {
      const deltaX = e.clientX - startPosRef.current.x;
      const deltaY = e.clientY - startPosRef.current.y;
      
      let newWidth = startPosRef.current.width + deltaX;
      let newHeight = startPosRef.current.height + deltaY;
      
      newWidth = Math.max(MIN_SIZE, Math.min(MAX_SIZE, snapToGrid(newWidth)));
      newHeight = Math.max(MIN_SIZE, Math.min(MAX_SIZE, snapToGrid(newHeight)));
      
      setResizeWidth(newWidth);
      setResizeHeight(newHeight);
    };

    const handleMouseUp = () => {
      setIsDraggingResize(false);
      // 拖拽结束后，如果尺寸有变化，触发弹窗
      const finalWidth = resizeWidth;
      const finalHeight = resizeHeight;
      
      if (finalWidth !== width || finalHeight !== height) {
        // 调用回调触发弹窗确认
        if (onResizeComplete && room.id) {
          onResizeComplete(room.id, finalWidth, finalHeight, room.room_type);
        }
      }
      setIsResizing(false);
    };

    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);

    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, [isDraggingResize, resizeWidth, resizeHeight, width, height, onResizeComplete, room.id, room.room_type]);

  // 进入调整大小模式
  const handleEnterResizeMode = (e: React.MouseEvent) => {
    e.stopPropagation();
    setIsResizing(true);
    setResizeWidth(width);
    setResizeHeight(height);
  };

  const currentWidth = isResizing ? resizeWidth : width;
  const currentHeight = isResizing ? resizeHeight : height;

  return (
    <div
      ref={isResizing ? undefined : drag as any}
      className="room-card-container"
      style={{
        position: 'absolute',
        left,
        top,
        opacity: isDragging ? 0 : 1,
        cursor: isResizing ? 'default' : 'move',
        width: currentWidth,
        height: currentHeight,
      }}
    >
      <Card
        hoverable={!isResizing}
        size="small"
        bodyStyle={{
          padding: '8px',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
          height: '100%',
        }}
        style={{
          boxShadow: isResizing 
            ? '0 0 0 2px #1890ff, 0 8px 16px rgba(24,144,255,0.3)' 
            : isDragging 
              ? '0 8px 16px rgba(0,0,0,0.2)' 
              : '0 2px 8px rgba(0,0,0,0.1)',
          height: '100%',
          border: isResizing 
            ? '2px solid #1890ff'
            : `2px solid ${
                room.status === 'available' ? '#52c41a' :
                room.status === 'occupied' ? '#ff4d4f' :
                room.status === 'maintenance' ? '#faad14' : '#d9d9d9'
              }`,
          backgroundColor: 
            room.status === 'available' ? '#f6ffed' :
            room.status === 'occupied' ? '#fff1f0' :
            room.status === 'maintenance' ? '#fffbe6' : '#fafafa',
        }}
      >
        <div style={{ textAlign: 'center', width: '100%' }}>
          {/* 床位图标 */}
          <Iconfont 
            name="bed" 
            size={26} 
            color={
              room.status === 'available' ? '#52c41a' :
              room.status === 'occupied' ? '#ff4d4f' :
              room.status === 'maintenance' ? '#faad14' : '#8c8c8c'
            }
          />
          <div style={{ 
            fontSize: '18px', 
            fontWeight: 'bold',
            marginTop: 6,
            color: '#000',
            lineHeight: 1,
          }}>
            {room.room_number}
          </div>
          {/* 调整大小模式下显示尺寸 */}
          {isResizing && (
            <div style={{ 
              marginTop: 6, 
              fontSize: '10px', 
              color: '#1890ff',
              fontWeight: 'bold',
            }}>
              {currentWidth} × {currentHeight}
            </div>
          )}
        </div>

        {/* 调整大小模式下的拖拽手柄 */}
        {isResizing && (
          <div
            onMouseDown={handleResizeMouseDown}
            style={{
              position: 'absolute',
              right: -4,
              bottom: -4,
              width: 16,
              height: 16,
              cursor: 'se-resize',
              background: '#1890ff',
              borderRadius: '50%',
              border: '2px solid #fff',
              boxShadow: '0 2px 4px rgba(0,0,0,0.2)',
              zIndex: 10,
            }}
          />
        )}
      </Card>
      
      {/* 悬浮时显示操作按钮 - 根据位置在卡片上方或下方 */}
      <div
        className="room-card-actions"
        style={{
          position: 'absolute',
          ...(top < 30 
            ? { bottom: -26 }  // 靠近顶部时，工具栏显示在下方
            : { top: -26 }     // 否则显示在上方
          ),
          left: '50%',
          transform: 'translateX(-50%)',
          display: 'none',
          background: 'rgba(255, 255, 255, 0.95)',
          borderRadius: 4,
          padding: '2px 6px',
          boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
          zIndex: 10,
          whiteSpace: 'nowrap',
        }}
      >
        <Space size={0}>
          {onEdit && (
            <Button 
              type="text" 
              size="small" 
              icon={<FormOutlined />} 
              onClick={(e) => {
                e.stopPropagation();
                onEdit(room);
              }}
              style={{ fontSize: 12, padding: '0 4px', height: 22 }}
              title="编辑"
            />
          )}
          {onResizeComplete && (
            <Button 
              type="text" 
              size="small" 
              icon={<GatewayOutlined />} 
              onClick={handleEnterResizeMode}
              style={{ fontSize: 12, padding: '0 4px', height: 22 }}
              title="调整大小"
            />
          )}
          <Popconfirm
            title="确定要删除这个房间吗？"
            onConfirm={() => room.id && onDelete(room.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button 
              type="text" 
              size="small" 
              danger 
              icon={<DeleteOutlined />}
              style={{ fontSize: 12, padding: '0 4px', height: 22 }}
            />
          </Popconfirm>
        </Space>
      </div>
      
      <style>{`
        .room-card-container:hover .room-card-actions {
          display: flex !important;
        }
      `}</style>
    </div>
  );
};

export default DraggableRoomCard;

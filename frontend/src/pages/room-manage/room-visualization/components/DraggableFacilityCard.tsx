import React, { useState, useCallback, useRef, useEffect } from 'react';
import { Card, Button, Popconfirm } from 'antd';
import { DeleteOutlined, GatewayOutlined, RedoOutlined } from '@ant-design/icons';
import { useDrag } from 'react-dnd';
import Iconfont, { IconName } from '@/components/Iconfont';

// 设施类型定义
export const FacilityTypes = {
  ELEVATOR: 'elevator',        // 电梯井
  CORRIDOR: 'corridor',        // 走廊
  LAUNDRY: 'laundry',          // 洗衣房
  STAIRS: 'stairs',            // 楼梯
  BATHROOM: 'bathroom',        // 公共卫生间
  STORAGE: 'storage',          // 储物间
  REST: 'rest',                // 休息区
  RECEPTION: 'reception',      // 前台
  BED: 'bed',                  // 床位区
  MICROWAVE: 'microwave',      // 微波炉/餐饮区
} as const;

export type FacilityType = typeof FacilityTypes[keyof typeof FacilityTypes];

// 设施配置
export const FacilityConfig: Record<FacilityType, {
  name: string;
  icon: IconName;
  color: string;
  bgColor: string;
  defaultWidth: number;
  defaultHeight: number;
}> = {
  elevator: {
    name: '电梯',
    icon: 'elevator',
    color: '#1890ff',
    bgColor: '#e6f7ff',
    defaultWidth: 80,
    defaultHeight: 80,
  },
  corridor: {
    name: '走廊',
    icon: 'corridor',
    color: '#722ed1',
    bgColor: '#f9f0ff',
    defaultWidth: 200,
    defaultHeight: 60,
  },
  laundry: {
    name: '洗衣房',
    icon: 'laundry',
    color: '#13c2c2',
    bgColor: '#e6fffb',
    defaultWidth: 100,
    defaultHeight: 100,
  },
  stairs: {
    name: '楼梯',
    icon: 'stairs',
    color: '#fa8c16',
    bgColor: '#fff7e6',
    defaultWidth: 80,
    defaultHeight: 120,
  },
  bathroom: {
    name: '卫生间',
    icon: 'bathroom',
    color: '#52c41a',
    bgColor: '#f6ffed',
    defaultWidth: 80,
    defaultHeight: 80,
  },
  storage: {
    name: '储物间',
    icon: 'storage',
    color: '#8c8c8c',
    bgColor: '#fafafa',
    defaultWidth: 80,
    defaultHeight: 80,
  },
  rest: {
    name: '休息区',
    icon: 'rest',
    color: '#eb2f96',
    bgColor: '#fff0f6',
    defaultWidth: 160,
    defaultHeight: 120,
  },
  reception: {
    name: '前台',
    icon: 'reception',
    color: '#faad14',
    bgColor: '#fffbe6',
    defaultWidth: 140,
    defaultHeight: 80,
  },
  bed: {
    name: '床位区',
    icon: 'bed',
    color: '#597ef7',
    bgColor: '#f0f5ff',
    defaultWidth: 100,
    defaultHeight: 80,
  },
  microwave: {
    name: '餐饮区',
    icon: 'microwave',
    color: '#f5222d',
    bgColor: '#fff1f0',
    defaultWidth: 100,
    defaultHeight: 100,
  },
};

export const ItemTypes = {
  ROOM_CARD: 'room_card',
  FACILITY_CARD: 'facility_card',
};

// 设施数据接口
export interface Facility {
  id: string | number;  // 支持临时字符串 id（新建未保存）和数字 id（已保存）
  type: FacilityType;
  floor: number;
  left: number;
  top: number;
  width: number;
  height: number;
  rotation?: number;
  label?: string;
}

interface DraggableFacilityCardProps {
  facility: Facility;
  onDelete: (id: string | number) => void;
  onDrop: (id: string | number, left: number, top: number) => void;
  onResizeComplete?: (id: string | number, newWidth: number, newHeight: number, facilityType: FacilityType, floor: number) => void;
  onRotate?: (id: string | number) => void;
}

// 网格大小
const GRID_SIZE = 20;
const MIN_SIZE = 40;
const MAX_SIZE = 400;

// 对齐到网格
const snapToGrid = (value: number): number => {
  return Math.round(value / GRID_SIZE) * GRID_SIZE;
};

const DraggableFacilityCard: React.FC<DraggableFacilityCardProps> = ({
  facility,
  onDelete,
  onDrop,
  onResizeComplete,
  onRotate,
}) => {
  const config = FacilityConfig[facility.type];
  const [isResizing, setIsResizing] = useState(false);
  const [resizeWidth, setResizeWidth] = useState(facility.width);
  const [resizeHeight, setResizeHeight] = useState(facility.height);
  const [isDraggingResize, setIsDraggingResize] = useState(false);
  const startPosRef = useRef({ x: 0, y: 0, width: 0, height: 0 });

  // 同步外部尺寸变化
  useEffect(() => {
    if (!isResizing) {
      setResizeWidth(facility.width);
      setResizeHeight(facility.height);
    }
  }, [facility.width, facility.height, isResizing]);

  const [{ isDragging }, drag, preview] = useDrag(
    () => ({
      type: ItemTypes.FACILITY_CARD,
      item: {
        id: facility.id,
        left: facility.left,
        top: facility.top,
        type: 'facility',
        facility,
      },
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
    [facility, onDrop, isResizing],
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

      if (finalWidth !== facility.width || finalHeight !== facility.height) {
        // 调用回调触发弹窗确认
        if (onResizeComplete) {
          onResizeComplete(facility.id, finalWidth, finalHeight, facility.type, facility.floor);
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
  }, [isDraggingResize, resizeWidth, resizeHeight, facility.width, facility.height, facility.id, facility.type, facility.floor, onResizeComplete]);

  // 进入调整大小模式
  const handleEnterResizeMode = (e: React.MouseEvent) => {
    e.stopPropagation();
    setIsResizing(true);
    setResizeWidth(facility.width);
    setResizeHeight(facility.height);
  };

  // 计算图标大小
  const currentWidth = isResizing ? resizeWidth : facility.width;
  const currentHeight = isResizing ? resizeHeight : facility.height;
  const iconSize = Math.min(currentWidth, currentHeight) > 80 ? 28 : 22;

  return (
    <div
      ref={isResizing ? undefined : drag as any}
      style={{
        position: 'absolute',
        left: facility.left,
        top: facility.top,
        width: currentWidth,
        height: currentHeight,
        opacity: isDragging ? 0.3 : 1,
        cursor: isResizing ? 'default' : 'move',
      }}
    >
      <Card
        hoverable={!isResizing}
        size="small"
        bodyStyle={{
          padding: '4px',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
          height: '100%',
          overflow: 'hidden',
        }}
        style={{
          height: '100%',
          border: isResizing
            ? `2px solid ${config.color}`
            : `2px dashed ${config.color}`,
          backgroundColor: config.bgColor,
          borderRadius: 4,
          boxShadow: isResizing
            ? `0 0 0 2px ${config.color}40, 0 8px 16px ${config.color}30`
            : undefined,
        }}
      >
        <div style={{
          textAlign: 'center',
          width: '100%',
          overflow: 'hidden',
        }}>
          <div style={{ marginBottom: 2 }}>
            <Iconfont
              name={config.icon}
              size={iconSize}
              color={config.color}
            />
          </div>
          <div style={{
            fontSize: '11px',
            fontWeight: 'bold',
            color: config.color,
            whiteSpace: 'nowrap',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
          }}>
            {facility.label || config.name}
          </div>
          {/* 调整大小模式下显示尺寸 */}
          {isResizing && (
            <div style={{
              marginTop: 2,
              fontSize: '10px',
              color: config.color,
              fontWeight: 'bold',
            }}>
              {currentWidth} × {currentHeight}
            </div>
          )}
        </div>

        {/* 悬浮操作按钮 */}
        <div
          className="facility-card-actions"
          style={{
            position: 'absolute',
            top: 2,
            right: 2,
            display: 'none',
            gap: 2,
          }}
        >
          {onResizeComplete && (
            <Button
              type="text"
              size="small"
              icon={<GatewayOutlined />}
              onClick={handleEnterResizeMode}
              style={{ fontSize: 10, padding: 2, minWidth: 20, height: 20 }}
              title="调整大小"
            />
          )}
          {onRotate && (
            <Button
              type="text"
              size="small"
              icon={<RedoOutlined />}
              onClick={(e) => {
                e.stopPropagation();
                onRotate(facility.id);
              }}
              style={{ fontSize: 10, padding: 2, minWidth: 20, height: 20 }}
              title="旋转90°"
            />
          )}
          <Popconfirm
            title="确定要删除这个设施吗？"
            onConfirm={() => onDelete(facility.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="text"
              size="small"
              danger
              icon={<DeleteOutlined />}
              onClick={(e) => e.stopPropagation()}
              style={{ fontSize: 10, padding: 2, minWidth: 20, height: 20 }}
            />
          </Popconfirm>
        </div>

        {/* 调整大小模式下的拖拽手柄 */}
        {isResizing && (
          <div
            onMouseDown={handleResizeMouseDown}
            style={{
              position: 'absolute',
              right: -4,
              bottom: -4,
              width: 14,
              height: 14,
              cursor: 'se-resize',
              background: config.color,
              borderRadius: '50%',
              border: '2px solid #fff',
              boxShadow: '0 2px 4px rgba(0,0,0,0.2)',
              zIndex: 10,
            }}
          />
        )}
      </Card>

      <style>{`
        .ant-card:hover .facility-card-actions {
          display: flex !important;
        }
      `}</style>
    </div>
  );
};

export default DraggableFacilityCard;

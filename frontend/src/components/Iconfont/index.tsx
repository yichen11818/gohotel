import React from 'react';

/**
 * Iconfont 图标映射表
 * 来源: https://www.iconfont.cn/ Project id 5077084
 */
export const IconMap = {
  // 设施图标
  'rest': '\ue663',        // 休息页-休息一下
  'reception': '\ue628',   // 前台
  'storage': '\ue632',     // 箱子/储物间
  'stairs': '\ue60d',      // 楼梯
  'bathroom': '\ue662',    // 卫生间
  'corridor': '\ue600',    // 过道/走廊
  'bed': '\ue623',         // 床
  'laundry': '\ue625',     // 洗衣房
  'microwave': '\ue626',   // 微波炉
  'elevator': '\ue627',    // 电梯
} as const;

export type IconName = keyof typeof IconMap;

interface IconfontProps {
  name: IconName;
  size?: number;
  color?: string;
  style?: React.CSSProperties;
  className?: string;
  onClick?: () => void;
}

/**
 * Iconfont 图标组件
 * 
 * @example
 * <Iconfont name="elevator" size={24} color="#1890ff" />
 */
const Iconfont: React.FC<IconfontProps> = ({
  name,
  size = 16,
  color,
  style,
  className,
  onClick,
}) => {
  return (
    <i
      className={`iconfont ${className || ''}`}
      style={{
        fontSize: size,
        color,
        ...style,
      }}
      onClick={onClick}
    >
      {IconMap[name]}
    </i>
  );
};

export default Iconfont;





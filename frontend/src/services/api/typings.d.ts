declare namespace API {
  type AddUserRequest = {
    email: string;
    is_admin?: boolean;
    phone?: string;
    real_name?: string;
    username: string;
  };

  type Booking = {
    /** 预订单号（唯一） */
    booking_number?: number;
    /** 取消原因 */
    cancel_reason?: string;
    /** 入住日期（有索引） */
    check_in?: string;
    /** 退房日期（有索引） */
    check_out?: string;
    /** 创建时间 */
    created_at?: string;
    /** 入住人身份证号 */
    guest_id_card?: string;
    /** 入住人姓名 */
    guest_name?: string;
    /** 入住人电话 */
    guest_phone?: string;
    /** 主键 */
    id?: number;
    /** 支付方式：wechat, alipay, card */
    payment_method?: string;
    /** 支付状态：unpaid, paid, refunded */
    payment_status?: string;
    /** 关联的房间 */
    room?: Room;
    /** 房间 ID（有索引） */
    room_id?: number;
    /** 特殊要求 */
    special_request?: string;
    /** 状态：pending, confirmed, checkin, checkout, cancelled */
    status?: string;
    /** 总天数 */
    total_days?: number;
    /** 总价 */
    total_price?: number;
    /** 更新时间 */
    updated_at?: string;
    /** 关联查询（可选）
当查询 Booking 时，可以同时加载 User 和 Room 的信息 */
    user?: User;
    /** 用户 ID（有索引） */
    user_id?: number;
  };

  type CreateBookingRequest = {
    /** 格式: "2024-01-01" */
    check_in: string;
    /** 格式: "2024-01-05" */
    check_out: string;
    guest_id_card?: string;
    guest_name: string;
    guest_phone: string;
    room_id: number;
    special_request?: string;
  };

  type CreateRoomRequest = {
    area?: number;
    bed_type?: string;
    capacity: number;
    description?: string;
    facilities?: string;
    floor: number;
    images?: string;
    original_price?: number;
    price: number;
    room_number: string;
    room_type: string;
  };

  type ErrorInfo = {
    code?: string;
    message?: string;
  };

  type ErrorResponse = {
    error?: ErrorInfo;
    success?: boolean;
  };

  type getAdminBookingsParams = {
    /** 页码 */
    page?: number;
    /** 每页数量 */
    page_size?: number;
  };

  type getAdminUsersIdParams = {
    /** 用户 ID */
    id: number;
  };

  type getAdminUsersParams = {
    /** 页码 */
    page?: number;
    /** 每页数量 */
    page_size?: number;
    /** 用户名 */
    username?: string;
    /** 邮箱 */
    email?: string;
    /** 手机号 */
    phone?: string;
    /** 真实姓名 */
    real_name?: string;
    /** 角色 */
    role?: string;
    /** 状态 */
    status?: string;
  };

  type getBookingsIdParams = {
    /** 预订 ID */
    id: number;
  };

  type getBookingsMyParams = {
    /** 页码 */
    page?: number;
    /** 每页数量 */
    page_size?: number;
  };

  type getRoomsAvailableParams = {
    /** 页码 */
    page?: number;
    /** 每页数量 */
    page_size?: number;
  };

  type getRoomsFloorFloorParams = {
    /** 楼层号 */
    floor: number;
    /** 页码 */
    page?: number;
    /** 每页数量 */
    page_size?: number;
  };

  type getRoomsIdParams = {
    /** 房间 ID */
    id: number;
  };

  type getRoomsParams = {
    /** 页码 */
    page?: number;
    /** 每页数量 */
    page_size?: number;
  };

  type getRoomsSearchTypeParams = {
    /** 房型 */
    type: string;
    /** 页码 */
    page?: number;
    /** 每页数量 */
    page_size?: number;
  };

  type LoginRequest = {
    password: string;
    username: string;
  };

  type LoginResponse = {
    token?: string;
    user?: User;
  };

  type postBookingsIdCancelParams = {
    /** 预订 ID */
    id: number;
  };

  type postBookingsIdCheckinParams = {
    /** 预订 ID */
    id: number;
  };

  type postBookingsIdCheckoutParams = {
    /** 预订 ID */
    id: number;
  };

  type postBookingsIdConfirmParams = {
    /** 预订 ID */
    id: number;
  };

  type postRoomsId_openAPI_deleteParams = {
    /** 房间 ID */
    id: number;
  };

  type postRoomsIdParams = {
    /** 房间 ID */
    id: number;
  };

  type RegisterRequest = {
    email: string;
    password: string;
    phone?: string;
    real_name?: string;
    username: string;
  };

  type Room = {
    /** 面积（平方米） */
    area?: number;
    /** 床型：单人床、双人床、大床 */
    bed_type?: string;
    /** 可住人数 */
    capacity?: number;
    /** 创建时间 */
    created_at?: string;
    /** 房间描述 */
    description?: string;
    /** 设施（JSON 字符串） */
    facilities?: string;
    /** 楼层 */
    floor?: number;
    /** 主键 */
    id?: number;
    /** 图片 URL（JSON 数组） */
    images?: string;
    /** 原价 */
    original_price?: number;
    /** 价格（每晚） */
    price?: number;
    /** 房间号（唯一，有索引） */
    room_number?: string;
    /** 房间类型（有索引） */
    room_type?: string;
    /** 状态：available, occupied, maintenance */
    status?: string;
    /** 更新时间 */
    updated_at?: string;
  };

  type UpdateRoomRequest = {
    area?: number;
    bed_type?: string;
    capacity?: number;
    description?: string;
    facilities?: string;
    floor?: number;
    images?: string;
    original_price?: number;
    price?: number;
    room_type?: string;
    status?: string;
  };

  type User = {
    /** 头像 URL */
    avatar?: string;
    /** 创建时间 */
    created_at?: string;
    /** 邮箱（唯一） */
    email?: string;
    /** 是否首次登录 */
    first_login?: boolean;
    /** 主键（使用雪花算法生成） */
    id?: number;
    /** 手机号（唯一，可为空） */
    phone?: string;
    /** 真实姓名 */
    real_name?: string;
    /** 角色：user, admin */
    role?: string;
    status?: string;
    /** 更新时间 */
    updated_at?: string;
    /** 用户名（唯一） */
    username?: string;
  };
}

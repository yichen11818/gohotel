package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// JSONInt64 自定义类型，用于将 int64 序列化为字符串
// 解决 JavaScript 中大整数精度丢失的问题
type JSONInt64 int64

// MarshalJSON 实现 json.Marshaler 接口
// 将 int64 序列化为字符串
func (i JSONInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(i), 10))
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
// 支持从字符串或数字反序列化
func (i *JSONInt64) UnmarshalJSON(data []byte) error {
	// 尝试解析为字符串
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		*i = JSONInt64(val)
		return nil
	}

	// 尝试解析为数字
	var num int64
	if err := json.Unmarshal(data, &num); err != nil {
		return err
	}
	*i = JSONInt64(num)
	return nil
}

// Value 实现 driver.Valuer 接口，用于数据库写入
func (i JSONInt64) Value() (driver.Value, error) {
	return int64(i), nil
}

// Scan 实现 sql.Scanner 接口，用于数据库读取
func (i *JSONInt64) Scan(value interface{}) error {
	if value == nil {
		*i = 0
		return nil
	}

	switch v := value.(type) {
	case int64:
		*i = JSONInt64(v)
		return nil
	case []byte:
		val, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return err
		}
		*i = JSONInt64(val)
		return nil
	case string:
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		*i = JSONInt64(val)
		return nil
	default:
		return fmt.Errorf("无法将 %T 类型转换为 JSONInt64", value)
	}
}

// Int64 返回 int64 值
func (i JSONInt64) Int64() int64 {
	return int64(i)
}

// String 返回字符串值
func (i JSONInt64) String() string {
	return strconv.FormatInt(int64(i), 10)
}




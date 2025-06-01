package pkkit

import "fmt"

var DefaultSnowflake *Snowflake

func init() {
	DefaultSnowflake, _ = NewSnowflake(1, 1)
}
func UseSnowflakeID() string {
	id, _ := DefaultSnowflake.NextID()
	return fmt.Sprintf("%d", id)
}
func UseSnowflakeIntID() int64 {
	id, _ := DefaultSnowflake.NextID()
	return id
}

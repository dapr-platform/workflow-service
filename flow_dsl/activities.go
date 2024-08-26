package flow_dsl

type FlowDslActivities struct {
}

// 设备
// 事件状态或状态更新
var Type_DeviceEventOrStateChange = "DeviceEventOrStateChange"

// 查询当前状态
var Type_DeviceStateQuery = "DeviceStateQuery"

// 执行操作
var Type_DeviceAction = "DeviceAction"

// 设备采集
var Type_DeviceDataCollection = "DeviceDataCollection"

// 时间
// 定时
var Type_TimeFixed = "TimeFixed"

// 时间段
var Type_TimeSpan = "TimeSpan"

// 延时
var Type_TimeDelay = "TimeDelay"

// 状态维持了一段时间
var Type_TimeStateKeeping = "TimeStateKeeping"

// 流程控制
// 当-如果-就
var Type_FlowWhenIfThen = "FlowWhenIfThen"

// 循环
var Type_FlowLoop = "FlowLoop"

// 最多触发次数
var Type_FlowMaxCountTrigger = "FlowMaxCountTrigger"

// 达到指定次数时
var Type_FlowFixedCountTrigger = "FlowFixedCountTrigger"

// 逻辑
// 当任意事件发生
var Type_LogicAnyEventOccured = "LogicAnyEventOccured"

// 满足任意条件
var Type_LogicAnyConditionMet = "LogicAnyConditionMet"

// 满足全部条件
var Type_LogicFullConditionMet = "LogicFullConditionMet"

// 其他
// 本规则启用时
var Type_OtherRuleEnable = "OtherRuleEnable"

// 事件先后发生
var Type_OtherEventsOccuredSequence = "OtherEventsOccuredSequence"

// 自定义状态
var Type_OtherCustomState = "OtherCustomState"

// 模式切换
var Type_OtherModSwitch = "OtherModSwitch"

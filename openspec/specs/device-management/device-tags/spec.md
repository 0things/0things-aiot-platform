# device-management/device-tags Specification

## Purpose
让用户在设备详情页查看并编辑（新增、删除）设备的自定义标签，复用已存在的 `device_tags` 后端接口。

## Requirements

### Requirement: 设备详情展示自定义标签
系统 SHALL 在设备详情页的「标签信息」区域展示该设备的全部自定义标签（key=value）。

#### Scenario: 设备有标签时展示列表
- **WHEN** 用户打开某设备的详情页且设备存在自定义标签
- **THEN** 「标签信息」区域以标签形式展示每个标签的 key=value

#### Scenario: 设备无标签时展示空态
- **WHEN** 用户打开某设备的详情页且设备无自定义标签
- **THEN** 「标签信息」区域展示无标签的空态文案

### Requirement: 设备详情编辑自定义标签
系统 SHALL 允许用户在设备详情页新增与删除自定义标签，操作后即时反映到列表。

#### Scenario: 新增标签
- **WHEN** 用户输入 key（及可选 value）并提交
- **THEN** 系统调用后端接口新增该标签并刷新列表

#### Scenario: 删除标签
- **WHEN** 用户对已存在的某个标签执行删除
- **THEN** 系统调用后端接口删除该标签并刷新列表

#### Scenario: 拒绝纯数字键
- **WHEN** 用户输入的 key 为纯数字（仅由 0-9 组成）
- **THEN** 系统拒绝该标签并返回错误提示，不写入后端

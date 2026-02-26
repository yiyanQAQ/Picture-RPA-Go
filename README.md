# Picture RPA

[中文](#中文) | [English](#english)

---

## 中文

### 简介
基于图像识别的轻量级RPA工具 (Go 语言版)

### 功能
- **内置截图工具**
- **图像识别**
- **逻辑控制**
- **流程编辑**

### 目录结构
- `main.go`: 程序入口
- `gui/`: UI与事件逻辑
- `automation/`: 核心自动化逻辑
- `constants/`: 配置与常量
- `snips/`: 截图保存目录

### 运行
1. **克隆**:
   ```bash
   git clone https://github.com/yiyanQAQ/Picture-RPA-Go.git
   ```
2. **构建**:
   ```bash
   go build -o picture-rpa main.go
   ```
3. **启动程序**:
   ```bash
   ./picture-rpa
   ```

### 许可协议
本项目采用 [GPLv3](LICENSE) 开源协议。

---

## English

### Introduction
A lightweight RPA tool based on image recognition (Go version).

### Features
- **Built-in Snipping Tool**
- **Image Recognition**
- **Logical Flow Control**
- **Workflow Editing**

### Directory Structure 
- `main.go`: Entry point
- `gui/`: UI & logic
- `automation/`: Core automation logic
- `constants/`: Config & constants
- `snips/`: Screenshot directory

### Run
1. **Clone**:
   ```bash
   git clone https://github.com/yiyanQAQ/Picture-RPA-Go.git
   ```

2. **Build**:
   ```bash
   go build -o picture-rpa main.go
   ```
3. **Run**:
   ```bash
   ./picture-rpa
   ```
   
### License
This project is licensed under the [GPLv3](LICENSE) License.

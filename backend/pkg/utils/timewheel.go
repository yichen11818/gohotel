package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"gohotel/pkg/logger"

	"go.uber.org/zap"
)

// Task 表示一个时间轮任务
type Task struct {
	ID        string                 // 任务ID
	ExecTime  time.Time              // 执行时间
	Callback  func()                 // 任务执行函数
	Cancelled bool                   // 任务是否被取消
	Meta      map[string]interface{} // 任务元数据，用于存储业务相关信息
}

// PersistTask 表示持久化的任务信息
type PersistTask struct {
	ID       string                 `json:"id"`        // 任务ID
	ExecTime time.Time              `json:"exec_time"` // 执行时间
	Meta     map[string]interface{} `json:"meta"`      // 任务元数据
}

// ToPersistTask 将Task转换为PersistTask
func (t *Task) ToPersistTask() *PersistTask {
	return &PersistTask{
		ID:       t.ID,
		ExecTime: t.ExecTime,
		Meta:     t.Meta,
	}
}

// PersistStore 定义持久化存储接口
type PersistStore interface {
	SaveTasks(tasks []*PersistTask) error
	LoadTasks() ([]*PersistTask, error)
}

// TaskExecutor 定义任务执行器接口
type TaskExecutor interface {
	Execute(task *Task)  // 执行任务
	GetTaskType() string // 获取任务类型
}

// TimeSlot 表示时间轮中的一个时间槽
type TimeSlot struct {
	Tasks map[string]*Task // 存储该时间槽的任务列表
	Mutex sync.RWMutex     // 保护任务列表的互斥锁
}

// NewTimeSlot 创建一个新的时间槽
func NewTimeSlot() *TimeSlot {
	return &TimeSlot{
		Tasks: make(map[string]*Task),
	}
}

// TimeWheel 表示时间轮
type TimeWheel struct {
	Interval      time.Duration           // 时间槽间隔
	Slots         int                     // 时间槽数量
	SlotArray     []*TimeSlot             // 时间槽数组
	CurrentSlot   int                     // 当前时间槽指针
	Mutex         sync.RWMutex            // 保护时间轮的互斥锁
	Timer         *time.Ticker            // 定时器，用于推动时间轮
	StopCh        chan struct{}           // 停止信号通道
	TaskIDMap     map[string]bool         // 记录任务ID是否存在
	TaskIDMutex   sync.RWMutex            // 保护TaskIDMap的互斥锁
	PersistStore  PersistStore            // 持久化存储实例
	Level         int                     // 层级标识（1:秒轮, 2:分轮, 3:时轮, 4:天轮）
	ParentWheel   *TimeWheel              // 父时间轮（上层）
	ChildWheel    *TimeWheel              // 子时间轮（下层）
	MultiWheel    *MultiTimeWheel         // 多层时间轮实例引用
	ExecutorMap   map[string]TaskExecutor // 任务执行器映射，key: taskType, value: TaskExecutor
	ExecutorMutex sync.RWMutex            // 保护ExecutorMap的互斥锁
}

// MultiTimeWheel 表示多层时间轮
type MultiTimeWheel struct {
	SecondWheel      *TimeWheel              // 秒轮（60格，1秒间隔）
	MinuteWheel      *TimeWheel              // 分轮（60格，1分钟间隔）
	HourWheel        *TimeWheel              // 时轮（24格，1小时间隔）
	DayWheel         *TimeWheel              // 天轮（365格，1天间隔）
	PersistStore     PersistStore            // 持久化存储实例
	StopCh           chan struct{}           // 停止信号通道
	wg               sync.WaitGroup          // 等待组，用于优雅关闭
	TaskIDMap        map[string]bool         // 记录所有任务ID
	TaskIDMutex      sync.RWMutex            // 保护TaskIDMap的互斥锁
	SnapshotInterval time.Duration           // 快照间隔时间
	ExecutorMap      map[string]TaskExecutor // 任务执行器映射
	ExecutorMutex    sync.RWMutex            // 保护ExecutorMap的互斥锁
}

// NewTimeWheel 创建一个新的时间轮实例
func NewTimeWheel(interval time.Duration, slots int) *TimeWheel {
	if interval <= 0 {
		interval = time.Second
	}
	if slots <= 0 {
		slots = 60
	}

	// 基于实际时间计算初始CurrentSlot，确保与实际时间同步
	currentTime := time.Now()
	currentUnix := currentTime.Unix()
	intervalSec := int64(interval.Seconds())
	// 计算初始槽位
	initialSlot := int(currentUnix/intervalSec) % slots

	tw := &TimeWheel{
		Interval:    interval,
		Slots:       slots,
		SlotArray:   make([]*TimeSlot, slots),
		CurrentSlot: initialSlot,
		StopCh:      make(chan struct{}),
		TaskIDMap:   make(map[string]bool),
		ExecutorMap: make(map[string]TaskExecutor),
	}

	// 初始化时间槽
	for i := range tw.SlotArray {
		tw.SlotArray[i] = NewTimeSlot()
	}

	return tw
}

// NewMultiTimeWheel 创建一个新的多层时间轮实例
func NewMultiTimeWheel() *MultiTimeWheel {
	mtw := &MultiTimeWheel{
		StopCh:           make(chan struct{}),
		TaskIDMap:        make(map[string]bool),
		SnapshotInterval: 30 * time.Second, // 默认30秒快照间隔
		ExecutorMap:      make(map[string]TaskExecutor),
	}

	// 创建各层时间轮
	secondWheel := NewTimeWheel(time.Second, 60)
	secondWheel.Level = 1
	secondWheel.MultiWheel = mtw

	minuteWheel := NewTimeWheel(time.Minute, 60)
	minuteWheel.Level = 2
	minuteWheel.MultiWheel = mtw

	hourWheel := NewTimeWheel(time.Hour, 24)
	hourWheel.Level = 3
	hourWheel.MultiWheel = mtw

	dayWheel := NewTimeWheel(24*time.Hour, 365)
	dayWheel.Level = 4
	dayWheel.MultiWheel = mtw

	// 建立层级关系
	dayWheel.ChildWheel = hourWheel
	hourWheel.ParentWheel = dayWheel
	hourWheel.ChildWheel = minuteWheel
	minuteWheel.ParentWheel = hourWheel
	minuteWheel.ChildWheel = secondWheel
	secondWheel.ParentWheel = minuteWheel

	// 设置到多层时间轮实例
	mtw.SecondWheel = secondWheel
	mtw.MinuteWheel = minuteWheel
	mtw.HourWheel = hourWheel
	mtw.DayWheel = dayWheel

	return mtw
}

// SetPersistStore 设置持久化存储实例（多层时间轮）
func (mtw *MultiTimeWheel) SetPersistStore(store PersistStore) {
	mtw.PersistStore = store
	// 为所有层级设置持久化存储
	mtw.DayWheel.PersistStore = store
	mtw.HourWheel.PersistStore = store
	mtw.MinuteWheel.PersistStore = store
	mtw.SecondWheel.PersistStore = store
}

// SetSnapshotInterval 设置快照间隔时间（多层时间轮）
func (mtw *MultiTimeWheel) SetSnapshotInterval(interval time.Duration) {
	if interval > 0 {
		mtw.SnapshotInterval = interval
	}
}

// RegisterExecutor 注册任务执行器（多层时间轮）
func (mtw *MultiTimeWheel) RegisterExecutor(executor TaskExecutor) {
	mtw.ExecutorMutex.Lock()
	defer mtw.ExecutorMutex.Unlock()
	mtw.ExecutorMap[executor.GetTaskType()] = executor
}

// GetExecutor 获取任务执行器（多层时间轮）
func (mtw *MultiTimeWheel) GetExecutor(taskType string) TaskExecutor {
	mtw.ExecutorMutex.RLock()
	defer mtw.ExecutorMutex.RUnlock()
	return mtw.ExecutorMap[taskType]
}

// AddTask 添加一个定时任务（多层时间轮）
func (mtw *MultiTimeWheel) AddTask(execTime time.Time, callback func(), meta map[string]interface{}) string {
	// 生成唯一任务ID
	taskID := fmt.Sprintf("task_%d", time.Now().UnixNano())

	// 创建任务
	task := &Task{
		ID:       taskID,
		ExecTime: execTime,
		Callback: callback,
		Meta:     meta,
	}

	// 计算延迟时间
	currentTime := time.Now()
	delay := execTime.Sub(currentTime)
	if delay < 0 {
		// 如果任务已经过期，立即执行
		go func() {
			callback()
			logger.Info("任务已过期，立即执行完成",
				zap.String("task_id", taskID),
				zap.Time("exec_time", execTime),
				zap.Any("meta", meta),
			)
		}()
		return taskID
	}

	// 根据延迟时间选择合适的层级
	var targetWheel *TimeWheel
	switch {
	case delay <= 60*time.Second:
		targetWheel = mtw.SecondWheel
	case delay < 60*time.Minute:
		targetWheel = mtw.MinuteWheel
	case delay < 24*time.Hour:
		targetWheel = mtw.HourWheel
	default:
		targetWheel = mtw.DayWheel
	}

	// 计算目标层级的精确时间槽索引
	targetWheel.Mutex.RLock()
	// 直接计算任务执行时间对应的槽位，简化计算逻辑
	taskUnix := execTime.Unix()
	intervalSec := int64(targetWheel.Interval.Seconds())
	// 任务执行时间在目标时间轮上的槽位
	index := int(taskUnix/intervalSec) % targetWheel.Slots
	targetWheel.Mutex.RUnlock()

	// 添加任务到对应时间槽
	targetWheel.SlotArray[index].Mutex.Lock()
	targetWheel.SlotArray[index].Tasks[taskID] = task
	targetWheel.SlotArray[index].Mutex.Unlock()

	// 记录任务ID
	targetWheel.TaskIDMutex.Lock()
	targetWheel.TaskIDMap[taskID] = true
	targetWheel.TaskIDMutex.Unlock()

	mtw.TaskIDMutex.Lock()
	mtw.TaskIDMap[taskID] = true
	mtw.TaskIDMutex.Unlock()

	// 持久化更新
	mtw.saveTasks()

	return taskID
}

// AddDelayTask 添加一个延迟任务（多层时间轮）
func (mtw *MultiTimeWheel) AddDelayTask(delay time.Duration, callback func(), meta map[string]interface{}) string {
	return mtw.AddTask(time.Now().Add(delay), callback, meta)
}

// RemoveTask 删除一个任务（多层时间轮）
func (mtw *MultiTimeWheel) RemoveTask(taskID string) bool {
	// 检查任务是否存在
	mtw.TaskIDMutex.RLock()
	_, exists := mtw.TaskIDMap[taskID]
	mtw.TaskIDMutex.RUnlock()
	if !exists {
		return false
	}

	// 遍历所有层级查找并删除任务
	wheels := []*TimeWheel{mtw.SecondWheel, mtw.MinuteWheel, mtw.HourWheel, mtw.DayWheel}
	for _, wheel := range wheels {
		// 遍历当前层级的所有时间槽
		for _, slot := range wheel.SlotArray {
			slot.Mutex.Lock()
			if _, ok := slot.Tasks[taskID]; ok {
				delete(slot.Tasks, taskID)
				slot.Mutex.Unlock()
				// 更新当前层级的任务ID映射
				wheel.TaskIDMutex.Lock()
				delete(wheel.TaskIDMap, taskID)
				wheel.TaskIDMutex.Unlock()
				// 更新多层时间轮的任务ID映射
				mtw.TaskIDMutex.Lock()
				delete(mtw.TaskIDMap, taskID)
				mtw.TaskIDMutex.Unlock()
				// 持久化更新
				mtw.saveTasks()
				return true
			}
			slot.Mutex.Unlock()
		}
	}

	return false
}

// LoadTasks 从持久化存储加载任务（多层时间轮）
func (mtw *MultiTimeWheel) LoadTasks() error {
	if mtw.PersistStore == nil {
		return nil
	}

	// 从持久化存储加载任务
	persistTasks, err := mtw.PersistStore.LoadTasks()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	currentTime := time.Now()
	restoredCount := 0

	// 恢复每个任务
	for _, pt := range persistTasks {
		// 检查任务是否已经过期
		if pt.ExecTime.Before(currentTime) {
			// 任务已经过期，创建任务并执行
			task := &Task{
				ID:       pt.ID,
				ExecTime: pt.ExecTime,
				Meta:     pt.Meta,
			}

			// 尝试获取任务执行器并执行任务
			executorType, ok := task.Meta["executor_type"].(string)
			if ok {
				executor := mtw.GetExecutor(executorType)
				if executor != nil {
					go func() {
						executor.Execute(task)
						logger.Info("加载的任务已过期，通过执行器执行完成",
							zap.String("task_id", task.ID),
							zap.Time("exec_time", task.ExecTime),
							zap.String("executor_type", executorType),
							zap.Any("meta", task.Meta),
						)
					}()
					restoredCount++
					continue
				}
			}

			continue
		}

		// 创建任务
		task := &Task{
			ID:       pt.ID,
			ExecTime: pt.ExecTime,
			Meta:     pt.Meta,
		}

		// 保存任务ID和元数据，用于回调函数
		taskID := task.ID
		taskMeta := task.Meta
		taskExecTime := task.ExecTime

		// 设置任务回调函数
		task.Callback = func() {
			// 任务执行时，先尝试使用任务执行器
			executorType, ok := taskMeta["executor_type"].(string)
			if ok {
				executor := mtw.GetExecutor(executorType)
				if executor != nil {
					// 重新构建任务对象，包含完整信息
					currentTask := &Task{
						ID:       taskID,
						ExecTime: taskExecTime,
						Meta:     taskMeta,
					}
					executor.Execute(currentTask)
					logger.Info("恢复的任务通过执行器执行完成",
						zap.String("task_id", taskID),
						zap.Time("exec_time", taskExecTime),
						zap.String("executor_type", executorType),
						zap.Any("meta", taskMeta),
					)
				}
			}
		}

		// 计算延迟时间
		delay := pt.ExecTime.Sub(currentTime)

		// 根据延迟时间选择合适的层级
		var targetWheel *TimeWheel
		switch {
		case delay <= 60*time.Second:
			targetWheel = mtw.SecondWheel
		case delay < 60*time.Minute:
			targetWheel = mtw.MinuteWheel
		case delay < 24*time.Hour:
			targetWheel = mtw.HourWheel
		default:
			targetWheel = mtw.DayWheel
		}

		// 计算目标层级的精确时间槽索引
		targetWheel.Mutex.RLock()
		// 直接计算任务执行时间对应的槽位，简化计算逻辑
		taskUnix := pt.ExecTime.Unix()
		intervalSec := int64(targetWheel.Interval.Seconds())
		// 任务执行时间在目标时间轮上的槽位
		index := int(taskUnix/intervalSec) % targetWheel.Slots
		targetWheel.Mutex.RUnlock()

		// 添加任务到对应时间槽
		targetWheel.SlotArray[index].Mutex.Lock()
		targetWheel.SlotArray[index].Tasks[pt.ID] = task
		targetWheel.SlotArray[index].Mutex.Unlock()

		// 记录任务ID到目标层级和多层时间轮
		targetWheel.TaskIDMutex.Lock()
		targetWheel.TaskIDMap[pt.ID] = true
		targetWheel.TaskIDMutex.Unlock()

		mtw.TaskIDMutex.Lock()
		mtw.TaskIDMap[pt.ID] = true
		mtw.TaskIDMutex.Unlock()

		restoredCount++
	}

	fmt.Printf("Restored %d tasks from persist store\n", restoredCount)
	return nil
}

// Start 启动多层时间轮
func (mtw *MultiTimeWheel) Start() {
	// 启动所有层级的时间轮
	mtw.wg.Add(4) // 四个层级的时间轮

	// 启动各层级的时间轮
	for _, wheel := range []*TimeWheel{mtw.SecondWheel, mtw.MinuteWheel, mtw.HourWheel, mtw.DayWheel} {
		// 所有层级都使用基于绝对时间的定时器
		go func(wheel *TimeWheel) {
			defer mtw.wg.Done()
			for {
				// 计算距离下一个触发点的时间
				now := time.Now()
				var next time.Time

				switch wheel.Level {
				case 1: // 秒轮：每秒0毫秒触发
					// 计算下一个整秒
					next = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+1, 0, now.Location())
				case 2: // 分轮：每分钟0秒触发
					next = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+1, 0, 0, now.Location())
				case 3: // 时轮：每小时0分0秒触发
					next = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
				case 4: // 天轮：每天0时0分0秒触发
					next = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
				default:
					// 意外情况，使用普通定时器
					timer := time.NewTicker(wheel.Interval)
					select {
					case <-timer.C:
						timer.Stop()
						wheel.advance()
					case <-mtw.StopCh:
						timer.Stop()
						return
					}
					continue
				}

				// 计算距离下一个触发点的时间
				duration := next.Sub(now)
				if duration < 0 {
					// 理论上不会发生，但为了安全
					duration = wheel.Interval
				}

				// 创建定时器
				timer := time.NewTimer(duration)

				select {
				case <-timer.C:
					// 定时器触发，执行任务
					wheel.advance()
				case <-mtw.StopCh:
					// 停止信号，退出
					timer.Stop()
					return
				}
			}
		}(wheel)
	}

	// 启动定期快照定时器
	if mtw.PersistStore != nil {
		snapshotTimer := time.NewTicker(mtw.SnapshotInterval)
		mtw.wg.Add(1)

		// 快照goroutine
		go func() {
			defer mtw.wg.Done()
			for {
				select {
				case <-snapshotTimer.C:
					mtw.saveTasks()
				case <-mtw.StopCh:
					snapshotTimer.Stop()
					return
				}
			}
		}()
	}
}

// Stop 停止多层时间轮
func (mtw *MultiTimeWheel) Stop() {
	close(mtw.StopCh)
	mtw.wg.Wait()
	// 执行最终快照
	if mtw.PersistStore != nil {
		mtw.saveTasks()
	}
}

// advance 推动时间轮前进一个槽
func (tw *TimeWheel) advance() {
	tw.Mutex.Lock()
	// 基于实际时间计算当前应该在哪个槽位，确保CurrentSlot与实际时间同步
	currentTime := time.Now()
	currentUnix := currentTime.Unix()
	intervalSec := int64(tw.Interval.Seconds())
	// 计算当前时间应该在哪个槽位
	actualSlot := int(currentUnix/intervalSec) % tw.Slots

	// 直接使用实际槽位，不移动指针，确保执行的是当前时间对应的任务
	currentSlot := actualSlot
	tw.Mutex.Unlock()

	// 执行当前槽的所有任务
	slot := tw.SlotArray[currentSlot]
	slot.Mutex.Lock()
	// 创建任务副本，避免在遍历过程中修改映射
	tasks := make(map[string]*Task, len(slot.Tasks))
	for id, task := range slot.Tasks {
		tasks[id] = task
	}
	// 清空当前槽的任务
	clear(slot.Tasks)
	slot.Mutex.Unlock()

	// 如果有任务需要处理
	if len(tasks) > 0 {
		if tw.ChildWheel != nil {
			// 如果有子时间轮，将任务迁移到子时间轮
			tw.migrateTasksToChildWheel(tasks)
		} else {
			// 没有子时间轮，执行所有任务
			// 从任务ID映射中删除所有执行的任务
			for id := range tasks {
				tw.TaskIDMutex.Lock()
				delete(tw.TaskIDMap, id)
				tw.TaskIDMutex.Unlock()
				// 从多层时间轮的映射中删除
				if tw.MultiWheel != nil {
					tw.MultiWheel.TaskIDMutex.Lock()
					delete(tw.MultiWheel.TaskIDMap, id)
					tw.MultiWheel.TaskIDMutex.Unlock()
				}
			}

			// 更新持久化存储
			if tw.MultiWheel != nil {
				tw.MultiWheel.saveTasks()
			}

			// 执行所有任务
			currentTime := time.Now()
			for _, task := range tasks {
				if !task.Cancelled {
					// 检查任务是否到了执行时间
					if currentTime.Before(task.ExecTime) {
						// 任务还没到执行时间，重新添加到正确的时间槽
						// 直接计算任务执行时间对应的槽位，简化计算逻辑
						taskUnix := task.ExecTime.Unix()
						intervalSec := int64(tw.Interval.Seconds())
						// 任务执行时间在目标时间轮上的槽位
						index := int(taskUnix/intervalSec) % tw.Slots
						// 重新添加到时间槽
						tw.SlotArray[index].Mutex.Lock()
						tw.SlotArray[index].Tasks[task.ID] = task
						tw.SlotArray[index].Mutex.Unlock()
						// 恢复任务ID
						tw.TaskIDMutex.Lock()
						tw.TaskIDMap[task.ID] = true
						tw.TaskIDMutex.Unlock()
						tw.MultiWheel.TaskIDMutex.Lock()
						tw.MultiWheel.TaskIDMap[task.ID] = true
						tw.MultiWheel.TaskIDMutex.Unlock()
						continue
					}
					// 尝试获取任务执行器并执行任务
					executorType, ok := task.Meta["executor_type"].(string)
					if ok && tw.MultiWheel != nil {
						executor := tw.MultiWheel.GetExecutor(executorType)
						if executor != nil {
							go func() {
								executor.Execute(task)
								logger.Info("任务通过执行器执行完成",
									zap.String("task_id", task.ID),
									zap.Time("exec_time", task.ExecTime),
									zap.String("executor_type", executorType),
									zap.Any("meta", task.Meta),
								)
							}()
							continue
						}
					}
					// 没有执行器，直接调用回调函数
					go func() {
						task.Callback()
						logger.Info("任务直接执行完成",
							zap.String("task_id", task.ID),
							zap.Time("exec_time", task.ExecTime),
							zap.Any("meta", task.Meta),
						)
					}()
				}
			}
		}
	}
}

// migrateTasksToChildWheel 将任务迁移到子时间轮
func (tw *TimeWheel) migrateTasksToChildWheel(tasks map[string]*Task) {
	childWheel := tw.ChildWheel
	if childWheel == nil {
		return
	}

	// 当前时间
	currentTime := time.Now()

	// 将每个任务迁移到子时间轮
	for _, task := range tasks {
		// 计算任务执行时间与当前时间的差值
		delay := task.ExecTime.Sub(currentTime)
		// 添加调试日志
		logger.Info("开始迁移任务到子时间轮",
			zap.String("task_id", task.ID),
			zap.Time("exec_time", task.ExecTime),
			zap.Time("current_time", currentTime),
			zap.Duration("delay", delay),
			zap.String("parent_wheel_level", fmt.Sprintf("%d", tw.Level)),
			zap.String("child_wheel_level", fmt.Sprintf("%d", childWheel.Level)),
			zap.Any("meta", task.Meta),
		)
		if delay < 0 {
			// 任务已经过期，立即执行
			// 尝试获取任务执行器并执行任务
			executorType, ok := task.Meta["executor_type"].(string)
			if ok && tw.MultiWheel != nil {
				executor := tw.MultiWheel.GetExecutor(executorType)
				if executor != nil {
					go func() {
						executor.Execute(task)
						logger.Info("迁移任务已过期，通过执行器执行完成",
							zap.String("task_id", task.ID),
							zap.Time("exec_time", task.ExecTime),
							zap.String("executor_type", executorType),
							zap.Any("meta", task.Meta),
						)
					}()
					// 从任务ID映射中删除
					tw.TaskIDMutex.Lock()
					delete(tw.TaskIDMap, task.ID)
					tw.TaskIDMutex.Unlock()
					// 从多层时间轮的映射中删除
					if tw.MultiWheel != nil {
						tw.MultiWheel.TaskIDMutex.Lock()
						delete(tw.MultiWheel.TaskIDMap, task.ID)
						tw.MultiWheel.TaskIDMutex.Unlock()
					}
					continue
				}
			}
			// 没有执行器，直接调用回调函数
			go func() {
				task.Callback()
				logger.Info("迁移任务已过期，直接执行完成",
					zap.String("task_id", task.ID),
					zap.Time("exec_time", task.ExecTime),
					zap.Any("meta", task.Meta),
				)
			}()
			// 从任务ID映射中删除
			tw.TaskIDMutex.Lock()
			delete(tw.TaskIDMap, task.ID)
			tw.TaskIDMutex.Unlock()
			// 从多层时间轮的映射中删除
			if tw.MultiWheel != nil {
				tw.MultiWheel.TaskIDMutex.Lock()
				delete(tw.MultiWheel.TaskIDMap, task.ID)
				tw.MultiWheel.TaskIDMutex.Unlock()
			}
			continue
		}

		// 计算在子时间轮中的精确时间槽索引
		childWheel.Mutex.RLock()
		// 直接计算任务执行时间对应的槽位，简化计算逻辑
		taskUnix := task.ExecTime.Unix()
		intervalSec := int64(childWheel.Interval.Seconds())
		// 任务执行时间在子时间轮上的槽位
		index := int(taskUnix/intervalSec) % childWheel.Slots
		childWheel.Mutex.RUnlock()

		// 添加调试日志
		logger.Info("迁移任务到子时间轮槽位",
			zap.String("task_id", task.ID),
			zap.Time("exec_time", task.ExecTime),
			zap.String("child_wheel_interval", fmt.Sprintf("%v", childWheel.Interval)),
			zap.Int("child_wheel_slots", childWheel.Slots),
			zap.Int("target_slot", index),
		)

		// 添加任务到子时间轮
		childWheel.SlotArray[index].Mutex.Lock()
		childWheel.SlotArray[index].Tasks[task.ID] = task
		childWheel.SlotArray[index].Mutex.Unlock()

		// 更新子时间轮的任务ID映射
		childWheel.TaskIDMutex.Lock()
		childWheel.TaskIDMap[task.ID] = true
		childWheel.TaskIDMutex.Unlock()

		// 从当前时间轮的任务ID映射中删除
		tw.TaskIDMutex.Lock()
		delete(tw.TaskIDMap, task.ID)
		tw.TaskIDMutex.Unlock()
	}

	// 更新持久化存储
	if tw.MultiWheel != nil {
		tw.MultiWheel.saveTasks()
	}
}

// saveTasks 保存所有任务到持久化存储（多层时间轮）
func (mtw *MultiTimeWheel) saveTasks() {
	if mtw.PersistStore == nil {
		return
	}

	// 收集所有未执行的任务
	var persistTasks []*PersistTask

	// 遍历所有层级的时间轮
	wheels := []*TimeWheel{mtw.SecondWheel, mtw.MinuteWheel, mtw.HourWheel, mtw.DayWheel}
	for _, wheel := range wheels {
		for _, slot := range wheel.SlotArray {
			slot.Mutex.RLock()
			for _, task := range slot.Tasks {
				if !task.Cancelled {
					persistTasks = append(persistTasks, task.ToPersistTask())
				}
			}
			slot.Mutex.RUnlock()
		}
	}

	// 保存任务到持久化存储
	if err := mtw.PersistStore.SaveTasks(persistTasks); err != nil {
		fmt.Printf("Failed to save tasks: %v\n", err)
	}
}

// FilePersistStore 实现基于文件系统的持久化存储
type FilePersistStore struct {
	FilePath string     // 持久化文件路径
	Mutex    sync.Mutex // 保护文件读写的互斥锁
}

// NewFilePersistStore 创建一个新的文件持久化存储实例
func NewFilePersistStore(filePath string) *FilePersistStore {
	return &FilePersistStore{
		FilePath: filePath,
	}
}

// SaveTasks 保存任务到文件
func (fs *FilePersistStore) SaveTasks(tasks []*PersistTask) error {
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()

	// 将任务转换为JSON格式
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	// 写入文件
	err = os.WriteFile(fs.FilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write tasks to file: %w", err)
	}

	return nil
}

// LoadTasks 从文件加载任务
func (fs *FilePersistStore) LoadTasks() ([]*PersistTask, error) {
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()

	// 检查文件是否存在
	if _, err := os.Stat(fs.FilePath); os.IsNotExist(err) {
		// 文件不存在，返回空任务列表
		return []*PersistTask{}, nil
	}

	// 读取文件内容
	data, err := os.ReadFile(fs.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks from file: %w", err)
	}

	// 如果文件为空，返回空任务列表
	if len(data) == 0 {
		return []*PersistTask{}, nil
	}

	// 解析JSON数据
	var tasks []*PersistTask
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	return tasks, nil
}

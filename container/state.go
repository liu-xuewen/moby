package container // import "github.com/docker/docker/container"

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	units "github.com/docker/go-units"
)

// State holds the current container state, and has methods to get and
// set the state. Container has an embed, which allows all of the
// functions defined against State to run against Container.
//
// State保存当前容器状态，并具有获取和设置状态的方法。
// CONTAINER有一个嵌入式，它允许针对State定义的所有函数针对Container运行。
type State struct {
	sync.Mutex
	// Note that `Running` and `Paused` are not mutually exclusive:
	// When pausing a container (on Linux), the freezer cgroup is used to suspend
	// all processes in the container. Freezing the process requires the process to
	// be running. As a result, paused containers are both `Running` _and_ `Paused`.
	//
	// 需要注意的是，`Running`和`Paused`并不是互斥的：在Linux上暂停容器时，会使用Freezer cgroup来暂停容器中的所有进程。
	// 冻结该进程需要该进程处于运行状态。
	// 因此，暂停的容器都是`Running`_和_`Paused`。
	Running           bool
	Paused            bool
	Restarting        bool
	OOMKilled         bool
	RemovalInProgress bool // Not need for this to be persistent on disk. // 不需要将其持久保存在磁盘上。
	Dead              bool
	Pid               int
	ExitCodeValue     int    `json:"ExitCode"`
	ErrorMsg          string `json:"Error"` // contains last known error during container start, stop, or remove
	StartedAt         time.Time
	FinishedAt        time.Time
	Health            *Health

	waitStop   chan struct{}
	waitRemove chan struct{}
}

// StateStatus is used to return container wait results.
// Implements exec.ExitCode interface.
// This type is needed as State include a sync.Mutex field which make
// copying it unsafe.
//
// StateStatus用于返回容器等待结果。
// 实现exec.ExitCode接口。
// 此类型是必需的，因为State包含使复制不安全的sync.Mutex字段。
type StateStatus struct {
	exitCode int
	err      error
}

// ExitCode returns current exitcode for the state.
func (s StateStatus) ExitCode() int {
	return s.exitCode
}

// Err returns current error for the state. Returns nil if the container had
// exited on its own.
func (s StateStatus) Err() error {
	return s.err
}

// NewState creates a default state object with a fresh channel for state changes.
func NewState() *State {
	return &State{
		waitStop:   make(chan struct{}),
		waitRemove: make(chan struct{}),
	}
}

// String returns a human-readable description of the state
func (s *State) String() string {
	if s.Running {
		if s.Paused {
			return fmt.Sprintf("Up %s (Paused)", units.HumanDuration(time.Now().UTC().Sub(s.StartedAt)))
		}
		if s.Restarting {
			return fmt.Sprintf("Restarting (%d) %s ago", s.ExitCodeValue, units.HumanDuration(time.Now().UTC().Sub(s.FinishedAt)))
		}

		if h := s.Health; h != nil {
			return fmt.Sprintf("Up %s (%s)", units.HumanDuration(time.Now().UTC().Sub(s.StartedAt)), h.String())
		}

		return fmt.Sprintf("Up %s", units.HumanDuration(time.Now().UTC().Sub(s.StartedAt)))
	}

	if s.RemovalInProgress {
		return "Removal In Progress"
	}

	if s.Dead {
		return "Dead"
	}

	if s.StartedAt.IsZero() {
		return "Created"
	}

	if s.FinishedAt.IsZero() {
		return ""
	}

	return fmt.Sprintf("Exited (%d) %s ago", s.ExitCodeValue, units.HumanDuration(time.Now().UTC().Sub(s.FinishedAt)))
}

// IsValidHealthString checks if the provided string is a valid container health status or not.
// IsValidHealthString检查提供的字符串是否为有效的容器健康状态。
func IsValidHealthString(s string) bool {
	return s == types.Starting ||
		s == types.Healthy ||
		s == types.Unhealthy ||
		s == types.NoHealthcheck
}

// StateString returns a single string to describe state
// StateString返回描述状态的单个字符串
func (s *State) StateString() string {
	if s.Running {
		if s.Paused {
			return "paused"
		}
		if s.Restarting {
			return "restarting"
		}
		return "running"
	}

	if s.RemovalInProgress {
		return "removing"
	}

	if s.Dead {
		return "dead"
	}

	if s.StartedAt.IsZero() {
		return "created"
	}

	return "exited"
}

// IsValidStateString checks if the provided string is a valid container state or not.
// IsValidStateString检查提供的字符串是否为有效的容器状态。
func IsValidStateString(s string) bool {
	if s != "paused" &&
		s != "restarting" &&
		s != "removing" &&
		s != "running" &&
		s != "dead" &&
		s != "created" &&
		s != "exited" {
		return false
	}
	return true
}

// WaitCondition is an enum type for different states to wait for.
// WaitCondition是不同州等待的枚举类型。
type WaitCondition int

// Possible WaitCondition Values.
//
// WaitConditionNotRunning (default) is used to wait for any of the non-running
// states: "created", "exited", "dead", "removing", or "removed".
//
// WaitConditionNextExit is used to wait for the next time the state changes
// to a non-running state. If the state is currently "created" or "exited",
// this would cause Wait() to block until either the container runs and exits
// or is removed.
//
// WaitConditionRemoved is used to wait for the container to be removed.
//
// 可能的WaitCondition值。
//
// WaitConditionNotRunning(默认值)用于等待任何非运行状态：“Created”、“Exted”、“Dead”、“Removing”或“Remove”。
// WaitConditionNextExit用于等待下次状态变为非运行状态。
// 如果状态当前为“CREATED”或“EXITED”，这将导致WAIT()阻塞，直到容器运行并退出或被删除。
//
// WaitConditionRemoved用于等待容器被移除。
const (
	WaitConditionNotRunning WaitCondition = iota
	WaitConditionNextExit
	WaitConditionRemoved
)

// Wait waits until the container is in a certain state indicated by the given
// condition. A context must be used for cancelling the request, controlling
// timeouts, and avoiding goroutine leaks. Wait must be called without holding
// the state lock. Returns a channel from which the caller will receive the
// result. If the container exited on its own, the result's Err() method will
// be nil and its ExitCode() method will return the container's exit code,
// otherwise, the results Err() method will return an error indicating why the
// wait operation failed.
//
// 等待，直到容器处于给定条件指示的特定状态。
// 必须使用上下文来取消请求、控制超时和避免Goroutine泄漏。
// 必须在不持有状态锁的情况下调用Wait。
// 返回调用方将从中接收结果的通道。
// 如果容器自己退出，则结果的err()方法将为空，其ExitCode()方法将返回容器的退出代码，否则，结果err()方法将返回错误，指示等待操作失败的原因。
//
func (s *State) Wait(ctx context.Context, condition WaitCondition) <-chan StateStatus {
	s.Lock()
	defer s.Unlock()

	if condition == WaitConditionNotRunning && !s.Running {
		// Buffer so we can put it in the channel now.
		// 缓冲器，这样我们现在就可以把它放进频道了。
		resultC := make(chan StateStatus, 1)

		// Send the current status.
		resultC <- StateStatus{
			exitCode: s.ExitCode(),
			err:      s.Err(),
		}

		return resultC
	}

	// If we are waiting only for removal, the waitStop channel should
	// remain nil and block forever.
	// 如果我们只等待删除，waitStop通道应该永远保持为空并阻塞。
	var waitStop chan struct{}
	if condition < WaitConditionRemoved {
		waitStop = s.waitStop
	}

	// Always wait for removal, just in case the container gets removed
	// while it is still in a "created" state, in which case it is never
	// actually stopped.
	//
	// 始终等待移除，以防容器在仍处于“已创建”状态时被移除，在这种情况下，它实际上永远不会停止。
	waitRemove := s.waitRemove

	resultC := make(chan StateStatus)

	go func() {
		select {
		case <-ctx.Done():
			// Context timeout or cancellation.
			resultC <- StateStatus{
				exitCode: -1,
				err:      ctx.Err(),
			}
			return
		case <-waitStop:
		case <-waitRemove:
		}

		s.Lock()
		result := StateStatus{
			exitCode: s.ExitCode(),
			err:      s.Err(),
		}
		s.Unlock()

		resultC <- result
	}()

	return resultC
}

// IsRunning returns whether the running flag is set. Used by Container to check whether a container is running.
func (s *State) IsRunning() bool {
	s.Lock()
	res := s.Running
	s.Unlock()
	return res
}

// GetPID holds the process id of a container.
func (s *State) GetPID() int {
	s.Lock()
	res := s.Pid
	s.Unlock()
	return res
}

// ExitCode returns current exitcode for the state. Take lock before if state
// may be shared.
func (s *State) ExitCode() int {
	return s.ExitCodeValue
}

// SetExitCode sets current exitcode for the state. Take lock before if state
// may be shared.
func (s *State) SetExitCode(ec int) {
	s.ExitCodeValue = ec
}

// SetRunning sets the state of the container to "running".
func (s *State) SetRunning(pid int, initial bool) {
	s.ErrorMsg = ""
	s.Paused = false
	s.Running = true
	s.Restarting = false
	if initial {
		s.Paused = false
	}
	s.ExitCodeValue = 0
	s.Pid = pid
	if initial {
		s.StartedAt = time.Now().UTC()
	}
}

// SetStopped sets the container state to "stopped" without locking.
func (s *State) SetStopped(exitStatus *ExitStatus) {
	s.Running = false
	s.Paused = false
	s.Restarting = false
	s.Pid = 0
	if exitStatus.ExitedAt.IsZero() {
		s.FinishedAt = time.Now().UTC()
	} else {
		s.FinishedAt = exitStatus.ExitedAt
	}
	s.ExitCodeValue = exitStatus.ExitCode
	s.OOMKilled = exitStatus.OOMKilled
	close(s.waitStop) // fire waiters for stop
	s.waitStop = make(chan struct{})
}

// SetRestarting sets the container state to "restarting" without locking.
// It also sets the container PID to 0.
func (s *State) SetRestarting(exitStatus *ExitStatus) {
	// we should consider the container running when it is restarting because of
	// all the checks in docker around rm/stop/etc
	s.Running = true
	s.Restarting = true
	s.Paused = false
	s.Pid = 0
	s.FinishedAt = time.Now().UTC()
	s.ExitCodeValue = exitStatus.ExitCode
	s.OOMKilled = exitStatus.OOMKilled
	close(s.waitStop) // fire waiters for stop
	s.waitStop = make(chan struct{})
}

// SetError sets the container's error state. This is useful when we want to
// know the error that occurred when container transits to another state
// when inspecting it
func (s *State) SetError(err error) {
	s.ErrorMsg = ""
	if err != nil {
		s.ErrorMsg = err.Error()
	}
}

// IsPaused returns whether the container is paused or not.
func (s *State) IsPaused() bool {
	s.Lock()
	res := s.Paused
	s.Unlock()
	return res
}

// IsRestarting returns whether the container is restarting or not.
func (s *State) IsRestarting() bool {
	s.Lock()
	res := s.Restarting
	s.Unlock()
	return res
}

// SetRemovalInProgress sets the container state as being removed.
// It returns true if the container was already in that state.
func (s *State) SetRemovalInProgress() bool {
	s.Lock()
	defer s.Unlock()
	if s.RemovalInProgress {
		return true
	}
	s.RemovalInProgress = true
	return false
}

// ResetRemovalInProgress makes the RemovalInProgress state to false.
func (s *State) ResetRemovalInProgress() {
	s.Lock()
	s.RemovalInProgress = false
	s.Unlock()
}

// IsRemovalInProgress returns whether the RemovalInProgress flag is set.
// Used by Container to check whether a container is being removed.
func (s *State) IsRemovalInProgress() bool {
	s.Lock()
	res := s.RemovalInProgress
	s.Unlock()
	return res
}

// IsDead returns whether the Dead flag is set. Used by Container to check whether a container is dead.
func (s *State) IsDead() bool {
	s.Lock()
	res := s.Dead
	s.Unlock()
	return res
}

// SetRemoved assumes this container is already in the "dead" state and
// closes the internal waitRemove channel to unblock callers waiting for a
// container to be removed.
func (s *State) SetRemoved() {
	s.SetRemovalError(nil)
}

// SetRemovalError is to be called in case a container remove failed.
// It sets an error and closes the internal waitRemove channel to unblock
// callers waiting for the container to be removed.
func (s *State) SetRemovalError(err error) {
	s.SetError(err)
	s.Lock()
	close(s.waitRemove) // Unblock those waiting on remove.
	// Recreate the channel so next ContainerWait will work
	s.waitRemove = make(chan struct{})
	s.Unlock()
}

// Err returns an error if there is one.
func (s *State) Err() error {
	if s.ErrorMsg != "" {
		return errors.New(s.ErrorMsg)
	}
	return nil
}

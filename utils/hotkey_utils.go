package utils

import (
	"fmt"
	goruntime "runtime"
	"sync"
	"time"

	"github.com/gonutz/w32/v2"
)

const (
	WM_HOTKEY = 0x0312

	ModAlt   = 0x0001
	ModCtrl  = 0x0002
	ModShift = 0x0004
	ModWin   = 0x0008

	VK_Q = 0x51
	VK_L = 0x4C
)

var (
	procRegisterHotKey   = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey = user32.NewProc("UnregisterHotKey")

	hotkeyMu       sync.RWMutex
	hotkeyRegistry       = make(map[int32]HotkeyEntry)
	hotkeyNextID   int32 = 1

	hotkeyListenerTID uint32

	registerCh   chan hotkeyRegReq
	unregisterCh chan int32
	stopCh       chan struct{}
)

type HotkeyCallback func()

type HotkeyEntry struct {
	Modifiers uint32
	VK        uint32
	Callback  HotkeyCallback
}

type hotkeyRegReq struct {
	entry    HotkeyEntry
	resultCh chan hotkeyRegResult
}

type hotkeyRegResult struct {
	id  int32
	err error
}

func RegisterHotkey(modifiers uint32, vk uint32, callback HotkeyCallback) (int32, error) {
	if registerCh == nil {
		return 0, fmt.Errorf("热键监听器未启动，请先调用 StartHotkeyListener")
	}

	resultCh := make(chan hotkeyRegResult, 1)
	registerCh <- hotkeyRegReq{
		entry:    HotkeyEntry{Modifiers: modifiers, VK: vk, Callback: callback},
		resultCh: resultCh,
	}
	result := <-resultCh
	return result.id, result.err
}

func UnregisterHotkey(id int32) error {
	if unregisterCh == nil {
		return fmt.Errorf("热键监听器未启动，请先调用 StartHotkeyListener")
	}

	unregisterCh <- id
	return nil
}

func UnregisterAllHotkeys() {
	hotkeyMu.Lock()
	defer hotkeyMu.Unlock()

	for id := range hotkeyRegistry {
		procUnregisterHotKey.Call(0, uintptr(id))
		delete(hotkeyRegistry, id)
	}
}

// doRegisterHotkey 在监听线程上实际调用 Windows API 注册热键
func doRegisterHotkey(modifiers uint32, vk uint32, callback HotkeyCallback) (int32, error) {
	hotkeyMu.Lock()
	defer hotkeyMu.Unlock()

	id := hotkeyNextID
	hotkeyNextID++

	ret, _, _ := procRegisterHotKey.Call(0, uintptr(id), uintptr(modifiers), uintptr(vk))
	if ret == 0 {
		hotkeyNextID--
		return 0, fmt.Errorf("注册热键失败: modifiers=0x%X, vk=0x%X", modifiers, vk)
	}

	hotkeyRegistry[id] = HotkeyEntry{
		Modifiers: modifiers,
		VK:        vk,
		Callback:  callback,
	}
	return id, nil
}

// doUnregisterHotkey 在监听线程上实际调用 Windows API 注销热键
func doUnregisterHotkey(id int32) {
	hotkeyMu.Lock()
	defer hotkeyMu.Unlock()

	if _, ok := hotkeyRegistry[id]; ok {
		procUnregisterHotKey.Call(0, uintptr(id))
		delete(hotkeyRegistry, id)
	}
}

func StartHotkeyListener(entries ...HotkeyEntry) {
	registerCh = make(chan hotkeyRegReq)
	unregisterCh = make(chan int32)
	stopCh = make(chan struct{})

	go func() {
		goruntime.LockOSThread()

		for _, e := range entries {
			doRegisterHotkey(e.Modifiers, e.VK, e.Callback)
		}

		hotkeyListenerTID = getCurrentThreadId()

		var msg w32.MSG
		for {
			select {
			case req := <-registerCh:
				id, err := doRegisterHotkey(req.entry.Modifiers, req.entry.VK, req.entry.Callback)
				req.resultCh <- hotkeyRegResult{id: id, err: err}

			case id := <-unregisterCh:
				doUnregisterHotkey(id)

			case <-stopCh:
				UnregisterAllHotkeys()
				return

			default:
				if w32.PeekMessage(&msg, 0, 0, 0, w32.PM_REMOVE) {
					if msg.Message == WM_HOTKEY {
						hotkeyMu.RLock()
						entry, ok := hotkeyRegistry[int32(msg.WParam)]
						hotkeyMu.RUnlock()
						if ok && entry.Callback != nil {
							entry.Callback()
						}
					}
					w32.TranslateMessage(&msg)
					w32.DispatchMessage(&msg)
				} else {
					time.Sleep(1 * time.Millisecond)
				}
			}
		}
	}()
}

func StopHotkeyListener() {
	if stopCh != nil {
		close(stopCh)
	}
}

func BuildModifiers(ctrl, alt, shift, win bool) uint32 {
	var mods uint32
	if ctrl {
		mods |= ModCtrl
	}
	if alt {
		mods |= ModAlt
	}
	if shift {
		mods |= ModShift
	}
	if win {
		mods |= ModWin
	}
	return mods
}

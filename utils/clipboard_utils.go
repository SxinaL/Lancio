package utils

import (
	"fmt"
	"unicode/utf16"
	"unsafe"

	"github.com/gonutz/w32/v2"
)

// ReadClipboardText 读取剪贴板中的 Unicode 文本
func ReadClipboardText() string {
	if !w32.OpenClipboard(0) {
		return ""
	}
	defer w32.CloseClipboard()

	h := w32.GetClipboardData(w32.CF_UNICODETEXT)
	if h == 0 {
		return ""
	}

	p := w32.GlobalLock(w32.HGLOBAL(h))
	if p == nil {
		return ""
	}
	defer w32.GlobalUnlock(w32.HGLOBAL(h))

	ptr := (*uint16)(p)
	var chars []uint16
	for {
		c := *ptr
		if c == 0 {
			break
		}
		chars = append(chars, c)
		ptr = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + unsafe.Sizeof(uint16(0))))
	}
	return string(utf16.Decode(chars))
}

// writeClipboardText 将文本写入剪贴板
func writeClipboardText(s string) error {
	if !w32.OpenClipboard(0) {
		return fmt.Errorf("打开剪贴板失败")
	}
	defer w32.CloseClipboard()
	w32.EmptyClipboard()

	u16 := utf16.Encode([]rune(s))
	u16 = append(u16, 0) // 以 0 结尾

	h := w32.GlobalAlloc(w32.GMEM_MOVEABLE, uint32(len(u16)*2))
	if h == 0 {
		return fmt.Errorf("分配剪贴板内存失败")
	}

	p := w32.GlobalLock(h)
	if p != nil {
		dst := unsafe.Slice((*uint16)(p), len(u16))
		copy(dst, u16)
		w32.GlobalUnlock(h)
	}

	ret := w32.SetClipboardData(w32.CF_UNICODETEXT, w32.HANDLE(h))
	if ret == 0 {
		return fmt.Errorf("写入剪贴板失败")
	}
	return nil
}

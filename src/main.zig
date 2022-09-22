pub const UNICODE = true;

const std = @import("std");
const WINAPI = std.os.windows.WINAPI;
const win32 = struct {
    usingnamespace @import("win32").zig;
    usingnamespace @import("win32").foundation;
    usingnamespace @import("win32").system.system_services;
    usingnamespace @import("win32").ui.windows_and_messaging;
    usingnamespace @import("win32").graphics.gdi;
    usingnamespace @import("win32").foundation;
};

const HWND = win32.HWND;
const WPARAM = win32.WPARAM;
const LPARAM = win32.LPARAM;
const LRESULT = win32.LRESULT;
const HINSTANCE = win32.HINSTANCE;
const L = win32.L;
const GetLastError = std.os.windows.kernel32.GetLastError;

const window_class_str = L("tmenu_class");
const window_title = L("tmenu title");

fn wndProc(hwnd: HWND, msg: c_uint, wparam: WPARAM, lparam: LPARAM) callconv(WINAPI) LRESULT {
    switch (msg) {
        win32.WM_DESTROY => {
            win32.PostQuitMessage(win32.WM_QUIT);
        },
        win32.WM_PAINT => {},
        else => {
            return win32.DefWindowProc(hwnd, msg, wparam, lparam);
        },
    }
    return 0;
}

fn makeWindowClassEx(hInstance: HINSTANCE) void {
    var wnd_class_ex = std.mem.zeroes(win32.WNDCLASSEX);
    wnd_class_ex.cbSize = @sizeOf(win32.WNDCLASSEX);
    wnd_class_ex.style = win32.WNDCLASS_STYLES.initFlags(.{
        .HREDRAW = 1,
        .VREDRAW = 1,
    });
    wnd_class_ex.lpfnWndProc = wndProc;
    wnd_class_ex.hInstance = hInstance;
    wnd_class_ex.hbrBackground = win32.GetStockObject(win32.WHITE_BRUSH);
    wnd_class_ex.lpszClassName = window_class_str;
    wnd_class_ex.hIconSm = win32.LoadIcon(null, win32.IDI_APPLICATION);
    wnd_class_ex.hCursor = win32.LoadCursor(null, win32.IDC_ARROW);

    _ = win32.RegisterClassEx(&wnd_class_ex);
}

pub export fn wWinMain(hInstance: HINSTANCE, hPrevInstance: ?HINSTANCE, pCmdLine: [*:0]u16, nCmdShow: u32) callconv(WINAPI) c_int {
    _ = hPrevInstance;
    _ = pCmdLine;

    var msg: win32.MSG = undefined;
    makeWindowClassEx(hInstance);

    var hwnd = win32.CreateWindowExW(
        win32.WS_EX_OVERLAPPEDWINDOW,
        window_class_str,
        window_title,
        win32.WS_OVERLAPPEDWINDOW,
        100,
        120,
        600,
        100,
        null,
        null,
        hInstance,
        null,
    );

    _ = win32.ShowWindow(hwnd, @intToEnum(win32.SHOW_WINDOW_CMD, nCmdShow));
    _ = win32.UpdateWindow(hwnd);

    std.log.info("{}", .{GetLastError()});

    while (win32.GetMessage(&msg, null, 0, 0) != 0) {
        _ = win32.TranslateMessage(&msg);
        _ = win32.DispatchMessage(&msg);
    }

    return 0;
}

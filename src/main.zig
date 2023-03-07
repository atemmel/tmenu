const std = @import("std");
const c = @cImport(@cInclude("SDL2/SDL.h"));
const assert = std.debug.assert;

pub fn main() !void {
    assert(c.SDL_Init(c.SDL_INIT_EVERYTHING) == 0);
    defer c.SDL_Quit();

    var window = c.SDL_CreateWindow("Test", 200, 100, 800, 600, 0);
    assert(window != null);
    defer c.SDL_DestroyWindow(window);

    var renderer = c.SDL_CreateRenderer(window, 0, c.SDL_RENDERER_ACCELERATED);
    assert(renderer != null);
    defer c.SDL_DestroyRenderer(renderer);

    loop: while (true) {
        var event: c.SDL_Event = undefined;
        while (c.SDL_PollEvent(&event) > 0) {
            switch (event.type) {
                c.SDL_KEYDOWN => {
                    if (event.key.keysym.sym == c.SDLK_ESCAPE) {
                        break :loop;
                    }
                },
                c.SDL_QUIT => {
                    break :loop;
                },
                else => {},
            }
        }

        const rect = c.SDL_Rect{
            .x = 20,
            .y = 20,
            .w = 100,
            .h = 100,
        };

        _ = c.SDL_SetRenderDrawColor(renderer, 0, 0, 0, 255);
        _ = c.SDL_RenderClear(renderer);
        _ = c.SDL_SetRenderDrawColor(renderer, 255, 0, 255, 255);
        _ = c.SDL_RenderFillRect(renderer, &rect);
        _ = c.SDL_RenderPresent(renderer);
        std.time.sleep(std.time.ns_per_ms * 16);
    }
}

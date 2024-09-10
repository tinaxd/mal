const std = @import("std");
const re = @cImport(@cInclude("../lib/regez/regez.h"));

pub fn READ(s: []u8) []u8 {
    return s;
}

fn read_form(s: []u8) []u8 {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();

    var slice = try allocator.alignedAlloc(u8, re.)
}
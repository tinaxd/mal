const std = @import("std");

pub fn main() !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    const stdin_file = std.io.getStdIn().reader();
    var br = std.io.bufferedReader(stdin_file);
    const stdin = br.reader();

    while (true) {
        try stdout.print("user> ", .{});
        try bw.flush();
        var buf: [1000]u8 = undefined;
        const result = try stdin.readUntilDelimiterOrEof(&buf, '\n');
        if (result) |line| {
            const pr = rep(line);
            try stdout.print("{s}\n", .{pr});
            try bw.flush();
        } else {
            break;
        }
    }
}

fn READ(s: []u8) []u8 {
    return s;
}

fn EVAL(s: []u8) []u8 {
    return s;
}

fn PRINT(s: []u8) []u8 {
    return s;
}

fn rep(line: []u8) []u8 {
    const r = READ(line);
    const e = EVAL(r);
    return PRINT(e);
}

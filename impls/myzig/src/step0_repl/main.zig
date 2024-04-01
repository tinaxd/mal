const std = @import("std");

pub fn read(input: []const u8) []const u8 {
    return input;
}

pub fn eval(input: []const u8) []const u8 {
    return input;
}

pub fn print(input: []const u8) []const u8 {
    return input;
}

pub fn repl_loop() !void {
    var general_purpose_allocator = std.heap.GeneralPurposeAllocator(.{}){};
    var gpa = general_purpose_allocator.allocator();

    const stdout = std.io.getStdOut().writer();
    const stdin = std.io.getStdIn().reader();
    while (true) {
        _ = try stdout.write("user> ");
        const bare_line = stdin.readUntilDelimiterAlloc(gpa, '\n', 8192) catch |err| {
            if (err == error.EndOfStream) {
                return;
            } else {
                return err;
            }
        };
        defer gpa.free(bare_line);

        const line = std.mem.trim(u8, bare_line, "\n");
        if (line.len == 0) {
            continue;
        }

        const step1 = read(line);
        const step2 = eval(step1);
        const step3 = print(step2);
        _ = try stdout.write(step3);
        _ = try stdout.write("\n");
    }
}

pub fn main() !void {
    try repl_loop();
}

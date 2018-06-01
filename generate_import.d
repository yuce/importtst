/*
Copyright 2017 Pilosa Corp.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
contributors may be used to endorse or promote products derived
from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH
DAMAGE.
*/


// Download a D compiler from: https://dlang.org/download.html
// LDC is recommended, since it produces more performant binaries.
// Compiling with LDC: ldc2 -release -O generate_import.d
// Compiling with dmd: dmd -release -O generate_import.d


import core.stdc.stdio: fprintf, stdout, FILE, fopen, fclose;
import std.exception: enforce;

int main(string[] args) {
    import std.stdio: writeln;
    import std.format: format;

    if (args.length < 2) {
        writeln(format("Usage: %s STRATEGY [PARAMS]", args[0]));
        return 1;
    }

    const strategy = args[1];
    const invalidStrategyMsg = "Invalid strategy: %s.\nShould be one of random, row-linear, reverse-row-linear, row-linear-gap, col-linear";

    BitWriter writer;
    try {
        switch (strategy) {
            case "random":
                writer = RandomBitWriter.fromArgv(args[2..$]);
                break;
            case "row-linear":
                writer = LinearRowBitWriter.fromArgv(args[2..$]);
                break;
            case "reverse-row-linear":
                writer = ReverseLinearRowBitWriter.fromArgv(args[2..$]);
                break;
            case "col-linear":
                writer = LinearColumnBitWriter.fromArgv(args[2..$]);
                break;
            case "row-linear-gap":
                writer = LinearGapRowBitWriter.fromArgv(args[2..$]);
                break;
            default:
                throw new Exception(format(invalidStrategyMsg, strategy));
        }
    }
    catch (Exception ex) {
        writeln(ex.msg);
        return 1;
    }

    writer.writeBits(stdout);

    return 0;
}

template IntParser(T) {
    T parseInt(string s) pure {
        import std.conv: to;
        import std.array: replace;
        return to!T(s.replace("_", ""));
    }
}

interface BitWriter {
    void writeBits(FILE *f);
}

class RandomBitWriter : BitWriter {
    import std.random: Random, uniform;

    static RandomBitWriter fromArgv(string[] args) {
        enforce(args.length >= 4, "Required params: RANDOM_SEED MAX_ROW_ID MAX_COL_ID BIT_COUNT");
        const seed = IntParser!int.parseInt(args[0]);
        const maxRowID = IntParser!long.parseInt(args[1]);
        const maxColID = IntParser!long.parseInt(args[2]);
        const bitCount = IntParser!long.parseInt(args[3]);
        return new RandomBitWriter(seed, maxRowID, maxColID, bitCount);
    }
    
    void writeBits(FILE *f) {
        foreach (i; 0 .. this.bitCount) {
            const rowID = uniform(0L, this.maxRowID, this.rnd);
            const colID = uniform(0L, this.maxColumnID, this.rnd);
            f.fprintf("%d,%d\n", rowID, colID);
        }
    }

    private this(int seed, long maxRowID, long maxColumnID, long bitCount) {
        this.rnd = Random(seed);
        this.maxRowID = maxRowID;
        this.maxColumnID = maxColumnID;
        this.bitCount = bitCount;
    }
    
    private Random rnd;
    private const long maxRowID;
    private const long maxColumnID;
    private const long bitCount;
}

class LinearRowBitWriter : BitWriter {
    static LinearRowBitWriter fromArgv(string[] args) {
        enforce(args.length >= 2, "Required params: ROW_ID_COUNT COLUMN_ID_COUNT");
        const rowIDCount = IntParser!ulong.parseInt(args[0]);
        const columnIDCount = IntParser!ulong.parseInt(args[1]);
        return new LinearRowBitWriter(rowIDCount, columnIDCount);
        
    }

    void writeBits(FILE *f) {
        foreach (rowID; 0 .. this.rowIDCount) {
            foreach (colID; 0 .. this.columnIDCount) {
                f.fprintf("%d,%d\n", rowID, colID);
            }
        }
    }

    private this(ulong rowIDCount, ulong columnIDCount) {
        this.rowIDCount = rowIDCount;
        this.columnIDCount = columnIDCount;
    }

    private ulong rowIDCount;
    private ulong columnIDCount;
}

class ReverseLinearRowBitWriter : BitWriter {
    static ReverseLinearRowBitWriter fromArgv(string[] args) {
        enforce(args.length >= 2, "Required params: ROW_ID_COUNT COLUMN_ID_COUNT");
        const rowIDCount = IntParser!long.parseInt(args[0]);
        const columnIDCount = IntParser!long.parseInt(args[1]);
        return new ReverseLinearRowBitWriter(rowIDCount, columnIDCount);
        
    }

    void writeBits(FILE *f) {
        foreach_reverse (rowID; 0 .. this.rowIDCount) {
            foreach_reverse (colID; 0 .. this.columnIDCount) {
                f.fprintf("%d,%d\n", rowID, colID);
            }
        }
    }

    private this(long rowIDCount, long columnIDCount) {
        this.rowIDCount = rowIDCount;
        this.columnIDCount = columnIDCount;
    }

    private long rowIDCount;
    private long columnIDCount;
}

class LinearColumnBitWriter : BitWriter {
    static LinearColumnBitWriter fromArgv(string[] args) {
        enforce(args.length >= 2, "Required params: ROW_ID_COUNT COLUMN_ID_COUNT");
        const rowIDCount = IntParser!ulong.parseInt(args[0]);
        const columnIDCount = IntParser!ulong.parseInt(args[1]);
        return new LinearColumnBitWriter(rowIDCount, columnIDCount);
        
    }

    void writeBits(FILE *f) {
        foreach (colID; 0 .. this.columnIDCount) {
            foreach (rowID; 0 .. this.rowIDCount) {
                f.fprintf("%d,%d\n", rowID, colID);
            }
        }
    }

    private this(ulong rowIDCount, ulong columnIDCount) {
        this.rowIDCount = rowIDCount;
        this.columnIDCount = columnIDCount;
    }

    private ulong rowIDCount;
    private ulong columnIDCount;
}

class LinearGapRowBitWriter : BitWriter {
    static LinearGapRowBitWriter fromArgv(string[] args) {
        enforce(args.length >= 2, "Required params: ROW_ID_COUNT COLUMN_ID_COUNT");
        const rowIDCount = IntParser!ulong.parseInt(args[0]);
        const columnIDCount = IntParser!ulong.parseInt(args[1]);
        return new LinearGapRowBitWriter(rowIDCount, columnIDCount);
        
    }

    void writeBits(FILE *f) {
        import std.range: iota;
        import std.algorithm: filter;
        foreach (rowID; 0 .. this.rowIDCount) {
            foreach (colID; this.columnIDCount.iota.filter!(n => n % 2)) {
                f.fprintf("%d,%d\n", rowID, colID);
            }
        }
    }

    private this(ulong rowIDCount, ulong columnIDCount) {
        this.rowIDCount = rowIDCount;
        this.columnIDCount = columnIDCount;
    }

    private ulong rowIDCount;
    private ulong columnIDCount;
}
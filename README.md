# Disassembly Unsplitter

A small tool to "unsplit" a Sonic disassembly. Inspired by [Hivebrain's ASMUnsplitter](https://github.com/cvghivebrain/ASMUnsplitter), which worked well but lacked manual file/folder exclusion.

## Usage
Drag `sonic.asm` onto `Unsplitter.exe`. After the process is done, a `sonic.unsplit.asm` and an `Unsplitter.log` file will be created.

## Ignore
If an `Unsplitter_ignore.txt` is found in the same path, it will be used as a list of files and folders to ignore.

Here's an example for the [Sonic 1 GitHub disassembly](https://github.com/sonicretro/s1disasm), which will effectively exclude everything from the unsplitting process except the files in `_incObj` and `_inc`:

```
; Use this file to define files or folders to exclude in the unsplitting process.
; Lines prefixed by semicolons are ignored.
_anim/
_maps/
sound/
Constants.asm
Variables.asm
Macros.asm
MacroSetup.asm
s1.sounddriver.asm
s1.sounddriver.ram.asm
```

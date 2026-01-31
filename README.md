# Disassembly Unsplitter

A small tool to "unsplit" a Sonic disassembly. Inspired by [Hivebrain's ASMUnsplitter](https://github.com/cvghivebrain/ASMUnsplitter), which worked well but lacked manual file/folder exclusion.

## Usage
Drag `sonic.asm` into `Unsplitter.exe`. After the process is done, a `sonic.asm.unsplit.asm` and an `Unsplitter.log` file will be created.

## Ignore
If an `Unsplitter_ignore.txt` is found in the same path, it will be used as a list of files and folders to ignore.

Here's an example for the Sgithub.com/sonicretro/s1disasm, which will effectively exclude everything from the unsplitting process except the files in _incObj and _inc:

```
_anim\
_maps\
sound\
s1.sounddriver.asm
s1.sounddriver.ram.asm
Variables.asm
Macros.asm
MacroSetup.asm
Constants.asm
```

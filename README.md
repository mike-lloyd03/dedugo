# Dedugo

## Image Duplicate Finder
#### De-duplicate in Go (get it?)

### Summary

This simple program evaluates two directories of images and finds images that are similar. It takes two arguments, the reference directory (a directory of images considered to be originals) and the evaluation directory (a directory of images that may contain duplicates of images in the reference directory).

Both directories are searched recursively for any compatible image formats (`.jpg`, `.png`, `.heic`). The path of images found to be similar is written to a file `duplicateImages.txt`.

### Usage
```
dedugo ./reference/image/directory ./evaluation/image/directory
```
Things will happen. Silicon will get hot. Fans will spin.

### To Do
[ ] Allow user to visually confirm if paired images are indeed duplicates or are actually just very similar
[ ] Convert this to use [Cobra](https://github.com/spf13/cobra)
[ ] A GUI would be neat

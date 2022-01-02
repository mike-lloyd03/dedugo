# Dedugo

## Image Duplicate Finder
#### De-duplicate in Go (get it?)

### Summary

This simple program evaluates two directories of images and finds images that are similar. It takes two arguments, the reference directory (a directory of images considered to be originals) and the evaluation directory (a directory of images that may contain duplicates of images in the reference directory).

Both directories are searched recursively for any compatible image formats (`.jpg`, `.png`, `.heic`).

### Usage
There are three phases to using this tool: finding duplicate images, confirming the detected duplicates, and deleting the confirmed duplicates.

#### Finding Duplicates
The first argument is the directory for the images to match against. You can consider this directory the "originals" which you don't want to be deleted. The second argument is the directory containing images which may be duplicated and which you may want to delete files from.
```
dedugo find-duplicates ./reference/image/directory ./evaluation/image/directory
```
Things will happen. Silicon will get hot. Fans will spin.

#### Checking Results
The `check-results` subcommand allows the user to visually confirm if detected duplicates are actually duplicate images. Because no algorithm is perfect, false positives are likely to happen. This will allow the user to confirm if a pair of images is a duplicate or not.
```
dedugo check-results
```

#### Deleting Duplicates
Once duplicate images are confirmed, they can be deleted in one fell swoop by running:
```
dedugo delete-duplicates
```

### To Do
- [x] Allow user to visually confirm if paired images are indeed duplicates or are actually just very similar
- [x] Convert this to use [Cobra](https://github.com/spf13/cobra)
- [x] A GUI would be neat
- [x] Add `delete` command to remove confirmed duplicates
- [ ] Make image loading faster
- [ ] Allow user to specify output filename for `find-duplicates` and input filename for `check-results`
- [ ] I probably need to incorporate the idea of similar image clusters rather than just image pairs.

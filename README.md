# Dedugo
__Duplicate Image Finder__

### Summary

This simple program evaluates two directories of images and finds images that are similar. It aims to be really fast and simple to use.

Both directories are searched recursively for any compatible image formats (`.jpg`, `.png`, `.heic`).

### Usage
There are three phases to using this tool: finding duplicate images, confirming the detected duplicates, and deleting the confirmed duplicates.

#### Finding Duplicates
The first argument is the directory for the images to match against. You can consider this directory the "originals" which you don't want to be deleted. The second argument is the directory containing images which may be duplicated and which you may want to delete files from.
```bash
dedugo find-duplicates ./reference/image/directory ./evaluation/image/directory
```
Things will happen. Silicon will get hot. Fans will spin.

#### Checking Results
The `check-results` subcommand allows the user to visually confirm if detected duplicates are actually duplicate images. Because no algorithm is perfect, false positives are likely to happen. This will allow the user to confirm if a pair of images is a duplicate or not.
```bash
dedugo check-results
```

#### Deleting Duplicates
Once duplicate images are confirmed, they can be deleted in one fell swoop by running:
```bash
dedugo delete-duplicates
```

### To Do
- [x] Allow user to visually confirm if paired images are indeed duplicates or are actually just very similar
- [x] Convert this to use [Cobra](https://github.com/spf13/cobra)
- [x] A GUI would be neat
- [x] Add `delete` command to remove confirmed duplicates
- [ ] Make image loading faster
- [x] Allow user to specify output filename for `find-duplicates` and input filename for `check-results`
- [ ] I probably need to incorporate the idea of similar image clusters rather than just image pairs.
- [ ] Write tests...

### Thanks
A special thanks to [Vitali Fedulov](https://github.com/vitali-fedulov) for writing the [Go package](https://github.com/vitali-fedulov/images) upon which this tool is built.

And also, thanks to [jdeng](https://github.com/jdeng) for writing the [heif decoder](https://github.com/jdeng/goheif) and [adrium](https://github.com/adrium/) for maintaining [a fork of it](https://github.com/adrium/goheif) that runs on my Linux install.
